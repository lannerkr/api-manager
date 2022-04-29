package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type UserInterface interface{}

type Payload struct {
	User []Userdata `json:"user"`
}
type Userdata struct {
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Enabled   string `json:"enabled"`
	PwdChange string `json:"change-password-at-signin"`
}

type PPSuser struct {
	AttributesTables AttributesTable `json:"attributes-table"`
}

type AttributesTable struct {
	AttributeTable []struct {
		Attr     string `json:"Attr"`
		StaticIP string `json:"Value"`
	} `json:"attribute-table"`
}

type PPSPayload struct {
	User []PPSUserdata `json:"user"`
}

type PPSUserdata struct {
	Username         string          `json:"username"`
	Fullname         string          `json:"fullname"`
	Enabled          string          `json:"enabled"`
	AttributesTables AttributesTable `json:"attributes-table"`
}

func getUsersdata(realm string) (u []Userdata, d map[string]string, err error) {

	url := "/local/users"

	resp, err := pulseReq(realm, "GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	user := Payload{}
	if err := decoder.Decode(&user); err != nil {
		fmt.Println(err)
	}
	users := user.User
	daysFromLastlogin := getDays(users)

	return users, daysFromLastlogin, nil
}

func getDays(u []Userdata) map[string]string {

	uri := configuration.DBUri
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
	}

	//insert or update user history
	coll := client.Database("ldapDB").Collection("user_history")

	var daymap map[string]string = make(map[string]string)
	for i := 0; i < len(u); i++ {
		var result bson.M
		if err := coll.FindOne(context.TODO(), bson.M{"user_name": u[i].Username}).Decode(&result); err != nil {
			continue
		}
		for k, v := range result {
			if k == "days" {
				vs := fmt.Sprint(v)
				daymap[u[i].Username] = vs
				break
			}
		}
	}
	return daymap
}

func getUserdata(realm, userid string) (u Userdata, d map[string]string, err error) {

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

	daysFromLastlogin := getDay(user)

	return user, daysFromLastlogin, nil
}

func getDay(u Userdata) map[string]string {

	uri := configuration.DBUri
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
	}

	//insert or update user history
	coll := client.Database("ldapDB").Collection("user_history")

	var daymap map[string]string = make(map[string]string)

	var result bson.M
	if err := coll.FindOne(context.TODO(), bson.M{"user_name": u.Username}).Decode(&result); err != nil {
		fmt.Println(err)
	}

	for k, v := range result {
		if k == "days" {
			vs := fmt.Sprint(v)
			daymap[u.Username] = vs
			break
		}
	}

	return daymap
}

func resetPW(realm, user string) error {

	updata := map[string]string{
		"password-cleartext":        "Fashion2022!",
		"change-password-at-signin": "true",
	}
	pbytes, _ := json.Marshal(updata)
	buff := bytes.NewBuffer(pbytes)

	url := "/local/users/user/" + user + "/password-cleartext"
	resp, err := pulseReq(realm, "PUT", url, buff)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("user: " + user + " password reset is done.")
		return nil
	} else {
		err := fmt.Errorf("user: " + user + " password reset is failed")
		return err
	}
}

func updateStatus(r, user string, status bool) error {

	st := strconv.FormatBool(status)
	stup := map[string]string{
		"enabled": st,
	}
	pbytes, _ := json.Marshal(stup)
	buff := bytes.NewBuffer(pbytes)

	url := "/local/users/user/" + user + "/enabled"
	// fmt.Println(url)
	resp, err := pulseReq(r, "PUT", url, buff)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("user: " + user + " enabled status is changed to " + st)
		return nil
	} else {
		err := fmt.Errorf("user: " + user + " status change is failed")
		return err
	}
}

func getStaticIP(realm, userid string) (staticip string) {
	url := "/local/users/user/" + userid

	resp, err := ppsReq(realm, "GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	user := PPSuser{}
	if err := decoder.Decode(&user); err != nil {
		fmt.Println(err)
	}

	// fmt.Println(user)

	if len(user.AttributesTables.AttributeTable) != 0 {
		staticip = user.AttributesTables.AttributeTable[0].StaticIP
	}

	return staticip
}

func getStaticUsers(realm string) (users []PPSUserdata, err error) {
	url := "/local/users"

	resp, err0 := ppsReq(realm, "GET", url, nil)
	if err0 != nil {
		fmt.Println(err0)
		defer resp.Body.Close()
		return nil, err0
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	user := PPSPayload{}
	if err := decoder.Decode(&user); err != nil {
		fmt.Println(err)
	}
	users = user.User

	return users, nil
}

func deleteUser(realm, userid string) error {

	url := "/local/users/user/" + userid

	resp, err := pulseReq(realm, "DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
	} else if resp.StatusCode == 202 || resp.StatusCode == 200 {
		defer resp.Body.Close()

		respGet, errGet := pulseReq(realm, "GET", url, nil)
		if respGet.StatusCode == 404 {
			defer respGet.Body.Close()
			return nil
		}
		defer respGet.Body.Close()
		fmt.Println(errGet)
		fmt.Println(respGet.Body)
		return errGet
	}
	defer resp.Body.Close()
	return err
}
