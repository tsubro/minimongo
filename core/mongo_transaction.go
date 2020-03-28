package core

import (
	"context"
	"fmt"
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

	for _, v := range md {
		log.Info(v)
		m.insert(v.SData, v.CollectionName)
	}
	return nil
}

func (m MongoTx) Get(o interface{}, collectionName string, searchKeys ...interface{}) {

	r, err := m.get(collectionName, searchKeys)
	if err != nil {
		log.Info("Error While Fetching Data ", err)
	}

	
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

func (m MongoTx) get(collectionName string, searchKeys ...interface{}) ([]*map[string]interface{}, error) {
	collection := m.getMongoConnection().Database(m.DatabaseName).Collection(collectionName)

	var query bson.D
	for k, v := range searchKeys {
		query = append(query, bson.E{fmt.Sprintf("%v", k), v})
	}

	cur, err := collection.Find(context.TODO(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var results []*map[string]interface{}
	for cur.Next(context.TODO()) {
    
		// create a value into which the single document can be decoded
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
	
		results = append(results, &elem)
	}

	log.Info("Fetched document succefully")
	return results, nil
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
