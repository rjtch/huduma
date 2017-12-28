package mongo

import (
	"log"

	"gopkg.in/mgo.v2"
)

//Configure sets a new mongo session
func Configure() {
	session, err := mgo.Dial("localhost")

	if err != nil {
		log.Panic("Failed when creating a new session")
	}
	defer session.Clone()
}
