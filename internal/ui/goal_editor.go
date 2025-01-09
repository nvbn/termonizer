package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
	"log"
)

const goalEditorPlaceholder = `* a things to do
* a thing to achieve

--
Some notes. This is just a placeholder in some opinionated format.
`

type GoalEditorProps struct {
	app             *tview.Application
	goalsRepository goalsRepository
	goal            model.Goal
	onFocus         func()
}

type GoalEditor struct {
	GoalEditorProps

	Primitive *tview.TextArea
}

func NewGoalEditor(ctx context.Context, props GoalEditorProps) *GoalEditor {
	e := &GoalEditor{GoalEditorProps: props}
	e.initPrimitive(ctx)
	return e
}

func (e *GoalEditor) initPrimitive(ctx context.Context) {
	p := tview.NewTextArea()
	p.SetTitle(e.goal.Title())
	p.SetBorder(true)
	p.SetText(e.goal.Content, false)
	p.SetPlaceholder(goalEditorPlaceholder)
	p.SetChangedFunc(func() {
		e.goal.Content = p.GetText()
		if err := e.goalsRepository.Update(ctx, e.goal); err != nil {
			log.Fatalf("failed to update goal: %s", err)
		}
	})
	p.SetFocusFunc(e.onFocus)
	p.SetInputCapture(e.handleHotkeys)
	e.Primitive = p
}

// manual ctrl+c / ctrl+v / ctrl + x
func (e *GoalEditor) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		log.Println("hotkey: ctrl c")
		selected, _, _ := e.Primitive.GetSelection()
		clipboard.Write(clipboard.FmtText, []byte(selected))
		return nil
	}

	if event.Key() == tcell.KeyCtrlV {
		log.Println("hotkey: ctrl v")
		text := clipboard.Read(clipboard.FmtText)
		e.Primitive.PasteHandler()(string(text), nil)
		return nil
	}

	if event.Key() == tcell.KeyCtrlX {
		log.Println("hotkey: ctrl x")
		selected, start, end := e.Primitive.GetSelection()
		e.Primitive.Replace(start, end, "")
		clipboard.Write(clipboard.FmtText, []byte(selected))
		return nil
	}

	return event
}

func (e *GoalEditor) Focus() {
	e.app.SetFocus(e.Primitive)
}
