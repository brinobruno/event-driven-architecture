package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	wg.Done()
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event1          TestEvent
	event2          TestEvent
	handler1        TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	EventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.EventDispatcher = NewEventDispatcher()
	suite.event1 = TestEvent{Name: "Test 1"}
	suite.event2 = TestEvent{Name: "Test 2"}
	suite.handler1 = TestEventHandler{
		ID: 1,
	}
	suite.handler2 = TestEventHandler{
		ID: 2,
	}
	suite.handler3 = TestEventHandler{
		ID: 3,
	}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcherRegister() {
	err := suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	assert.Equal(suite.T(), &suite.handler1, suite.EventDispatcher.handlers[suite.event1.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.EventDispatcher.handlers[suite.event1.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcherRegisterWithSameHandler() {
	err := suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.NotNil(err)
	suite.Error(err, "handler already registered")
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcherClear() {
	// Event 1
	err := suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	// Event 2
	err = suite.EventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event2.GetName()]))

	err = suite.EventDispatcher.Clear()
	suite.Nil(err)
	suite.Equal(0, len(suite.EventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcherHas() {
	err := suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	assert.True(suite.T(), suite.EventDispatcher.Has(suite.event1.GetName(), &suite.handler1))
	assert.True(suite.T(), suite.EventDispatcher.Has(suite.event1.GetName(), &suite.handler2))
	assert.False(suite.T(), suite.EventDispatcher.Has(suite.event1.GetName(), &suite.handler3))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcherDispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suite.event1)
	suite.EventDispatcher.Register(suite.event1.GetName(), eh)
	suite.EventDispatcher.Dispatch(&suite.event1)

	time.Sleep(10 * time.Millisecond)
	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcherRemove() {
	// Event 1
	err := suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.EventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	// Event 2
	err = suite.EventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event2.GetName()]))

	suite.EventDispatcher.Remove(suite.event1.GetName(), &suite.handler1)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))
	suite.Equal(&suite.handler2, suite.EventDispatcher.handlers[suite.event1.GetName()][0])

	suite.EventDispatcher.Remove(suite.event1.GetName(), &suite.handler2)
	suite.Equal(0, len(suite.EventDispatcher.handlers[suite.event1.GetName()]))

	suite.EventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.Equal(0, len(suite.EventDispatcher.handlers[suite.event2.GetName()]))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
