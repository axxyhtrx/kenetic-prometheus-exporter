package zyx

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"
)

// ungzip used for unpacking http payload to text
func ungzip(input *bytes.Reader) (string, error) {
	gzreader, e1 := gzip.NewReader(input)
	if e1 != nil {
		return "", e1
	}

	output, e2 := io.ReadAll(gzreader)
	if e2 != nil {
		return "", e2
	}

	result := string(output)
	return result, nil
}

// GetMD5Hash calculate md5 hash from string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetSha256Hash calculate sha256 hash from string
func GetSha256Hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Poller poll statistics from router each 10 seconds and update stats
func Poller() {
	for range time.Tick(time.Second * 10) {
		Zyx()
		SystemInf()
	}
}

// UnmarshallAndUpdateSpeed is used to unmarshal json data from router about uplink bandwidth utilization
func UnmarshallAndUpdateSpeed(Zmon *Zymon, jsoninput []byte) {
	var statistics InterfacesStatistics

	if err := json.Unmarshal(jsoninput, &statistics); err != nil {
		panic(err)
	}
	Zmon.TXSpeed = statistics.Show.Interface.Stat[0].Txspeed
	Zmon.RXSpeed = statistics.Show.Interface.Stat[0].Rxspeed
}

// UnmarshallAndUpdateSpeed is used to unmarshal json data from router about system cpu/ram resources utilization
func UnmarshallAndUpdateSystem(Zmon *Zymon, jsoninput []byte) {
	var system SystemInfo

	if err := json.Unmarshal(jsoninput, &system); err != nil {
		panic(err)
	}
	Zmon.CPU = system.Cpuload
	Zmon.FreeRam = system.Memfree
	Zmon.UsedRam = system.Memtotal - system.Memfree
}
