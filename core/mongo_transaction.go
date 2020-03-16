package core

import (
	"context"
	"minimongo/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
)

type MongoTx struct {
	URL          string
	DatabaseName string
}

func (m MongoTx) Save(o interface{}, collectionName string) error {

	md, err := utils.Parse(o, collectionName, "", nil)
	if err != nil {
		return err
	}

	for _,v := range md {
		log.Info(v)
		m.insert(v.SData, v.CollectionName)
	}
	return nil
}

func (m MongoTx) Get(o interface{}, collectionName string, searchKeys ...interface{}) {

}

func (m MongoTx) Commit() {

}

func (m MongoTx) Rollback() {

}

func (m MongoTx) getMongoConnection() *mongo.Client {

	clientOptions := options.Client().ApplyURI(m.URL) //"mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to MongoDB!")
	return client
}

func (m MongoTx) insert(o interface{}, collectionName string) {

	collection := m.getMongoConnection().Database(m.DatabaseName).Collection(collectionName)

	insertResult, err := collection.InsertOne(context.TODO(), o)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Inserted a single document: ", insertResult.InsertedID)
}

func (m MongoTx) get(collectionName string, searchKeys ...interface{}) (map[string]interface{}, error) {
	collection := m.getMongoConnection().Database(m.DatabaseName).Collection(collectionName)

	var o map[string]interface{}

	err := collection.FindOne(context.TODO(), bson.D{{"jobid", jobId}}).Decode(&o)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("Fetched a single document succefully")
	return o, nil
}

// func update(job *models.RenderJob) {
// 	collection := getMongoConnection().Database("document-renderer").Collection("render-job")

// 	update := bson.M{
// 		"$set": bson.M{
// 			"state":  job.State,
// 			"inputs": job.Inputs,
// 		},
// 	}
// 	_, err := collection.UpdateOne(context.TODO(), bson.D{{"jobid", job.JobId}}, update)
// 	if err != nil {
// 		log.Error(err)
// 	}
// }
