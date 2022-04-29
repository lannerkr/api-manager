package actions

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Response struct {
	ActiveUsers ActiveUserRecords `json:"active-users"`
}

type ActiveUserRecords struct {
	ActiveUserRecord UserRecords `json:"active-user-records"`
	TotalMatch       int         `json:"total-matched-record-number"`
	TotalRuturn      int         `json:"total-returned-record-number"`
	UserPermission   bool        `json:"user-login-permission"`
}

type UserRecords struct {
	UserRecord []Records `json:"active-user-record"`
}

type Records struct {
	Username      string `json:"active-user-name"`
	Realm         string `json:"authentication-realm"`
	ConnectIP     string `json:"network-connect-ip"`
	ClientVersion string `json:"pulse-client-version"`
	Role          string `json:"user-roles"`
	LoginTime     string `json:"user-sign-in-time"`
	SessionID     string `json:"session-id"`
}

func getActiveUsers(user, count string) (record []Records, err error) {
	var url string
	if user == "all" {
		url = "system/active-users?number=" + count
	} else {
		url = "system/active-users?name=" + user
	}
	// fmt.Println("url = " + url)
	resp, err := pulseSysReq("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	respRecords := Response{}
	if err := decoder.Decode(&respRecords); err != nil {
		fmt.Println(err)
	}
	record = respRecords.ActiveUsers.ActiveUserRecord.UserRecord
	return record, nil
}

func deleteActiveUser(sid string) int {
	url := "system/active-users/session/" + sid
	//fmt.Println("delete session url = " + url)
	resp, err := pulseSysReq("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func totpReset(user string) error {
	totp := configuration.TotpServer
	url := "totp/" + totp + "/users/" + user + "?operation=reset"
	// fmt.Println(url)
	resp, err := pulseSysReq("PUT", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	respStrings := respString(resp, "OTP")
	//fmt.Println(respStrings)
	if resp.StatusCode == 200 && strings.Contains(respStrings, "has been reset") {
		fmt.Println(respStrings)
		return nil
	} else {
		err = fmt.Errorf(respStrings)
	}

	return err
}

func totpUnlock(user string) error {
	totp := configuration.TotpServer
	url := "totp/" + totp + "/users/" + user + "?operation=unlock"
	// fmt.Println(url)
	resp, err := pulseSysReq("PUT", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	respStrings := respString(resp, "OTP")
	//fmt.Println(respStrings)
	if resp.StatusCode == 200 && strings.Contains(respStrings, "has been unlocked") {
		fmt.Println(respStrings)
		return nil
	} else {
		err = fmt.Errorf(respStrings)
	}

	return err
}
