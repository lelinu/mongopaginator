# mongodpaginator

Mongodb paginator is a wrapper package to load records via mongodb with pagination

Requires Go 1.12 or newer.

## Usage

### Simple initialization

```go
import "github.com/lelinu/mongopaginator"

    var docs []TestDocument

    paginator := Init("OrderByField", 1, 1)
    collection := dbSession.C("testcollection")
    output := paginator.Paginate(collection, nil, &docs)

    returnedDocs := output.Records.(*[]TestDocument)
    fmt.Printf("Returned docs are %v", returnedDocs)
```
