# 常用命令
## 设置环境
- 设置gopath `export GOPATH="/Users/bytedance/Desktop/my_proj/golang_review"`
- 安装自己的包的时候需要设置 `export GO111MODULE=off` 关闭111模块模式 [关于go1111module 模式的问题](https://learnku.com/go/t/39086#449e69) 
- 简单来说就是go111module允许go get 选择仓库的特定分支来编译

## 常用命令 
- 只编译 `go build -o main_run main.go`  生成main_run 的可执行文件
- 编译&运行 `go run xx.go`
- 从远端安装 go get xxx
- 安装自己写的包 `go install draw` 需要确保 draw 文件夹位于go path的src下面

## 注意事项
- 在不使用 modules 时，需要把goland中的modules选项关掉
![img.png](img.png)

# 内容说明
这里的代码主要是学习go语言使用 都是一些基本的小样例练手用的