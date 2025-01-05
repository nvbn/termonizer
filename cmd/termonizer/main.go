package main

import (
	"context"
	"flag"
	"github.com/nvbn/termonizer/internal/repository"
	"github.com/nvbn/termonizer/internal/storage"
	"github.com/nvbn/termonizer/internal/ui"
	"io"
	"log"
	"os"
	"time"
)

var dbPath = flag.String("db", "${HOME}/.termonizer.db", "path to the database")
var debug = flag.String("debug", "", "debug output path")

func main() {
	flag.Parse()

	if *debug == "" {
		log.SetOutput(io.Discard)
	} else {
		f, err := os.OpenFile(*debug, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	}

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
