package cache

import (
	"errors"
	"github.com/Back1ng/wbtech-0/internal/entity"
	"sync"
)

type InMemoryOrderCache struct {
	orders map[string]entity.Order
	mu     *sync.RWMutex
}

func New() InMemoryOrderCache {
	return InMemoryOrderCache{
		orders: make(map[string]entity.Order),
		mu:     &sync.RWMutex{},
	}
}

func (c InMemoryOrderCache) Has(uid string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.orders[uid]

	return ok
}

func (c InMemoryOrderCache) Get(uid string) (*entity.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.Has(uid) {
		//log.Println("cache hit:", uid)
		order := c.orders[uid]
		return &order, nil
	}

	return nil, errors.New("order not found")
}

func (c InMemoryOrderCache) GetAll() []entity.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	orders := make([]entity.Order, 0, len(c.orders))

	for _, order := range c.orders {
		orders = append(orders, order)
	}

	return orders
}

func (c InMemoryOrderCache) Store(order entity.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders[order.OrderUID] = order
}

func (c InMemoryOrderCache) StoreAll(orders []entity.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}
}

func (c InMemoryOrderCache) Delete(uid string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.orders, uid)
}
