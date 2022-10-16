package web_run

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

func HandlerFunc(conn net.Conn) {
	for {
		fmt.Printf("[INFO][Server] TCP HandlerFunc Start %d\n", time.Now().Unix())
		readBuffer := make([]byte, 100)
		fmt.Printf("[INFO][Server] TCP:Server:waiting data\n")
		readBytes, readErr := conn.Read(readBuffer)
		//readbytes, readErr := ioutil.ReadAll(conn.re)
		if readErr != nil {
			fmt.Printf("[ERROR][Server] TCP:Server:read Error %v\n", readErr)
			break
		}
		// process data
		time.Sleep(time.Second)
		clientString := string(readBuffer[:readBytes])
		fmt.Printf("[INFO][Server] TCP:Server:reviced=%s\n", clientString)
		wrBytes, wrErr := conn.Write([]byte("Echo:" + clientString))
		if wrErr != nil {
			fmt.Printf("[ERROR][Server] TCP:Server:write fail\n", wrBytes)
		} else {
			fmt.Printf("[INFO][Server] TCP:Server:write back %d\n", wrBytes)
		}
		fmt.Printf("[INFO][Server] TCP HandlerFunc End %d\n", time.Now().Unix())
	}
}

func TCP_Server_Run() {
	// 解析ip地址
	tcpAddr, rsvErr := net.ResolveTCPAddr("tcp", "127.0.0.1:12335")
	if rsvErr != nil {
		fmt.Printf("[ERROR][Server] TCP resolve laddr Err\n")
	}
	// 监听 对应的端口 这里用的是laddr 因为Listen不可能监听远端的ip 端口 所以这里能用到的只有端口
	listener, lisErr := net.ListenTCP("tcp", tcpAddr)
	if lisErr != nil {
		fmt.Printf("[ERROR][Server] TCP:Server: listen Err\n")
	}

	fmt.Printf("[INFO][Server] Enter Server Loop \n")
	for {
		// 每次accept 实际上就是建立了一个新的TCP连接
		//就会生成一个conn对象 只要对着这个连接当做文件一样进行读写 即可完成TCP两端的数据交互
		conn, connErr := listener.Accept()
		fmt.Printf("[INFO][Server] Connecttion Accecpted\n")
		if connErr != nil {
			fmt.Printf("[INFO][Server] TCP:Server:connErr Error")
		}
		// 通过go原生的并发 对每个连接都使用一个go程 这样就实现了服务端的多并发。
		go HandlerFunc(conn)
	}
}

func TCP_Client_Run() {
	testDataArr := []string{"songziheng.666", "wuyu.wy01", "zhaoyanbin.zyb", "liuyang.x"}
	// 解析ip字符串 或者域名 解析为 ip:port
	tcpAddr, rsvErr := net.ResolveTCPAddr("tcp", "127.0.0.1:12335")
	if rsvErr != nil {
		fmt.Printf("[ERROR][Client] TCP: rsvErr\n")
	}
	wg := sync.WaitGroup{}
	clientNum := 5
	readWriteTime := 10
	// Step0: 并发启动 clientNum 个Client
	for j := 0; j < clientNum; j++ {
		// Step1: 每个Client 会建立一个到ServerAddr的TCP连接 然后基于该链接
		myTcpConn, dialErr := net.DialTCP("tcp", nil, tcpAddr)
		if dialErr != nil {
			fmt.Printf("[ERROR][Client] TCP: dialErr\n")
		}
		// 需要设置等待每个链接读写完成 多个client都读写完成之后才能退出
		wg.Add(1)
		go func(clientIdx int) {
			// Step2: 每个client内 和TCP连接进行 ${readWriteTime} 次读写
			for i := 0; i < readWriteTime; i++ {
				writeN, err := myTcpConn.Write([]byte(fmt.Sprintf("client[%d]=%s", clientIdx, testDataArr[rand.Intn(4)])))
				if err != nil {
					fmt.Printf("[ERROR][Client] TCP: write err\n")
				}
				fmt.Printf("[INFO][Client] Client[%d]: writeN=%d\n", clientIdx, writeN)
				readBuffer := make([]byte, 200)
				resultLen, readErr := myTcpConn.Read(readBuffer)
				if readErr != nil {
					fmt.Printf("[ERROR][Client] TCP: readErr err\n")
				}
				fmt.Printf("[INFO][Client] Client[%d] TCP: readData %s\n", clientIdx, string(readBuffer[:resultLen]))
			}
			// Step3: 读写完成后 client主动 关闭TCP连接 这时Server read时 会收到EOF ERROR 需要在server端中处理
			closeErr := myTcpConn.Close()
			if closeErr != nil {
				fmt.Printf("[ERROR][Client] Client[%d] TCP: close fail\n", clientIdx)
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
}

func UDP_Server_Run() {

}

func UDP_Client_Run() {

}
