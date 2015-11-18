package breeze

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

// PushContext represents a context that can be used to send push messages. This
// means an app token and a user key.
type PushContext struct {
	appToken string
	userKey  string
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

func (pushContext *PushContext) push(url string, values url.Values) (bool, error) {
	response, err := http.PostForm(url, values)
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

func (pushContext *PushContext) addValues(values url.Values) {
	values.Add("token", pushContext.appToken)
	values.Add("user", pushContext.userKey)
}

func (pushContext *PushContext) Push(message *Message) (bool, error) {
	parameters := url.Values{}

	pushContext.addValues(parameters)
	message.addValues(parameters)

	return pushContext.push(pushURL, parameters)
}

func (pushContext *PushContext) validatePushContext(device string) (bool, error) {
	parameters := url.Values{}

	pushContext.addValues(parameters)

	if device != "" {
		parameters.Add("device", device)
	}

	return pushContext.push(validateURL, parameters)
}

// NewPushContext is the primary interface for receiving a new PushContext
// struct that can be used to interact with the Pushover API.
func NewPushContext(appToken string, userKey string) (*PushContext, error) {
	var pushContext = new(PushContext)

	pushContext.appToken = appToken
	pushContext.userKey = userKey

	ok, err := pushContext.validatePushContext("")
	if !ok {
		return nil, err
	}

	return pushContext, nil
}
