# Ika

Priority queue proxy manager

# Endpoints

## /proxy

Retrieves a proxy according to the queue corresponding to the passed channel & domain

- takes a channel & domain
- returns a proxy

## /proxy/meta

Informs the priority queue to reprioritize according to the error (currently only considers if the error is populated or not, and lowers proxy priority according to domain, channel, & address)

- takes a channel, domain, proxy address, & error string
- returns nothing

## /proxy/refresh

Updates all proxies in the service (currently only the master list is updated, restart server to rebuild all)

- takes nothing
- returns nothing

# Key Routing

A key router resolves a proxy queue key by the passed channel & domain. (Routes are currently loaded at server start by redis keys -- `proxy:route:channel` and `proxy:route:domain`. This should be updated to refresh on some API hit)

# Priority Queue

Each proxy queue is loaded from redis at a key such as `proxy:queue:<queue_name>`. Every time a proxy fails, the priority is lowered for that proxy.
