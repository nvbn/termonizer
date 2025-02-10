package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"log"
	"time"
)

const editorsCacheSize = 256

type GoalsListProps struct {
	app                *tview.Application
	timeNow            func() time.Time
	period             model.Period
	goalsRepository    goalsRepository
	settingsRepository settingsRepository
	onFocus            func()
}

type GoalsList struct {
	GoalsListProps

	Primitive *tview.Flex

	inView       []*GoalEditor
	idToPosition map[string]int
	offset       int
	currentFocus int

	editorsCache *lru.Cache[string, *GoalEditor] // rendered editors cache to persist editor state
}

func NewGoalsList(ctx context.Context, props GoalsListProps) *GoalsList {
	editorsCache, err := lru.New[string, *GoalEditor](editorsCacheSize)
	if err != nil {
		log.Fatalf("failed to init editors cache: %v", err)
	}

	l := &GoalsList{
		GoalsListProps: props,
		editorsCache:   editorsCache,
		offset:         1,
	}

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
	if l.offset >= 1 {
		l.offset -= 1
	}

	l.render(ctx) // always re-render as it's the only way to get a new day
}

func (l *GoalsList) ScrollNow(ctx context.Context) {
	l.offset = 1
	l.currentFocus = 0
	l.render(ctx)
}

func (l *GoalsList) ScrollPast(ctx context.Context) {
	amount, err := l.goalsRepository.CountForPeriod(ctx, l.period)
	if err != nil {
		log.Fatalf("failed to count goals: %v", err)
	}

	if l.offset >= (amount - l.amountToShow()) {
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

	return goals[l.offset:min(l.offset+l.amountToShow(), len(goals))]
}

func (l *GoalsList) initPrimitive(ctx context.Context) {
	p := tview.NewFlex().SetDirection(tview.FlexRow)
	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return l.handleHotkeys(ctx, event) })
	p.SetFocusFunc(l.onFocus)
	l.Primitive = p
}

func (l *GoalsList) zoomIn(ctx context.Context) {
	amountToShow := l.amountToShow()

	if amountToShow == 1 {
		return
	}

	if amountToShow > len(l.inView) && len(l.inView) > 1 {
		amountToShow = len(l.inView) - 1
	} else {
		amountToShow -= 1
	}

	if l.currentFocus >= amountToShow {
		l.offset += 1
		l.currentFocus -= 1
	}

	if err := l.setAmountToShow(ctx, amountToShow); err != nil {
		log.Fatalf("failed to set amount to show: %v", err)
	}

	l.render(ctx)
}

func (l *GoalsList) zoomOut(ctx context.Context) {
	amountToShow := l.amountToShow()

	amountToShow += 1

	if amountToShow <= len(l.getVisibleGoals(ctx)) && l.offset > 0 {
		l.offset -= 1
		l.currentFocus += 1
	}

	if err := l.setAmountToShow(ctx, amountToShow); err != nil {
		log.Fatalf("failed to set amount to show: %v", err)
	}

	l.render(ctx)
}

func (l *GoalsList) amountToShow() int {
	return l.settingsRepository.GetAmountForPeriod(l.period)
}

func (l *GoalsList) setAmountToShow(ctx context.Context, amount int) error {
	return l.settingsRepository.SetAmountForPeriod(ctx, l.period, amount)
}

func (l *GoalsList) handleHotkeys(ctx context.Context, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyUp && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
		log.Println("hotkey: option shift up")
		l.ScrollNow(ctx)
		return nil
	}

	if event.Key() == tcell.KeyUp && event.Modifiers()&tcell.ModAlt != 0 {
		log.Println("hotkey: option up")
		l.focusFuture(ctx)
		return nil
	}

	if event.Key() == tcell.KeyDown && event.Modifiers()&tcell.ModAlt != 0 {
		log.Println("hotkey: option down")
		l.focusPast(ctx)
		return nil
	}

	// option + =
	if event.Key() == tcell.KeyRune && event.Rune() == '≠' {
		log.Println("hotkey: option +")
		l.zoomIn(ctx)
		return nil
	}

	// option + -
	if event.Key() == tcell.KeyRune && event.Rune() == '–' {
		log.Println("hotkey: option -")
		l.zoomOut(ctx)
		return nil
	}

	return event
}

func (l *GoalsList) render(ctx context.Context) {
	l.Primitive.Clear()

	goals := l.getVisibleGoals(ctx)

	nextIdToPosition := make(map[string]int)
	nextInView := make([]*GoalEditor, 0, len(goals))
	idToFocusNow := ""
	for n, goal := range goals {
		nextIdToPosition[goal.ID] = n

		var editor *GoalEditor
		if existingEditor, ok := l.editorsCache.Get(goal.ID); ok {
			editor = existingEditor
		} else {
			editor = NewGoalEditor(ctx, GoalEditorProps{
				app:             l.app,
				timeNow:         l.timeNow,
				goalsRepository: l.goalsRepository,
				goal:            goal,
				onFocus: func() {
					// could be called during the first rendering
					if pos, ok := l.idToPosition[goal.ID]; ok {
						l.currentFocus = pos
						l.onFocus()
					}
				},
			})

			l.editorsCache.Add(goal.ID, editor)
		}

		nextInView = append(nextInView, editor)
		l.Primitive.AddItem(editor.Primitive, 0, 1, false)

		if l.currentFocus == n {
			idToFocusNow = goal.ID
		}
	}

	l.app.SetFocus(nextInView[nextIdToPosition[idToFocusNow]].Primitive)

	l.inView = nextInView
	l.idToPosition = nextIdToPosition
}
