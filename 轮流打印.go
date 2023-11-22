package main

import (
	"fmt"
	"time"
)

func printJi(ch1, ch2 chan int) {
	for i := 1; i <= 100; i += 2 {
		<-ch1
		fmt.Println(i)
		ch2 <- 1
	}
}
func printOu(ch1, ch2 chan int) {
	for i := 2; i <= 100; i += 2 {
		<-ch2
		fmt.Println(i)
		ch1 <- 1
	}
}
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go printJi(ch1, ch2)
	go printOu(ch1, ch2)
	ch1 <- 1
	time.Sleep(10 * time.Second)
}
