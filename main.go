package main

import (
	"minimongo/core"
	// "minimongo/test"
)

type Employee struct {
	Id      int     `odm.field_name:"id"`
	Name    string  `odm.field_name:"name"`
	Vehicle map[string]Vehicle `odm.field_name:"emp_id" odm.collection_name:"odm-vehicle" odm.reference_key:"Id"`
}

type Vehicle struct {
	Id    int		`odm.field_name:"vehicle_id"`
	Name  string
	Model string
}

func main() {

	v := Vehicle{1, "Honda City", "Sedan"}
	v2 := Vehicle{2, "Honda Civic", "Sedan"}

	ma := map[string]Vehicle{"v1":v, "v2":v2}

	e := Employee{1, "Subro", ma}
	
	// test.Ref(&e)
	m := core.MongoTx{"mongodb://localhost:27017", "odm-test"}
	m.Save(&e, "odm-job")

	emp := Employee{}

	m.Get(&emp, "odj-job", nil)
}
