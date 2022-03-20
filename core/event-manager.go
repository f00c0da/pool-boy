package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	DataFile   = "pool-boy-data.json"
	ConfigFile = "pool-boy-conf.json"

	Pool1 = "https://pretix.eu/Baeder/79/"
	Pool2 = "https://pretix.eu/Baeder/15/"
)

func GetPool() {
	poolEventsUrl := getPoolNumberFromUserInput()
	ok := saveConfig(Config{PoolUrl: poolEventsUrl})
	if ok {
		fetchAndSaveData()
	}
}

func GetConfig() Config {
	config := loadConfig()
	return config
}

func GetEvents() []PoolEvent {
	events, ok := loadData()
	if !ok {
		events = fetchAndSaveData()
	}
	return FilterValidEvents(events)
}

func FilterValidEvents(events []PoolEvent) []PoolEvent {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})
	var validEvents []PoolEvent
	currentTime := time.Now()
	for _, event := range events {
		if event.Time.After(currentTime) {
			validEvents = append(validEvents, event)
		}
	}
	return validEvents
}

func saveConfig(config Config) bool {
	if config.IsUrlNotValid() {
		return false
	}
	content, err := json.MarshalIndent(config, "", "")
	if err != nil {
		fmt.Printf("error during save conf file %s\n", err.Error())
		return false
	}
	_ = ioutil.WriteFile(ConfigFile, content, 0644)
	return true
}

func fetchAndSaveData() []PoolEvent {
	config := GetConfig()
	events, err := ScrapeTickets(config.PoolUrl)
	if err != nil {
		fmt.Printf("error during scraping %s\n", err.Error())
	}
	content, _ := json.MarshalIndent(events, "", " ")
	_ = ioutil.WriteFile(DataFile, content, 0644)
	return events
}

func loadConfig() Config {
	jsonFile, err := os.Open(ConfigFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	_ = json.Unmarshal(byteValue, &config)
	if err != nil {
		_ = errors.New(fmt.Sprintf("error during load %s", ConfigFile))
	}
	_ = jsonFile.Close()
	return config
}

func loadData() ([]PoolEvent, bool) {
	jsonFile, err := os.Open(DataFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var poolEvents []PoolEvent
	_ = json.Unmarshal(byteValue, &poolEvents)
	if err != nil || len(poolEvents) == 0 {
		_ = errors.New(fmt.Sprintf("error during load %s", DataFile))
		return poolEvents, false
	}
	_ = jsonFile.Close()
	return poolEvents, true
}

func getPoolNumberFromUserInput() string {
	fmt.Printf("1 - Schwimm- und Sprunghalle im Europasportpark\n")
	fmt.Printf("2 - Wellenbad am Spreewaldplatz\n")
	fmt.Printf("Use Pool (Nr.): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	poolNumber, err := strconv.Atoi(text)
	if err != nil {
		fmt.Printf("could not read the number of the pool\n")
	}

	var poolEventsUrl string
	switch poolNumber {
	case 1:
		poolEventsUrl = Pool1
		break
	case 2:
		poolEventsUrl = Pool2
		break
	}

	return poolEventsUrl
}
