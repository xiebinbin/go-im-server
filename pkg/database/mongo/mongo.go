// Package mongo 三方 mgo 太长时间不维护了
//
//	官方 mongo 驱动很不友好
//	所以这里稍微对常用方法做了处理,可以直接调用这里的方法进行一些常规操作
//	复杂的操作,调用这里的 Collection 之后可获取里边的 Database 属性 和 Collection 属性操作
//	这里的添加和修改操作将会自动补全 create_time update_time 和 _id
package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"imsdk/pkg/app"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/unique"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mode = readpref.Mode

// Mode constants
const (
	PrimaryMode            = readpref.PrimaryMode
	PrimaryPreferredMode   = readpref.PrimaryPreferredMode
	SecondaryMode          = readpref.SecondaryMode
	SecondaryPreferredMode = readpref.SecondaryPreferredMode
	NearestMode            = readpref.NearestMode

	ServiceTypeSelfBuild = "self"
	ServiceTypeAws       = "aws"
)

var (
	client *mongo.Client
	conf   config
	isRepl bool
)

type (
	// CollectionInfo 集合包含的连接信息和查询等操作信息
	CollectionInfo struct {
		Database        *mongo.Database
		DatabasePrimary *mongo.Database
		DatabaseReader  *mongo.Database
		Collection      *mongo.Collection
		filter          bson.M
		limit           int64
		skip            int64
		sort            bson.M
		fields          bson.M
	}

	config struct {
		URL             string `toml:"url"`
		Database        string `toml:"database"`
		MaxConnIdleTime int    `toml:"max_conn_idle_time"`
		MaxPoolSize     int    `toml:"max_pool_size"`
		Username        string `toml:"username"`
		Password        string `toml:"password"`
		ReplicaSet      string `toml:"replicaSet"`
		PrimaryUrl      string `toml:"primary_url"`
		ReaderUrl       string `toml:"reader_url"`
		ServiceType     string `toml:"service_type"` // self(self build) | aws
	}
)

// Start 启动 mongo
func Start() {
	err1 := app.Config().Bind("db", "mongo", &conf)
	fmt.Println("mongo conf: ", conf, err1)
	if errors.Is(err1, app.ErrNodeNotExists) {
		// 配置节点不存在, 不启动服务
		return
	}
	var err error
	mongoOptions := options.Client()
	mongoOptions.SetMaxConnIdleTime(time.Duration(conf.MaxConnIdleTime) * time.Second)
	mongoOptions.SetMaxPoolSize(uint64(conf.MaxPoolSize))
	mongoOptions.SetRetryReads(true)
	if conf.ReplicaSet != "" {
		mongoOptions.SetReplicaSet(conf.ReplicaSet)
		isRepl = true
	}
	if conf.Username != "" && conf.Password != "" {
		mongoOptions.SetAuth(options.Credential{Username: conf.Username, Password: conf.Password})
	}
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "startMongo"})
	client, err = mongo.NewClient(mongoOptions.ApplyURI(conf.URL))
	if err != nil {
		log.Logger().Error(logCtx, "failed to new client , err: ", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Logger().Error(logCtx, "failed to connect , err: ", err)
	}
}

// Database 获取数据库连接
func Database(name ...string) *CollectionInfo {
	var db *mongo.Database
	if len(name) == 1 {
		db = client.Database(name[0])
	} else {
		db = client.Database(conf.Database)
	}
	collection := &CollectionInfo{
		Database: db,
		filter:   make(bson.M),
	}
	return collection
}

// SetTable 设置集合名称
func (collection *CollectionInfo) SetTable(name string, mode ...Mode) *CollectionInfo {
	//fmt.Println("SetTable: ", conf.URL)
	if isRepl {
		rpMode := PrimaryMode
		if len(mode) > 0 && mode[0] > 0 {
			rpMode = mode[0]
		}

		rp, _ := readpref.New(rpMode)
		collection.Collection = collection.Database.Collection(name, options.Collection().SetReadPreference(rp))
	} else {
		collection.Collection = collection.Database.Collection(name)
	}
	return collection
}

// Where 条件查询, bson.M{"field": "value"}
func (collection *CollectionInfo) Where(m bson.M) *CollectionInfo {
	collection.filter = m
	return collection
}

// Limit 限制条数
func (collection *CollectionInfo) Limit(n int64) *CollectionInfo {
	collection.limit = n
	return collection
}

// Skip 跳过条数
func (collection *CollectionInfo) Skip(n int64) *CollectionInfo {
	collection.skip = n
	return collection
}

// Sort 排序 bson.M{"create_time":-1}
func (collection *CollectionInfo) Sort(sorts bson.M) *CollectionInfo {
	collection.sort = sorts
	return collection
}

// Fields 指定查询字段
func (collection *CollectionInfo) Fields(fields interface{}) *CollectionInfo {
	kind := reflect.TypeOf(fields).Kind()
	if kind == reflect.String {
		fieldStr := fields.(string)
		if fieldStr != "" {
			collection.fields = GetMongoFieldsBsonByString(fieldStr)
		}
	} else if kind == reflect.Map && fields != nil {
		collection.fields = fields.(bson.M)
	}

	return collection
}

// InsertOne 写入单条数据
func (collection *CollectionInfo) InsertOne(document interface{}) (*mongo.InsertOneResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Collection.InsertOne(ctx, BeforeCreate(document))
}

func (collection *CollectionInfo) InsertOneOrigin(document interface{}) (*mongo.InsertOneResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Collection.InsertOne(ctx, document)
}

// InsertMany 写入多条数据
func (collection *CollectionInfo) InsertMany(documents interface{}) (*mongo.InsertManyResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var data []interface{}
	data = BeforeCreate(documents).([]interface{})
	return collection.Collection.InsertMany(ctx, data)
}

// UpdateOrInsert 存在更新,不存在写入, documents 里边的文档需要有 _id 的存在
func (collection *CollectionInfo) UpdateOrInsert(documents interface{}) (*mongo.UpdateResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	upsert := true
	return collection.Collection.UpdateMany(ctx, bson.M{}, documents, &options.UpdateOptions{Upsert: &upsert})
}

func (collection *CollectionInfo) Upsert(document interface{}) *mongo.SingleResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	upsert := true
	var isReturn options.ReturnDocument = 1
	opt := options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &isReturn}
	result := collection.Collection.FindOneAndUpdate(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)}, &opt)
	return result
}

func (collection *CollectionInfo) UpsertByBson(document interface{}) *mongo.SingleResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	upsert := true
	var isReturn options.ReturnDocument = 1
	opt := options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &isReturn}
	return collection.Collection.FindOneAndUpdate(ctx, collection.filter, document, &opt)
}

// UpdateOne 更新一条
func (collection *CollectionInfo) UpdateOne(document interface{}) (*mongo.UpdateResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Collection.UpdateOne(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)})
}

func (collection *CollectionInfo) UpdateOneBit(document interface{}) (*mongo.UpdateResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Collection.UpdateOne(ctx, collection.filter, bson.M{"$bit": document})
}

func (collection CollectionInfo) UpByID(id interface{}, document interface{}) (*mongo.UpdateResult, error) {
	return collection.Where(bson.M{"_id": id}).UpdateOne(document)
}

// UpdateMany 更新多条
func (collection *CollectionInfo) UpdateMany(document interface{}) (*mongo.UpdateResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Collection.UpdateMany(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)})
}

func (collection *CollectionInfo) UpdateManyBit(document interface{}) (*mongo.UpdateResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Collection.UpdateMany(ctx, collection.filter, bson.M{"$bit": BeforeUpdate(document)})
}

func (collection *CollectionInfo) UpsertMany(document interface{}) (*mongo.UpdateResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	return collection.Collection.UpdateMany(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)}, &opt)
}

// FindOne 查询一条数据
func (collection *CollectionInfo) FindOne(document interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := collection.Collection.FindOne(ctx, collection.filter, &options.FindOneOptions{
		Skip:       &collection.skip,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	return result.Decode(document)
}

func (collection *CollectionInfo) FindByID(id interface{}, document interface{}) error {
	return collection.Where(bson.M{"_id": id}).FindOne(document)
}

// FindMany 查询多条数据
func (collection *CollectionInfo) FindMany(documents interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	opt := options.Find().SetProjection(collection.fields).SetLimit(collection.limit).
		SetSort(collection.sort).SetSkip(collection.skip)
	result, err := collection.Collection.Find(ctx, collection.filter, opt)
	if err != nil {
		return err
	}
	defer result.Close(ctx)
	val := reflect.ValueOf(documents)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}
	slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)

	itemTyp := val.Elem().Type().Elem()
	for result.Next(ctx) {
		item := reflect.New(itemTyp)
		err := result.Decode(item.Interface())
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, reflect.Indirect(item))
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	if resErr := result.Err(); resErr != nil {
		return resErr
	}

	val.Elem().Set(slice)
	return nil
}

// Delete 删除数据,并返回删除成功的数量
func (collection *CollectionInfo) Delete() (int64, error) {
	if collection.filter == nil || len(collection.filter) == 0 {
		return 0, errors.New("you can't delete all documents, it's very dangerous")
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Collection.DeleteMany(ctx, collection.filter)
	return result.DeletedCount, err
}

// Count 根据指定条件获取总条数
func (collection *CollectionInfo) Count() int64 {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Collection.CountDocuments(ctx, collection.filter)
	if err != nil {
		return 0
	}
	return result
}

// BeforeCreate 创建数据前置操作
func BeforeCreate(document interface{}) interface{} {
	millis := funcs.GetMillis()
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)
	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeCreate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeCreate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			data[typ.Field(i).Tag.Get("bson")] = val.Field(i).Interface()
		}
		if val.FieldByName("ID").Type() == reflect.TypeOf(primitive.ObjectID{}) {
			data["_id"] = primitive.NewObjectID()
		}

		if val.FieldByName("ID").Kind() == reflect.String && val.FieldByName("ID").Interface() == "" {
			data["_id"] = primitive.NewObjectID().Hex()
		}

		if IsIntn(val.FieldByName("ID").Kind()) && val.FieldByName("ID").Interface() == 0 {
			data["_id"] = unique.ID()
		}

		if data["create_time"] == 0 {
			data["create_time"] = millis
		}
		if data["update_time"] == 0 {
			data["update_time"] = millis
		}
		return data

	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			if !val.MapIndex(reflect.ValueOf("_id")).IsValid() {
				val.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
			}
			val.SetMapIndex(reflect.ValueOf("create_time"), reflect.ValueOf(millis))
			val.SetMapIndex(reflect.ValueOf("update_time"), reflect.ValueOf(millis))
		}
		return val.Interface()
	}
}

// BeforeUpdate 更新数据前置操作
func BeforeUpdate(document interface{}) interface{} {
	millis := funcs.GetMillis()
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)
	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeUpdate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeUpdate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			if !isZero(val.Field(i)) {
				tag := strings.Split(typ.Field(i).Tag.Get("bson"), ",")[0]
				data[tag] = val.Field(i).Interface()
				if tag != "_id" {
					data[tag] = val.Field(i).Interface()
				}
			}
		}
		//time.Now().Unix()
		if data["update_time"] == 0 {
			data["update_time"] = time.Now().UnixNano() / 1e6
		}

		return data

	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			val.SetMapIndex(reflect.ValueOf("update_time"), reflect.ValueOf(millis))
		}
		return val.Interface()
	}
}

// IsIntn 是否为整数
func IsIntn(p reflect.Kind) bool {
	return p == reflect.Int || p == reflect.Int64 || p == reflect.Uint64 || p == reflect.Uint32
}

func isZero(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func (collection *CollectionInfo) LBS(pipeline interface{}, documents interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Aggregate()
	result, _ := collection.Collection.Aggregate(ctx, pipeline, opts)
	defer result.Close(ctx)
	val := reflect.ValueOf(documents)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
	}
	slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)
	////fmt.Println(slice)
	itemTyp := val.Elem().Type().Elem()
	for result.Next(ctx) {
		item := reflect.New(itemTyp)
		err := result.Decode(item.Interface())
		if err != nil {
			continue
		}
		if err != nil {
			continue
		}
		slice = reflect.Append(slice, reflect.Indirect(item))
	}
	ctxErr := ctx.Err()
	resErr := result.Err()
	if ctxErr != nil || resErr != nil {
	}

	val.Elem().Set(slice)
	//return result
}

func (collection *CollectionInfo) Aggregate(pipeline interface{}, documents interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Aggregate()
	result, _ := collection.Collection.Aggregate(ctx, pipeline, opts)
	defer result.Close(ctx)
	val := reflect.ValueOf(documents)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}
	slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)
	itemTyp := val.Elem().Type().Elem()
	for result.Next(ctx) {
		item := reflect.New(itemTyp)
		err := result.Decode(item.Interface())
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, reflect.Indirect(item))
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	if resErr := result.Err(); resErr != nil {
		return resErr
	}

	val.Elem().Set(slice)
	return nil
}

func GetMongoFieldsBsonByString(fields string) bson.M {
	fieldsSlice := strings.Split(fields, ",")
	var res = make(bson.M)
	for _, f := range fieldsSlice {
		f = strings.Replace(f, " ", "", -1)
		res[f] = 1
	}
	return res
}

func IsFiledDuplicateKeyError(err error, filed string) bool {
	if mongo.IsDuplicateKeyError(err) {
		return strings.Contains(err.Error(), "index: "+filed+"_")
	}
	return false
}
