package main

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/storage"
	"github.com/nvbn/termonizer/internal/ui"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	storagePath := os.ExpandEnv("${HOME}/.termonizer.db")
	goalsStorage, err := storage.NewSQLite(ctx, storagePath)
	if err != nil {
		panic(err)
	}
	defer goalsStorage.Close()
	goalsRepository, err := model.NewGoalsRepository(ctx, time.Now, goalsStorage)
	if err != nil {
		panic(err)
	}
	if err = ui.Show(ctx, goalsRepository); err != nil {
		panic(err)
	}
}
