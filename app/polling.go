package app

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/rivo/tview"
	"log"
	"os/exec"
	"pool-boy/core"
	"runtime"
	"time"
)

func StartEventPolling(appData *AppData, config core.Config, uiApplication *tview.Application, uiElements *map[string]tview.Primitive, events *[]core.PoolEvent) {
	processTicker := time.NewTicker(60 * time.Second)
	internalTicker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-internalTicker.C:
				uiApplication.QueueUpdateDraw(func() {
					updateLeftStatusView(events, uiElements)
					updateRightStatusView(appData, uiElements)
				})

			case <-processTicker.C:
				uiApplication.QueueUpdateDraw(func() {
					updatePoolEventData(appData, config, events, uiElements)
					notifyUserIfAvailableTickets(appData, events)
				})

			case <-quit:
				processTicker.Stop()
				internalTicker.Stop()
				return
			}
		}
	}()
}

func updatePoolEventData(appData *AppData, config core.Config, events *[]core.PoolEvent, uiElements *map[string]tview.Primitive) {
	tmp, err := core.ScrapeTickets(config.PoolUrl)
	if err != nil {
		updateRightStatusViewWithMessage(err.Error(), uiElements)
		return
	}

	appData.RequestCount++
	appData.LastRequest = time.Now()

	fetchedEvents := core.FilterValidEvents(tmp)
	for _, fetchedEvent := range fetchedEvents {
		for index, _ := range *events {
			(*events)[index].UpdateStatusIfSameEvent(fetchedEvent)
		}
	}
	appData.LastMessage = fmt.Sprintf("%d events fetched", len(fetchedEvents))
}

func notifyUserIfAvailableTickets(appData *AppData, events *[]core.PoolEvent) {
	var shouldSaySomething = false
	for _, event := range *events {
		if event.IsActive() && event.IsAvailable() {
			shouldSaySomething = true
			title := "PoolBoy"
			message := fmt.Sprintf("Order ticket now, Juuuuuuuunge!\n%s, %s", event.PoolTimeLabel(), event.PoolDateLabel())
			err := beeep.Notify(title, message, "")
			if err != nil {
				panic(err)
			}
		}
	}

	if appData.SpeechActive && shouldSaySomething {
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("/usr/bin/say", "-v", "Anna", appData.TextToSpeech)
			err := cmd.Run()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
		}
	}
}

func updateLeftStatusView(events *[]core.PoolEvent, uiElements *map[string]tview.Primitive) {
	targetEvents := filterTicketsToWatch(events)
	availableEvents := filterAvailableTickets(events)
	message := fmt.Sprintf(
		"%d watching, %d available",
		len(targetEvents),
		len(availableEvents),
	)
	(*uiElements)[core.LeftStatus].(*tview.TextView).SetText(message)
}

func updateRightStatusView(appData *AppData, uiElements *map[string]tview.Primitive) {
	message := fmt.Sprintf(
		"%s, last update %s, %d total",
		appData.LastMessage,
		appData.LastRequest.Format("03:04:05PM"),
		appData.RequestCount,
	)
	(*uiElements)[core.RightStatus].(*tview.TextView).SetText(message)
}

func updateRightStatusViewWithMessage(message string, uiElements *map[string]tview.Primitive) {
	(*uiElements)[core.RightStatus].(*tview.TextView).SetText(message)
}

func filterAvailableTickets(events *[]core.PoolEvent) []core.PoolEvent {
	var availableEvents []core.PoolEvent
	for _, event := range *events {
		if event.IsAvailable() {
			availableEvents = append(availableEvents, event)
		}
	}
	return availableEvents
}

func filterTicketsToWatch(events *[]core.PoolEvent) []core.PoolEvent {
	var activeEvents []core.PoolEvent
	eventCount := len(*events)
	for r := 0; r < eventCount; r++ {
		event := (*events)[r]
		if event.IsActive() {
			activeEvents = append(activeEvents, event)
		}
	}
	return activeEvents
}
