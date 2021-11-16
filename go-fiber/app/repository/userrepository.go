package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/fusesuphasin/go-fiber/app/domain"
	"github.com/fusesuphasin/go-fiber/app/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Ctx context.Context
}

func userCollecttion() *mongo.Collection{
	coll, err := infrastructure.GetMongoDbCollection("myLearning","users")
	if err != nil {
		panic(err)
	}
	return coll
}

func (ur *UserRepository) Insert(User *domain.User) (user *domain.User) {
	coll := userCollecttion()

	updateResult, err :=  coll.InsertOne(context.TODO(), User)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	fmt.Println("Insert ID :", updateResult.InsertedID)
	// defer conn.Close()
	return User
}

func (ur *UserRepository) CountByUsername(input string) (res int64) {
	var count int64

	coll := userCollecttion()
	
	filter := bson.D{{"username", input}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	// defer conn.Close()
	return count
}

func (ur *UserRepository) FindByUsername(input string) (res *domain.User) {
	var user *domain.User

	coll := userCollecttion()
	filter := bson.M{"username": input}
	
	coll.FindOne(context.TODO(), filter).Decode(&user)
	
	return user
}

func (ur *UserRepository) FindById(input uint64) (res *[]domain.User) {
	var user *[]domain.User

	coll := userCollecttion()
	filter := bson.M{"id": input}

	coll.FindOne(context.TODO(), filter).Decode(&user)
	// defer conn.Close()
	return user
}

func (ur *UserRepository) FindByIdWithRelation(input string) (res *domain.User) {
	var user *domain.User
	coll := userCollecttion()
	
	filter := bson.M{"username": input}
	coll.FindOne(context.TODO(), filter).Decode(&user)

	return user
}

func (ur *UserRepository) InsertRedis(key string, value interface{}, expires time.Duration) error {
	redisClient := infrastructure.RedisInit()
	set := redisClient.Set(ur.Ctx, key, value, expires).Err()
	defer redisClient.Close()
	return set
}

func (ur *UserRepository) GettRedis(key string) (res string, err error) {
	redisClient := infrastructure.RedisInit()
	
	get, err := redisClient.Get(ur.Ctx, key).Result()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer redisClient.Close()
	return get, err 
}