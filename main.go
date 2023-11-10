package main

import (
	"fmt"
	"os"
	"errors"
	"github.com/glendc/go-external-ip"
	"github.com/ip2location/ip2location-go"
	"github.com/pschou/go-suncalc"
	"time"
)

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

func getPublicIpAddr() string {
	if _, err := os.Stat("/tmp/zcock_ip_cache"); err == nil {
		ip, err := os.ReadFile("/tmp/zcock_ip_cache")
		if err != nil {
			panic(err)
		}
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

	f, err := os.Create("/tmp/zcock_ip_cache")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(ip)
	if err != nil {
		panic(err)
	}

	return ip
}

func getGeolocation(ip string) (float64, float64) {
	db, err := ip2location.OpenDB("IP2LOCATION-LITE-DB5.BIN")
	if err != nil {
		panic(err)
	}
	results, err := db.Get_all(ip)
	if err != nil {
		panic(err)
	}

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
	hours := []rune{'🐀', '🐂', '🐅', '🐇', '🐉', '🐍', '🐎', '🐐', '🐒', '🐓', '🐕', '🐖'}

	ip := getPublicIpAddr() 

	lat, long := getGeolocation(ip)
	hour := currentSolarHour(lat, long)
	idx := getAnimalIdx(hour)
//	fmt.Println(idx)
	fmt.Printf("%c", hours[idx])
}
