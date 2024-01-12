package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type NanoTime struct {
	time.Time
}

func (nt NanoTime) MarshalBSONValue() (bsontype bsontype.Type, data []byte, err error) {
	return bson.MarshalValue(nt.UnixNano())
}

func (nt *NanoTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var unixNano int64
	err := bson.UnmarshalValue(t, data, &unixNano)
	if err != nil {
		return err
	}
	*nt = NanoTime{time.Unix(0, unixNano).UTC()}
	return nil
}
