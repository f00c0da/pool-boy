package core

import (
	"fmt"
	"strings"
	"time"
)

const (
	LeftStatus  = "huntview"
	RightStatus = "lastseenview"
	TicketsView = "ticketsview"
	GridLayout  = "gridlayout"

	StatusSoon      = "soon"
	StatusReserved  = "reserved"
	StatusSoldout   = "soldout"
	StatusOver      = "over"
	StatusAvailable = "available"
)

type Config struct {
	PoolUrl string
}

func (config *Config) IsUrlNotValid() bool {
	return len(config.PoolUrl) == 0
}

type PoolEvent struct {
	Time       time.Time
	OrderTime  time.Time
	TicketLink string
	Status     string
	active     bool
}

func (event *PoolEvent) PoolDateLabel() string {
	return event.Time.Format("Monday, 02 Jan 2006")
}

func (event *PoolEvent) PoolTimeLabel() string {
	start := event.Time.Add(time.Hour * 1).Format("15:04")
	end := event.Time.Add(time.Hour * 3).Format("15:04")
	return fmt.Sprintf("%s - %s", start, end)
}

func (event *PoolEvent) OrderDateLabel() string {
	return fmt.Sprintf("from %s", event.OrderTime.Format("Monday, 02 Jan, 15:04"))
}

func (event *PoolEvent) ToggleActive() {
	event.active = !event.active
}

func (event *PoolEvent) IsActive() bool {
	return event.active
}

func (event *PoolEvent) IsAvailable() bool {
	return strings.Compare(event.Status, StatusAvailable) == 0
}

func (event *PoolEvent) UpdateStatusIfSameEvent(newEvent PoolEvent) {
	if event.Time.Equal(newEvent.Time) {
		event.Status = newEvent.Status
	}
}
