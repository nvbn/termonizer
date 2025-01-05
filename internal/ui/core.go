package ui

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
)

type goalsRepository interface {
	FindForPeriod(ctx context.Context, period model.Period) ([]model.Goal, error)
	CountForPeriod(ctx context.Context, period model.Period) (int, error)
	Update(ctx context.Context, goals model.Goal) error
}

func Show(ctx context.Context, goalsRepository goalsRepository) error {
	app := tview.NewApplication()

	container := tview.NewFlex().SetDirection(tview.FlexColumn)

	for _, period := range model.Periods {
		panel := newPanel(ctx, app, period, goalsRepository)
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
