package main

import (
	"fmt"
	"github.com/ip2location/ip2location-go"
	"github.com/pschou/go-suncalc"
	"github.com/glendc/go-external-ip"
	"time"
)

func getAnimalIdx(hour int) int {
	for i := 23; i != 1; i -= 2 {
		if hour >= i {
			return 12 - ((i + 1) / 2)
		}
	}
	
	return 0
	
}

func main() {
	hours := []rune{'🐀', '🐂', '🐅', '🐇', '🐉', '🐍', '🐎', '🐐', '🐒', '🐓', '🐕', '🐖'}
	db, err := ip2location.OpenDB("IP2LOCATION-LITE-DB5.BIN")

	if err != nil {
		return
	}

	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(4)

	 ip, err := consensus.ExternalIP()
	 if err != nil {
		 panic(err)
	 }

//	fmt.Println(ip.String())
	results, err := db.Get_all(ip.String())

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

	tzNoon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	solarNoon := times[suncalc.SolarNoon].Value
	diffNoon := tzNoon.Sub(solarNoon)
	hour := now.Add(diffNoon).Hour()
//	fmt.Println(diffNoon)
//	fmt.Println(hour)
	
	idx := getAnimalIdx(hour)
	fmt.Printf("%c", hours[idx])
}
