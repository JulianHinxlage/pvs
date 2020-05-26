package main_v1

import (
	"log"
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

var changeActiveChannel1 = make(chan string, 3)
var changeActiveChannel2 = make(chan string, 3)
var changeActiveChannel3 = make(chan string, 3)
var changeActiveChannel4 = make(chan string, 3)

var changePhaseChannel1 = make(chan string, 1)
var changePhaseChannel2 = make(chan string, 1)

func changeActive(d CardinalDirection) {
	switch d {
	case north:
		changeActiveChannel2 <- "a"
		changeActiveChannel3 <- "a"
		changeActiveChannel4 <- "a"
		<-changeActiveChannel1
		<-changeActiveChannel1
		<-changeActiveChannel1
	case east:
		changeActiveChannel1 <- "a"
		changeActiveChannel3 <- "a"
		changeActiveChannel4 <- "a"
		<-changeActiveChannel2
		<-changeActiveChannel2
		<-changeActiveChannel2
	case south:
		changeActiveChannel1 <- "a"
		changeActiveChannel2 <- "a"
		changeActiveChannel4 <- "a"
		<-changeActiveChannel3
		<-changeActiveChannel3
		<-changeActiveChannel3
	case west:
		changeActiveChannel1 <- "a"
		changeActiveChannel2 <- "a"
		changeActiveChannel3 <- "a"
		<-changeActiveChannel4
		<-changeActiveChannel4
		<-changeActiveChannel4
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

func TrafficLight1(d CardinalDirection){
	if d == north || d == south {
		TrafficLight2(d, none)
	}else{
		changeActive(d)
		TrafficLight2(d, none)
	}
}

func TrafficLight2(d CardinalDirection, c Colour){
	for true{
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
	go TrafficLight1(north)
	go TrafficLight1(east)
	go TrafficLight1(south)
	go TrafficLight1(west)
	<-quitChannel
}
