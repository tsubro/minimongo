package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type MetaData struct {
	CollectionName string
	data           map[string]interface{}
	SData          interface{}
}

func Parse(o interface{}, collectionName string, name string, val interface{}) ([]MetaData, error) {

	var res []MetaData
	m := MetaData{}
	m.CollectionName = collectionName
	m.data = map[string]interface{}{}

	if name != "" && val != nil {
		m.data[name] = val
	}

	fmt.Println(reflect.TypeOf(o).Kind())
	s := reflect.ValueOf(o).Elem()

	for i := 0; i < s.NumField(); i++ {

		f := s.Field(i)
		tf := s.Type().Field(i)

		v, err := parseTag(tf.Tag)
		if err != nil {
			m.data[tf.Name] = f.Interface()
			continue
		}

		iTag := v[ignoreTag]
		nTag := v[fieldNameTag]
		cTag := v[collectionNameTag]
		rTag := v[referenceKeyTag]

		if iTag != "" {
			continue
		}

		if nTag == "" {
			m.data[tf.Name] = f.Interface()
			continue
		}

		if cTag == "" || rTag == "" {
			m.data[nTag] = f.Interface()
			continue
		}

		fv := s.FieldByName(rTag).Interface()
		if f.Kind() == reflect.Slice || f.Kind() == reflect.Array {

			for i := 0; i < f.Len(); i++ {
				item := f.Index(i)
				if item.Kind() == reflect.Struct {
					md, err := Parse(reflect.Indirect(item).Addr().Interface(), cTag, (nTag + "_ref"), fv)
					if err != nil {
						return nil, err
					}
					res = append(res, md...)

				} else if item.Kind() == reflect.Map {
					return nil, errors.New("Odm Tags Not Supported For Maps")
				}
			}

		} else if f.Kind() == reflect.Struct {

			md, err := Parse(f.Addr().Interface(), cTag, (nTag + "_ref"), fv)
			if err != nil {
				return nil, err
			}
			res = append(res, md...)

		} else if f.Kind() == reflect.Map {
			return nil, errors.New("Odm Tags Not Supported For Maps")
		}
	}

	m.SData = convertMapToStruct(m.data)
	res = append(res, m)
	return res, nil
}

func parseTag(tag reflect.StructTag) (map[string]string, error) {
	if tag == "" {
		return nil, errors.New("Tag Not Found")
	}

	iTag := tag.Get(ignoreTag)
	nTag := tag.Get(fieldNameTag)
	cTag := tag.Get(collectionNameTag)
	rTag := tag.Get(referenceKeyTag)

	return map[string]string{
		ignoreTag:         iTag,
		fieldNameTag:      nTag,
		collectionNameTag: cTag,
		referenceKeyTag:   rTag,
	}, nil
}

func convertMapToStruct(d map[string]interface{}) interface{} {

	b, _ := json.Marshal(d)
	var a interface{}
	json.Unmarshal(b, &a)
	return a
}

func Unparse(o interface{}, results []map[string]interface{}) []interface{} {

	var r []interface{}

	var tagMap map[string]string
	s := reflect.ValueOf(o).Elem()
	for i := 0; i < s.NumField(); i++ {

		tf := s.Type().Field(i)
		n := tf.Name

		if tf.Tag != "" {
			nTag := tf.Tag.Get(fieldNameTag)

			if nTag != "" {
				n = nTag
			}
		}

		tagMap[n] = tf.Name
	}

	for _, result := range results {
		ele := o
		rEle := reflect.ValueOf(ele).Elem()

		for k, v := range result {

			f := rEle.FieldByName(tagMap[k])

			if f.Kind() == reflect.Int {
				f.SetInt(v.(int64))
			} else if f.Kind() == reflect.String {
				f.SetString(v.(string))
			}
		}
		r = append(r, ele)
	}
	return r
}
