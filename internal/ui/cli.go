package ui

import (
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/rivo/tview"
	"log"
	"time"
)

type CLI struct {
	app                *tview.Application
	timeNow            func() time.Time
	goalsRepository    goalsRepository
	settingsRepository settingsRepository
	container          *tview.Flex
	panels             []*PeriodPanel
	currentFocus       int
	lastEscapePress    time.Time
}

func NewCLI(
	ctx context.Context,
	timeNow func() time.Time,
	goalsRepository goalsRepository,
	settingsRepository settingsRepository,
) *CLI {
	c := &CLI{
		goalsRepository:    goalsRepository,
		settingsRepository: settingsRepository,
		timeNow:            timeNow,
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
			app:                c.app,
			timeNow:            c.timeNow,
			period:             period,
			goalsRepository:    c.goalsRepository,
			settingsRepository: c.settingsRepository,
			onFocus:            func() { c.currentFocus = n },
		})
		c.container.AddItem(panel.Primitive, 0, 1, false)
		c.panels[n] = panel
	}
}

func (c *CLI) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		log.Printf("hotkey: esc")

		now := c.timeNow()
		if now.Sub(c.lastEscapePress) < time.Second {
			c.app.Stop()
			return nil
		}

		c.lastEscapePress = now
	}

	if event.Key() == tcell.KeyCtrlC {
		return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
	}

	if event.Key() == tcell.KeyLeft && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
		log.Printf("hotkey: option shift left")
		c.focusLeft()
		return nil
	}

	if event.Key() == tcell.KeyRight && event.Modifiers()&tcell.ModShift != 0 && event.Modifiers()&tcell.ModAlt != 0 {
		log.Printf("hotkey: option shift right")
		c.focusRight()
		return nil
	}

	return event
}

func (c *CLI) focusLeft() {
	if c.currentFocus == 0 {
		return
	}

	c.panels[c.currentFocus-1].Focus()
}

func (c *CLI) focusRight() {
	if c.currentFocus == len(c.panels)-1 {
		return
	}

	c.panels[c.currentFocus+1].Focus()
}

func (c *CLI) Run() error {
	return c.app.Run()
}
