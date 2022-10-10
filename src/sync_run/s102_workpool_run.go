package sync_run

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func calcFunc(args float64) float64 {
	time.Sleep(time.Millisecond * 50)
	return math.Pow(args, 1)
}

// 注意单向通道的写法
func WorkerFunc(jobs <-chan float64, result chan<- float64, workerIdx int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("worker[%d] ready\n", workerIdx)
	// channel 可以用range 来迭代
	for jobArgs := range jobs {
		result <- calcFunc(jobArgs)
	}
	fmt.Printf("worker[%d] done\n", workerIdx)
}
func WorkPoolRun(concurrentNum int) {
	wg := sync.WaitGroup{}

	jobs := make(chan float64, 100)
	// result 必须有缓冲区 否则直接锁死了
	result := make(chan float64, 1000)

	for i := 0; i < concurrentNum; i++ {
		wg.Add(1)
		go WorkerFunc(jobs, result, i, &wg)
	}
	start := time.Now().Unix()

	// 输入参数
	for j := 0; j < 1000; j++ {
		jobs <- (float64(j) + 0)
	}

	// 关闭通道
	close(jobs)
	fmt.Printf("jobs closed\n")

	finalRes := 0.0
	fmt.Println("result_len1=", len(result))
	// 需要等待worker线程都退出了 说明结果都算完了 才可以退出
	wg.Wait()
	fmt.Println("result_len2=", len(result))

	for len(result) > 0 {
		v := <-result
		finalRes += v
	}
	end := time.Now().Unix()
	fmt.Printf("finalRes=%f, totaltimeCost=%d\n", finalRes, end-start)
}
