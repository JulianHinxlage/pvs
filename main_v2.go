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

func changeActive(d CardinalDirection) {
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

func changePhase(d CardinalDirection){
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
	var c = none
	if d == west || d == east {
		changeActive(d)
	}

	for true{
		time.Sleep(100 * time.Millisecond)
		show(d, next(c))
		if next(c) == red {
			changeActive(d)
			changeActive(d)
		}else{
			changePhase(d)
		}
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
