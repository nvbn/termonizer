package ui

import (
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
	period    model.Period
	goals     []model.Goal
	onChange  func(goals model.Goal)
	container *tview.Flex
	offset    int
}

func newPanel(period model.Period, goals []model.Goal, onChange func(goals model.Goal)) *Panel {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).SetTitle(model.PeriodName(period))

	panel := &Panel{
		goals:     goals,
		onChange:  onChange,
		container: container,
		offset:    0,
		period:    period,
	}

	panel.render()

	return panel
}

func (p *Panel) render() {
	p.container.Clear()

	after := tview.NewButton("after")
	after.SetSelectedFunc(func() {
		if p.offset-1 >= 0 {
			p.offset -= 1
		}
		p.render()
	})
	p.container.AddItem(after, 1, 1, false)

	goals := p.goals
	if p.offset+periodToAmount[p.period] <= len(goals) {
		goals = goals[p.offset : p.offset+periodToAmount[p.period]]
	}

	for _, goal := range goals {
		input := tview.NewTextArea().SetText(goal.Content, false)
		input.SetTitle(goal.Title()).SetTitle(goal.Title()).SetBorder(true)
		input.SetChangedFunc(func() {
			goal.Content = input.GetText()
			goal.Updated = time.Now()
			p.onChange(goal)
		})
		p.container.AddItem(input, 0, 1, false)
	}

	before := tview.NewButton("before")
	before.SetSelectedFunc(func() {
		if p.offset+1 <= (len(p.goals) - periodToAmount[p.period]) {
			p.offset += 1
		}
		p.render()
	})
	p.container.AddItem(before, 1, 1, false)
}
