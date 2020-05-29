package main

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

var channelNorth = make(chan string, 1)
var channelEast = make(chan string, 1)
var channelSouth = make(chan string, 1)
var channelWest = make(chan string, 1)


func show(d CardinalDirection, c Colour){
	log.Print("show(" + string(d) + ", " + string(c) + ")")
}

func TrafficLight(d CardinalDirection, own chan string, opposite chan string, right chan string){
	//Colour c
	var c = none

	for{
		select{
			case msg := <-own:
				if msg == "p" {
					//time.Sleep(100 * time.Millisecond)
					show(d, next(c))
					if next(c) == red {
						opposite <- "b"
					}else{
						opposite <- "p"
					}
					c = next(c)
				}else if msg == "a" {
					own <- "p"
				}else if msg == "b" {
					right <- "a"
				}
		}
	}

}

func main(){
	var quitChannel = make(chan string)
	channelNorth <- "p"
	channelSouth <- "p"
	go TrafficLight(north, channelNorth, channelSouth, channelEast)
	go TrafficLight(east, channelEast, channelWest, channelSouth)
	go TrafficLight(south, channelSouth, channelNorth, channelWest)
	go TrafficLight(west, channelWest, channelEast, channelNorth)
	<-quitChannel
}
