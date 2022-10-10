package sync_run

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func WithCancelCtxRun(t int) {
	// 获取全局最高context
	ctx := context.Background()
	// 通过with 函数创建下一级的context对象
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// 500 ms 后主动调用Context的cancel方法 主动取消
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	select {
	case <-ctx.Done():
		fmt.Println("testWCancel.Done:", ctx.Err())
	case e := <-time.After(time.Duration(t) * time.Millisecond):
		fmt.Println("testWCancel:", e)
	}
	return
}

// 多层级context 演示
// 层级越深 超时时间越短
// 当func0中的超时时间短于 下级函数中最短的时候 会一次性取消下级所有函数的等待
func WithCancelMultiCtxRun(t int) {
	wg := sync.WaitGroup{}
	// 获取全局最高context
	rootCtx := context.Background()
	// 通过with 函数创建下一级的context对象
	ctx0, cancel0 := context.WithCancel(rootCtx)
	ctx1, cancel1 := context.WithCancel(ctx0)
	ctx2, cancel2 := context.WithCancel(ctx1)
	ctx3, cancel3 := context.WithCancel(ctx2)

	// 500 ms 后主动调用Context的cancel方法 主动取消
	func0 := func() {
		defer wg.Done()
		defer cancel0()
		// root context  在5s 后自动取消
		time.Sleep(time.Duration(t) * time.Second)
		cancel0()
		fmt.Printf("func 0 call cancled at %d\n", time.Now().Unix())
	}

	func1 := func() {
		// 下级context通过defer调用主要是为了避免context泄漏
		defer wg.Done()
		defer cancel1()
		select {
		case <-ctx1.Done():
			fmt.Printf("func 1 cancled at %d\n", time.Now().Unix())
		case <-time.After(time.Second * 10):
			fmt.Printf("func 1 timeup at %d\n", time.Now().Unix())
		}
	}

	func2 := func() {
		// 下级context通过defer调用主要是为了避免context泄漏
		defer wg.Done()
		defer cancel2()
		select {
		case <-ctx2.Done():
			fmt.Printf("func 2 cancled at %d\n", time.Now().Unix())
		case <-time.After(time.Second * 8):
			fmt.Printf("func 2 timeup at %d\n", time.Now().Unix())
		}
	}

	func3 := func() {
		// 下级context通过defer调用主要是为了避免context泄漏
		defer wg.Done()
		defer cancel3()
		select {
		case <-ctx3.Done():
			fmt.Printf("func 3 cancled at %d\n", time.Now().Unix())
		case <-time.After(time.Second * 6):
			fmt.Printf("func 3 timeup at %d\n", time.Now().Unix())
		}
	}
	wg.Add(4)
	go func1()
	go func2()
	go func3()
	go func0()
	wg.Wait()
}

// deadline context 用time.Time 结构作为输入
func WithDeadlineCtxRun(t, timeoutOffset int) {
	// 获取全局最高context
	ctx := context.Background()
	// 通过with 函数创建下一级的context对象
	dl := time.Now().Add(time.Duration(1*t) * time.Second)
	fmt.Printf("deadline %+v\n", dl)
	// 这里的deadline 是一个绝对 时间 所以用time.Time 结构作为输入
	ctx, cancel := context.WithDeadline(ctx, dl)
	defer cancel()
	// 500 ms 后主动调用Context的cancel方法 主动取消
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	select {
	case <-ctx.Done():
		fmt.Printf("ContextDeadline.Done: %+v\n", time.Now())
		// ctx.Err 返回 ctx.Done 发生的原因 是因为超时 还是因为主动取消
		fmt.Printf("ContextDeadline.Err: %+v\n", ctx.Err())

	case e := <-time.After(time.Duration(t+timeoutOffset) * time.Second):
		fmt.Println("contextWaitTimeOut:", e)
	}
	return
}

// deadline context 用time.Time 结构作为输入
func WithTimeOutCtxRun(t, timeoutOffset, cancel_time int) {
	// 获取全局最高context
	ctx := context.Background()
	// 这里的deadline 是一个绝对 时间 所以用time.Time 结构作为输入
	ctx, cancel := context.WithTimeout(ctx, time.Duration(t)*time.Second)
	defer cancel()
	// 500 ms 后主动调用Context的cancel方法 主动取消
	go func() {
		time.Sleep(time.Duration(cancel_time) * time.Second)
		cancel()
	}()
	select {
	case <-ctx.Done():
		fmt.Printf("ContextDeadline.Done: %+v\n", time.Now())
		// ctx.Err 返回 ctx.Done 发生的原因 是因为超时 还是因为主动取消
		fmt.Printf("ContextDeadline.Err: %+v\n", ctx.Err())

	case e := <-time.After(time.Duration(t+timeoutOffset) * time.Second):
		fmt.Println("contextWaitTimeOut:", e)
	}
	return
}

/*
func testWDeadline(t int) {
	ctx := context.Background()
	dl := time.Now().Add(time.Duration(1*t) * time.Second)
	ctx, cancel := context.WithDeadline(ctx, dl)
	defer cancel()
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	select {
	case <-ctx.Done():
		fmt.Println("testWDeadline.Done:", ctx.Err())
	case e := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("testWDeadline:", e)
	}
	return
}

func testWTimeout(t int) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(t)*time.Second)
	defer cancel()

	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	select {
	case <-ctx.Done():
		fmt.Println("testWTimeout.Done:", ctx.Err())
	case e := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("testWTimeout:", e)
	}
	return
}

func main_() {
	var t = 4
	testWCancel(t)
	testWDeadline(t)
	testWTimeout(t)
}
*/
