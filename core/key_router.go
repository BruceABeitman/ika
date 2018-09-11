package core

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

const (
	// DomainRouteKey key to domain routes
	domainRouteKey = "proxy:route:domain"

	// ChannelRouteKey key to channel routes
	channelRouteKey = "proxy:route:channel"

	defaultKey = "proxy:queue:master"
)

// ErrorRetrievingDomainRoute ...
var ErrorRetrievingDomainRoute = errors.New("Could not retrieve domain route")

// ErrorRetrievingChannelRoute ...
var ErrorRetrievingChannelRoute = errors.New("Could not retrieve channel route")

// KeyRouter interface for resolving key from channel & domain
type KeyRouter interface {
	GetKey(queueID QueueID) string
}

// KeyRouterImpl implements KeyRouter interface
type KeyRouterImpl struct {
	pool       *redis.Pool
	domainMap  map[string]string
	channelMap map[string]string
}

// NewKeyRouter constructor
func NewKeyRouter(pool *redis.Pool) (KeyRouter, error) {
	// Initialize keyRouter
	keyRouter := &KeyRouterImpl{
		pool:       pool,
		domainMap:  make(map[string]string),
		channelMap: make(map[string]string),
	}
	// Refresh all routes to initialize domain & channel keys
	err := keyRouter.RefreshRoutes()
	if err != nil {
		return nil, err
	}
	return keyRouter, nil
}

// GetKey resolves a key from the queueID
func (k KeyRouterImpl) GetKey(queueID QueueID) string {
	// Check if domain matches first, if so return key
	domainKey, ok := k.domainMap[queueID.Domain]
	if ok {
		return domainKey
	}
	// Check if channel matches second, if so return key
	channelKey, ok := k.channelMap[queueID.Channel]
	if ok {
		return channelKey
	}
	// If none of above were found, return default key
	return defaultKey
}

// RefreshRoutes updates routing configs
func (k *KeyRouterImpl) RefreshRoutes() error {
	conn := k.pool.Get()
	defer conn.Close()

	domainMap, err := redis.StringMap(conn.Do("HGETALL", domainRouteKey))
	if err != nil {
		return ErrorRetrievingDomainRoute
	}
	channelMap, err := redis.StringMap(conn.Do("HGETALL", channelRouteKey))
	if err != nil {
		return ErrorRetrievingChannelRoute
	}

	k.domainMap = domainMap
	k.channelMap = channelMap
	return nil
}
