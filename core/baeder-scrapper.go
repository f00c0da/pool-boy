package core

import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
	"time"
)

const (
	TargetDomain = "pretix.eu"
)

func ScrapeTickets(poolUrl string) ([]PoolEvent, error) {
	var events []PoolEvent

	collector := colly.NewCollector(
		colly.AllowedDomains(TargetDomain),
	)

	collector.OnHTML(".available, .reserved, .soon, .soldout, .over", func(element *colly.HTMLElement) {
		rawOrderTimes := element.ChildAttrs("time", "datetime")
		orderTime, _ := time.Parse(time.RFC3339, rawOrderTimes[len(rawOrderTimes)-1])
		eventTime, _ := time.Parse(time.RFC3339, element.ChildAttr("span", "data-time"))
		link := fmt.Sprintf("https://%s%s", TargetDomain, element.Attr("href"))
		status := createStatus(element.Attr("class"))

		event := PoolEvent{
			Time:       eventTime, // +1h for real opening time
			OrderTime:  orderTime, // time you can order the ticket
			TicketLink: link,
			Status:     status,
		}
		events = append(events, event)
	})

	poolEventsUrl := poolUrl
	poolEventsNextMonthUrl := fmt.Sprintf("%s?date=%s-%s", poolUrl, getCurrentYear(), getNextMonth())

	var err error
	err = collector.Visit(poolEventsUrl)
	err = collector.Visit(poolEventsNextMonthUrl)

	return events, err
}

func getNextMonth() string {
	firstDayOfMonth := calculateBeginningOfMonth(time.Now())
	_, month, _ := firstDayOfMonth.AddDate(0, 1, 0).Date() // next month
	return fmt.Sprintf("%02d", int(month))
}

func getCurrentYear() string {
	year, _, _ := time.Now().Date()
	return strconv.Itoa(year)
}

func calculateBeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func createStatus(classes string) string {
	var state string
	if strings.Contains(classes, StatusAvailable) {
		state = StatusAvailable
	}
	if strings.Contains(classes, StatusSoon) {
		state = StatusSoon
	}
	if strings.Contains(classes, StatusReserved) {
		state = StatusReserved
	}
	if strings.Contains(classes, StatusSoldout) {
		state = StatusSoldout
	}
	if strings.Contains(classes, StatusOver) {
		state = StatusOver
	}
	return state
}
