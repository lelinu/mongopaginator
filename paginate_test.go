package mongopaginator

import (
	"math/rand"
	"testing"
	"time"

	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/lelinu/mongopaginator/db"

	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
)

//constants
var TotalDocsNum = 20
var CollectionName = "test_collection"

//dbSession to be used throughout all test cases
var dbSession *mgo.Database

//TestDocument object
type TestDocument struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	Age  int           `bson:"age"`
}

//init function
func init() {

	var dbHost = "localhost:27017"
	var dbName = "testdb"
	var dbUsername = ""
	var dbPassword = ""
	var dbTimeout time.Duration = 5
	var err error

	//init database
	var database = db.NewDatabase(dbHost, dbName, dbUsername, dbPassword, dbTimeout)
	dbSession, err = database.Init()
	check(err)

	//get collection instance and generate test records
	collection := dbSession.C(CollectionName)
	generateTestRecords(collection)
}

//generateRandomChars will generate random alphabetical characters
func generateRandomChars(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//generateTestRecords will generate test records for testing purposes
func generateTestRecords(collection *mgo.Collection) {

	//get total number of documents
	totalDocuments, err := collection.Find(nil).Count()
	check(err)

	if totalDocuments >= TotalDocsNum {
		return
	}

	//create a waitgroup and guard
	var wg sync.WaitGroup
	guard := make(chan struct{}, 100)
	round := 10
	for i := totalDocuments; i < TotalDocsNum; i += round {
		guard <- struct{}{}
		wg.Add(1)
		go func(x int) {
			docs := make([]interface{}, round)
			for total := x + round; x < total && x < TotalDocsNum; x++ {
				docs = append(docs, TestDocument{ID: bson.NewObjectId(), Name: generateRandomChars(10), Age: rand.Intn(1000)})
			}
			b := collection.Bulk()
			b.Unordered()
			b.Insert(docs...)
			_, err = b.Run()
			check(err)
			<-guard
			wg.Done()
		}(i)
	}
	wg.Wait()
}

//removeTestRecords will remove the test records
func removeTestRecords(collection *mgo.Collection) {
	err := collection.DropCollection()
	check(err)
}

//TestPaginationWithOneRecord
func TestPaginationWithOneRecord(t *testing.T) {

	collection := dbSession.C(CollectionName)
	var docs []TestDocument

	// should return 1 record
	paginator := Init("", 1, 1)
	output := paginator.Paginate(collection, nil, &docs)

	returnedDocs := output.Records.(*[]TestDocument)
	assert.Equal(t, 1, len(*returnedDocs))
	assert.Equal(t, TotalDocsNum, output.TotalRecords)
	// should return 1 record
}

//TestPaginationWithTenRecords
func TestPaginationWithTenRecords(t *testing.T) {

	collection := dbSession.C(CollectionName)
	var docs []TestDocument

	// should return 10 records
	paginator := Init("", 1, 10)
	output := paginator.Paginate(collection, nil, &docs)

	returnedDocs := output.Records.(*[]TestDocument)
	assert.Equal(t, 10, len(*returnedDocs))
	assert.Equal(t, TotalDocsNum, output.TotalRecords)
	// should return 10 records
}

//TestTearDown will remove all the setup
func TestTearDown(t *testing.T) {

	//get collection instance and generate test records
	collection := dbSession.C(CollectionName)
	removeTestRecords(collection)
}

//check will be used to check the error
func check(err error) {
	if err != nil {
		panic(err)
	}
}
