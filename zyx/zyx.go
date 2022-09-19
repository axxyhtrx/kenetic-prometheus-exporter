package zyx

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	IPaddr   = "192.168.1.1"
	User     = "admin"
	Password = "admin"
)

//Zymon is a main struct with collected metrics

type Zymon struct {
	TXSpeed int
	RXSpeed int
	CPU     int
	FreeRam int
	UsedRam int
}

var Zmon Zymon

// Auth is used for authentication to kinetic web api
func Auth() string {

	client := &http.Client{}
	url := "http://" + IPaddr + "/auth"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	referer := "http://" + IPaddr + "/dashboard"

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Host", IPaddr)
	req.Header.Set("Referer", referer)
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	realm, err := io.ReadAll(bytes.NewReader([]byte(resp.Header.Get("X-NDM-Realm"))))
	if err != nil {
		log.Fatal(err)
	}
	md5string := User + ":" + string(realm) + ":" + Password

	challenge, err := io.ReadAll(bytes.NewReader([]byte(resp.Header.Get("X-NDM-Challenge"))))
	if err != nil {
		log.Fatal(err)
	}
	challengeString := string(challenge) + GetMD5Hash(md5string)
	token := GetSha256Hash(challengeString)

	cookie, err := io.ReadAll(bytes.NewReader([]byte(resp.Header.Get("Set-Cookie"))))
	if err != nil {
		log.Fatal(err)
	}
	leftCookie := strings.Split(string(cookie), ";")[0]

	if err != nil {
		log.Fatal(err)
	}

	Jsondata := `{"login":"","password":""}`

	type Credentials struct {
		Login    string `json:"login,omitempty"`
		Password string `json:"password,omitempty"`
	}
	var cred Credentials
	json.Unmarshal([]byte(Jsondata), &cred)

	cred.Login = User
	cred.Password = token

	data, _ := json.Marshal(cred)
	jdata := strings.NewReader(string(data))

	req, err = http.NewRequest("POST", url, jdata)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Cookie", leftCookie)
	req.Header.Set("Host", IPaddr)
	req.Header.Set("Referer", referer)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		log.Println("Authenticated as user " + User + " on host " + IPaddr)
		return leftCookie
	} else {
		log.Fatal("Incorrect credentials for user " + User + " on host " + IPaddr)
	}

	return ""
}
