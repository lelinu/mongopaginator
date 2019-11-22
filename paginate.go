package mongopaginator

import (
	"errors"
	"log"

	mgo "gopkg.in/mgo.v2"
)

//Pagination struct
type Pagination struct {
	OrderByList string
	Page        int
	PerPage     int
}

//Result struct
type Result struct {
	TotalRecords int         `json:"totalRecords"`
	Records      interface{} `json:"records"`
	Error        error       `json:"error"`
}

//Init will init the pagination
func Init(orderByList string, page int, perPage int) *Pagination {
	if perPage == 0 {
		log.Fatalf("Pagination: PerPage must be greater than 0")
	}

	//return pointer to paginator
	return &Pagination{
		OrderByList: orderByList,
		Page:        page,
		PerPage:     perPage,
	}
}

//Paginate this method will paginate the documents
func (p *Pagination) Paginate(collection *mgo.Collection, query, dataSource interface{}) *Result {

	done := make(chan bool, 1)

	var output Result
	var count int

	//if it's empty default it to primary key
	if len(p.OrderByList) == 0 {
		p.OrderByList = "_id"
	}

	//validations
	if p.Page <= 0 {
		output.Error = errors.New("Pagination: Page must be greater than zero")
		return &output
	}

	//get count in different thread
	go func() {
		count, _ = collection.Find(query).Count()
		done <- true
	}()

	//calculate skip
	skip := (p.Page - 1) * p.PerPage

	//build query
	err := collection.Find(query).Sort(p.OrderByList).Skip(skip).Limit(p.PerPage).All(dataSource)
	<-done

	if err != nil {
		output.Error = errors.New("Pagination: Page must be greater than zero")
		return &output
	}

	//set output parameters
	output.Records = dataSource
	output.TotalRecords = count

	return &output
}
