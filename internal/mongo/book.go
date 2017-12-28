package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

type (

	//Book represents the structure of the ressource
	Book struct {
		ID      bson.ObjectId `json:"id" bson:"id"`
		Title   string        `json:"title" bson:"title"`
		Authors []string      `json:"authors"  bson:"authors"`
		Price   string        `json:"price" bson:"price"`
	}
)
