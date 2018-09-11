package core

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

// Store handles updating master lists from a persistent store
type Store interface {
	GetProxies(key string) ([]Proxy, error)
}

// StoreImpl implements the store interface
type StoreImpl struct {
	pool *redis.Pool
}

// NewStore constructor
func NewStore(pool *redis.Pool) StoreImpl {
	return StoreImpl{pool: pool}
}

// GetProxies gets a proxy list from redis using the key passed
func (s StoreImpl) GetProxies(key string) ([]Proxy, error) {
	conn := s.pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	proxies := make([]Proxy, 0)
	if err := json.Unmarshal(data, &proxies); err != nil {
		return nil, err
	}
	return proxies, nil
}
