package zyx

import (
	"io"
	"log"
	"net/http"
)

type SystemInfo struct {
	Hostname   string `json:"hostname"`
	Domainname string `json:"domainname"`
	Cpuload    int    `json:"cpuload"`
	Memory     string `json:"memory"`
	Swap       string `json:"swap"`
	Memtotal   int    `json:"memtotal"`
	Memfree    int    `json:"memfree"`
	Membuffers int    `json:"membuffers"`
	Memcache   int    `json:"memcache"`
	Swaptotal  int    `json:"swaptotal"`
	Swapfree   int    `json:"swapfree"`
	Uptime     string `json:"uptime"`
}

func SystemInf() {
	client := &http.Client{}

	url := "http://" + IPaddr + "/rci/show/system"
	origin := "http://" + IPaddr
	referer := "http://" + IPaddr + "/dashboard"
	authcookie := "_authorized=" + User + ";" + Auth() + "; sysmode=router"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Origin", origin)
	req.Header.Set("Cookie", authcookie)
	req.Header.Set("Host", IPaddr)
	req.Header.Set("Referer", referer)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		log.Println("Not Authorized")
	}

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		UnmarshallAndUpdateSystem(&Zmon, bodyText)

	} else {
		log.Fatal("Request to API return an error: ", resp.Status)
	}

}
