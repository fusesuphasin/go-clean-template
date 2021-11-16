package pagination

import (
	"context"
	"fmt"
	"time"

	"github.com/fusesuphasin/go-fiber/app/domain"
	"github.com/fusesuphasin/go-fiber/app/infrastructure"
	. "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

func Paginate(page int, limit int) []domain.Role{
	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}

	filter := bson.M{}
	var glimit int64 = int64(page)
    var gpage int64 = int64(limit)

	if page == 0 {
		page = 1
	}
	
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coll, err := infrastructure.GetMongoDbCollection("myLearning","phonebooks")
	if err != nil {
		/* panic(err) */
		fmt.Println("paging 1: ", err)
	}
	projection := bson.D{
		{"name", 1},
		/* {"qty", 1}, */
	}

	// Querying paginated data
    // Sort and select are optional
    // Multiple Sort chaining is also allowed
    // If you want to do some complex sort like sort by score(weight) for full text search fields you can do it easily
    // sortValue := bson.M{
    //		"$meta" : "textScore",
    //	}
    // aggPaginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("score", sortValue)...
    
	var produce []domain.Role
	paginatedData, err := New(coll).Context(ctx).Limit(glimit).Page(gpage)/* .Sort("price", -1) */.Select(projection).Filter(filter).Decode(&produce).Find()
	if err != nil {
		/* panic(err) */
		fmt.Println("paging 2: ", err)
	}
	// paginated data or paginatedData.Data will be nil because data is already decoded on through Decode function
	// pagination info can be accessed in  paginatedData.Pagination
	// print ProductList
	fmt.Printf("Normal Find Data: %+v\n", produce)

	// print pagination data
	fmt.Printf("Normal find pagination info: %+v\n", paginatedData.Pagination)
	return produce
}
