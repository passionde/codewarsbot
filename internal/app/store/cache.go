package store

import (
	"github.com/Yarik-xxx/CodeWarsRestApi/pkg/codewars"
	"sync"
	"time"
)

type Cache struct {
	m                 *sync.RWMutex
	items             map[string]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

type Item struct {
	Value      codewars.Kata
	Expiration int64
	Created    time.Time
}

func (c *Cache) Set(key string, value codewars.Kata) {
	expiration := time.Now().Add(c.defaultExpiration).UnixNano()

	c.m.Lock()
	defer c.m.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

func (c *Cache) Get(key string) (codewars.Kata, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	item, found := c.items[key]

	// Не менять, иначе выходит паника
	if !found {
		return codewars.Kata{}, false
	}

	return item.Value, true
}

func (c *Cache) expiredKeys() (keys []string) {
	c.m.Lock()
	defer c.m.Unlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}
