package tgbot

import (
	"sync"
	"time"
)

type Cache struct {
	m      *sync.RWMutex
	memory map[string]Item
}

type Item struct {
	data       map[string]string
	keys       []string
	expiration int64
}

func NewCache() Cache {
	return Cache{memory: make(map[string]Item), m: &sync.RWMutex{}}
}

func (c *Cache) Set(username string, data map[string]string, keys []string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.memory[username] = Item{
		data: data,
		keys: keys,
		//expiration: time.Now().Add(time.Second * 10).UnixNano(),
		expiration: time.Now().Add(time.Minute * 10).UnixNano(),
	}
}

func (c *Cache) Get(username string) (Item, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	res, ok := c.memory[username]
	return res, ok
}

func (c *Cache) Remove(username string) {
	if _, ok := c.Get(username); ok {
		c.m.Lock()
		delete(c.memory, username)
		c.m.Unlock()
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.memory == nil || len(c.memory) == 0 {
		return
	}

	for k, i := range c.memory {
		if time.Now().UnixNano() > i.expiration {
			keys = append(keys, k)
		}
	}
	return
}
