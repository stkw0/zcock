package main

import (
	"fmt"
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

func main() {
	hours := []rune{'🐀', '🐂', '🐅', '🐇', '🐉', '🐍', '🐎', '🐐', '🐒', '🐓', '🐕', '🐖'}
	db, err := ip2location.OpenDB("IP2LOCATION-LITE-DB5.BIN")

	if err != nil {
		return
	}

	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(4)

	tip, err := consensus.ExternalIP()
	if err != nil {
		panic(err)
	}

	ip := tip.String()
	results, err := db.Get_all(ip)

	if err != nil {
		fmt.Print(err)
		return
	}

	//	fmt.Printf("latitude: %f\n", results.Latitude)
	//	fmt.Printf("longitude: %f\n", results.Longitude)

	var now = time.Now()

	// get today's sunlight times for London
	lat, long := float64(results.Latitude), float64(results.Longitude)
	//lat, long := 41.428437, 2.13821

	// get the times for today, latitude, longitude, height below or above the
	// horizon, and in timezone
	var times = suncalc.GetTimes(now, lat, long)
	// var times = suncalc.GetTimesWithObserver(now, suncalc.Observer{lat, long, 0, now.Location()})

	tzNoon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)

	solarNoon := times[suncalc.SolarNoon].Value
	diffNoon := tzNoon.Sub(solarNoon)
	//fmt.Println(tzNoon, " ", solarNoon)
	diffTime := now.Add(diffNoon)
	hour := diffTime.Hour()
	//fmt.Println(diffTime)
	//fmt.Println(diffNoon)
	//	fmt.Println(hour)
	//fmt.Println(now)

	idx := getAnimalIdx(hour)
	//	fmt.Println(idx)
	fmt.Printf("%c", hours[idx])
}
