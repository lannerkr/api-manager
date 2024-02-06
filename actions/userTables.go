package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type PayloadTable struct {
	User []UsersTable `json:"user"`
}
type UsersHistory struct {
	Username   string    `json:"username" bson:"user_name"`
	Enabled    string    `json:"enabled" bson:"enabled"`
	LastLogin  time.Time `json:"lastLogin" bson:"last_login"`
	Realm      string    `json:"realm,omitempty" bson:"realm"`
	Days       int64     `json:"days" bson:"days"`
	AccExpires time.Time `json:"accountExpires" bson:"accountExpires"`
	StaticIP   string    `json:"Static_IP" bson:"Static_IP"`
	FramedIP   string    `json:"framedip" bson:"framedip"`
}
type UsersTable struct {
	Username    string `json:"username" bson:"user_name"`
	Fullname    string `json:"fullname" bson:"fullname"`
	Enabled     string `json:"enabled" bson:"enabled"`
	PwdChange   string `json:"change-password-at-signin" bson:"cpas"`
	UserHistory UsersHistory
}

func getUsersPCS(realm string) (userTable []UsersTable, err error) {

	url := "/local/users"

	resp, err := pulseReq(realm, "GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	user := PayloadTable{}
	if err := decoder.Decode(&user); err != nil {
		fmt.Println(err)
	}
	users := user.User

	return users, nil
}

func getUsersHistory(realm string) (userMap map[string]UsersHistory, err error) {
	userMap = make(map[string]UsersHistory)

	uri := configuration.DBUri
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Println(err)
		}
	}()
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(err)
		return nil, err
	}
	coll := client.Database("ldapDB").Collection("user_history")

	var users []UsersHistory
	filter := bson.D{{Key: "realm", Value: realm}}
	u, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = u.All(context.TODO(), &users); err != nil {
		log.Println(err)
		return nil, err
	}
	for _, u := range users {
		userMap[u.Username] = u
	}

	return userMap, nil
}

func getUsersTable(realm string, update bool) (u []UsersTable, err error) {

	var userMap map[string]UsersHistory
	var userTable []UsersTable

TOP:
	if update {
		if userTable, err = getUsersPCS(realm); err != nil {
			log.Println(err)
			return nil, err
		}
		if userMap, err = getUsersHistory(realm); err != nil {
			log.Println(err)
			return nil, err
		}
		for i, u := range userTable {
			userTable[i].UserHistory = userMap[u.Username]
		}
		if err = putUserTableDB(realm, userTable); err != nil {
			log.Println(err)
			return nil, err
		}
	} else {
		userTable, err = getUserTableDB(realm)
		if err != nil {
			if err.Error() == "NeedUpdate" {
				update = true
				goto TOP
			}
			log.Println(err)
			return nil, err
		}
	}

	return userTable, nil
}

func putUserTableDB(realm string, users []UsersTable) error {

	uri := configuration.DBUri
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Println(err)
		}
	}()
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(err)
		return err
	}

	var realmColl string
	if realm == "EMP-GOTP" {
		realmColl = "userTable_" + "EMPGOTP"
	} else {
		realmColl = "userTable_" + realm
	}
	coll := client.Database("ldapDB").Collection(realmColl)
	if err := coll.Drop(context.TODO()); err != nil {
		log.Println(err)
		return err
	}

	var update []interface{}
	for _, v := range users {
		update = append(update, v)
	}

	if _, err := coll.InsertMany(context.TODO(), update); err != nil {
		log.Println(err)
		return err
	}

	collu := client.Database("ldapDB").Collection("updateTime")
	filter := bson.M{"realm": realm}
	newdate := bson.D{{Key: "$set", Value: bson.M{"updateTime": time.Now()}}}
	opts := options.Update().SetUpsert(true)
	if _, err = collu.UpdateOne(context.TODO(), filter, newdate, opts); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func getUserTableDB(realm string) ([]UsersTable, error) {

	uri := configuration.DBUri
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Println(err)
		}
	}()
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(err)
		return nil, err
	}

	if checkDBupdate(client, realm) {
		return nil, fmt.Errorf("NeedUpdate")
	}

	var users []UsersTable
	var realmColl string
	if realm == "EMP-GOTP" {
		realmColl = "userTable_" + "EMPGOTP"
	} else {
		realmColl = "userTable_" + realm
	}
	coll := client.Database("ldapDB").Collection(realmColl)

	filter := bson.D{}
	u, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = u.All(context.TODO(), &users); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

type Utime struct {
	//ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	//realm      string             `bson:"realm"`
	UpdateTime time.Time `bson:"updateTime"`
}

func checkDBupdate(client *mongo.Client, realm string) bool {
	collu := client.Database("ldapDB").Collection("updateTime")

	//log.Printf("checkDBupdate : %v\n", realm)
	var result Utime
	if err := collu.FindOne(context.TODO(), bson.D{{Key: "realm", Value: realm}}).Decode(&result); err != nil {
		log.Printf("check Find fail: %v\n", err)
		return false
	}

	t1 := result.UpdateTime.Add(time.Hour)
	//log.Println(result, ":", result.UpdateTime, ":", t1)

	return t1.Before(time.Now())
}

func getSingleUserdata(realm, userid string) (u Userdata, err error) {

	url := "/local/users/user/" + userid

	resp, err := pulseReq(realm, "GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	user := Userdata{}
	if err := decoder.Decode(&user); err != nil {
		fmt.Println(err)
	}

	return user, nil
}

func getSingleUserTable(realm, user_id string) (userTable UsersTable, err error) {
	uri := configuration.DBUri
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(err)
		return userTable, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
		return userTable, err
	}

	var realmColl string
	if realm == "EMP-GOTP" {
		realmColl = "userTable_" + "EMPGOTP"
	} else {
		realmColl = "userTable_" + realm
	}
	coll := client.Database("ldapDB").Collection(realmColl)

	users := coll.FindOne(context.TODO(), bson.M{"user_name": user_id})
	if err = users.Decode(&userTable); err != nil {
		return userTable, err
	}

	return userTable, nil
}
