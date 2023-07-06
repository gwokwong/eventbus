package eventbus

import (
	"sync"
)

// EventID 定义事件ID类型
type EventID string

// Event 定义事件接口
type Event interface {
	EventID() EventID
}

// EventHandler 定义事件处理函数类型
type EventHandler func(event Event)

// Subscription 定义订阅对象
type Subscription struct {
	eventID EventID
	id      uint64
}

// BusSubscriber 定义事件订阅者接口
type BusSubscriber interface {
	Subscribe(eventID EventID, cb EventHandler) Subscription
	Unsubscribe(sub Subscription)
}

// BusPublisher 定义事件发布者接口
type BusPublisher interface {
	Publish(event Event)
}

// Bus 定义事件总线接口
type Bus interface {
	BusSubscriber
	BusPublisher
}

// bus 实现了 Bus 接口
type bus struct {
	lock   sync.Mutex
	nextID uint64
	infos  map[EventID][]*subscriptionInfo
}

type subscriptionInfo struct {
	id uint64
	cb EventHandler
}

// New 创建一个新的事件总线
func New() Bus {
	return &bus{
		infos: make(map[EventID][]*subscriptionInfo),
	}
}

func (b *bus) Subscribe(eventID EventID, cb EventHandler) Subscription {
	b.lock.Lock()
	defer b.lock.Unlock()

	id := b.nextID
	b.nextID++

	sub := &subscriptionInfo{
		id: id,
		cb: cb,
	}
	b.infos[eventID] = append(b.infos[eventID], sub)

	return Subscription{
		eventID: eventID,
		id:      id,
	}
}

func (b *bus) Unsubscribe(subscription Subscription) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if infos, ok := b.infos[subscription.eventID]; ok {
		for idx, info := range infos {
			if info.id == subscription.id {
				infos = append(infos[:idx], infos[idx+1:]...)
				break
			}
		}

		if len(infos) == 0 {
			delete(b.infos, subscription.eventID)
		} else {
			b.infos[subscription.eventID] = infos
		}
	}
}

func (b *bus) Publish(event Event) {
	infos := b.copySubscriptions(event.EventID())

	for _, sub := range infos {
		sub.cb(event)
	}
}

func (b *bus) copySubscriptions(eventID EventID) []*subscriptionInfo {
	b.lock.Lock()
	defer b.lock.Unlock()

	infos, ok := b.infos[eventID]
	if !ok {
		return nil
	}

	// 复制订阅列表，以避免在迭代时修改原始列表
	copiedInfos := make([]*subscriptionInfo, len(infos))
	copy(copiedInfos, infos)

	return copiedInfos
}
