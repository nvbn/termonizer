package ui

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
)

type goalsRepository interface {
	FindByPeriod(period model.Period) ([]model.Goal, error)
	Update(ctx context.Context, goals model.Goal) error
}

func Show(ctx context.Context, goalsRepository goalsRepository) error {
	app := tview.NewApplication()

	container := tview.NewFlex().SetDirection(tview.FlexColumn)

	for _, period := range model.Periods {
		goals, err := goalsRepository.FindByPeriod(period)
		if err != nil {
			return err
		}
		panel := newPanel(period, goals, func(goals model.Goal) {
			if err := goalsRepository.Update(ctx, goals); err != nil {
				panic(err)
			}
		})
		container.AddItem(
			panel.container,
			0, 1, false,
		)
	}

	if err := app.SetRoot(container, true).EnableMouse(true).Run(); err != nil {
		return err
	}

	return nil
}
