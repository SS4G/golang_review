package basic_run

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 测试channel close 的影响
// 如果channel 已经close 写入会panic
// 如果channel有缓冲 缓冲中还有数据 即使写入端关闭了channel 那么读取端在没有将缓冲区的数据读完的情况下仍然可以读取数据
// channel 读取在返回两个值的情况下 第二个返回值是表示channel 是否还有更多数据 如果返回false 说明这个channel 没有更多数据(缓冲区没有数据 且channel已经生产端被关闭)
func ChannelCommunicateWithCloseRun() {
	fmt.Printf("Begin: =================ChannelCommunicateWithCloseRun()=========================\n")
	ch := make(chan int, 5)
	wg := sync.WaitGroup{}
	producerFunc := func(waitSec time.Duration) {
		defer wg.Done()
		for j := 0; j < 3; j++ {
			for i := 0; i < 3; i++ {
				ch <- rand.Intn(10) + 100*j
			}
		}
		time.Sleep(time.Millisecond * waitSec)
		close(ch)
		fmt.Printf("producer channel closed\n")
		// wg.Done()
	}

	consumerFunc := func(readWaitSec time.Duration) {
		defer wg.Done()
		var val int
		var has_more bool
		for {
			val, has_more = <-ch
			if has_more {
				fmt.Printf("get val=%d\n", val)
				time.Sleep(time.Millisecond * readWaitSec)
			} else {
				fmt.Printf("channle colsed\n")
				break
			}
		}
		// wg.Done()
	}
	// 读取完成后再关闭channel
	fmt.Println("读取完成后再关闭channel:")
	wg.Add(2)
	go producerFunc(1000)
	go consumerFunc(1)
	wg.Wait()
	fmt.Println("读取完成后再关闭channel: Done")

	ch = make(chan int, 5)
	// 未读取完成就关闭channle
	fmt.Println("未读取完成就关闭channle:")
	wg.Add(2)
	go producerFunc(1)
	go consumerFunc(100)
	wg.Wait()
	fmt.Printf("End: =================ChannelCommunicateWithCloseRun()=========================\n")
}

// 对于无缓冲区的channel 如果channel被关闭 读取channel的item会直接返回0
func ChannelCloseRun() {
	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		close(ch)
	}()

	go func() {
		r := <-ch
		fmt.Printf("channel closed %d\n", r)
	}()

	time.Sleep(4 * time.Second)
}

// 测试select 函数
// 生产端 每个Round下 在一个随机的等待时间下 从3个channel 中随机挑选一个
// 消费端 从四个channel 中通过select 获取数据 同时设置一个3ms超时的 time.After(time.Millisecond * 3)
func SelectChannelRun() {
	fmt.Printf("Begin: =================SelectChannelRun()=========================\n")

	channels := []chan int{make(chan int, 10), make(chan int, 10), make(chan int, 10), make(chan int, 10)}
	wg := sync.WaitGroup{}
	// !注意 wait group的add操作要从外面开始
	wg.Add(2)

	producerFunc := func() {
		for i := 0; i < 4; i++ {
			channel_idx := rand.Intn(4)
			channel_val := rand.Intn(100)
			sleep_second := time.Duration(rand.Int63() % 10)
			channels[channel_idx] <- channel_val
			fmt.Printf("Round %d: channel[%d]<-%d delay %d ms\n", i, channel_idx, channel_val, sleep_second)
			time.Sleep(time.Millisecond * sleep_second)
		}
		wg.Done()
	}

	consumerFunc := func() {
		var i1, i2, i3, i4 int
		for i := 0; i < 5; i++ {
			fmt.Printf("Wait Round %d\n", i)
			select {
			case i1 = <-channels[0]:
				fmt.Printf("CH[0]=%d\n", i1)
			case i2 = <-channels[1]:
				fmt.Printf("CH[1]=%d\n", i2)
			case i3 = <-channels[2]:
				fmt.Printf("CH[2]=%d\n", i3)
			case i4 = <-channels[3]:
				fmt.Printf("CH[3]=%d\n", i4)
			case <-time.After(time.Millisecond * 3):
				fmt.Printf("timeout\n")
				//default: fmt.Printf("defaule=")
			}
		}
		wg.Done()
	}

	go producerFunc()
	go consumerFunc()
	wg.Wait()
	fmt.Printf("End: =================SelectChannelRun()=========================\n")
}
