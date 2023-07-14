package store

import (
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/models"
	"github.com/Yarik-xxx/CodeWarsRestApi/pkg/codewars"
	"log"
	"strings"
	"sync"
	"time"
)

type ChallengeRepo struct {
	database *ChallengeDatabase
	cache    *Cache
}

func NewChallengeRepo(defaultExpiration, cleanupInterval time.Duration, store *Store) (*ChallengeRepo, error) {
	repo := &ChallengeRepo{}
	items := make(map[string]Item)

	if err := store.db.Ping(); err != nil {
		return nil, err
	}

	cache := &Cache{
		m:                 &sync.RWMutex{},
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	database := &ChallengeDatabase{store: store}

	repo.cache = cache
	repo.database = database

	if err := repo.InitCache(); err != nil {
		return nil, err
	}

	log.Println("Загружена информация о", len(repo.cache.items), "катах")

	if cleanupInterval > 0 {
		repo.StartGC()
	}

	return repo, nil
}

func (c *ChallengeRepo) InitCache() error {
	challenges, err := c.database.SelectAllRow()
	if err != nil {
		return err
	}

	for _, challenge := range challenges {
		expiration := challenge.LastUpdate.Add(c.cache.defaultExpiration).UnixNano()
		c.cache.m.Lock()

		c.cache.items[challenge.ID] = Item{
			Value:      challenge.Info,
			Expiration: expiration,
			Created:    challenge.LastUpdate,
		}

		c.cache.m.Unlock()
	}
	return nil
}

func (c *ChallengeRepo) Count() int {
	c.cache.m.Lock()
	defer c.cache.m.Unlock()
	return len(c.cache.items)
}

func (c *ChallengeRepo) Get(key string) (codewars.Kata, bool, error) {
	item, found := c.cache.Get(key)

	if !found {
		kata, err := codewars.GetKata(key)
		if err != nil {
			return kata, false, err
		}

		if err := c.database.AddDB(&models.Challenge{ID: kata.ID, Info: kata, LastUpdate: time.Now()}); err != nil {
			return kata, false, err
		}

		c.cache.Set(kata.ID, kata)

		return kata, false, nil
	}

	return item, true, nil
}

func (c *ChallengeRepo) StartGC() {
	go c.GC()
}

func (c *ChallengeRepo) GC() {
	for {
		if c.cache.items == nil {
			return
		}

		if keys := c.cache.expiredKeys(); len(keys) != 0 {
			c.updateItems(keys)
		}
		time.Sleep(c.cache.cleanupInterval)
	}

}

func (c *ChallengeRepo) updateItems(keys []string) {
	for _, k := range keys {
		kata, err := codewars.GetKata(k)

		for err != nil && (err.Error() == "429" || strings.TrimSpace(strings.ToLower(err.Error())) == "retry later") {
			time.Sleep(5 * time.Second)
			kata, err = codewars.GetKata(k)
		}
		if err != nil {
			log.Println(err, "updateItems") //todo
			continue
		}

		if err := c.database.AddDB(&models.Challenge{
			ID:         k,
			Info:       kata,
			LastUpdate: time.Now(),
		}); err != nil {
			return
		}

		c.cache.Set(k, kata)
	}
}
