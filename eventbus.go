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

type EventHandler func(event Event)

type Subscriber struct {
	handler EventHandler
}

type EventBus struct {
	subscribers map[string][]*Subscriber
	mu          sync.RWMutex
}

// NewEventBus 新建事件
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]*Subscriber),
	}
}

// Subscribe 订阅
func (eb *EventBus) Subscribe(eventName string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if _, ok := eb.subscribers[eventName]; !ok {
		eb.subscribers[eventName] = []*Subscriber{}
	}
	subscriber := &Subscriber{handler: handler}
	eb.subscribers[eventName] = append(eb.subscribers[eventName], subscriber)
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(eventName string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if subscribers, ok := eb.subscribers[eventName]; ok {
		for i, subscriber := range subscribers {
			if reflect.ValueOf(subscriber.handler).Pointer() == reflect.ValueOf(handler).Pointer() {
				// Remove the subscriber from the subscribers list
				eb.subscribers[eventName] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
	}
}

// Publish 发布
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if subscribers, ok := eb.subscribers[event.Name]; ok {
		for _, subscriber := range subscribers {
			subscriber.handler(event)
		}
	}
}
