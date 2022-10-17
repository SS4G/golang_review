export GOPATH="/Users/bytedance/go:/Users/bytedance/Desktop/my_proj/golang_review"
#go build -o main_run *.go
go build -o server_run web_server_main.go
go build -o client_run web_client_main.go
if test $? -eq 0;
then
    echo "compile done"
    #./main_run
else
    echo "compile error"
fi
