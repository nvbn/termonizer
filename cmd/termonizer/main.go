package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nvbn/termonizer/internal/ai"
	"github.com/nvbn/termonizer/internal/repository"
	"github.com/nvbn/termonizer/internal/storage"
	"github.com/nvbn/termonizer/internal/ui"
	"golang.design/x/clipboard"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var dbPath = flag.String("db", "${HOME}/.termonizer.db", "path to the database")
var debug = flag.String("debug", "", "debug output path")

var ollamaUrl = flag.String("ollama-url", "http://localhost:11434/api/generate", "ollama url")
var ollamaModel = flag.String("ollama-model", "llama3.2", "ollama model")

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

	sqlite, err := storage.NewSQLite(ctx, os.ExpandEnv(*dbPath))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := sqlite.Close(); err != nil {
			panic(err)
		}
	}()

	if err := sqlite.Vacuum(ctx); err != nil {
		panic(err)
	}

	goalsRepository := repository.NewGoalsRepository(time.Now, sqlite)

	settingsRepository, err := repository.NewSettings(ctx, time.Now, sqlite)
	if err != nil {
		panic(err)
	}

	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	aiClient := ai.NewClient(http.DefaultClient, *ollamaUrl, *ollamaModel)
	if err = ui.NewCLI(ctx, time.Now, goalsRepository, settingsRepository, aiClient).Run(); err != nil {
		panic(err)
	}
}
