package main

import (
	"fmt"
	"time"
)

type clock struct{}

type ClockIF interface {
	RunClock(int)
}

func Clock() ClockIF {
	return &clock{}
}

func (self *clock) RunClock(runctimeSecond int) {
	secondString := "tick"
	mintString := "tock"
	hourString := "bong"
	for i := 1; i <= runctimeSecond; i++ {
		time.Sleep(1 * time.Second)
		if i >= 36000 && i%36000 == 0 {
			if secondString == "tick" {
				secondString = "quack"
			} else {
				secondString = "tick"
			}
		}
		if i%3600 == 0 {
			fmt.Println(hourString)
		} else if i%60 == 0 {
			fmt.Println(mintString)
		} else {
			fmt.Println(secondString)
		}
	}

}
