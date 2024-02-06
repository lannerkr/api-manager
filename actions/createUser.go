package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type PPSuserAdd struct {
	Enabled         string    `json:"enabled"`
	Password        string    `json:"password-cleartext"`
	ConsoelAcc      string    `json:"console-access"`
	UserTimezone    string    `json:"user-timezone"`
	Guestcreatedby  string    `json:"guestcreatedby"`
	Username        string    `json:"username"`
	Fullname        string    `json:"fullname"`
	AttributesTable AttsTable `json:"attributes-table"`
}
type AttsTable struct {
	AttributeTable []AttTable `json:"attribute-table"`
}
type AttTable struct {
	Attr     string `json:"Attr"`
	StaticIP string `json:"Value"`
}

func userCreate(realm, user, ip, mac string) (err1, err2, err3 error) {
	// fmt.Println(realm, user, ip, mac)

	if user != "" && realm != "" {
		if err := addUserPulse(realm, user); err != nil {
			err1 = err
		}
	}

	if ip != "" {
		if err := addUserPPS(realm, user, ip); err != nil {
			err2 = err
		}
	} else {
		err2 = fmt.Errorf("static-ip is empty")
	}

	if mac != "" {
		if err := addDevicePPS(realm, mac); err != nil {
			err3 = err
		}

		// if err := macApprove(mac); err != nil {
		// 	fmt.Println(err)
		// }
	} else {
		err3 = fmt.Errorf("mac-address is empty")
	}
	return err1, err2, err3
}

func addUserPulse(realm, user string) error {
	//fmt.Println(realm, user)

	newuser := map[string]string{
		"enabled":                   "true",
		"change-password-at-signin": "true",
		"fullname":                  user,
		"password-cleartext":        "Fashion2022!",
		"username":                  user,
	}
	pbytes, _ := json.Marshal(newuser)
	buff := bytes.NewBuffer(pbytes)

	url := "/local/users/user/"

	resp, err := pulseReq(realm, "POST", url, buff)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	respStrings := respString(resp, "POST")

	if resp.StatusCode == 201 {

		if strings.Contains(respStrings, PostOK) {
			fmt.Printf("user: %v has been added successfully\n", user)
			return nil

		} else if strings.Contains(respStrings, PostExist) {
			err := fmt.Errorf("user: " + user + " already exist")
			return err
		} else {
			err := fmt.Errorf("user: " + user + " status change is failed")
			return err
		}
	}
	return nil
}

func addUserPPS(realm, user, ip string) error {
	//fmt.Println(realm, user)

	url := "/local/users/user"
	var newatt []AttTable
	newatt = append(newatt, AttTable{"Framed-IP-Address", ip})
	newatts := AttsTable{newatt}

	newuser := &PPSuserAdd{
		"true",
		"12345",
		"false",
		"seoul",
		"superman",
		user,
		ip,
		newatts,
	}

	pbytes, _ := json.Marshal(newuser)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsReq(realm, "POST", url, buff)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(resp.StatusCode)
	defer resp.Body.Close()

	respStrings := respString(resp, "POST")
	if resp.StatusCode == 201 {
		if strings.Contains(respStrings, PostOK) {
			fmt.Printf("static-ip for user : %v has been added successfully\n", user)
			return nil
		} else if strings.Contains(respStrings, PostExist) {
			err := fmt.Errorf("user: " + user + " already exist")
			return err
		} else {
			err := fmt.Errorf("user: " + user + " status change is failed")
			return err
		}
	}

	return nil
}

func addDevicePPS(realm, mac string) error {

	url := "profiler/endpoints?=addNewDevice"
	if realm == "EMP-GOTP" {
		realm = "emp"
	}
	// fmt.Println(url)
	setApprove := map[string]string{
		"macaddr":  mac,
		"category": "Windows",
		"notes":    realm,
		"status":   "approved",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("POST", url, buff)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	respStrings := respString(resp, "PPS")
	//fmt.Println(respStrings)
	if resp.StatusCode == 201 {
		if strings.Contains(respStrings, "operation: add") {
			fmt.Printf("mac-address for user : %v has been added successfully\n", mac)
			return nil
		} else {
			//fmt.Println("type 3" + PostExist + ":" + respStrings)
			err := fmt.Errorf("mac-address: " + mac + " adding is failed")
			return err
		}
	} else if resp.StatusCode == 409 {
		err := fmt.Errorf("mac-address: " + mac + " already exist")
		return err
	} else {
		err := fmt.Errorf(strconv.Itoa(resp.StatusCode))
		return err
	}
}

func macApprove(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	// fmt.Println(url)
	setApprove := map[string]string{
		"status": "approved",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		fmt.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func macUnApprove(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	// fmt.Println(url)
	setApprove := map[string]string{
		"status": "unapproved",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		fmt.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func macPermit(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	// fmt.Println(url)
	setApprove := map[string]string{
		"category": "permit",
		"override": "true",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		fmt.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func macProtect(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	// fmt.Println(url)
	setApprove := map[string]string{
		"category": "Windows",
		"override": "false",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		fmt.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}
