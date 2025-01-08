package ui

import (
	"context"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
)

type PeriodPanelProps struct {
	app             *tview.Application
	period          model.Period
	goalsRepository goalsRepository
	onFocus         func()
}

type PeriodPanel struct {
	PeriodPanelProps

	Primitive *tview.Flex

	goalsList *GoalsList
}

func NewPeriodPanel(ctx context.Context, props PeriodPanelProps) *PeriodPanel {
	p := &PeriodPanel{
		PeriodPanelProps: props,
		goalsList: NewGoalsList(ctx, GoalsListProps{
			app:             props.app,
			period:          props.period,
			goalsRepository: props.goalsRepository,
			onFocus:         props.onFocus,
		}),
	}

	p.initPrimitive(ctx)

	return p
}

func (p *PeriodPanel) Focus() {
	p.goalsList.Focus()
}

func (p *PeriodPanel) PrimitiveInFocus() tview.Primitive {
	return p.goalsList.EditorInFocus().Primitive
}

func (p *PeriodPanel) makeTopButtons(ctx context.Context) tview.Primitive {
	topButtons := tview.NewFlex().SetDirection(tview.FlexColumn)

	future := tview.NewButton(" future") // white space for centering
	future.SetSelectedFunc(func() { p.goalsList.ScrollFuture(ctx) })
	topButtons.AddItem(future, 0, 1, false)

	now := tview.NewButton("↑ ") // white space for centering
	now.SetSelectedFunc(func() { p.goalsList.ScrollNow(ctx) })
	topButtons.AddItem(now, 1, 0, false)

	return topButtons
}

func (p *PeriodPanel) makeBottomButton(ctx context.Context) tview.Primitive {
	past := tview.NewButton(" past") // white space for centering
	past.SetSelectedFunc(func() { p.goalsList.ScrollPast(ctx) })
	return past
}

func (p *PeriodPanel) initPrimitive(ctx context.Context) {
	c := tview.NewFlex().SetDirection(tview.FlexRow)
	c.SetFocusFunc(p.onFocus)
	c.SetBorder(true).SetTitle(model.PeriodName(p.period))

	c.AddItem(p.makeTopButtons(ctx), 1, 1, false)
	c.AddItem(p.goalsList.Primitive, 0, 1, false)
	c.AddItem(p.makeBottomButton(ctx), 1, 1, false)

	p.Primitive = c
}
