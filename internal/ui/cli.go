package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"log"
)

type CLI struct {
	app             *tview.Application
	goalsRepository goalsRepository
	container       *tview.Flex
	panels          []*PeriodPanel
}

func NewCLI(ctx context.Context, goalsRepository goalsRepository) *CLI {
	c := &CLI{
		goalsRepository: goalsRepository,
	}
	c.init(ctx)
	return c
}

func (c *CLI) init(ctx context.Context) {
	c.app = tview.NewApplication()
	c.app.SetInputCapture(c.handleHotkeys)
	c.render(ctx)
	c.app.SetRoot(c.container, true).
		EnableMouse(true).
		EnablePaste(true).
		SetFocus(c.panels[len(c.panels)-1].PrimitiveInFocus())
}

func (c *CLI) render(ctx context.Context) {
	c.container = tview.NewFlex().SetDirection(tview.FlexColumn)

	c.panels = make([]*PeriodPanel, len(model.Periods))

	for n, period := range model.Periods {
		panel := NewPeriodPanel(ctx, PeriodPanelProps{
			app:             c.app,
			period:          period,
			goalsRepository: c.goalsRepository,
			focusLeft:       func() { c.focusLeft(n) },
			focusRight:      func() { c.focusRight(n) },
		})
		c.container.AddItem(panel.Primitive, 0, 1, false)
		c.panels[n] = panel
	}
}

func (c *CLI) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		c.app.Stop()
		return nil
	}

	if event.Key() == tcell.KeyCtrlC {
		log.Printf("ctrl+c ignored")
		return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
	}

	return event
}

func (c *CLI) focusLeft(n int) {
	if n == 0 {
		return
	}

	c.panels[n-1].Focus()
}

func (c *CLI) focusRight(n int) {
	if n == len(c.panels)-1 {
		return
	}

	c.panels[n+1].Focus()
}

func (c *CLI) Run() error {
	return c.app.Run()
}
