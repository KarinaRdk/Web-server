package service

import (
	"TestWebServer/internal/cache"
	"TestWebServer/internal/model"
	"TestWebServer/internal/storage"
	"encoding/json"
	"log"
)

// Service struct encapsulates the service layer logic, including database and cache access.
type Service struct {
	pg *storage.Database // Database instance for persistent storage.
	c  *cache.InMemory   // In-memory cache for fast access.
}

// New creates a new instance of the Service with the provided database and cache instances.
func New(s *storage.Database, cache *cache.InMemory) *Service {
	service := Service{pg: s, c: cache}
	if s.IfDataExists() {
		service.populateCache()
	}
	return &service
}

// populateCache populates the cache with JSON data from the orders table.
// populateCache populates the cache with JSON data from all orders.
func (s *Service) populateCache() {
	jsonOrders, err := s.pg.FetchAllOrders()
	if err != nil {
		log.Printf("Failed to fetch orders for cache: %v\n", err)
		return
	}

	for _, jsonOrder := range jsonOrders {
		var order model.Order
		if err := json.Unmarshal(jsonOrder, &order); err != nil {
			log.Printf("Failed to unmarshal order JSON: %v\n", err)
			continue
		}

		// Use the order_uid as the key to store the order in the cache.
		// Ensure order.OrderUid is not nil before using it.
		if order.OrderUid != nil {
			err := s.c.Set(*order.OrderUid, jsonOrder)
			if err != nil {
				log.Printf("Failed to set order %s in cache: %v\n", *order.OrderUid, err)
			}
		} else {
			log.Println("Order UID is nil, skipping cache set.")
		}
	}
}

// Get retrieves a value by its key from the cache. If the value is not found in the cache,
// it logs an error and returns the error. Otherwise, it returns the value.
func (s *Service) Get(key string) (value []byte, err error) {

	// Attempt to retrieve the value from the cache.
	if value, err = s.c.Get(key); err != nil {
		log.Printf("Failed to retrieve from cache with key %s: %v\n", key, err)
		return value, err
	}
	return value, nil
}
