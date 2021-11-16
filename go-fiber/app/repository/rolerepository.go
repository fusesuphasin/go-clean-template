package repository

import (
	"context"
	"fmt"

	"github.com/fusesuphasin/go-fiber/app/domain"
	"github.com/fusesuphasin/go-fiber/app/infrastructure"
	"github.com/fusesuphasin/go-fiber/app/utils/pagination"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleRepository struct {
	Ctx context.Context
}

func roleCollecttion() *mongo.Collection{
	coll, err := infrastructure.GetMongoDbCollection("myLearning","rule")
	if err != nil {
		fmt.Println("role 0:" , err)
	}
	return coll
}

func (rr *RoleRepository) Insert(Role *domain.Role) (role *domain.Role, err error) {
	coll := roleCollecttion()

	updateResult, err :=  coll.InsertOne(context.TODO(), Role)
	if err != nil {
		fmt.Println("role 1: ", err)
		//panic(err)
	}
	fmt.Println("Insert ID :", updateResult.InsertedID)

	// defer conn.Close()
	return Role, err
}

func (rr *RoleRepository) 	GetAll(page int, limit int) (role *[]domain.Role, err error) {
	coll := roleCollecttion()
	
	var Role *[]domain.Role
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		var c *fiber.Ctx
		return nil, c.JSON(err)
	}
	if page == 0 && limit == 0 {
		cursor.All(context.TODO(), &Role)
	} else {
		Role := pagination.Paginate(page, limit)
		_ = Role
	}
	cursor.All(context.TODO(), &Role)
	fmt.Println(Role)
	// defer conn.Close()
	return Role, err
}

func (rr *RoleRepository) Update(role_id string, Role *domain.Role) (role *domain.Role, err error) {
	coll := roleCollecttion()
	filter := bson.D{{"_id", role_id}}

	update := bson.D{{Key: "$set", Value: Role.Name}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	// conn.Role.UpdateOneID(role_id).SetName(Role.Name).Save(rr.Ctx)
	_ = result
	return Role, err
}

func (rr *RoleRepository) Delete(role_id string) (err error) {
	coll := roleCollecttion()
	role, err := rr.FindById(role_id)
	if err != nil {
		return err
	}
	_ = role
	filter := bson.D{{"_id", role_id }}
	result, err := coll.DeleteOne(context.TODO(),filter)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return err
}

func (rr *RoleRepository) CountByName(input string) (res int64) {
	var count int64
	/* //conn, err := infrastructure.Open()
	if err != nil {
		panic(err)
	}

	//conn.Model(&domain.Role{}).Where("name = ?", input).Count(&count) */
	return count
}

func (rr *RoleRepository) FindById(role_id string) (roleData *domain.Role, err error) {
	var role *domain.Role
	objectId, err := primitive.ObjectIDFromHex(role_id)
	
	if err != nil {fmt.Println(err)}

	coll := roleCollecttion()
	filter := bson.M{"_id": objectId}
	var result domain.Role

	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil{
		fmt.Println(err)
	}
	role = &result
	return role, err
}