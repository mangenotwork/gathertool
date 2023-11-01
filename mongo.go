/*
	Description : mongoDB相关的操作
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	User        string
	Password    string
	Host        string
	Port        string
	Conn        *mongo.Client
	Database    *mongo.Database
	Collection  *mongo.Collection
	MaxPoolSize int
	TimeOut     time.Duration
}

// NewMongo 新建mongoDB客户端对象
func NewMongo(user, password, host, port string) (*Mongo, error) {
	m := &Mongo{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
	if m.Host == "" {
		m.Host = "127.0.0.1"
	}
	if m.Port == "" {
		m.Port = "27017"
	}
	err := m.GetConn()
	return m, err
}

// GetConn 建立mongodb 连接
func (m *Mongo) GetConn() (err error) {
	uri := fmt.Sprintf("mongodb://")
	if m.User != "" && m.Password != "" {
		uri = uri + fmt.Sprintf("%s:%s@", m.User, m.Password)
	}
	uri = uri + fmt.Sprintf("%s:%s", m.Host, m.Port)

	if m.TimeOut < 10*time.Second {
		m.TimeOut = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.TimeOut)
	defer cancel()

	o := options.Client().ApplyURI(uri)
	if m.MaxPoolSize > 0 {
		o.SetMaxPoolSize(uint64(m.MaxPoolSize))
	}

	m.Conn, err = mongo.Connect(ctx, o)
	if err != nil {
		return
	}

	err = m.Conn.Ping(context.Background(), readpref.Primary())
	return
}

// GetDB 连接mongodb 的db
// dbname:DB名
func (m *Mongo) GetDB(dbname string) {
	if m.Conn == nil {
		_ = m.GetConn()
	}
	m.Database = m.Conn.Database(dbname)
}

// GetCollection 连接mongodb 的db的集合
// dbname:DB名;  name:集合名
func (m *Mongo) GetCollection(dbname, name string) {
	if m.Conn == nil {
		_ = m.GetConn()
	}
	m.Collection = m.Conn.Database(dbname).Collection(name)
}

// Insert 插入数据
// document:可以是 Struct, 是 Slice
func (m *Mongo) Insert(document interface{}) error {
	if m.Collection == nil {
		return fmt.Errorf("Collection is nil;")
	}
	v := reflect.ValueOf(document)

	if reflect.ValueOf(document).Kind() == reflect.Struct {
		insertResult, err := m.Collection.InsertOne(context.TODO(), document)
		if err != nil {
			return err
		}
		Info("Inserted a single document: ", insertResult.InsertedID)
	}

	if v.Kind() == reflect.Slice {
		insertManyResult, err := m.Collection.InsertMany(context.TODO(), document.([]interface{}))
		if err != nil {
			return err
		}
		Info("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	}

	return nil
}

// MongoConn mongoDB客户端连接
func MongoConn() {
	m, err := NewMongo("", "", "", "")
	if err != nil {
		Error(err)
		return
	}
	m.GetCollection("test", "trainers")

	//数据结构体
	type Trainer struct {
		Name string
		Age  int
		City string
	}

	// ===== 插入一个单独的文档
	ash := Trainer{"aa", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}
	trainers := []interface{}{ash, misty, brock}
	_ = m.Insert(trainers)
}

// MongoConn1 mongoDB客户端连接
func MongoConn1() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		Error(err)
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		Error(err)
		return
	}
	Info("Connected to MongoDB!")
	// ===== 获得在test库里面Trainer集合的handle
	collection := client.Database("test").Collection("trainers")
	Info(collection)
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
	//	Error(err)
	//}
	//Info("Inserted a single document: ", insertResult.InsertedID)

	// ===== 插入多个文档 collection.InsertMany() 函数会采用一个slice对象
	//misty := Trainer{"Misty", 10, "Cerulean City"}
	//brock := Trainer{"Brock", 15, "Pewter City"}
	//trainers := []interface{}{misty, brock}
	//insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	//if err != nil {
	//	Error(err)
	//}
	//Info("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	// ===== 查找文档
	// create a value into which the result can be decoded
	var result Trainer
	filter := bson.D{{"name", "Ash"}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		Error("not find err : ", err)
	}
	Info("Found a single document: %+v\n", result)

	// ===== 要查询多个文档， 使用collection.Find()，这个函数返回一个游标
	findOptions := options.Find()
	//findOptions.SetLimit(100)
	var results []*Trainer
	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		Error(err)
	}
	for cur.Next(context.TODO()) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			Error(err)
		}
		Info("elem is : %v", elem)
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		Error(err)
	}
	_ = cur.Close(context.TODO())
	Info("Found multiple documents (array of pointers): %+v\n", results)

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
		Error(err)
	}
	Info("updateResult : %+v\n", updateResult)

	resultUpdate := &Trainer{}
	err = collection.FindOne(context.TODO(), filterUpdate).Decode(&resultUpdate)
	if err != nil {
		Error(err)
	}
	Info("resultUpdate : ", resultUpdate)

	// ===== 删除文档
	//可以使用collection.DeleteOne() 或者 collection.DeleteMany()来删除文档
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{"name", "Ash"}})
	if err != nil {
		Error(err)
	}
	Info("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	// ===== 关闭连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		Error(err)
	}
	Info("Connection to MongoDB closed.")
}

// TODO Mongo 连接池

// TODO  如何避免写入重复数据

// TODO  查看所有Database

// TODO  查看所有Collection
