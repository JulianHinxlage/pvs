package main

import (
	"log"
	"time"
)

type CardinalDirection string
const(
	north CardinalDirection = "north"
	east = "east"
	south = "south"
	west = "west"
)

type Colour string
const(
	none Colour = "none"
	red = "red"
	yellow = "yellow"
	green = "green"
)

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

var changeActiveChannelNorth = make(chan string, 3)
var changeActiveChannelEast = make(chan string, 3)
var changeActiveChannelSouth = make(chan string, 3)
var changeActiveChannelWest = make(chan string, 3)

var changePhaseChannel1 = make(chan string, 1)
var changePhaseChannel2 = make(chan string, 1)

//zur Synchronisierung aller Ampeln
func changeActive(d CardinalDirection) {
	//bei der Synchronisierung wird eine Nachricht
	//an alle anderen Ampeln(channels) gesendet.
	//dannach wird eine Nachricht für jede andere Ampel empfangen.
	//nur wenn alle Ampeln in dies Funktion eintreten, wird die Ausführung aller
	//Ampeln fortgesetzt, da dann alle nachricht vorhanden ist.
	switch d {
	case north:
		changeActiveChannelEast <- "a"
		changeActiveChannelSouth <- "a"
		changeActiveChannelWest <- "a"
		<-changeActiveChannelNorth
		<-changeActiveChannelNorth
		<-changeActiveChannelNorth
	case east:
		changeActiveChannelNorth <- "a"
		changeActiveChannelSouth <- "a"
		changeActiveChannelWest <- "a"
		<-changeActiveChannelEast
		<-changeActiveChannelEast
		<-changeActiveChannelEast
	case south:
		changeActiveChannelNorth <- "a"
		changeActiveChannelEast <- "a"
		changeActiveChannelWest <- "a"
		<-changeActiveChannelSouth
		<-changeActiveChannelSouth
		<-changeActiveChannelSouth
	case west:
		changeActiveChannelNorth <- "a"
		changeActiveChannelEast <- "a"
		changeActiveChannelSouth <- "a"
		<-changeActiveChannelWest
		<-changeActiveChannelWest
		<-changeActiveChannelWest
	}
}

//zwei Ampeln Synchronisieren
func changePhase(d CardinalDirection){
	//es wird eine Nachricht an dei andere Ampel gesendet
	//und die Nachricht von der anderen Ampel empfangen.
	if d == east || d == north {
		changePhaseChannel1<- "a"
		<-changePhaseChannel2
	}else{
		changePhaseChannel2<- "a"
		<-changePhaseChannel1
	}
}

func show(d CardinalDirection, c Colour){
	log.Print("show(" + string(d) + ", " + string(c) + ")")
}

func TrafficLight(d CardinalDirection){
	//Colour c
	var c = none
	if d == west || d == east {
		//west und east warten auf den ersten Richtungswechsel.
		changeActive(d)
	}

	for true{
		time.Sleep(100 * time.Millisecond)
		show(d, next(c))
		if next(c) == red {
			//zwei mal die Richtung Wechseln damit die Ampel wieder an der Reihe ist.
			changeActive(d)
			changeActive(d)
		}else{
			//den Farbwechsel Synchronisieren
			changePhase(d)
		}
		//Iteration anstatt Rekrusion
		c = next(c)
	}
}

func main(){
	var quitChannel = make(chan string)
	go TrafficLight(north)
	go TrafficLight(east)
	go TrafficLight(south)
	go TrafficLight(west)
	<-quitChannel
}
