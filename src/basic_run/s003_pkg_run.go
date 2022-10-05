package basic_run

import (
	Draw "draw_pkg"
	"fmt"
)

// 设置GOPATH=/Users/bytedance/Desktop/my_proj/golang_review"
// 执行 go install ./draw_pkg
func PackageInstalledRun() {
	p := &Draw.Round{R: 4.3}
	fmt.Println(p.ShowRound())
}
