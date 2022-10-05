export GOPATH="/Users/bytedance/go:/Users/bytedance/Desktop/my_proj/golang_review"
go build -o main_run *.go
if test $? -eq 0;
then
    ./main_run
else
    echo "compile error"
fi
