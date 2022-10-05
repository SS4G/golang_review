package main

import (
	"fmt"
	"time"
)

// 测试 channel
func communicate() {
	str_chan := make(chan string)

	fmt.Println("write !")

	go func() {
		str_chan <- "aaaa"
		fmt.Println("write done!")
	}()

	go func() {
		fmt.Printf("read done! %s", <-str_chan)
	}()

	time.Sleep(time.Second * 3)
}

// 测试select 函数
func select_run() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)
	var i1, i2, i3, i4 int

	go func() {
		time.Sleep(time.Second * 3)
		fmt.Println("ch1 in")
		ch1 <- 1
	}()

	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("ch2 in")
		ch2 <- 1
	}()

	go func() {
		select {
		case i1 = <-ch1:
			fmt.Printf("CH1=%d", i1)
		case i2 = <-ch2:
			fmt.Printf("CH2=%d", i2)
		case i3 = <-ch3:
			fmt.Printf("CH3=%d", i3)
		case i4 = <-ch4:
			fmt.Printf("CH4=%d", i4)
		case <-time.After(time.Second * 5):
			fmt.Printf("timeout")
			//default: fmt.Printf("defaule=")
		}
	}()

	time.Sleep(time.Second * 4)
}

func channel2_run() {

}
