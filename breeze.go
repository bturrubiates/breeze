package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	validateURL = "https://api.pushover.net/1/users/validate.json"
	pushURL     = "https://api.pushover.net/1/messages.json"
	receiptURL  = "https://api.pushover.net/1/receipts"
)

const (
	// Lowest priority does not generate any sort of notification.
	Lowest = -2 + iota

	// Low priority generates a pop-up notification but does not produce a sound.
	Low

	// Normal priority messages trigger a sound, vibration, and display a pop-up
	// notification.
	Normal

	// High priority messages produce the same notifications as normal priority,
	// but can bypass quiet-hours.
	High

	// Emergency priority messages are like high priority, but are repeated until
	// they are acknowledged by the receiver. Emergency priority messages require
	// the use of the retry and expire paramters. When a message with emergency
	// priority is posted, the Pushover API responds with a receipt that can be
	// used to query the delivery information.
	Emergency
)

const (
	// MaxTitleSize defines the maximum size of a message title as defined by the
	// Pushover API.
	MaxTitleSize = 250

	// MaxMessageLen defines the maximum length of a message as defined by the
	// Pushover API.
	MaxMessageLen = 1024

	// MaxSuppURLTitle defines the maximum length of a supplementary URL title as
	// defined by the Pushover API.
	MaxSuppURLTitle = 100

	// MaxSuppURLLen defined the maximum legnth of the supplementary URL as
	// defined by the Pushover API.
	MaxSuppURLLen = 512
)

// PushContext represents a context that can be used to send push messages. This
// means an app token and a user key.
type PushContext struct {
	AppToken string
	UserKey  string
}

// Message represents the message being sent to the push receiver. Only the
// "Message" field is required. If a device is given, it will be validated.
// Requesting emergency priority requires also providing retry and expire.
type Message struct {
	Message   string
	Title     string
	URL       string
	URLTitle  string
	Priority  int
	Retry     int
	Expire    int
	Timestamp int64
	Sound     string
	Device    string
}

// Reponse represents the response received from the Pushover API. If the
// "Status" field is set to 1, then it was successful. If the status field is
// not 1, then the "Errors" field will contain a list of errors received from
// the API. The devices are gathered when requesting validation of a device.
// The devices field will eventually be exposed as a discovery service.
type Response struct {
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Receipt string   `json:"receipt"`
	Devices []string `json:"devices"`
	Errors  []string `json:"errors"`
}

// PushError represents an error that occured while interacting with the API.
type PushError struct {
	RequestResponse Response
}

func (pushError *PushError) Error() string {
	requestKey := pushError.RequestResponse.Request
	errors := pushError.RequestResponse.Errors

	return fmt.Sprintf("Request: %s failed with errors: %s", requestKey,
		strings.Join(errors[:], ","))
}

func (pushContext *PushContext) validatePushContext(device string) (bool, error) {
	parameters := url.Values{}

	parameters.Add("token", pushContext.AppToken)
	parameters.Add("user", pushContext.UserKey)

	if device != "" {
		parameters.Add("device", device)
	}

	response, err := http.PostForm(validateURL, parameters)
	if err != nil {
		fmt.Println("Handle the error case.")
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var responseData Response
	decoder.Decode(&responseData)

	var success = responseData.Status == 1

	if !success {
		return success, &PushError{RequestResponse: responseData}
	}

	return success, nil
}

// NewPushContext is the primary interface for receiving a new PushContext
// struct that can be used to interact with the Pushover API.
func NewPushContext(appToken string, userKey string) (*PushContext, error) {
	var pushContext = new(PushContext)

	pushContext.AppToken = appToken
	pushContext.UserKey = userKey

	ok, err := pushContext.validatePushContext("")
	if !ok {
		return nil, err
	}

	return pushContext, nil
}

func main() {
}
