package sync_run

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

/**
	mutexRunTest
  1. 启动N个go程
  2. 每个go 程中 分两步 在flag 为false的情况下 对共享变量cnt++ 并设置flag为True 等待2ms后 将flag设置为False
  3. 有并发冲突的情况下 cnt的最终值 可能不等于N 如果不解决冲突 其他go程 设置flag=true的情况下 可能不会执行++操作就结束了
*/
func MutexRun(open_sync bool) {
	fmt.Printf("Begin: =================MutexRun() open_sync=%t =========================\n", open_sync)
	m := sync.Mutex{}
	cnt := 0
	flag := false
	// 添加 wait group
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
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
	fmt.Printf("End: =================MutexRun()=========================\n")
}

/**
  写入期间会将 临界变量 rw_val 设置为-1
  如果不加读写锁 读锁就会读取到-1这个值 如果加上就不会
  这个例子实际上写的不太好 应该比较的是读写锁和互斥锁的性能
*/
func RWMutexRun(open_sync bool) {
	fmt.Printf("Begin: =================RWMutexRun() open_sync=%t =========================\n", open_sync)
	rw_val := 0
	rwMutex := sync.RWMutex{}
	totalOperateTime := 100
	rwRecord := make([][100]int, 3)

	readerFunc := func(goIdx, callIdx int) {
		if open_sync {
			defer rwMutex.RUnlock()
			rwMutex.RLock()
		}
		rwRecord[goIdx][callIdx] = rw_val
		return
	}

	writerFunc := func(callIdx int) {
		if open_sync {
			defer rwMutex.Unlock()
			rwMutex.Lock()
		}
		tmp_rwval := rw_val
		rw_val = -1 // 中间会存在一段时间的 -1 正常加锁的情况下不会读取到这个值
		time.Sleep(time.Millisecond * 2)
		rw_val = tmp_rwval + 1
		return
	}

	for i := 0; i < totalOperateTime; i++ {
		go writerFunc(i)
		for j := 0; j < 3; j++ {
			go readerFunc(j, i)
		}
	}

	for goIdx := 0; goIdx < 3; goIdx++ {
		for callIdx := 0; callIdx < totalOperateTime; callIdx++ {
			if rwRecord[goIdx][callIdx] == -1 {
				fmt.Printf("-1 got")
			}
		}
	}
	fmt.Printf("End: =================RWMutexRun() open_sync=%t =========================\n", open_sync)
}

func SyncOnceRun(open_sync bool) {
	fmt.Printf("Begin: =================SyncOnceRun() open_sync=%t =========================\n", open_sync)
	once := sync.Once{}
	addCnt := 0
	addFunc := func() {
		addCnt++
	}
	for i := 0; i < 100; i++ {
		go func() {
			if open_sync {
				once.Do(addFunc)
			} else {
				addFunc()
			}
		}()
	}
	fmt.Printf("addCnt=%d\n", addCnt)
	fmt.Printf("End: =================SyncOnceRun() open_sync=%t =========================\n", open_sync)
}

// Cond 用于多个condition等待状态的情形
// Cond 实际上就是在一个go程 检测到条件不满足的情况下 就进入wait状态 然后就是等待 cond.Signal 或者 cond.BroadCast 唤醒
// cond.Signal 是随机唤醒一个等待的go程
// cond.BroadCast 是唤醒所有等待中的go程
// 在cond唤醒后 需要记得解锁 因为 如果不执行waid 是不会对cond.L 解锁的
func SyncCondRun() {
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	condP := sync.NewCond(&m) // *sync.Mutex 才实现了Lock接口
	cnt := 0

	readFunc := func(idx int) {
		defer wg.Done()
		// 首先上锁 在上锁状态下检查 condition是否满足
		condP.L.Lock()
		for cnt != 100 {
			fmt.Printf("proc %d cond not met wait..\n", idx)
			condP.Wait()
			fmt.Printf("proc %d cond check cond_cnt=%d\n", idx, cnt)
		}
		fmt.Printf("proc %d cond met success\n", idx)
		condP.L.Unlock()
	}
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go readFunc(j)
	}
	time.Sleep(time.Second * 2)
	cnt = 100

	// 通过 signal 逐个唤醒
	for j := 0; j < 10; j++ {
		cnt = 100
		condP.Signal()
		time.Sleep(time.Millisecond * 200)
		cnt = 0
		time.Sleep(time.Millisecond * 500)
	}
	//go readFunc()
	//go writeFunc()
	wg.Wait()

}

// 如果不打开同步开关 会报 "fatal error: concurrent map writes" 错误
// 所以这种情况下只能使用 sync.Map
func SyncMapRun(open_sync bool) {
	sMapx := sync.Map{}
	mapx := map[string]string{}

	writeFunc := func(open_sync bool, val int) {
		mkey := strconv.Itoa(val)
		mval := strconv.Itoa(val)
		if open_sync {
			sMapx.Store(mkey+"_sync", mval+"_sync")
		} else {
			mapx[mkey+"_asy"] = mval + "_asy"
		}
	}
	readFunc := func(open_sync bool, val int) {
		mkey := strconv.Itoa(val)
		//mval := strconv.Itoa(val)
		if open_sync {
			if v, ok := sMapx.Load(mkey + "_sync"); ok && v.(string) != mkey+"_sync" {
				fmt.Printf("sync_err")
			}
		} else {
			if v, ok := mapx[mkey+"_asy"]; ok && v != mkey+"_asy" {
				fmt.Printf("async_err")
			}
		}
	}
	for i := 0; i < 1000; i++ {
		go writeFunc(open_sync, i)
		go readFunc(open_sync, i)
	}
}

type Item struct {
	Buffer [10000]byte
	Bame   string
}

func NewItem() interface{} {
	return &Item{}
}

// Pool 实际上就是一个对象缓存 适合存放一些当前不用但是后面可能用的对象
// 一旦对象被put到Pool中 就处于一种可以被释放的状态
// 但是如果还没释放就Get的话 则可以打打节省创建对象的开销。
func SyncPoolRun(open_sync bool) {
	pool := sync.Pool{New: NewItem}
	start := time.Now().Unix()
	var res *Item
	for i := 0; i < 10000000; i++ {
		if open_sync {
			res = pool.Get().(*Item)
			// 需要注意的是 这里需要使用put 操作, 如果不put 回比直接new更慢
			// 但是对象一旦Put进来就处于一种可以被回收的状态 具体什么时候回收要看最终的内存问题
			pool.Put(res)
		} else {
			res = NewItem().(*Item)
		}
		if i%100000 == 0 {
			fmt.Printf("res_size=%d\n", len(res.Buffer))
		}
	}
	end := time.Now().Unix()
	fmt.Printf("time duration=%d\n", end-start)
}
