package main

import (
	"context"
	"flag"
	"github.com/nvbn/termonizer/internal/repository"
	"github.com/nvbn/termonizer/internal/storage"
	"github.com/nvbn/termonizer/internal/ui"
	"os"
	"time"
)

var dbPath = flag.String("db", "${HOME}/.termonizer.db", "path to the database")

func main() {
	flag.Parse()

	ctx := context.Background()
	goalsStorage, err := storage.NewSQLite(ctx, os.ExpandEnv(*dbPath))
	if err != nil {
		panic(err)
	}
	defer goalsStorage.Close()
	goalsRepository := repository.NewGoalsRepository(time.Now, goalsStorage)
	if err = ui.Show(ctx, goalsRepository); err != nil {
		panic(err)
	}
}
