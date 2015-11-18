package breeze

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

var appToken = os.Getenv("APP_TOKEN")
var userKey = os.Getenv("USER_KEY")
var pc *PushContext

func TestMain(m *testing.M) {
	var err error

	pc, err = NewPushContext(appToken, userKey)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(m.Run())
}

func expectError(t *testing.T, message *Message, what int) {
	ok, err := pc.Push(message)
	if ok {
		t.Errorf("Push should have failed.")
	}

	switch err := err.(type) {
	case *ValueError:
		fmt.Println(err)
		if err.What != what {
			t.Errorf("Expected %d, but got %d", what, err.What)
		}
	default:
		t.Errorf("Expected value error, got: %s", reflect.TypeOf(err))
	}
}

func expectOK(t *testing.T, message *Message) {
	ok, err := pc.Push(message)
	if !ok {
		t.Errorf("Not ok: %s", err.Error())
	}
}

func TestMessageCreation(t *testing.T) {
	const (
		testTitle   = "This is a title"
		testMessage = "This is a message"
	)

	var message = NewMessage(testMessage)
	if message.message != testMessage {
		t.Errorf("Expected %s, but got %s", testMessage, message.message)
	}

	message.AddTitle(testTitle)
	if message.title != testTitle {
		t.Errorf("Expected %s, but got %s", testTitle, message.title)
	}
}

func TestEmptyMessage(t *testing.T) {
	msg := NewMessage("")
	expectError(t, msg, ErrMessageBlank)
}

func TestOverflowingMessage(t *testing.T) {
	msg := NewMessage(strings.Repeat("a", MaxMessageLen+1))
	expectError(t, msg, ErrMessageTooLong)
}

func TestOverflowingTitle(t *testing.T) {
	msg := NewMessage("a").AddTitle(strings.Repeat("a", MaxTitleLen+1))
	expectError(t, msg, ErrTitleTooLong)
}

func TestOverflowingSuppURL(t *testing.T) {
	msg := NewMessage("a").AddURL(strings.Repeat("a", MaxSuppURLLen+1))
	expectError(t, msg, ErrSuppURLTooLong)
}

func TestOverflowingSuppURLTitle(t *testing.T) {
	msg := NewMessage("a").
		AddURL("a").
		AddURLTitle(strings.Repeat("a", MaxSuppURLTitleLen+1))
	expectError(t, msg, ErrSuppURLTitleTooLong)
}

func TestMissingSuppURL(t *testing.T) {
	msg := NewMessage("a").AddURLTitle("b")
	expectError(t, msg, ErrMissingParameter)
}

func TestInvalidPriority(t *testing.T) {
	msg := NewMessage("a").AddPriority(5)
	expectError(t, msg, ErrInvalidPriority)

	msg.AddPriority(-4)
	expectError(t, msg, ErrInvalidPriority)
}

func TestMissingParameter(t *testing.T) {
	msg := NewMessage("a").AddPriority(Emergency)
	expectError(t, msg, ErrMissingParameter)

	msg.AddRetry(MinRetryTime + 1)
	expectError(t, msg, ErrMissingParameter)
}

func TestRetryTooShort(t *testing.T) {
	msg := NewMessage("a").
		AddPriority(Emergency).
		AddRetry(MinRetryTime - 1).
		AddExpire(MaxExpireTime - 1)
	expectError(t, msg, ErrRetryTimeTooShort)
}

func TestExpireTooLong(t *testing.T) {
	msg := NewMessage("a").
		AddPriority(Emergency).
		AddRetry(MinRetryTime + 1).
		AddExpire(MaxExpireTime + 1)
	expectError(t, msg, ErrExpireTimeTooLong)
}
