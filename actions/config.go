package actions

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	TotpServer     string
	PulseUri       string
	PulseApiKey    string
	AuthStore      string
	AuthPartner    string
	AuthEmp        string
	PPSUri         string
	PPSApiKey      string
	PPSAuthStore   string
	PPSAuthPartner string
	PPSAuthEmp     string
	DBUri          string
	ApiLog         string
	MongoDBS       string
	Htpasswd       string
	Cert           string
	CertKey        string
	IPaccess       []string
	LogPath        string
}

var configuration Configuration

func Config(conf string) (cert, certkey, logPath string) {
	file, _ := os.Open(conf)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(configuration)

	return configuration.Cert, configuration.CertKey, configuration.LogPath
}
