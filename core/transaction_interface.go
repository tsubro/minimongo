package core

type Transaction interface {

	Save(o interface{}, collectionName string)
	Get(o interface{}, collectionName string)
	Commit()
	Rollback()
}