package eventbus

import (
	"reflect"
	"sync"
)

// Event 定义事件
type Event struct {
	Name    string
	Payload interface{}
}

// EventHandler 定义处理事件的处理函数类型
type EventHandler func(event Event)

// Subscriber 订阅者结构体
type Subscriber struct {
	handler EventHandler
}

// EventBus 事件总线
type EventBus struct {
	subscribers sync.Map
}

// NewEventBus 新建事件总线
func NewEventBus() *EventBus {
	return &EventBus{}
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventName string, handler EventHandler) {
	subscribers, _ := eb.subscribers.LoadOrStore(eventName, []*Subscriber{})
	subscriber := &Subscriber{handler: handler}
	eb.subscribers.Store(eventName, append(subscribers.([]*Subscriber), subscriber))
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(eventName string, handler EventHandler) {
	if subscribers, ok := eb.subscribers.Load(eventName); ok {
		oldSubscribers := subscribers.([]*Subscriber)
		newSubscribers := make([]*Subscriber, 0, len(oldSubscribers))
		for _, subscriber := range oldSubscribers {
			if reflect.ValueOf(subscriber.handler).Pointer() != reflect.ValueOf(handler).Pointer() {
				newSubscribers = append(newSubscribers, subscriber)
			}
		}
		if len(newSubscribers) == 0 {
			eb.subscribers.Delete(eventName)
		} else {
			eb.subscribers.Store(eventName, newSubscribers)
		}
	}
}

// Publish 发布事件
func (eb *EventBus) Publish(event Event) {
	wg := sync.WaitGroup{}
	if subscribers, ok := eb.subscribers.Load(event.Name); ok {
		for _, subscriber := range subscribers.([]*Subscriber) {
			wg.Add(1)
			go func(subscriber *Subscriber) {
				defer wg.Done()
				subscriber.handler(event)
			}(subscriber)
		}
	}
	wg.Wait()
}
