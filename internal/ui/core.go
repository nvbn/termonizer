package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"log"
)

type goalsRepository interface {
	FindForPeriod(ctx context.Context, period model.Period) ([]model.Goal, error)
	CountForPeriod(ctx context.Context, period model.Period) (int, error)
	Update(ctx context.Context, goals model.Goal) error
}

func Show(ctx context.Context, goalsRepository goalsRepository) error {
	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		}

		if event.Key() == tcell.KeyCtrlC {
			log.Printf("ctrl+c ignored")
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
		}

		return event
	})

	container := tview.NewFlex().SetDirection(tview.FlexColumn)

	panels := make([]*Panel, len(model.Periods))

	for n, period := range model.Periods {
		panel := newPanel(ctx, app, period, goalsRepository, func() {
			if n != 0 {
				panels[n-1].Focus()
			}
		}, func() {
			if n != len(model.Periods)-1 {
				panels[n+1].Focus()
			}
		})
		container.AddItem(
			panel.Container(),
			0, 1, false,
		)
		panels[n] = panel
	}

	lastPanel := panels[len(panels)-1]

	if err := app.SetRoot(container, true).
		EnableMouse(true).
		EnablePaste(true).
		SetFocus(lastPanel.FocusPrimitive()).
		Run(); err != nil {
		return err
	}

	return nil
}
