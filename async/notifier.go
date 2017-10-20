package async

import (
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	// A Notification is the message forwarded by the Notifier.
	Notification struct {
		Error           string           `json:"error,omitempty"`
		UpdatedAt       time.Time        `json:"updated_at"`
		OsStats         *OsStat          `json:"os_stats"`
		ContainersStats []*ContainerStat `json:"containers_stats"`
	}

	// A Notifier notifies registred subscribers with a Notification.
	Notifier struct {
		sync.RWMutex
		subscribers map[string]chan *Notification
	}
)

// WsNotifier is the golbal instance of a Notifier for websockets.
var WsNotifier *Notifier

func init() {
	WsNotifier = &Notifier{
		subscribers: make(map[string]chan *Notification, 0),
	}
}

// Len returns the number of subscribers.
func (n *Notifier) Len() int {
	n.RLock()
	defer n.RUnlock()

	return len(n.subscribers)
}

// Notify sends a Notification to all subscribers.
func (n *Notifier) Notify(notification *Notification) {
	n.RLock()
	defer n.RUnlock()

	for _, subscriber := range n.subscribers {
		subscriber <- notification
	}
}

// Subscribe returns a channel listening on Notifier events.
func (n *Notifier) Subscribe() (string, <-chan *Notification) {
	n.Lock()
	defer n.Unlock()

	id := uuid.NewV4().String()
	ch := make(chan *Notification, 500)

	n.subscribers[id] = ch

	return id, ch
}

// UnSubscribe returns a channel listening on Notifier events.
func (n *Notifier) UnSubscribe(id string) {
	n.Lock()
	defer n.Unlock()

	ch := n.subscribers[id]
	delete(n.subscribers, id)
	close(ch)
}
