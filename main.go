package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/glendc/go-external-ip"
	"github.com/ip2location/ip2location-go"
	"github.com/sixdouglas/suncalc"
	flag "github.com/spf13/pflag"
	"os"
	"time"
)

const geolocCacheFile = "/tmp/zcock_geoloc_cache"

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func getAnimalIdx(hour int) int {
	if hour == 23 || hour == 0 {
		return 0
	}

	for i := 23; i != 1; i -= 2 {
		if hour >= i {
			return (i + 1) / 2
		}

	}

	return 1
}

func floatToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func getPublicIpAddr() string {
	ipCacheFile := "/tmp/zcock_ip_cache"

	if _, err := os.Stat(ipCacheFile); err == nil {
		ip, err := os.ReadFile(ipCacheFile)
		handleError(err)
		return string(ip)
	} else if !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(4)
	tip, err := consensus.ExternalIP()
	if err != nil {
		panic(err)
	}

	ip := tip.String()

	f, err := os.Create(ipCacheFile)
	handleError(err)
	defer f.Close()

	_, err = f.WriteString(ip)
	handleError(err)

	return ip
}

func getCachedGeolocation() (float64, float64, bool) {
	if _, err := os.Stat(geolocCacheFile); err == nil {
		fd, err := os.Open(geolocCacheFile)
		handleError(err)
		defer fd.Close()

		var lat, long float32
		err = binary.Read(fd, binary.BigEndian, &lat)
		handleError(err)
		err = binary.Read(fd, binary.BigEndian, &long)
		handleError(err)

		return float64(lat), float64(long), true

	} else if !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	return 0, 0, false
}

func getGeolocation(ip string) (float64, float64) {
	dbPath, err := xdg.SearchDataFile("IP2LOCATION-LITE-DB5.BIN")
	handleError(err)

	db, err := ip2location.OpenDB(dbPath)
	handleError(err)

	results, err := db.Get_all(ip)
	handleError(err)

	f, err := os.Create(geolocCacheFile)
	handleError(err)
	defer f.Close()

	lat := floatToByte(results.Latitude)
	long := floatToByte(results.Longitude)

	_, err = f.Write(lat)
	handleError(err)

	_, err = f.Write(long)
	handleError(err)

	return float64(results.Latitude), float64(results.Longitude)
}

func getSolarNoon(now time.Time, lat, long float64) time.Time {
	var times = suncalc.GetTimes(now, lat, long)
	//var times = suncalc.GetTimesWithObserver(now, suncalc.Observer{lat, long, 0, now.Location()})

	solarNoon := times[suncalc.SolarNoon].Value
	return solarNoon
}

func currentSolarHour(lat, long float64) int {
	var now = time.Now()

	solarNoon := getSolarNoon(now, lat, long)
	tzNoon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)

	diffNoon := solarNoon.Sub(tzNoon)
	//fmt.Println(tzNoon, " ", solarNoon)
	diffTime := now.Add(diffNoon)
	//fmt.Println(diffTime)
	//fmt.Println(diffNoon)
	//fmt.Println(hour)
	//fmt.Println(now)

	return diffTime.Hour()
}

func main() {
	var printIp = flag.BoolP("ip", "i", false, "Get your public IP")
	var printGeolocation = flag.BoolP("geoloc", "g", false, "Get geolocation coordinates for your public IP")

	flag.Parse()

	cmdMode := false
	if *printIp {
		cmdMode = true
		fmt.Println(getPublicIpAddr())
	}

	if *printGeolocation {
		cmdMode = true
		ip := getPublicIpAddr()
		lat, long := getGeolocation(ip)
		fmt.Println("Latitude: ", lat)
		fmt.Println("Longitude: ", long)
	}

	if cmdMode {
		return
	}

	hours := []rune{'🐀', '🐂', '🐅', '🐇', '🐉', '🐍', '🐎', '🐐', '🐒', '🐓', '🐕', '🐖'}

	lat, long, b := getCachedGeolocation()
	if !b {
		ip := getPublicIpAddr()
		lat, long = getGeolocation(ip)
	}

	hour := currentSolarHour(lat, long)
	idx := getAnimalIdx(hour)
	fmt.Printf("%c", hours[idx])
}
