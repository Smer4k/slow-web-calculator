package main

import (
	"github.com/Smer4k/slow-web-calculator/internal/orchestrator"
)

func main() {
	o := orchestrator.NewOrchestrator()
	o.ExpressionParser("232+2+322*2/4")
}