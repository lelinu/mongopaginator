package db

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

//Database object
type Database struct {
	Host             string
	DbName           string
	Username         string
	Password         string
	TimeoutInSeconds time.Duration
}

//NewDatabase this method is used to init a Database object
func NewDatabase(host string, dbName string, username string, password string, timeoutInSeconds time.Duration) *Database {
	return &Database{host, dbName, username, password, timeoutInSeconds}
}

//Init this method is used to init and connect to a database
func (m *Database) Init() (*mgo.Database, error) {

	//init dial info
	info := &mgo.DialInfo{
		Addrs:    []string{m.Host},
		Timeout:  m.TimeoutInSeconds * time.Second,
		Database: m.DbName,
		Username: m.Username,
		Password: m.Password,
	}

	// connect with dial info
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	db := session.DB(m.DbName)
	return db, nil
}
