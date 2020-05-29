package main

import (
	"log"
	"time"
)

//Definition der Himmelsrichtung
type CardinalDirection int
const(
	north CardinalDirection = 0
	east = 1
	south = 2
	west = 3
)

//Definition der Ampelfarbe
type Colour int
const(
	none Colour = 0
	red = 1
	yellow = 2
	green = 3
)

func directionToString(d CardinalDirection) string {
	switch d {
	case north:
		return "none"
	case east:
		return "east"
	case south:
		return "south"
	case west:
		return "west"
	default:
		return ""
	}
}

func colourToString(c Colour) string {
	switch c {
	case none:
		return "none"
	case red:
		return "red"
	case yellow:
		return "yellow"
	case green:
		return "green"
	default:
		return ""
	}
}

//Definition der Phasen Reihenfolge
func next(c Colour)  Colour{
	switch c {
	case none:
		return red
	case red:
		return yellow
	case yellow:
		return green
	case green:
		return red
	default:
		return none
	}
}

//zwei goroutines werden über einen channel synchronisiert
func synchronize(channel chan bool){
	select {
	case _ = <-channel:
	case channel <- true:
	}
}

//zwei gegenüberliegende Ampeln Synchronisieren
func changePhase(phaseChannel chan bool){
	synchronize(phaseChannel)
}

//zur Synchronisierung aller Ampeln
func changeActive(phaseChannel chan bool, activeChannel chan bool) {
	synchronize(activeChannel)
	synchronize(phaseChannel)
}

func show(d CardinalDirection, c Colour){
	log.Print("show(" + directionToString(d) + ", " + colourToString(c) + ")")
}

func TrafficLight(d CardinalDirection, phaseChannel chan bool, activeChannel chan bool){
	//Colour c
	var c = none

	if d == east || d == west {
		//west und east warten auf den ersten Richtungswechsel.
		changeActive(phaseChannel, activeChannel)
	}

	for true{
		time.Sleep(time.Millisecond * 100)
		show(d, next(c))
		if next(c) == red {
			//zwei mal die Richtung Wechseln damit die Ampel wieder an der Reihe ist.
			changeActive(phaseChannel, activeChannel)
			changeActive(phaseChannel, activeChannel)
		}else{
			//den Farbwechsel Synchronisieren
			changePhase(phaseChannel)
		}
		//Iteration anstatt Rekrusion
		c = next(c)
	}
}

func main(){
	var quitChannel = make(chan bool)

	//Synchronisations Channel anlegen
	//ein channel ist immer zwischen zwei Ampeln
	var channelNorthSouth = make(chan bool)
	var channelEastWest = make(chan bool)
	var channelNorthEast = make(chan bool)
	var channelSouthWest = make(chan bool)
	go TrafficLight(north, channelNorthSouth, channelNorthEast)
	go TrafficLight(east, channelEastWest, channelNorthEast)
	go TrafficLight(south, channelNorthSouth, channelSouthWest)
	go TrafficLight(west, channelEastWest, channelSouthWest)
	<-quitChannel
}
