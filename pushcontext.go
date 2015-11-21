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
	soundURL    = "https://api.pushover.net/1/sounds.json"
)

// PushContext represents a context that can be used to send push messages. This
// means an app token and a user key.
type PushContext struct {
	appToken string
	userKey  string

	// A list of sounds supported by the push context.
	SupportedSounds map[string]string

	// A list of devices supported by the push context.
	SupportedDevices []string
}

// Response represents the response received from the Pushover API. If the
// "Status" field is set to 1, then it was successful. If the status field is
// not 1, then the "Errors" field will contain a list of errors received from
// the API. The devices are gathered when requesting validation of a device.
// The devices field will eventually be exposed as a discovery service.
type Response struct {
	Status  int               `json:"status"`
	Request string            `json:"request"`
	Receipt string            `json:"receipt"`
	Devices []string          `json:"devices"`
	Errors  []string          `json:"errors"`
	Sounds  map[string]string `json:"sounds"`
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

func (pushContext *PushContext) push(pushURL string, values url.Values) (Response, error) {
	response, err := http.PostForm(pushURL, values)
	if err != nil {
		fmt.Println("Handle the error case.")
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var responseData Response
	decoder.Decode(&responseData)

	if responseData.Status != 1 {
		return Response{}, &PushError{RequestResponse: responseData}
	}

	return responseData, nil
}

func (pushContext *PushContext) addValues(values url.Values) {
	values.Add("token", pushContext.appToken)
	values.Add("user", pushContext.userKey)
}

// Push will use the Pushover API to push the given message using the given push
// context.
func (pushContext *PushContext) Push(message *Message) (Response, error) {
	parameters := url.Values{}

	err := message.validateMessage(pushContext)
	if err != nil {
		return Response{}, err
	}

	pushContext.addValues(parameters)
	message.addValues(parameters)

	return pushContext.push(pushURL, parameters)
}

func (pushContext *PushContext) validatePushContext() error {
	parameters := url.Values{}
	pushContext.addValues(parameters)

	response, err := pushContext.push(validateURL, parameters)
	if err != nil {
		return err
	}
	pushContext.SupportedDevices = response.Devices

	sounds, err := pushContext.availableSounds()
	if err != nil {
		return err
	}
	pushContext.SupportedSounds = sounds

	return err
}

func (pushContext *PushContext) get(getURL string, values url.Values) (Response, error) {
	var encodedURL *url.URL
	encodedURL, err := url.Parse(getURL)
	encodedURL.RawQuery = values.Encode()

	response, err := http.Get(encodedURL.String())
	if err != nil {
		return Response{}, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var responseData Response
	decoder.Decode(&responseData)

	return responseData, nil
}

func (pushContext *PushContext) availableSounds() (map[string]string, error) {
	parameters := url.Values{}
	pushContext.addValues(parameters)

	response, err := pushContext.get(soundURL, parameters)
	if err != nil {
		return make(map[string]string), err
	}

	return response.Sounds, nil
}

// IsValidDevice can be used to check if a requested device is valid with the
// given push context.
func (pushContext *PushContext) IsValidDevice(device string) bool {
	for _, suppDev := range pushContext.SupportedDevices {
		if suppDev == device {
			return true
		}
	}

	return false
}

// IsValidSound can be used to check if a requested sound is valid with the
// given push context.
func (pushContext *PushContext) IsValidSound(sound string) bool {
	_, ok := pushContext.SupportedSounds[sound]
	return ok
}

// NewPushContext is the primary interface for receiving a new PushContext
// struct that can be used to interact with the Pushover API.
func NewPushContext(appToken string, userKey string) (*PushContext, error) {
	var pushContext = new(PushContext)

	pushContext.appToken = appToken
	pushContext.userKey = userKey

	err := pushContext.validatePushContext()
	if err != nil {
		return nil, err
	}

	return pushContext, nil
}
