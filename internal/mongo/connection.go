package mongo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

type (

	//BooksDB represents the controller for operating on the book ressource
	//Only mongo db will be implemented
	BooksDB struct {
		database *mgo.Database
		session  *mgo.Session
	}
)

//ErrInvalidDBprovided is returned when an uninitialized db is used to perform actions aigainst
var ErrInvalidDBprovided = errors.New("invalid DB provided")

//NewCollection creates a new session with timeout and url into mongo
func NewCollection(url string, timeout time.Duration) (*BooksDB, error) {

	//Default timeout for one session
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	//sess is a socket connection to mongo
	sess, err := mgo.DialWithTimeout(url, timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "mgo.DialWithTimeout: %s,%v", url, timeout)
	}

	//In the Monotonic consistency mode reads may not be entirely up-to-date,
	//but they will always see the history of changes moving forward,
	//the data read will be consistent across sequential queries in the same session,
	//and modifications made within the session will be observed in following queries (read-your-writes).
	//For more details read https://godoc.org/labix.org/v2/mgo#Session.SetMode
	sess.SetMode(mgo.Monotonic, true)

	db := BooksDB{
		database: sess.DB(""),
		session:  sess,
	}
	return &db, nil
}

//Execute executes mongo command
func (datB *BooksDB) Execute(ctx context.Context, name string, f func(*mgo.Collection) error) error {
	if datB == nil || datB.session == nil {
		return errors.Wrap(ErrInvalidDBprovided, "Requesed db not available")
	}
	return f(datB.database.C(name))
}

//ExecuteWithTimeout executes mongo command with timeout
func (datB *BooksDB) ExecuteWithTimeout(ctx context.Context, name string, f func(*mgo.Collection) error, timeout time.Duration) error {
	if datB == nil || datB.session == nil {
		return errors.Wrap(ErrInvalidDBprovided, "Requested db not available")
	}

	datB.session.SetSocketTimeout(timeout)

	return f(datB.database.C(name))
}

//Query provides a sting version of the v
func (datB *BooksDB) Query(v interface{}) string {
	json, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(json)
}

//Close close the db value used in mongo by a session
//It operates like mongo close method.
func (datB *BooksDB) Close() {
	datB.session.Close()
}

//Copy works like New but it returns the original session
func (datB *BooksDB) Copy() (*BooksDB, error) {

	sess := datB.session.Copy()

	//If no database was specified like mentioned in the docu then it will return
	//the default one, or the one that the connection was initated with.
	newDB := BooksDB{
		database: sess.DB(""),
		session:  sess,
	}
	return &newDB, nil
}
