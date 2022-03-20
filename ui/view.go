package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"pool-boy/core"
)

func InitializeUI(events *[]core.PoolEvent) (*tview.Application, *map[string]tview.Primitive) {
	application := tview.NewApplication()

	uiElements := make(map[string]tview.Primitive)

	uiElements[core.LeftStatus] = createTextView("No tickets to hunt", tview.AlignLeft, tcell.ColorYellow)

	uiElements[core.RightStatus] = createTextView("No request made", tview.AlignRight, tcell.ColorYellow)

	uiElements[core.TicketsView] = createTicketTable(application)
	uiElements[core.TicketsView].(*tview.Table).SetSelectedFunc(func(row int, column int) {
		contentRowIndex := row - 1
		(*events)[contentRowIndex].ToggleActive()
	})
	uiElements[core.TicketsView].(*tview.Table).SetContent(&EventTable{events: events})

	uiElements[core.GridLayout] = tview.NewGrid().
		SetColumns(50, 0, 150).
		SetRows(0, 1). // 0:max height for table, 1:one height for the status bar
		AddItem(uiElements[core.TicketsView], 0, 0, 1, 3, 0, 0, true).
		AddItem(uiElements[core.LeftStatus], 1, 0, 1, 1, 0, 0, false).
		AddItem(uiElements[core.RightStatus], 1, 2, 1, 1, 0, 0, false)

	application.SetRoot(uiElements[core.GridLayout], true)

	return application, &uiElements
}

func createTextView(text string, align int, color tcell.Color) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(align).
		SetText(text).
		SetTextColor(color)
}

func createTicketTable(application *tview.Application) *tview.Table {
	table := tview.NewTable()
	table.SetFixed(1, 0)
	table.SetSelectable(true, false)
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			application.Stop()
		}
	})
	return table
}
