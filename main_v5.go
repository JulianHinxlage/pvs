package main

import (
	"log"
)

//Definition der Himmelsrichtung
type CardinalDirection string
const(
	north CardinalDirection = "north"
	east = "east"
	south = "south"
	west = "west"
)

//Definition der Ampelfarbe
type Colour string
const(
	none Colour = "none"
	red = "red"
	yellow = "yellow"
	green = "green"
)

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
	for i := 0; i < 2; i++ {
		select {
		case _ = <-channel:
		case channel <- true:
		}
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
	log.Print("show(" + string(d) + ", " + string(c) + ")")
}

func TrafficLight(d CardinalDirection, phaseChannel chan bool, activeChannel chan bool){
	//Colour c
	var c = none

	if d == east || d == west {
		//west und east warten auf den ersten Richtungswechsel.
		changeActive(phaseChannel, activeChannel)
	}

	for true{
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
