package cfevents

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const eventUrl = "http://hottopic.apps.bogata.cf-app.com/map"

type eventHandler func(payload map[string]interface{})

type CfEvent struct {
	handler eventHandler
	appName string
}

func NewCfEvent(handler eventHandler) CfEvent {
	appName := getAppName()

	return CfEvent{
		handler,
		appName,
	}
}

func (e *CfEvent) Run() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {})

	go func() {
		fmt.Println("listening...")
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			panic(err)
		}
	}()

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

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response while getting topic %d", res.StatusCode)
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

func getAppName() string {
	vcapApp := os.Getenv("VCAP_APPLICATION")
	if vcapApp == "" {
		fmt.Println("App name cannot be found in vcap")
		os.Exit(1)
	}

	var vcap VcapApplication

	err := json.Unmarshal([]byte(vcapApp), &vcap)
	if err != nil {

	}

	return vcap.AppName
}

type VcapApplication struct {
	AppName string `json:"application_name"`
}
