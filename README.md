# 事件总线

## 安装

```
go get -u github.com/gowkwong/eventbus
```

## 用法示例

```golang
package main

import (
	"fmt"
	eventbus "github.com/gwokwong/eventbus"
)

func main() {

	eventBus := eventbus.NewEventBus()
	handler1 := func(event eventbus.Event) {
		fmt.Println("Handler 1 received event:", event)
	}
	handler2 := func(event eventbus.Event) {
		fmt.Println("Handler 2 received event:", event)
	}

	// Subscribe to events
	eventBus.Subscribe("event1", handler1)
	eventBus.Subscribe("event2", handler2)

	// Publish events
	eventBus.Publish(eventbus.Event{Name: "event1", Payload: "Hello, World!"})
	eventBus.Publish(eventbus.Event{Name: "event2", Payload: 42})

	// Unsubscribe from an event
	eventBus.Unsubscribe("event1", handler1)

	// Publish events again
	eventBus.Publish(eventbus.Event{Name: "event1", Payload: "This event should not be handled"})
	eventBus.Publish(eventbus.Event{Name: "event2", Payload: "Another event"})

}



 ```
