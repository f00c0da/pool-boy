package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"pool-boy/core"
	"strconv"
	"strings"
	"time"
)

type EventTable struct {
	tview.TableContentReadOnly
	events *[]core.PoolEvent
}

// GetCell Return the cell at the given position or nil if there is no cell. The
//row and column arguments start at 0 and end at what GetRowCount() and
//GetColumnCount() return, minus 1.
func (eventTable *EventTable) GetCell(row, column int) *tview.TableCell {
	var cell *tview.TableCell

	if row == 0 {

		// header
		if column == 0 {
			cell = tview.NewTableCell("Time")
		} else if column == 1 {
			cell = tview.NewTableCell("Date")
		} else if column == 2 {
			cell = tview.NewTableCell("Link")
		} else if column == 3 {
			cell = tview.NewTableCell("Action").SetAlign(tview.AlignRight)
		}
		cell.SetSelectable(false).SetTextColor(tcell.ColorYellow)
	} else {

		//content
		eventIndex := row - 1
		event := (*eventTable.events)[eventIndex]
		if column == 0 {
			cell = tview.NewTableCell(event.PoolTimeLabel())
		} else if column == 1 {
			cell = tview.NewTableCell(event.PoolDateLabel())
		} else if column == 2 {
			cell = tview.NewTableCell(event.TicketLink)
		} else if column == 3 {
			cell = tview.NewTableCell(createActionColumnContent(event)).SetAlign(tview.AlignRight)
		}

		// color by action
		if event.IsAvailable() && event.IsActive() {
			cell.SetBackgroundColor(tcell.ColorDarkGreen)
		} else if event.IsAvailable() {
			cell.SetBackgroundColor(tcell.ColorSlateGray)
		}

	}
	return cell
}

// GetRowCount Return the total number of rows in the table.
func (eventTable *EventTable) GetRowCount() int {
	return len(*eventTable.events)
}

// GetColumnCount Return the total number of columns in the table.
func (eventTable *EventTable) GetColumnCount() int {
	return 4
}

func createActionColumnContent(event core.PoolEvent) string {
	var content string

	if event.IsActive() {
		content, _ = unquoteCodePoint("\\U0001F916")
	} else if event.OrderTime.After(time.Now()) {
		content = event.OrderDateLabel()
	} else {
		content = event.Status
	}

	return content
}

func unquoteCodePoint(s string) (string, error) {
	r, err := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	return string(r), err
}
