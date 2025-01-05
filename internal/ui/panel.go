package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"log"
	"time"
)

type focus int

const (
	focusNone focus = iota
	focusPreserveOrFirst
	focusPreserveOrLast
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
	inFocus         string
}

func newPanel(ctx context.Context, app *tview.Application, period model.Period, goalsRepository goalsRepository) *Panel {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).SetTitle(model.PeriodName(period))

	panel := &Panel{
		app:             app,
		container:       container,
		offset:          0,
		period:          period,
		goalsRepository: goalsRepository,
		inFocus:         "",
	}

	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log.Printf("event: %+v\n", event)
		if event.Key() == tcell.KeyUp && event.Modifiers()&tcell.ModAlt != 0 {
			panel.scrollAfter(ctx)
			return nil
		}

		if event.Key() == tcell.KeyDown && event.Modifiers()&tcell.ModAlt != 0 {
			panel.scrollBefore(ctx)
			return nil
		}

		return event
	})

	panel.render(ctx)

	return panel
}

func (p *Panel) scrollAfter(ctx context.Context) {
	if p.offset-1 >= 0 {
		p.offset -= 1
	}
	if err := p.renderGoals(ctx, focusPreserveOrLast); err != nil {
		panic(err)
	}
}

func (p *Panel) scrollNow(ctx context.Context) {
	p.offset = 0
	if err := p.renderGoals(ctx, focusPreserveOrFirst); err != nil {
		panic(err)
	}
}

func (p *Panel) scrollBefore(ctx context.Context) {
	amount, err := p.goalsRepository.CountForPeriod(ctx, p.period)
	if err != nil {
		panic(err)
	}

	if p.offset+1 <= (amount - periodToAmount[p.period]) {
		p.offset += 1
	}
	if err := p.renderGoals(ctx, focusPreserveOrFirst); err != nil {
		panic(err)
	}
}

func (p *Panel) renderGoals(ctx context.Context, focus focus) error {
	p.goalsContainer.Clear()

	goals, err := p.goalsRepository.FindForPeriod(ctx, p.period)
	if err != nil {
		return err
	}

	if p.offset+periodToAmount[p.period] <= len(goals) {
		goals = goals[p.offset : p.offset+periodToAmount[p.period]]
	}

	alreadyFocusedById := false
	nextIdToFocus := ""
	for n, goal := range goals {
		input := tview.NewTextArea().SetText(goal.Content, false)
		input.SetTitle(goal.Title()).SetTitle(goal.Title()).SetBorder(true)
		input.SetChangedFunc(func() {
			goal.Content = input.GetText()
			goal.Updated = time.Now()
			p.goalsRepository.Update(ctx, goal)
		})
		input.SetFocusFunc(func() {
			p.inFocus = goal.ID
		})
		p.goalsContainer.AddItem(input, 0, 1, false)

		// that logic sucks
		if focus == focusPreserveOrFirst && n == 0 {
			nextIdToFocus = goal.ID
			p.app.SetFocus(input)
		} else if focus == focusPreserveOrLast && n == len(goals)-1 && !alreadyFocusedById {
			nextIdToFocus = goal.ID
			p.app.SetFocus(input)
		} else if p.inFocus == goal.ID {
			nextIdToFocus = goal.ID
			alreadyFocusedById = true
			p.app.SetFocus(input)
		}
	}

	p.inFocus = nextIdToFocus

	return nil
}

func (p *Panel) render(ctx context.Context) error {
	topButtons := tview.NewFlex().SetDirection(tview.FlexColumn)
	after := tview.NewButton("after")
	after.SetSelectedFunc(func() { p.scrollAfter(ctx) })
	topButtons.AddItem(after, 0, 1, false)
	now := tview.NewButton("↑")
	now.SetSelectedFunc(func() { p.scrollNow(ctx) })
	topButtons.AddItem(now, 1, 0, false)
	p.container.AddItem(topButtons, 1, 1, false)

	p.goalsContainer = tview.NewFlex().SetDirection(tview.FlexRow)
	if err := p.renderGoals(ctx, focusNone); err != nil {
		return err
	}
	p.container.AddItem(p.goalsContainer, 0, 1, false)

	before := tview.NewButton("before")
	before.SetSelectedFunc(func() { p.scrollBefore(ctx) })
	p.container.AddItem(before, 1, 1, false)

	return nil
}
