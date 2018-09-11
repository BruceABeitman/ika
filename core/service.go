package core

// Service blah blah blah
type Service interface {
	GetProxy(queueID QueueID) (*Proxy, error)
	UpdateProxyMeta(queueID QueueID, meta ProxyMeta) error
	RefreshProxies(key string) error
}

// ProxyMeta struct holds meta data about proxy success
type ProxyMeta struct {
	Addr  string `json:"addr"`
	Error string `json:"error"`
}

// Ignore is true if we can safely ignore this meta object
// when considering queue updates
func (meta ProxyMeta) Ignore() bool {
	return len(meta.Error) == 0
}

// QueueID struct holds data to identify a queue
type QueueID struct {
	Channel string `json:"channel"`
	Domain  string `json:"domain"`
}

// Proxy struct
type Proxy struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

// ServiceImpl implements the service interface
type ServiceImpl struct {
	KeyRouter    KeyRouter
	QueueManager QueueManager
}

// GetProxy retrieves a proxy from the service
func (s ServiceImpl) GetProxy(queueID QueueID) (*Proxy, error) {
	// Resolve key from channel & domain
	key := s.KeyRouter.GetKey(queueID)
	return s.QueueManager.GetProxy(key)
}

// UpdateProxyMeta updates meta-data concerning a proxies success
func (s ServiceImpl) UpdateProxyMeta(queueID QueueID, proxyMeta ProxyMeta) error {
	// Resolve key from channel & domain
	key := s.KeyRouter.GetKey(queueID)
	return s.QueueManager.UpdateProxyMeta(key, proxyMeta)
}

// RefreshProxies updates queue corresponding with key
func (s ServiceImpl) RefreshProxies(key string) error {
	return s.QueueManager.RefreshProxies(key)
}
