package main

import (
	Draw "draw"
	"fmt"
)

func draw_main() {
	p := &Draw.Round{R: 4.3}
	fmt.Println(p.ShowRound())
}
