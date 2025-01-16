package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nvbn/termonizer/internal/repository"
	"github.com/nvbn/termonizer/internal/storage"
	"github.com/nvbn/termonizer/internal/ui"
	"golang.design/x/clipboard"
	"io"
	"log"
	"os"
	"time"
)

var dbPath = flag.String("db", "${HOME}/.termonizer.db", "path to the database")
var debug = flag.String("debug", "", "debug output path")

var hotkeysDoc = `
Esc Esc - exit

Navigation:
  ⌥↑	future/up goal
  ⇧⌥↑	current/first goal
  ⌥↓	past/down goal
  ⇧⌥←	longer/left period
  ⇧⌥→	shorter/right period

Zooming:
  ⌥+	zoom in / decrease the amount of visible goals
  ⌥-	zoom out / increase the amount of visible goals

Text editing:
  ⌃C	copy
  ⌃X	cut
  ⌃V	paste
  ⌃A	select all
  Esc	remove selection
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of termonizer:\n")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), hotkeysDoc)
	}

	flag.Parse()

	if *debug == "" {
		log.SetOutput(io.Discard)
	} else {
		f, err := os.OpenFile(*debug, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

	if err := goalsStorage.Cleanup(ctx); err != nil {
		panic(err)
	}

	goalsRepository := repository.NewGoalsRepository(time.Now, goalsStorage)

	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	if err = ui.NewCLI(ctx, goalsRepository).Run(); err != nil {
		panic(err)
	}
}
