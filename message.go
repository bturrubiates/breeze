package breeze

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	// Lowest priority does not generate any sort of notification.
	Lowest = -2 + iota

	// Low priority generates a pop-up notification but does not produce a sound.
	Low

	// Normal priority messages trigger a sound, vibration, and display a pop-up
	// notification. This is the default.
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
	minRetryTime  = 30
	maxExpireTime = 86400
)

const (
	// MaxTitleLen defines the maximum size of a message title as defined by the
	// Pushover API.
	MaxTitleLen = 250

	// MaxMessageLen defines the maximum length of a message as defined by the
	// Pushover API.
	MaxMessageLen = 1024

	// MaxSuppURLTitleLen defines the maximum length of a supplementary URL title
	// as defined by the Pushover API.
	MaxSuppURLTitleLen = 100

	// MaxSuppURLLen defined the maximum legnth of the supplementary URL as
	// defined by the Pushover API.
	MaxSuppURLLen = 512
)

const (
	// ErrMessageBlank indicates that the message was blank. The Pushover API
	// requires that the message have content.
	ErrMessageBlank = -(1 + iota)

	// ErrMessageTooLong indicates the message was longer than MaxMessageLen.
	ErrMessageTooLong

	// ErrTitleTooLong indicates the message title was longer than MaxTitleLen.
	ErrTitleTooLong

	// ErrSuppURLTitleTooLong indicates that the supplementary url title was
	// longer than MaxSuppURLTitleLen.
	ErrSuppURLTitleTooLong

	// ErrSuppURLTooLong indicates that the supplementary url was longer than
	// MaxSuppURLLen.
	ErrSuppURLTooLong

	// ErrInvalidPriority indicates that the priority is invalid.
	ErrInvalidPriority

	// ErrMissingParameter indicates that some parameter is missing. This is often
	// due to retry or expire missing on a priority message, or url missing but
	// url title is given.
	ErrMissingParameter

	// ErrRetryTimeTooShort indicates that the retry time is less than
	// MinRetryTime
	ErrRetryTimeTooShort

	// ErrExpireTimeTooLong indicates that the expire time is longer than
	// MaxExpireTime.
	ErrExpireTimeTooLong

	// ErrNoDevice indicates that a device validation failed.
	ErrNoDevice
)

// ValueError represents some error in the parameters passed in by the user.
type ValueError struct {
	What int
	Why  string
}

func (errValue *ValueError) Error() string {
	return fmt.Sprintf("%d: %s", errValue.What, errValue.Why)
}

// Message represents the message being sent to the push receiver. Only the
// "Message" field is required. If a device is given, it will be validated.
// Requesting emergency priority requires also providing retry and expire.
type Message struct {
	message   string
	title     string
	url       string
	urlTitle  string
	priority  int
	retry     int
	expire    int
	timestamp int64
	sound     string
	device    string
}

func (message *Message) addValues(values url.Values) {
	if message.message != "" {
		values.Add("message", message.message)
	}
	if message.title != "" {
		values.Add("title", message.title)
	}
	if message.url != "" {
		values.Add("url", message.url)
	}
	if message.urlTitle != "" {
		values.Add("url_title", message.urlTitle)
	}
	values.Add("priority", strconv.Itoa(message.priority))
	if message.priority == Emergency {
		values.Add("retry", strconv.Itoa(message.retry))
		values.Add("expire", strconv.Itoa(message.expire))
	}
	if message.timestamp != 0 {
		values.Add("timestamp", strconv.FormatInt(message.timestamp, 10))
	}
	if message.sound != "" {
		values.Add("sound", message.sound)
	}
	if message.device != "" {
		values.Add("device", message.device)
	}
}

// AddTitle can be used to add a title to a message. The title is limited to a
// maximum length of MaxTitleLen, and will be verified before sending.
// TODO: Verify that the length is less than MaxTitleLen.
func (message *Message) AddTitle(title string) *Message {
	message.title = title
	return message
}

// AddURL can be used to add a URL to a message. The URL is limited to a maximum
// length of MaxSuppURLLen, and will be verified before sending.
// TODO: Verify that the length is less than MaxSuppURLLen.
func (message *Message) AddURL(url string) *Message {
	message.url = url
	return message
}

// AddURLTitle can be used to add a url title. The url title is limited to a
// maximum length of MaxSuppURLTitleLen, and will be verified before sending.
// TODO: Verify that the length is less than MaxSuppURLTitleLen.
func (message *Message) AddURLTitle(title string) *Message {
	message.urlTitle = title
	return message
}

// AddPriority can be used to associate a priority with the message. The
// priority is limited to Lowest, Low, Normal, High, and Emergency. The default
// is Normal. If Emergency is set, then retry and expire must also be provided.
// TODO: Verify that retry and expire are provided.
func (message *Message) AddPriority(priority int) *Message {
	message.priority = priority
	return message
}

// AddRetry can be used to establish an associated retry time and is only
// important when the message is sent with Emergency priority. The retry time
// dictates how many seconds the API will wait before re-pushing the message.
// The value is expressed in seconds. This value must be at least 30 seconds.
// TODO: Verify retry time is less than 30 seconds.
func (message *Message) AddRetry(retry int) *Message {
	message.retry = retry
	return message
}

// AddExpire can be used to establish an associated expire time and is only
// important when the message is sent with Emergency priority. The expire time
// determines the length of the window in which the API will attempt retries.
// The value is expressed in seconds. This value must be at most 86,400 seconds.
// TODO: Verify expire time is less than 86,400.
func (message *Message) AddExpire(expire int) *Message {
	message.expire = expire
	return message
}

// AddTimestamp can be used to associate a time stamp with a message. The
// timestamp expresses the time that the Pushover API received the push request.
// The timestamp given will be reflected in the receiving application. This
// value is expressed in epoch time.
func (message *Message) AddTimestamp(timestamp int64) *Message {
	message.timestamp = timestamp
	return message
}

// AddSound can be used to associate a sound with a message. The options are:
// pushover, bike, bugle, cashregister, classical, cosmic, falling, gamelan,
// incoming, intermission, magic, mechanical, pianobar, siren, spacealarm,
// tugboat, alien, climb, persistent, echo, updown, none.
// TODO: Verify sound is one of the above.
func (message *Message) AddSound(sound string) *Message {
	message.sound = sound
	return message
}

// AddDevice can be used to filter which device the message gets sent to. This
// will be validated using the Pushover API before the message is sent.
func (message *Message) AddDevice(device string) *Message {
	message.device = device
	return message
}

// NewMessage is the primary interface for receiving a new Message struct that
// can be given to the Push function of a PushContext.
func NewMessage(msg string) *Message {
	var message = new(Message)

	message.message = msg
	return message
}
