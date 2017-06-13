package cfevents

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const eventUrl = "http://hottopic.apps.bogata.cf-app.com/map"

type eventHandler func(payload map[string]interface{})

type CfEvent struct {
	handler eventHandler
	appName string
}

func NewCfEvent(handler eventHandler, appName string) CfEvent {
	return CfEvent{
		handler,
		appName,
	}
}

func (e *CfEvent) Run() {
	for {
		payLoad, err := e.GetTopic()
		if err != nil {
			fmt.Printf("err getting from topic: %v", payLoad)
		} else if payLoad == nil {
			fmt.Println("nothing on the topic")
		} else {
			e.handler(payLoad)
		}
		time.Sleep(3)
	}
}

func (e *CfEvent) GetTopic() (map[string]interface{}, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", eventUrl, e.appName))

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 204 {
		fmt.Println("no content")
		return nil, nil
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var payLoad map[string]interface{}
	err = json.Unmarshal(body, &payLoad)
	if err != nil {
		return nil, err
	}

	return payLoad, nil
}
