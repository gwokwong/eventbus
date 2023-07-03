package eventBus

import (
	"sync"
)

type EventID string

type Event interface {
	EventID() EventID
}

type EventHandler func(event Event)

type Subscription struct {
	eventID EventID
	id      uint64
}

type BusSubscriber interface {
	Subscribe(eventID EventID, cb EventHandler) Subscription
	Unsubscribe(id Subscription)
}

type BusPublisher interface {
	Publish(event Event)
}

type Bus interface {
	BusSubscriber
	BusPublisher
}

func New() Bus {
	b := &bus{
		infos: make(map[EventID]subscriptionInfoList),
	}
	return b
}

type subscriptionInfo struct {
	id uint64
	cb EventHandler
}

type subscriptionInfoList []*subscriptionInfo

type bus struct {
	lock   sync.Mutex
	nextID uint64
	infos  map[EventID]subscriptionInfoList
}

func (bus *bus) Subscribe(eventID EventID, cb EventHandler) Subscription {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	id := bus.nextID
	bus.nextID++
	sub := &subscriptionInfo{
		id: id,
		cb: cb,
	}
	bus.infos[eventID] = append(bus.infos[eventID], sub)
	return Subscription{
		eventID: eventID,
		id:      id,
	}
}

func (bus *bus) Unsubscribe(subscription Subscription) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if infos, ok := bus.infos[subscription.eventID]; ok {
		for idx, info := range infos {
			if info.id == subscription.id {
				infos = append(infos[:idx], infos[idx+1:]...)
				break
			}
		}
		if len(infos) == 0 {
			delete(bus.infos, subscription.eventID)
		} else {
			bus.infos[subscription.eventID] = infos
		}
	}
}

func (bus *bus) Publish(event Event) {
	infos := bus.copySubscriptions(event.EventID())
	for _, sub := range infos {
		sub.cb(event)
	}
}

func (bus *bus) copySubscriptions(eventID EventID) subscriptionInfoList {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if infos, ok := bus.infos[eventID]; ok {
		return infos
	}
	return subscriptionInfoList{}
}
