package main

import (
	"fmt"
)

var (
	complete = make(chan int)
)

func loop(name string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d_%s \n", i, name)

	}

	complete <- 0

}

func main22() {
	go loop(" go ")
	//loop(" main ")
	//time.Sleep(time.Second)

	var messages chan string = make(chan string)
	go func(message string) {
		messages <- message
	}("ping")

	fmt.Println(<-messages)

	//for e := range complete {
	//	fmt.Println(" complete : ", e)
	//}
	<-complete
}
