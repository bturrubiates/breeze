package breeze

import (
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

func (message *Message) AddTitle(title string) *Message {
	message.title = title
	return message
}

func (message *Message) AddURL(url string) *Message {
	message.url = url
	return message
}

func (message *Message) AddURLTitle(title string) *Message {
	message.urlTitle = title
	return message
}

func (message *Message) AddPriority(priority int) *Message {
	message.priority = priority
	return message
}

func (message *Message) AddRetry(retry int) *Message {
	message.retry = retry
	return message
}

func (message *Message) AddExpire(expire int) *Message {
	message.expire = expire
	return message
}

func (message *Message) AddTimestamp(timestamp int64) *Message {
	message.timestamp = timestamp
	return message
}

func (message *Message) AddSound(sound string) *Message {
	message.sound = sound
	return message
}

func (message *Message) AddDevice(device string) *Message {
	message.device = device
	return message
}

func NewMessage(msg string) *Message {
	var message = new(Message)

	message.message = msg
	return message
}
