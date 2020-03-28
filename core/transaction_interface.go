package core

import "go.mongodb.org/mongo-driver/bson"

type Transaction interface {
	Save(o interface{}, collectionName string)
	Get(o interface{}, collectionName string, query bson.D)
	
}
