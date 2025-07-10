package {{.PKG}}

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound        = mongo.ErrNoDocuments
	ErrInvalidObjectId = errors.New("invalid objectId")
)
