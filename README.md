# 事件总线

## 安装

```
go get -u github.com/gowkwong/eventbus
```

## 用法示例

```golang


func main() {
	eventBus := NewEventBus()
	handler1 := func(event Event) {
		fmt.Println("Handler 1 received event:", event)
	}
	handler2 := func(event Event) {
		fmt.Println("Handler 2 received event:", event)
	}

	// Subscribe to events
	eventBus.Subscribe("event1", handler1)
	eventBus.Subscribe("event2", handler2)

	// Publish events
	eventBus.Publish(Event{Name: "event1", Payload: "Hello, World!"})
	eventBus.Publish(Event{Name: "event2", Payload: 42})

	// Unsubscribe from an event
	eventBus.Unsubscribe("event1", handler1)

	// Publish events again
	eventBus.Publish(Event{Name: "event1", Payload: "This event should not be handled"})
	eventBus.Publish(Event{Name: "event2", Payload: "Another event"})
}


 ```
