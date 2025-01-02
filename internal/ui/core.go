package ui

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"time"
)

func makePanel(name string, goals []model.Goals, onChange func(goals model.Goals)) tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)

	container.SetBorder(true).SetTitle(name)

	for _, goal := range goals {
		input := tview.NewTextArea().SetText(goal.Content, false)
		input.SetTitle(goal.Title()).SetTitle(goal.Title()).SetBorder(true)
		input.SetChangedFunc(func() {
			goal.Content = input.GetText()
			goal.Updated = time.Now()
			onChange(goal)
		})
		container.AddItem(input, 0, 1, false)
	}

	return container
}

type goalsProvider interface {
	FindByPeriod(period model.Period) ([]model.Goals, error)
	Update(ctx context.Context, goals model.Goals) error
}

func Show(ctx context.Context, goalsRepository goalsProvider) error {
	app := tview.NewApplication()

	container := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	for _, period := range model.Periods {
		goals, err := goalsRepository.FindByPeriod(period)
		if err != nil {
			return err
		}
		container.AddItem(
			makePanel(model.PeriodName(period), goals, func(goals model.Goals) {
				if err := goalsRepository.Update(ctx, goals); err != nil {
					panic(err)
				}
			}),
			0, 1, false,
		)
	}

	if err := app.SetRoot(container, true).EnableMouse(true).Run(); err != nil {
		return err
	}

	return nil
}
