package gathertool

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func MongoConn(){
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// ===== 获得在test库里面Trainer集合的handle
	collection := client.Database("test").Collection("trainers")
	log.Println(collection)


	//数据结构体
	type Trainer struct {
		Name string
		Age  int
		City string
	}


	// ===== 插入一个单独的文档
	//ash := Trainer{"aa", 10, "Pallet Town"}
	//insertResult, err := collection.InsertOne(context.TODO(), ash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Inserted a single document: ", insertResult.InsertedID)


	// ===== 插入多个文档 collection.InsertMany() 函数会采用一个slice对象
	//misty := Trainer{"Misty", 10, "Cerulean City"}
	//brock := Trainer{"Brock", 15, "Pewter City"}
	//trainers := []interface{}{misty, brock}
	//insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)


	// ===== 查找文档
	// create a value into which the result can be decoded
	var result Trainer
	filter := bson.D{{"name", "Ash"}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println("not find err : ", err)
	}
	log.Printf("Found a single document: %+v\n", result)



	// ===== 要查询多个文档， 使用collection.Find()，这个函数返回一个游标
	findOptions := options.Find()
	//findOptions.SetLimit(100)
	var results []*Trainer
	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("elem is : %v", elem)
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	log.Printf("Found multiple documents (array of pointers): %+v\n", results)



	// ===== 更新文档
	//collection.UpdateOne()函数允许你更新单一的文档
	filterUpdate := bson.D{{"name", "aa"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filterUpdate, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("updateResult : %+v\n", updateResult)

	resultUpdate := &Trainer{}
	err = collection.FindOne(context.TODO(), filterUpdate).Decode(&resultUpdate)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("resultUpdate : ", resultUpdate)


	// ===== 删除文档
	//可以使用collection.DeleteOne() 或者 collection.DeleteMany()来删除文档
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{"name", "Ash"}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)


	// ===== 关闭连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}

// TODO Mongo 连接池

// TODO ??  如何避免写入重复数据

// TODO ??  查看所有Database

// TODO ??  查看所有Collection