package main

import (
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/storage"
	"github.com/nvbn/termonizer/internal/ui"
	"os"
	"time"
)

func main() {
	storagePath := os.ExpandEnv("${HOME}/.termonizer.json")
	goalsStorage := storage.NewJSON(storagePath)
	goalsRepository, err := model.NewGoalsRepository(time.Now, goalsStorage)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-ticker.C
			if err := goalsRepository.Sync(); err != nil {
				panic(err)
			}
		}
	}()
	if err = ui.Show(goalsRepository); err != nil {
		panic(err)
	}
}
