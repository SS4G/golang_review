package main

import (
	"fmt"
	"sync"
	"time"
)

/**
	mutexRunTest
  1. 启动N个go程
  2. 每个go 程中 分两步 在flag 为false的情况下 对共享变量cnt++ 并设置flag为True 等待2ms后 将flag设置为False
  3. 有并发冲突的情况下 cnt的最终值 可能不等于N 如果不解决冲突 其他go程 设置flag=true的情况下 可能不会执行++操作就结束了
*/
func mutexRun(open_sync bool) {
	m := sync.Mutex{}
	cnt := 0
	flag := false
	// 添加 wait group
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		fmt.Println(i)
		go func() {
			wg.Add(1)
			defer wg.Done()
			if open_sync {
				m.Lock()
				defer m.Unlock()
			}
			if !flag {
				cnt++
				time.Sleep(time.Millisecond * 2)
				flag = true
			}
			time.Sleep(time.Millisecond * 2)
			flag = false
		}()
	}
	wg.Wait()
	fmt.Printf("add success cnt=%d\n", cnt)
}

func rwMutexRun() {
	m := sync.RWMutex{}
	cnt := 0
	//m.Lock()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		fmt.Println(i)
		go func() {
			wg.Add(1)
			m.RLock()
			defer m.Unlock()
			defer wg.Done()
			cnt += 1
			time.Sleep(time.Millisecond * 2)
			//m.Unlock()
			//wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("final cnt=%d\n", cnt)
}
