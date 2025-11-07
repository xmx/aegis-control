package repository

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"path"
	"strings"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FS interface {
	fs.ReadDirFS
	Repository[bson.ObjectID, model.FS, model.FSs]

	OpenContext(ctx context.Context, name string) (fs.File, error)
	List(ctx context.Context, name string) (model.FSs, error)
	Create(ctx context.Context, name string, rd io.Reader) (*model.FS, error)
	Upload(ctx context.Context, name string, rd io.Reader) (*model.FS, error)
	Mkdir(ctx context.Context, name string) error
	Remove(ctx context.Context, name string) error
}

func NewFS(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) FS {
	name := "fs"

	coll := db.Collection(name+".infos", opts...)
	repo := NewRepository[bson.ObjectID, model.FS, model.FSs](coll)
	opt := options.GridFSBucket().SetName(name)
	bucket := db.GridFSBucket(opt)

	return &fsRepo{
		Repository: repo,
		bucket:     bucket,
	}
}

type fsRepo struct {
	Repository[bson.ObjectID, model.FS, model.FSs]
	bucket *mongo.GridFSBucket
}

func (r *fsRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "full_path", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

// Open opens the named file.
// [fs.File.Close] must be called to release any associated resources.
//
// When Open returns an error, it should be of type *PathError
// with the Op field set to "open", the Path field set to name,
// and the Err field describing the problem.
//
// Open should reject attempts to open names that do not satisfy
// fs.ValidPath(name), returning a *fs.PathError with Err set to
// fs.ErrInvalid or fs.ErrNotExist.
func (r *fsRepo) Open(name string) (fs.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return r.OpenContext(ctx, name)
}

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
func (r *fsRepo) ReadDir(name string) ([]fs.DirEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fss, err := r.List(ctx, name)
	if err != nil {
		return nil, err
	}
	entries := make([]fs.DirEntry, 0, len(fss))
	for _, f := range fss {
		entries = append(entries, &fsDirEntry{file: f})
	}

	return entries, nil
}

func (r *fsRepo) OpenContext(ctx context.Context, name string) (fs.File, error) {
	fp := r.normalization(name)
	info, err := r.FindOne(ctx, bson.M{"full_path": fp})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fs.ErrNotExist
		}
		return nil, err
	}
	if info.IsDir() {
		return &fsFile{file: info}, nil
	}
	fid := info.FileID
	stm, err := r.bucket.OpenDownloadStream(ctx, fid)
	if err != nil {
		return nil, err
	}

	return &fsFile{file: info, stm: stm}, nil
}

func (r *fsRepo) List(ctx context.Context, name string) (model.FSs, error) {
	fp := r.normalization(name)
	info, err := r.FindOne(ctx, bson.M{"full_path": fp})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fs.ErrNotExist
		}
	} else if !info.IsDir() {
		pe := errors.New("not a directory")
		return nil, &fs.PathError{Op: "read", Path: fp, Err: pe}
	}
	filter := bson.M{"parent_path": fp}
	fss, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	fss.Sort()

	return fss, nil
}

func (r *fsRepo) Remove(ctx context.Context, name string) error {
	fp := r.normalization(name)
	filter := bson.M{"full_path": fp}
	dat, err := r.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fs.ErrNotExist
		}
		return err
	}
	if !dat.IsDir() {
		_, _ = r.DeleteOne(ctx, filter)
		err = r.bucket.Delete(ctx, dat.FileID)
		return err
	}

	// 如果是目录就检查目录下是否存在子目录或文件
	if cnt, _ := r.CountDocuments(ctx, bson.M{"parent_path": fp}); cnt != 0 {
		pe := errors.New("directory not empty")
		return &fs.PathError{Op: "remove", Path: fp, Err: pe}
	}
	_, err = r.DeleteOne(ctx, filter)

	return err
}

func (r *fsRepo) Create(ctx context.Context, name string, rd io.Reader) (*model.FS, error) {
	createdAt := time.Now()

	fp := r.normalization(name)
	parentDir := path.Dir(fp)
	basename := path.Base(fp)
	ext := strings.ToLower(path.Ext(basename))

	// 检查父目录是否存在，不存在就报错。
	// 如果是根目录就自动创建，其他目录不会自动创建，因为需要逐级检查。
	if parent, err := r.FindOne(ctx, bson.M{"full_path": parentDir}); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		} else if parentDir != "/" {
			return nil, fs.ErrNotExist
		}
		// 根目录不存在就自动创建
		root := &model.FS{
			FullPath:  "/",
			Name:      "/",
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}
		if _, err = r.InsertOne(ctx, root); err != nil {
			return nil, err
		}
	} else if !parent.IsDir() { // 检查父节点是否是目录
		pe := errors.New("not a directory")
		return nil, &fs.PathError{Op: "read", Path: fp, Err: pe}
	}

	// 先插入数据
	dat := &model.FS{
		FullPath:   fp,
		ParentPath: parentDir,
		Name:       basename,
		Extension:  ext,
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
	}
	ret, err := r.InsertOne(ctx, dat)
	if err != nil {
		if exc, ok := err.(mongo.WriteException); ok && exc.HasErrorCode(11000) {
			return nil, fs.ErrExist
		}
		return nil, err
	}

	id := ret.InsertedID.(bson.ObjectID)
	fileID, sum, cnt, err1 := r.upload(ctx, fp, rd)
	if err1 != nil {
		_, _ = r.DeleteOne(ctx, id)
		return nil, err1
	}

	updatedAt := time.Now()
	update := bson.M{"$set": bson.M{
		"size": cnt, "file_id": fileID, "checksum": sum, "updated_at": updatedAt,
	}}
	dat.ID, dat.FileID, dat.Checksum, dat.UpdatedAt = id, fileID, sum, updatedAt
	if _, err = r.UpdateByID(ctx, id, update); err == nil {
		return dat, nil
	}

	// 出错就删除数据
	_, _ = r.DeleteOne(ctx, id)
	_ = r.bucket.Delete(ctx, id)

	return nil, err
}

func (r *fsRepo) Upload(ctx context.Context, name string, rd io.Reader) (*model.FS, error) {
	updatedAt := time.Now()

	// 检查文件是否存在
	fp := r.normalization(name)
	dat, err := r.FindOne(ctx, bson.M{"full_path": fp})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fs.ErrNotExist
		}
		return nil, err
	} else if dat.IsDir() {
		pe := errors.New("bad file descriptor")
		return nil, &fs.PathError{Op: "write", Path: fp, Err: pe}
	}

	fileID, sum, cnt, err1 := r.upload(ctx, fp, rd)
	if err1 != nil {
		return nil, err1
	}
	update := bson.M{"$set": bson.M{
		"size": cnt, "file_id": fileID, "checksum": sum, "updated_at": updatedAt,
	}}
	if _, err = r.UpdateByID(ctx, dat.ID, update); err != nil {
		_ = r.bucket.Delete(ctx, fileID) // 出错就删除掉刚上传的文件。
		return nil, err
	}
	// 删除掉老的文件
	_ = r.bucket.Delete(ctx, dat.FileID)
	dat.FileID, dat.Checksum, dat.Size, dat.UpdatedAt = fileID, sum, cnt, updatedAt

	return dat, nil
}

// upload 上传文件流，并返回文件ID、文件哈希、文件大小。
func (r *fsRepo) upload(ctx context.Context, name string, rd io.Reader) (bson.ObjectID, model.Checksum, int64, error) {
	stm, err := r.bucket.OpenUploadStream(ctx, name)
	if err != nil {
		return bson.NilObjectID, model.Checksum{}, 0, err
	}
	defer stm.Close()

	hsw := model.NewHashWriter()
	cnt, err1 := io.Copy(io.MultiWriter(hsw, stm), rd)
	if err1 != nil {
		return bson.NilObjectID, model.Checksum{}, 0, err1
	}
	sum := hsw.Sum()
	fileID, _ := stm.FileID.(bson.ObjectID)

	return fileID, sum, cnt, nil
}

func (r *fsRepo) Mkdir(ctx context.Context, name string) error {
	now := time.Now()
	fp := r.normalization(name)
	parentDir := path.Dir(fp)
	basename := path.Base(fp)
	// 检查文件名是否存在
	_, err := r.FindOne(ctx, bson.M{"full_path": fp})
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return fs.ErrExist
	}

	// 父目录
	parent, err1 := r.FindOne(ctx, bson.M{"full_path": parentDir})
	if err1 != nil {
		if !errors.Is(err1, mongo.ErrNoDocuments) {
			return err1
		} else if parentDir != "/" {
			return fs.ErrNotExist
		}
		// 根目录不存在就自动创建
		root := &model.FS{
			FullPath:  "/",
			Name:      "/",
			CreatedAt: now,
			UpdatedAt: now,
		}
		if _, err = r.InsertOne(ctx, root); err != nil {
			return err
		}
	} else if !parent.IsDir() {
		pe := errors.New("not a directory")
		return &fs.PathError{Op: "mkdir", Path: fp, Err: pe}
	}

	// 创建
	dir := &model.FS{
		FullPath:   fp,
		ParentPath: parentDir,
		Name:       basename,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = r.InsertOne(ctx, dir)

	return err
}

func (r *fsRepo) normalization(name string) string {
	return path.Join("/", name)
}
