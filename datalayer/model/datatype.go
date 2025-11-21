package model

import (
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Operator struct {
	ID   bson.ObjectID `json:"id"   bson:"id"`   // 用户 ID
	Name string        `json:"name" bson:"name"` // 用户名
}

type Duration time.Duration

func (d *Duration) UnmarshalText(raw []byte) error {
	du, err := time.ParseDuration(string(raw))
	if err != nil {
		return err
	}
	*d = Duration(du)

	return nil
}

func (d Duration) MarshalText() ([]byte, error) {
	du := time.Duration(d)
	s := du.String()

	return []byte(s), nil
}

func (d Duration) String() string {
	du := time.Duration(d)
	return du.String()
}

type HTTPHeader map[string]string

func (h HTTPHeader) Canonical() HTTPHeader {
	hm := make(HTTPHeader, len(h))
	for k, v := range h {
		key := http.CanonicalHeaderKey(strings.TrimSpace(k))
		hm[key] = v
	}

	return hm
}
