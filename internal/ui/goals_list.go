package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"log"
)

// TODO: make configurable
var periodToAmount = map[model.Period]int{
	model.Year:    4,
	model.Quarter: 4,
	model.Week:    4,
	model.Day:     5,
}

type GoalsListProps struct {
	app             *tview.Application
	period          model.Period
	goalsRepository goalsRepository
}

type GoalsList struct {
	GoalsListProps

	Primitive *tview.Flex

	inView       []*GoalEditor
	offset       int
	currentFocus int
}

func NewGoalsList(ctx context.Context, props GoalsListProps) *GoalsList {
	l := &GoalsList{GoalsListProps: props}

	l.initPrimitive(ctx)
	l.render(ctx)

	return l
}

func (l *GoalsList) EditorInFocus() *GoalEditor {
	return l.inView[l.currentFocus]
}

func (l *GoalsList) Focus() {
	l.EditorInFocus().Focus()
}

func (l *GoalsList) ScrollFuture(ctx context.Context) {
	if l.offset < 1 {
		return
	}

	l.offset -= 1
	l.render(ctx)
}

func (l *GoalsList) ScrollNow(ctx context.Context) {
	l.offset = 0
	l.currentFocus = 0
	l.render(ctx)
}

func (l *GoalsList) ScrollPast(ctx context.Context) {
	amount, err := l.goalsRepository.CountForPeriod(ctx, l.period)
	if err != nil {
		log.Fatalf("failed to count goals: %v", err)
	}

	if l.offset+1 == (amount - periodToAmount[l.period]) {
		return
	}

	l.offset += 1
	l.render(ctx)
}

func (l *GoalsList) focusFuture(ctx context.Context) {
	if l.currentFocus == 0 {
		l.ScrollFuture(ctx)
	} else {
		l.currentFocus -= 1
		l.Focus()
	}
}

func (l *GoalsList) focusPast(ctx context.Context) {
	if l.currentFocus == len(l.inView)-1 {
		l.ScrollPast(ctx)
	} else {
		l.currentFocus += 1
		l.Focus()
	}
}

func (l *GoalsList) getVisibleGoals(ctx context.Context) []model.Goal {
	goals, err := l.goalsRepository.FindForPeriod(ctx, l.period)
	if err != nil {
		log.Fatalf("failed to find goals: %v", err)
	}

	if l.offset+periodToAmount[l.period] <= len(goals) {
		return goals[l.offset : l.offset+periodToAmount[l.period]]
	}

	return goals
}

func (l *GoalsList) initPrimitive(ctx context.Context) {
	p := tview.NewFlex().SetDirection(tview.FlexRow)
	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return l.handleHotkeys(ctx, event) })
	l.Primitive = p
}

func (l *GoalsList) handleHotkeys(ctx context.Context, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyUp && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
		l.ScrollNow(ctx)
		return nil
	}

	if event.Key() == tcell.KeyUp && event.Modifiers()&tcell.ModAlt != 0 {
		l.focusFuture(ctx)
		return nil
	}

	if event.Key() == tcell.KeyDown && event.Modifiers()&tcell.ModAlt != 0 {
		l.focusPast(ctx)
		return nil
	}

	return event
}

func (l *GoalsList) render(ctx context.Context) {
	l.Primitive.Clear()

	goals := l.getVisibleGoals(ctx)

	nextInView := make([]*GoalEditor, 0, len(goals))
	for n, goal := range goals {
		editor := NewGoalEditor(ctx, GoalEditorProps{
			app:             l.app,
			goalsRepository: l.goalsRepository,
			goal:            goal,
			onFocus: func() {
				l.currentFocus = n
			},
		})

		nextInView = append(nextInView, editor)
		l.Primitive.AddItem(editor.Primitive, 0, 1, false)

		if l.currentFocus == n {
			l.app.SetFocus(editor.Primitive)
		}
	}
	l.inView = nextInView
}
