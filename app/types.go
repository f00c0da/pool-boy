package app

import "time"

type AppData struct {
	StartTime    time.Time
	LastRequest  time.Time
	RequestCount int
	LastMessage  string
	SpeechActive bool
	TextToSpeech string
}
