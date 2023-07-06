package main

import (
	eventBus "eventbus/eventbus"
	"fmt"
)

type testEvent struct {
	duration int
}

func (e *testEvent) EventID() eventBus.EventID {
	return testEventId
}

const (
	testEventId = eventBus.EventID("test_event")
)

func main() {
	bus := eventBus.New()
	duration := 10000
	id := bus.Subscribe(testEventId, func(e eventBus.Event) {
		se := e.(*testEvent)
		fmt.Println("收到的推送消息-------->")
		fmt.Println(se.duration)
	})
	bus.Publish(&testEvent{
		duration: duration,
	})
	bus.Unsubscribe(id)
}
