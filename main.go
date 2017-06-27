package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	prefix        string
	eventTitle    string
	eventMessage  string
	eventTag      string
	datadogApiKey string
	datadogAppKey string
)

func init() {
	flag.StringVar(&prefix, "prefix", "", "Datadog event prefix")
	flag.StringVar(&eventTitle, "event-title", "", "Event title")
	flag.StringVar(&eventMessage, "event-message", "", "Event Message")
	flag.StringVar(&eventTag, "event-tag", "", "Event Tag")
	flag.StringVar(&datadogApiKey, "datadog-api-key", "", "Datadog API key")
	flag.StringVar(&datadogAppKey, "datadog-app-key", "", "Datadog APP key")

	flag.Parse()
}

func main() {
	if err := publishEvent(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type event struct {
	Title        string   `json:"title"`
	Text         string   `json:"text"`
	Tags         []string `json:"tags"`
	DateHappened int64    `json:"date_happened"`
}

func publishEvent() error {
	event := event{
		Title:        eventTitle,
		Text:         eventMessage,
		Tags:         []string{eventTag},
		DateHappened: time.Now().Unix(),
	}

	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(event)
	if err != nil {
		return err
	}

	fmt.Println("publishing event")
	eventsEndpoint := "https://app.datadoghq.com/api/v1/events?api_key=" + datadogApiKey + "&application_key=" + datadogAppKey
	response, err := http.Post(eventsEndpoint, "application/json", buffer)
	if err != nil {
		return fmt.Errorf("Submit event: %s", err)
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Submit event returned status code %d", response.StatusCode)
	}

	return nil
}
