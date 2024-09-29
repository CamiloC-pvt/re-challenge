package order

import pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"

type IOrderBusiness interface {
	CalculatePackaging(availablePacks []pack_domain.Pack, size int32) []OrderPack
	Cancel(orderID int32) error
	Create(size int32) (Order, error)
	GetAll() ([]Order, error)
}
