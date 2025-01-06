package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"log"
)

var periodToAmount = map[model.Period]int{
	model.Year:    4,
	model.Quarter: 4,
	model.Week:    4,
	model.Day:     5,
}

type Panel struct {
	app             *tview.Application
	period          model.Period
	goalsRepository goalsRepository
	container       *tview.Flex
	focusLeft       func()
	focusRight      func()
	goalsList       *GoalsList
}

func newPanel(
	ctx context.Context,
	app *tview.Application,
	period model.Period,
	goalsRepository goalsRepository,
	focusLeft func(),
	focusRight func(),
) *Panel {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).SetTitle(model.PeriodName(period))

	goalsList := NewGoalsList(ctx, GoalsListProps{
		app:             app,
		period:          period,
		goalsRepository: goalsRepository,
	})

	panel := &Panel{
		app:             app,
		container:       container,
		period:          period,
		goalsRepository: goalsRepository,
		focusLeft:       focusLeft,
		focusRight:      focusRight,
		goalsList:       goalsList,
	}

	panel.setupHotkeys(ctx)

	if err := panel.render(ctx); err != nil {
		panic(err)
	}

	return panel
}

func (p *Panel) Container() tview.Primitive {
	return p.container
}

func (p *Panel) EditorInFocus() *GoalEditor {
	return p.goalsList.EditorInFocus()
}

func (p *Panel) Focus() {
	p.EditorInFocus().Focus()
}

func (p *Panel) setupHotkeys(ctx context.Context) {
	p.container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log.Printf("event: %#v\n", event)

		if event.Key() == tcell.KeyLeft && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
			p.focusLeft()
			return nil
		}

		if event.Key() == tcell.KeyRight && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
			p.focusRight()
			return nil
		}

		return event
	})
}

func (p *Panel) render(ctx context.Context) error {
	topButtons := tview.NewFlex().SetDirection(tview.FlexColumn)
	future := tview.NewButton("future")
	future.SetSelectedFunc(func() { p.goalsList.ScrollFuture(ctx) })
	topButtons.AddItem(future, 0, 1, false)
	now := tview.NewButton("↑")
	now.SetSelectedFunc(func() { p.goalsList.ScrollNow(ctx) })
	topButtons.AddItem(now, 1, 0, false)
	p.container.AddItem(topButtons, 1, 1, false)

	p.container.AddItem(p.goalsList.Primitive, 0, 1, false)

	past := tview.NewButton("past")
	past.SetSelectedFunc(func() { p.goalsList.ScrollPast(ctx) })
	p.container.AddItem(past, 1, 1, false)

	return nil
}
