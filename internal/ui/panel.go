package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"time"
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
	goalsContainer  *tview.Flex
	offset          int
	inView          []tview.Primitive
	currentFocus    int
	focusLeft       func()
	focusRight      func()
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

	panel := &Panel{
		app:             app,
		container:       container,
		offset:          0,
		period:          period,
		goalsRepository: goalsRepository,
		currentFocus:    0,
		focusLeft:       focusLeft,
		focusRight:      focusRight,
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

func (p *Panel) FocusPrimitive() tview.Primitive {
	return p.inView[p.currentFocus]
}

func (p *Panel) Focus() {
	p.app.SetFocus(p.FocusPrimitive())
}

func (p *Panel) setupHotkeys(ctx context.Context) {
	p.container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
			p.focusLeft()
			return nil
		}

		if event.Key() == tcell.KeyRight && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
			p.focusRight()
			return nil
		}

		if event.Key() == tcell.KeyUp && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
			p.currentFocus = 0
			p.scrollNow(ctx)
			return nil
		}

		if event.Key() == tcell.KeyUp && event.Modifiers()&tcell.ModAlt != 0 {
			if p.currentFocus == 0 {
				p.scrollToPast(ctx)
			} else {
				p.currentFocus -= 1
				p.Focus()
			}

			return nil
		}

		if event.Key() == tcell.KeyDown && event.Modifiers()&tcell.ModAlt != 0 {
			if p.currentFocus == len(p.inView)-1 {
				p.scrollToFuture(ctx)
			} else {
				p.currentFocus += 1
				p.Focus()
			}

			return nil
		}

		return event
	})
}

func (p *Panel) scrollToPast(ctx context.Context) {
	if p.offset-1 >= 0 {
		p.offset -= 1
	}
	if err := p.renderGoals(ctx); err != nil {
		panic(err)
	}
}

func (p *Panel) scrollNow(ctx context.Context) {
	p.offset = 0
	if err := p.renderGoals(ctx); err != nil {
		panic(err)
	}
}

func (p *Panel) scrollToFuture(ctx context.Context) {
	amount, err := p.goalsRepository.CountForPeriod(ctx, p.period)
	if err != nil {
		panic(err)
	}

	if p.offset+1 <= (amount - periodToAmount[p.period]) {
		p.offset += 1
	}
	if err := p.renderGoals(ctx); err != nil {
		panic(err)
	}
}

func (p *Panel) renderGoals(ctx context.Context) error {
	p.goalsContainer.Clear()

	goals, err := p.goalsRepository.FindForPeriod(ctx, p.period)
	if err != nil {
		return err
	}

	if p.offset+periodToAmount[p.period] <= len(goals) {
		goals = goals[p.offset : p.offset+periodToAmount[p.period]]
	}

	nextInView := make([]tview.Primitive, 0)
	for n, goal := range goals {
		input := tview.NewTextArea().SetText(goal.Content, false)
		input.SetTitle(goal.Title()).SetTitle(goal.Title()).SetBorder(true)
		input.SetChangedFunc(func() {
			goal.Content = input.GetText()
			goal.Updated = time.Now()
			p.goalsRepository.Update(ctx, goal)
		})
		input.SetFocusFunc(func() {
			p.currentFocus = n
		})
		nextInView = append(nextInView, input)
		p.goalsContainer.AddItem(input, 0, 1, false)

		if p.currentFocus == n {
			p.app.SetFocus(input)
		}
	}

	p.inView = nextInView

	return nil
}

func (p *Panel) render(ctx context.Context) error {
	topButtons := tview.NewFlex().SetDirection(tview.FlexColumn)
	future := tview.NewButton("future")
	future.SetSelectedFunc(func() { p.scrollToPast(ctx) })
	topButtons.AddItem(future, 0, 1, false)
	now := tview.NewButton("↑")
	now.SetSelectedFunc(func() { p.scrollNow(ctx) })
	topButtons.AddItem(now, 1, 0, false)
	p.container.AddItem(topButtons, 1, 1, false)

	p.goalsContainer = tview.NewFlex().SetDirection(tview.FlexRow)
	if err := p.renderGoals(ctx); err != nil {
		return err
	}
	p.container.AddItem(p.goalsContainer, 0, 1, false)

	past := tview.NewButton("past")
	past.SetSelectedFunc(func() { p.scrollToFuture(ctx) })
	p.container.AddItem(past, 1, 1, false)

	return nil
}
