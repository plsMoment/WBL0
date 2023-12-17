package service

import (
	"WBL0/internal/cache"
	"WBL0/internal/db"
	"WBL0/internal/model"
	"context"
	"fmt"
	"log"
)

type OrderService interface {
	GetOrder(orderID string) ([]byte, error)
}

// Service struct is using for manage cache and database data
type Service struct {
	dbStorage db.OrderStorage
	ch        chan model.Order
	memCache  *cache.Cache
}

// New create instance of service struct and refill cache
func New(storage db.OrderStorage, orderCh chan model.Order) *Service {
	ctx := context.Background()
	newCache := cache.InitCache()
	orders, err := storage.GetAllOrders(ctx)
	if err != nil {
		log.Printf("restoring cache failed: %v", err)
	}

	for _, order := range orders {
		newCache.Add(order.Id, order.Data)
	}

	service := &Service{
		dbStorage: storage,
		ch:        orderCh,
		memCache:  newCache,
	}

	go service.orderReceiver()
	return service
}

// GetOrder returns order with such id
func (s *Service) GetOrder(orderID string) ([]byte, error) {
	orderData, ok := s.memCache.Get(orderID)
	if !ok {
		return nil, fmt.Errorf("order with such id doesn't exists")
	}
	res, ok := orderData.([]byte)
	if !ok {
		return nil, fmt.Errorf("type assertion to []byte failed")
	}

	return res, nil
}

// OrderReceiver receives data from STAN producer
func (s *Service) orderReceiver() {
	ctx := context.Background()

	for {
		order, ok := <-s.ch
		if !ok {
			return
		}
		if err := s.dbStorage.CreateOrder(ctx, order); err != nil {
			log.Printf("can't create order: %v", err)
		} else {
			s.memCache.Add(order.Id, order.Data)
		}
	}
}
