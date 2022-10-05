package draw_pkg

import "fmt"

type Round struct {
	R float64
}

func (p *Round) ShowRound() string {
	return fmt.Sprintf("Round: r=%f", p.R)
}
