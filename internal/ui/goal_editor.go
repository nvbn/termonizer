package ui

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/utils"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
	"log"
	"strings"
	"time"
)

const goalEditorPlaceholder = `* a things to do
* a thing to achieve

--
Some notes. This is just a placeholder in some opinionated format.
`

const futureGoalEditorPlaceholder = `Goals and notes for the future`

type GoalEditorProps struct {
	app             *tview.Application
	timeNow         func() time.Time
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

	switch e.goal.CompareStart(e.timeNow()) {
	case 1:
		p.SetTitle(fmt.Sprintf("%s (future)", e.goal.FormatStart()))
		p.SetTitleColor(tcell.ColorBlue)
	case 0:
		p.SetTitle(fmt.Sprintf("%s (now)", e.goal.FormatStart()))
	case -1:
		p.SetTitle(e.goal.FormatStart())
	}

	p.SetBorder(true)
	p.SetText(e.goal.Content, false)

	if e.goal.CompareStart(e.timeNow()) == 1 {
		p.SetPlaceholder(futureGoalEditorPlaceholder)
	} else {
		p.SetPlaceholder(goalEditorPlaceholder)
	}

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

func (e *GoalEditor) handleList() bool {
	_, start, end := e.Primitive.GetSelection()
	if start != end {
		return false
	}

	content := e.Primitive.GetText()
	if content == "" {
		return false
	}

	lineStart := utils.FindLineStart(content, start)
	if lineStart == len(content) || start == lineStart || content[lineStart] != '*' {
		return false
	}

	lineEnd := utils.FindLineEnd(content, lineStart)
	lineContent := content[lineStart:lineEnd]
	if strings.TrimRight(lineContent, " \t") == "*" {
		e.Primitive.Replace(lineStart, lineEnd+1, "")
		return true
	}

	toInsert := "\n*"
	if !(start < len(content) && content[start] == ' ') {
		toInsert += " "
	}
	e.Primitive.PasteHandler()(toInsert, nil)

	return true
}

// manual ctrl+c / ctrl+v / ctrl + x
func (e *GoalEditor) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		log.Println("hotkey editor: ctrl c")
		selected, _, _ := e.Primitive.GetSelection()
		clipboard.Write(clipboard.FmtText, []byte(selected))
		return nil
	}

	if event.Key() == tcell.KeyCtrlV {
		log.Println("hotkey editor: ctrl v")
		text := clipboard.Read(clipboard.FmtText)
		e.Primitive.PasteHandler()(string(text), nil)
		return nil
	}

	if event.Key() == tcell.KeyCtrlX {
		log.Println("hotkey editor: ctrl x")
		selected, start, end := e.Primitive.GetSelection()
		e.Primitive.Replace(start, end, "")
		clipboard.Write(clipboard.FmtText, []byte(selected))
		return nil
	}

	if event.Key() == tcell.KeyCtrlA {
		log.Println("hotkey editor: ctrl A")
		e.Primitive.Select(0, len(e.Primitive.GetText()))
		return nil
	}

	if event.Key() == tcell.KeyEnter {
		log.Println("hotkey editor: enter")
		if e.handleList() {
			return nil
		}
	}

	if event.Key() == tcell.KeyEsc {
		log.Println("hotkey editor: escape")
		_, start, end := e.Primitive.GetSelection()
		if start != end {
			e.Primitive.Select(start, start)
		}
	}

	return event
}

func (e *GoalEditor) Focus() {
	e.app.SetFocus(e.Primitive)
}
