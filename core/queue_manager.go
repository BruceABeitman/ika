package core

import (
	"errors"
	"fmt"
)

// ErrorQueueUninitialized is thrown when we try to access a queue that's uninitialized
var ErrorQueueUninitialized = errors.New("Queue is uninitialized")

// QueueManager service manages all the priority queues, and handles
// serving/routing requests from/to queues
// fetching, building, and re-building proxy priority queues
type QueueManager interface {
	GetProxy(key string) (*Proxy, error)
	UpdateProxyMeta(key string, meta ProxyMeta) error
	RefreshProxies(key string) error
}

// QueueManagerImpl struct
type QueueManagerImpl struct {
	queues map[string]*Queue
	store  Store
}

// NewQueueManager constructor
func NewQueueManager(store Store) QueueManager {
	return &QueueManagerImpl{
		queues: make(map[string]*Queue),
		store:  store,
	}
}

// GetProxy resolves a proxy from a given channel & domain
func (qm QueueManagerImpl) GetProxy(key string) (*Proxy, error) {
	q, ok := qm.queues[key]
	// If corresponding queue does not exist, refresh queue
	if !ok {
		if err := qm.RefreshProxies(key); err != nil {
			return nil, err
		}
		q = qm.queues[key]
	}

	proxy, err := q.Pop()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, key)
	}
	return &proxy, err
}

// UpdateProxyMeta updates meta data and updates priority queue
func (qm QueueManagerImpl) UpdateProxyMeta(key string, meta ProxyMeta) error {
	if meta.Ignore() {
		return nil
	}
	q, ok := qm.queues[key]
	if !ok {
		return nil
	}
	// If we got an error, re-prioritize queue
	q.Update(meta.Addr)
	return nil
}

// RefreshProxies updates the proxies at the corresponding key
func (qm *QueueManagerImpl) RefreshProxies(key string) error {
	proxies, err := qm.store.GetProxies(key)
	if err != nil {
		return err
	}
	if _, ok := qm.queues[key]; !ok {
		qm.queues[key] = NewQueue()
	}
	qm.queues[key].Rebuild(proxies)
	return nil
}
