package events

import (
	"sync"
	"time"
)

// Load data
type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

// Operations that are executed upon event called
type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

// Event/Handler manager (register and dispatch/fire)
type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear() error
}
