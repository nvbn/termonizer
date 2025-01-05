package ui

import (
	"context"
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
	period          model.Period
	goalsRepository goalsRepository
	container       *tview.Flex
	goalsContainer  *tview.Flex
	offset          int
}

func newPanel(ctx context.Context, period model.Period, goalsRepository goalsRepository) *Panel {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).SetTitle(model.PeriodName(period))

	panel := &Panel{
		container:       container,
		offset:          0,
		period:          period,
		goalsRepository: goalsRepository,
	}

	panel.render(ctx)

	return panel
}

func (p *Panel) scrollAfterHandler(ctx context.Context) func() {
	return func() {
		if p.offset-1 >= 0 {
			p.offset -= 1
		}
		if err := p.renderGoals(ctx); err != nil {
			panic(err)
		}
	}
}

func (p *Panel) scrollNowHandler(ctx context.Context) func() {
	return func() {
		p.offset = 0
		if err := p.renderGoals(ctx); err != nil {
			panic(err)
		}
	}
}

func (p *Panel) scrollBeforeHandler(ctx context.Context) func() {
	return func() {
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

	for _, goal := range goals {
		input := tview.NewTextArea().SetText(goal.Content, false)
		input.SetTitle(goal.Title()).SetTitle(goal.Title()).SetBorder(true)
		input.SetChangedFunc(func() {
			goal.Content = input.GetText()
			goal.Updated = time.Now()
			p.goalsRepository.Update(ctx, goal)
		})
		p.goalsContainer.AddItem(input, 0, 1, false)
	}

	return nil
}

func (p *Panel) render(ctx context.Context) error {
	topButtons := tview.NewFlex().SetDirection(tview.FlexColumn)
	after := tview.NewButton("after")
	after.SetSelectedFunc(p.scrollAfterHandler(ctx))
	topButtons.AddItem(after, 0, 1, false)
	now := tview.NewButton("↑")
	now.SetSelectedFunc(p.scrollNowHandler(ctx))
	topButtons.AddItem(now, 1, 0, false)
	p.container.AddItem(topButtons, 1, 1, false)

	p.goalsContainer = tview.NewFlex().SetDirection(tview.FlexRow)
	if err := p.renderGoals(ctx); err != nil {
		return err
	}
	p.container.AddItem(p.goalsContainer, 0, 1, false)

	before := tview.NewButton("before")
	before.SetSelectedFunc(p.scrollBeforeHandler(ctx))
	p.container.AddItem(before, 1, 1, false)

	return nil
}
