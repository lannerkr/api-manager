package actions

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func pulseReq(realm, method, url string, buff io.Reader) (resp *http.Response, err error) {
	pulseUri := configuration.PulseUri
	pulseApiKey := configuration.PulseApiKey
	var auth string

	switch realm {
	case "store":
		auth = configuration.AuthStore
	case "partner":
		auth = configuration.AuthPartner
	case "emp":
		auth = configuration.AuthEmp
	case "EMP-GOTP":
		auth = configuration.AuthEmp
	default:
		err := fmt.Errorf("realm %v is not available", r)
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := pulseApiKey
	pUri := pulseUri + "/api/v1/configuration/authentication/auth-servers/auth-server/" + auth + url

	//fmt.Println(pUri)

	req, err := http.NewRequest(method, pUri, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// if method != "GET" {
	// 	defer resp.Body.Close()
	// }

	return resp, nil
}

func pulseSysReq(method, url string, buff io.Reader) (resp *http.Response, err error) {
	pulseUri := configuration.PulseUri
	pulseApiKey := configuration.PulseApiKey

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := pulseApiKey

	pUri := pulseUri + "/api/v1/" + url

	// fmt.Println("pulseSysReq URL = " + pUri)
	req, err := http.NewRequest(method, pUri, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	// if method != "GET" {
	// 	defer resp.Body.Close()
	// }

	return resp, nil
}

func ppsReq(realm, method, url string, buff io.Reader) (resp *http.Response, err error) {
	ppsUri := configuration.PPSUri
	ppsApiKey := configuration.PPSApiKey
	var auth string

	switch realm {
	case "store":
		auth = configuration.PPSAuthStore
	case "partner":
		auth = configuration.PPSAuthPartner
	case "emp":
		auth = configuration.PPSAuthEmp
	case "EMP-GOTP":
		auth = configuration.PPSAuthEmp
	default:
		err := fmt.Errorf("realm %v is not available", r)
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := ppsApiKey
	pUri := ppsUri + "/api/v1/configuration/authentication/auth-servers/auth-server/" + auth + url

	req, err := http.NewRequest(method, pUri, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ppsSysReq(method, url string, buff io.Reader) (resp *http.Response, err error) {
	ppsUri := configuration.PPSUri
	ppsApiKey := configuration.PPSApiKey

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := ppsApiKey
	pUri := ppsUri + "/api/v1/" + url

	req, err := http.NewRequest(method, pUri, buff)
	if err != nil {
		//fmt.Printf("ppsSysReq Error 1 : %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "")

	resp, err = client.Do(req)
	if err != nil {
		//fmt.Printf("ppsSysReq Error 2 : %v\n", err)
		return nil, err
	}
	//fmt.Printf("ppsSysReq resp : %v\n", resp)
	return resp, nil
}

type PostResp struct {
	Results Result `json:"result"`
}
type Result struct {
	Warnings []Warning `json:"warnings"`
}
type Warning struct {
	Messages string `json:"message"`
}

const (
	PostOK    string = "The configuration has been implicitly changed"
	PostExist string = "already exist"
)

func respString(resp *http.Response, method string) (respString string) {

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	//fmt.Println(bodyString)
	jsonResp := strings.NewReader(bodyString)

	if strings.Contains(bodyString, "admin_op") {
		type DeviceResp struct {
			Status  string `json:"status"`
			Mac     string `json:"macaddr"`
			AdminOp string `json:"admin_op"`
		}
		var postResp DeviceResp
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			fmt.Println(err)
		}
		respString = "mac-address: " + postResp.Mac + " , status: " + postResp.Status + " , operation: " + postResp.AdminOp
		return respString
	}

	if method == "POST" {
		var postResp PostResp
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			fmt.Println(err)
		}
		respString = postResp.Results.Warnings[0].Messages
	} else if method == "PPS" {
		var postResp Warning
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			fmt.Println(err)
		}
		respString = postResp.Messages
	} else if method == "OTP" {
		type TypeError struct {
			Messages string `json:"message"`
		}
		type Info struct {
			Messages string `json:"message"`
		}
		type Result struct {
			ErrorsR []TypeError `json:"errors"`
			Infos   []Info      `json:"info"`
		}
		type OTPresp struct {
			Results Result `json:"result"`
		}

		var postResp OTPresp
		//fmt.Println(postResp)
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			fmt.Println(err)
		}
		//fmt.Println(postResp)
		if postResp.Results.Infos != nil {
			//fmt.Println(postResp.Results.Infos)
			respString = postResp.Results.Infos[0].Messages
		} else if postResp.Results.ErrorsR != nil {
			//fmt.Println(postResp.Results.ErrorsR)
			respString = postResp.Results.ErrorsR[0].Messages
		}
	}

	return respString
}
