package model

type Checksum struct {
	MD5     string `json:"md5"     bson:"md5"`
	SHA1    string `json:"sha1"    bson:"sha1"`
	SHA256  string `json:"sha256"  bson:"sha256"`
	SHA512  string `json:"sha512"  bson:"sha512"`
	SHA3256 string `json:"sha3256" bson:"sha3256"`
}

func (chk Checksum) IsZero() bool {
	return chk.MD5 == "" &&
		chk.SHA1 == "" &&
		chk.SHA256 == "" &&
		chk.SHA512 == "" &&
		chk.SHA3256 == ""
}
