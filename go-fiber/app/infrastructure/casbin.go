package infrastructure

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Initmodel() model.Model {
	// Initialize the model from Go code.
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	return m
}

func CasbinLoad(driver *options.ClientOptions, dbName string) *casbin.Enforcer {
	
	a,err := mongodbadapter.NewAdapterWithClientOption(driver, dbName)
	// Or you can use NewAdapterWithCollectionName for custom collection name.
	if err != nil {
		/* panic(err) */
		fmt.Println("casnin 1: ", err)
	}
	m := Initmodel()
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		/* panic(err) */
		fmt.Println("casnin 2: ", err)
	}

	// Load the policy from DB.
	e.LoadPolicy()
	e.SavePolicy()
	return e
}