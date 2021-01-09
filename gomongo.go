package gomongo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Payload type
type Payload struct {
	Kind    string      `bson:"kind" json:"kind"`
	Values  interface{} `bson:"values" json:"values"`
	Options FindOptions `bson:"options,omitempty" json:"options,omitempty"`
}

//FindOptions type
type FindOptions struct {
	Limit      int64       `bson:"limit,omitempty" json:"limit,omitempty"`
	Projection interface{} `bson:"projection,omitempty" json:"projection,omitempty"`
	Sort       interface{} `bson:"sort,omitempty" json:"sort,omitempty"`
	Skip       int64       `bson:"skip,omitempty" json:"skip,omitempty"`
	Pagination interface{} `bson:"pagination,omitempty" json:"pagination,omitempty"`
}

//Identity type
type Identity struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Attributes        []Attributes       `bson:"attributes,omitempty" json:"attributes,omitempty"`
	CertificateNumber string             `bson:"certificate_number,omitempty" json:"certificate_number,omitempty"`
	CallSign          string             `bson:"call_sign" json:"call_sign"`
	EventID           string             `bson:"event_id" json:"event_id"`
	EventName         string             `bson:"event_name,omitempty" json:"event_name,omitempty"`
	Name              string             `bson:"name" json:"name"`
	IsFulfilled       bool               `bson:"is_fulfilled" json:"is_fulfilled"`
	DownloadCount     int32              `bson:"download_count" json:"download_count"`
	// Date              Date         `bson:"date" json:"date"`
}

// Date type
type Date struct {
	CreatedBy    string `bson:"created_by" json:"created_by"`
	DateCreated  string `bson:"date_created" json:"date_created"`
	DateModified string `bson:"date_modified" json:"date_modified"`
	ModifiedBy   string `bson:"modified_by" json:"modified_by"`
}

// Attributes type
type Attributes struct {
	Band      string `bson:"band" json:"band"`
	Frequency string `bson:"frequency" json:"frequency"`
	Date      string `bson:"date" json:"date"`
	RST int64 `bson:"rst,omitempty" json:"rst,omitempty"`
	Mode string `bson:"mode,omitempty" json:"mode,omitempty"`
}

// Image type
type Image struct {
	FileName string `bson:"file_name" json:"file_name"`
	B64      string `bson:"b64" json:"b64"`
}

type RGB struct {
	R int `bson:"r" json:"r"`
	G int `bson:"g" json:"g"`
	B int `bson:"b" json:"b"`
}

type TextPosition struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}

type TemplateProperty struct {
	TextPosition TextPosition `bson:"text_position" json:"text_position"`
	TextAlign    string       `bson:"text_align" json:"text_align"`
	FontColor    RGB          `bson:"font_color" json:"font_color"`
	FontName     string       `bson:"font_name" json:"font_name"`
	FontSize     float64      `bson:"font_size" json:"font_size"`
	FontDir      string       `bson:"font_dir" json:"font_dir"`
}

// ImageCertTemplate type
type ImageCertTemplate struct {
	FileName           string `bson:"file_name" json:"file_name"`
	B64                string `bson:"b64" json:"b64"`
	TemplateProperties struct {
		CallSign     TemplateProperty `bson:"call_sign" json:"call_sign"`
		IdentityName TemplateProperty `bson:"identity_name" json:"identity_name"`
		Frequency    TemplateProperty `bson:"frequency" json:"frequency"`
	} `bson:"template_properties" json:"template_properties"`
}

// CallSignPayload type
type CallSignPayload struct {
	Attributes        Attributes `bson:"attributes" json:"attributes"`
	CertificateNumber string     `bson:"certificate_number,omitempty" json:"certificate_number,omitempty"`
	CallSign          string     `bson:"call_sign" json:"call_sign"`
	EventID           string     `bson:"event_id" json:"event_id"`
	EventName         string     `bson:"event_name,omitempty" json:"event_name,omitempty"`
	Name              string     `bson:"name" json:"name"`
}

// EventCallSign type
type EventCallSign struct {
	Attributes          []Attributes       `bson:"attributes" json:"attributes"`
	CertificateTemplate string             `bson:"certificate_template" json:"certificate_template"`
	CertificateFormat   string             `bson:"certificate_format" json:"certificate_format"`
	Description         string             `bson:"description" json:"description"`
	Name                string             `bson:"name" json:"name"`
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	IsActive            bool               `bson:"is_active" json:"is_active"`
	IsHidden            bool               `bson:"is_hidden" json:"is_hidden"`
	CityID              int32              `bson:"city_id" json:"city_id"`
	//Date                string       `bson:"date" json:"date"`
}

// Adaptor Type
type Adaptor struct {
	Client mongo.Client
	DBName string
}

// Connect method
func (adaptor *Adaptor) Connect(ctx context.Context, uri string) {
	Client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	adaptor.Client = *Client

	err = adaptor.Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

// QueryUpdateDocument method
func (adaptor *Adaptor) QueryUpdateMany(ctx context.Context, collname string, filterQuery bson.M, updateQuery bson.M) error {
	var err error
	fmt.Println("filterQuery", filterQuery)
	fmt.Println("updateQuery", updateQuery)

	Collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	_, err = Collection.UpdateMany(ctx, filterQuery, updateQuery)

	return err
}

// QueryUpdateOne method
func (adaptor *Adaptor) QueryUpdateOne(ctx context.Context, collname string, filterQuery bson.M, updateQuery bson.M) (*mongo.UpdateResult, error) {
	var err error
	fmt.Println("filterQuery", filterQuery)
	fmt.Println("updateQuery", updateQuery)

	var result *mongo.UpdateResult
	result, err = adaptor.Client.Database(adaptor.DBName).Collection(collname).UpdateOne(ctx, filterQuery, updateQuery)

	return result, err
}

// QueryCreateCollection create collection in mongodb
func (adaptor *Adaptor) QueryCreateCollection(ctx context.Context, collname string) error {
	errCreateCollection := adaptor.Client.Database(adaptor.DBName).CreateCollection(ctx, collname)
	return errCreateCollection
}

// QueryInsert Query Insert to mongodb
func (adaptor *Adaptor) QueryInsert(ctx context.Context, collname string, byteQuery []byte) (interface{}, error) {
	var insertResult interface{}
	var errorInserting error

	var query bson.M
	bson.UnmarshalJSON(byteQuery, &query)
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	insertResult, errorInserting = collection.InsertOne(ctx, query)

	return insertResult, errorInserting
}

// QueryInsertV2 Query Insert to mongodb
func (adaptor *Adaptor) QueryInsertV2(ctx context.Context, collname string, query interface{}, result interface{}) error {
	result, errorInserting := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		InsertOne(ctx, query)

	if errorInserting != nil {
		log.Println(errorInserting)
	}

	return errorInserting
}

// QueryInsertV2 Query Insert to mongodb
func (adaptor *Adaptor) QueryInsertV3(ctx context.Context, collname string, query interface{}) (*mongo.InsertOneResult, error) {
	result, errorInserting := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		InsertOne(ctx, query)

	if errorInserting != nil {
		log.Println(errorInserting)
	}

	return result, errorInserting
}

// QueryFind query find to mongodb
func (adaptor *Adaptor) QueryFind(ctx context.Context, collname string, byteQuery []byte) ([]byte, error) {
	var query bson.M
	bson.UnmarshalJSON(byteQuery, &query)

	var received bson.M
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	errFinding := collection.FindOne(ctx, query).Decode(&received)
	jsonBytes, _ := bson.MarshalJSON(&received)

	return jsonBytes, errFinding
}

// QueryFindV2 query find to mongodb
func (adaptor *Adaptor) QueryFindV2(ctx context.Context, collName string, findOneOptions *options.FindOneOptions ,query interface{}, result interface{}) error {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collName)
	return collection.FindOne(ctx, query, findOneOptions).Decode(result)
}

// QueryFindMany query find many to mongodb
func (adaptor *Adaptor) QueryFindMany(ctx context.Context, collname string, byteQuery []byte, findOptions *options.FindOptions) ([]byte, error) {
	var query bson.M
	bson.UnmarshalJSON(byteQuery, &query)

	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	cursor, _ := collection.Find(ctx, query, findOptions)

	var received []bson.M
	var err error

	if err = cursor.All(ctx, &received); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(received)
	var results []byte
	results, err = bson.MarshalJSON(received)

	return results, err
}

// QueryFindManyV2 query find many to mongodb
func (adaptor *Adaptor) QueryFindManyV2(ctx context.Context, collname string, findOptions *options.FindOptions, query interface{}, result interface{}) error {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	cursor, _ := collection.Find(ctx, query, findOptions)

	err := cursor.All(ctx, result)
	if err != nil {
		panic(err)
	}

	return err
}

// QueryCount query find to mongodb
func (adaptor *Adaptor) QueryCount(ctx context.Context, collname string, query bson.M) (int64, error) {
	Count, err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		CountDocuments(ctx, query)

	return Count, err
}

// QueryFindAndUpdate method
func (adaptor *Adaptor) QueryFindAndUpdate(ctx context.Context, collname string, queryFilter bson.M, setQuery bson.M, setOnInsertQuery bson.M) (int64, error) {
	var err error
	var count int64

	updateQuery := bson.M{
		"$set":         setQuery,
		"$setOnInsert": setOnInsertQuery,
	}

	var updateOptions *options.FindOneAndUpdateOptions
	updateOptions.SetReturnDocument(1)
	updateOptions.SetUpsert(true)

	insertResult := adaptor.Client.Database(adaptor.DBName).Collection(collname).FindOneAndUpdate(ctx, queryFilter, updateQuery, updateOptions)
	fmt.Println(*insertResult)
	return count, err
}

// QueryFindAndUpdateV2 method
// updateQuery := bson.M{
//		"$set":         setQuery,
//		"$setOnInsert": setOnInsertQuery,
//	}
func (adaptor *Adaptor) QueryFindAndUpdateV2(ctx context.Context, collname string, filterQuery interface{}, updateQuery interface{}, result *bson.M) error {
	var updateOptions options.FindOneAndUpdateOptions
	updateOptions.SetReturnDocument(1)
	updateOptions.SetUpsert(true)

	err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		FindOneAndUpdate(ctx, filterQuery, updateQuery, &updateOptions).
		Decode(&result)

	fmt.Println("LAST INSERTED -->", result)

	return err
}

//QueryRemoveOne method
func (adaptor *Adaptor) QueryRemoveOne(ctx context.Context, collname string, queryFilter interface{}) (int64, error) {
	delResult, err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		DeleteOne(ctx, queryFilter)

	return delResult.DeletedCount, err
}

// QueryRemoveMany method
func (adaptor *Adaptor) QueryRemoveMany(ctx context.Context, collname string, queryFilter interface{}) (int64, error) {
	delResult, err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		DeleteMany(ctx, queryFilter)

	return delResult.DeletedCount, err
}

// QueryConfirm method
func (adaptor *Adaptor) QueryConfirm(ctx context.Context, collname, key, value string) bool {
	queryResult := bson.M{}
	errFindKey := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		FindOne(ctx, bson.M{"key": key}).
		Decode(&queryResult)
	if errFindKey != nil {
		panic(errFindKey)
	}

	if queryResult["value"].(string) == value {
		return true
	} else {
		return false
	}
}

///////////// PAYLOAD FILTER /////////////

// ParsePayload method
func (adaptor *Adaptor) ParsePayload(jsonByte []byte, out interface{}) {
	if isErr := bson.UnmarshalJSON(jsonByte, out); isErr != nil {
		fmt.Println(isErr)
	}
}

// Modeling filler
func (adaptor *Adaptor) Modeling(jsonByte *[]byte, collname string) error {
	var err error

	if collname == "identity" {
		identity := Identity{}
		err = bson.UnmarshalJSON(*jsonByte, &identity)
		*jsonByte, err = bson.MarshalJSON(&identity)

	} else if collname == "event" {
		event := EventCallSign{}
		err = bson.UnmarshalJSON(*jsonByte, &event)
		*jsonByte, err = bson.MarshalJSON(&event)
	}
	return err
}

// ParseOptions method
func (adaptor *Adaptor) ParseOptions(payload Payload, options *options.FindOptions) {
	// LIMIT
	limitVal := payload.Options.Limit
	if limitVal > 0 {
		if limitVal >= 100 {
			options.SetLimit(100)
		} else {
			options.SetLimit(limitVal)
		}
	} else {
		options.SetLimit(100)
	}

	// SORT
	if payload.Options.Sort != nil {
		options.SetSort(payload.Options.Sort)
	}

	// SKIP
	skipVal := payload.Options.Skip
	if skipVal >= 0 {
		options.SetSkip(skipVal)
	} else {
		options.SetSkip(0)
	}

	// PROJECTION
	if payload.Options.Projection != nil {
		options.SetProjection(payload.Options.Projection)
	}
}

// //GetIdentities method
// func (adaptor *Adaptor) GetIdentities(ctx context.Context, queryName bson.M, name string) []Identity {
// 	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
// 	cursor, err := collection.Find(ctx, queryName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(ctx)
//
// 	var result []Identity
// 	for cursor.Next(ctx) {
// 		received := bson.M{}
//
// 		if err := cursor.Decode(&received); err != nil {
// 			fmt.Println(err)
// 		}
//
// 		bsonBytes, _ := bson.Marshal(&received)
// 		var subIdentity Identity
// 		bson.Unmarshal(bsonBytes, &subIdentity)
//
// 		result = append(result, subIdentity)
// 	}
// 	return result
// }

// //DeleteCollection method
// func (adaptor *Adaptor) DeleteCollection(ctx context.Context, name string, deleteCode string) {
// 	if deleteCode == "AGREE TO DELETE "+adaptor.CollName {
// 		collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
// 		delResult, err := collection.DeleteMany(ctx, bson.M{"name": name})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		fmt.Println(delResult)
// 	} else {
// 		fmt.Println("ACCESS DENIED")
// 	}
// }

// // LegalizePayload method
// func (adaptor *Adaptor) LegalizePayload(bsonPayload bson.M, out interface{}) {
// 	for key, value := range bsonPayload {
// 		fmt.Printf("Key: %v Value: %v\n", key, value)
// 	}
// }

// GetDate method
func (adaptor *Adaptor) GetDate() string {
	Time := time.Now().UnixNano()
	dateRune := []rune(strconv.Itoa(int(Time)))
	parsedDate := string(dateRune[0:13])

	return parsedDate
}
