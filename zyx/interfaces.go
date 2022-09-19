package zyx

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
)

type InterfacesStatistics struct {
	Show struct {
		Interface struct {
			Stat []struct {
				Rxpackets          int    `json:"rxpackets"`
				RxMulticastPackets int    `json:"rx-multicast-packets"`
				RxBroadcastPackets int    `json:"rx-broadcast-packets"`
				Rxbytes            int64  `json:"rxbytes"`
				Rxerrors           int    `json:"rxerrors"`
				Rxdropped          int    `json:"rxdropped"`
				Txpackets          int    `json:"txpackets"`
				TxMulticastPackets int    `json:"tx-multicast-packets"`
				TxBroadcastPackets int    `json:"tx-broadcast-packets"`
				Txbytes            int64  `json:"txbytes"`
				Txerrors           int    `json:"txerrors"`
				Txdropped          int    `json:"txdropped"`
				Timestamp          string `json:"timestamp"`
				LastOverflow       string `json:"last-overflow"`
				Rxspeed            int    `json:"rxspeed"`
				Txspeed            int    `json:"txspeed"`
			} `json:"stat"`
		} `json:"interface"`
	} `json:"show"`
}

func Zyx() {
	client := &http.Client{}
	var data = strings.NewReader(`{"show":{"interface":{"stat":[{"name":"GigabitEthernet1"},{"name":"UsbModem0"}]}}}`)

	url := "http://" + IPaddr + "/rci/"
	origin := "http://" + IPaddr
	referer := "http://" + IPaddr + "/dashboard"
	authcookie := "_authorized=" + User + ";" + Auth() + "; sysmode=router"

	req, err := http.NewRequest("POST", url, data)
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
		reader := bytes.NewReader([]byte(bodyText))
		responsePayload, err := ungzip(reader)
		if err != nil {
			log.Println(err)
		}
		UnmarshallAndUpdateSpeed(&Zmon, []byte(responsePayload))

	} else {
		log.Fatal("Request to API return an error: ", resp.Status)
	}

}
