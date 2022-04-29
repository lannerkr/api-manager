package actions

import (
	"encoding/json"
	"fmt"
)

type ReponseUserStats struct {
	UserStats UserStat `json:"user-stats"`
}
type UserStat struct {
	LicensedUserCount string `json:"allocated-user-count"`
	CurrentUserCount  string `json:"current-user-count"`
	MaxUserLast24hrs  string `json:"max-active-user-count-24hrs"`
}

type HealthCheck struct {
	CpuUtil  string `json:"CPU-UTILIZATION"`
	SwapUtil string `json:"SWAP-UTILIZATION"`
	DiskUtil string `json:"DISK-UTILIZATION"`
}

type ResponseStatus struct {
	ClusterStatus    []ClusterMember  `json:"cluster-member-status"`
	LastConfigUpdate LastConfigUpdate `json:"last-config-update"`
	SystemDate       string           `json:"system-date-and-time"`
	UpTime           UpTime           `json:"uptime"`
}
type ClusterMember struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
	Notes   string `json:"notes"`
	Status  string `json:"status-code"`
}
type LastConfigUpdate struct {
	Device string `json:"device"`
}
type UpTime struct {
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

func getStatus() (sysdate string, uptime string, configdate string, cluster []ClusterMember, err error) {

	url := "system/status/overview"

	// fmt.Println("url = " + url)
	resp, err := pulseSysReq("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	respRecords := ResponseStatus{}
	if err := decoder.Decode(&respRecords); err != nil {
		fmt.Println(err)
	}

	recordMember := respRecords.ClusterStatus
	//recordMember1 := respRecords.ClusterStatus[1]
	recordConfigUpdate := respRecords.LastConfigUpdate.Device
	recordSystemDate := respRecords.SystemDate

	upTimeD := fmt.Sprint(respRecords.UpTime.Days)
	upTimeH := fmt.Sprint(respRecords.UpTime.Hours)
	upTimeM := fmt.Sprint(respRecords.UpTime.Minutes)
	upTimeS := fmt.Sprint(respRecords.UpTime.Seconds)
	recordUpTime := upTimeD + " days, " + upTimeH + ":" + upTimeM + ":" + upTimeS

	//fmt.Println(respRecords.UpTime)
	//fmt.Println(recordSystemDate, recordUpTime, recordConfigUpdate)

	return recordSystemDate, recordUpTime, recordConfigUpdate, recordMember, nil
}

func getStats() (licensed, current, lastmax string, err error) {

	url := "system/user-stats"

	// fmt.Println("url = " + url)
	resp, err := pulseSysReq("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	respRecords := ReponseUserStats{}
	if err := decoder.Decode(&respRecords); err != nil {
		fmt.Println(err)
	}

	recordLicensed := respRecords.UserStats.LicensedUserCount
	recordCurrent := respRecords.UserStats.CurrentUserCount
	recordLastMax := respRecords.UserStats.MaxUserLast24hrs

	//fmt.Println(recordLicensed, recordCurrent, recordLastMax)

	return recordLicensed, recordCurrent, recordLastMax, nil
}

func getHealth() (cpu, swap, disk string, err error) {

	url := "system/healthcheck?status=all"

	// fmt.Println("url = " + url)
	resp, err := pulseSysReq("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	respRecords := HealthCheck{}
	if err := decoder.Decode(&respRecords); err != nil {
		fmt.Println(err)
	}

	recordCpu := respRecords.CpuUtil
	recordSwap := respRecords.SwapUtil
	recordDisk := respRecords.DiskUtil

	//fmt.Println(recordCpu, recordSwap, recordDisk)

	return recordCpu, recordSwap, recordDisk, nil
}
