package web_run

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {

}

type G55HttpHandler struct{}

func (m *G55HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTPFunc(w, r)
}

// 实际上 响应Http请求只需要两个参数 一个是response writer 另一个是 req
// 说白了就是一个请求 一个写回
func ServeHTTPFunc(w http.ResponseWriter, r *http.Request) {
	returnMsg := fmt.Sprintf("G55HttpHandler: Method[%s] RemoteAddr[%s] Host[%s] author[%s]", r.Method, r.RemoteAddr, r.Host, r.Header.Get("author"))
	byteNum, err := w.Write([]byte(returnMsg))
	if err != nil {
		fmt.Printf("[ERROR][Server] return write error\n")
	} else {
		fmt.Printf("[INFO][Server] byteNum=%d\n", byteNum)
	}
}

// http 服务端程序 默认的http只支持 精确匹配的路由
// http 实际上是对底层 TCP 服务的一个封装 因为http是应用层程序 所以在http内部已经 做了对于每个请求单独起一个go程来做的方式
// http 应该会复用底层的TCP连接
func Http_Server_Run() {
	// 方法1 注册一个实现了 HttpHandler 接口的类
	http.Handle("/hello_struct", &G55HttpHandler{})

	// 方法2: 一个普通的符合HttpFunc类型的函数  func(w http.ResponseWriter, r *http.Request)
	// 强制类型转换为 http.HandlerFunc 使用Handler注册
	http.Handle("/hello_func", http.HandlerFunc(ServeHTTPFunc))

	// 方法3: 直接使用Handle Func方法注册
	http.HandleFunc("/hello_handle_func", ServeHTTPFunc)

	// 不同的注册方式实际上是给不同的路由路径注册了不同的handle函数
	// 底层应该是 一个叫ServerMux的东西 pattern -> handleFunc
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("[ERROR][Server] httpFunc Failed\n")
	}
}

// 自定义 请求参数 包括请求方法以及 请求Header 等参数
func Http_Option_Client_Run() {
	url := "http://127.0.0.1:8080/hello_struct"
	reader := strings.NewReader(url)
	// 把new req中的reader换成一个 file 的reader 应该就可以发送整个文件了
	req, _ := http.NewRequest("GET", url, reader)
	req.Header.Add("author", "songziheng.666")
	rsp, _ := http.DefaultClient.Do(req)
	defer rsp.Body.Close()
	//rsp.Body.Read(readBuffer)
	//readByteNum, readErr
	readBodyBuffer, _ := io.ReadAll(rsp.Body)
	fmt.Printf("rsp=%s", readBodyBuffer)
}

// 默认发送Http请求的方式
func Http_Client_Run() {
	//rsp, err := http.Get("http://httpbin.org/get")
	url := "http://127.0.0.1:8080/hello_struct"
	// reader := strings.NewReader(url)
	// req := http.NewRequest("GET", url, reader)

	// rsp, err := http.Get(url)
	rsp, err := http.Get(url)

	if err != nil {
		fmt.Printf("[ERROR][Client] http get Failed\n")
	}
	defer rsp.Body.Close()
	//rsp.Body.Read(readBuffer)
	//readByteNum, readErr
	readBodyBuffer, readErr := io.ReadAll(rsp.Body)

	if readErr != nil {
		fmt.Printf("[ERROR][Client] Body read Failed err=%v\n", readErr)
	}
	fmt.Printf("read bytes=%s\n", string(readBodyBuffer))
}
