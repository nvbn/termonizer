package main

import (
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/ui"
	"time"
)

func main() {
	goalsRepository := model.NewGoalsRepository(time.Now)
	ui.Show(goalsRepository)
}
