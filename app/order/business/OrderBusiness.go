package business

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	order_domain "github.com/CamiloC-pvt/re-challenge/app/order/domain"
	pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"
)

var (
	orderBusinessInstance *OrderBusiness
	orderBusinessOnce     sync.Once
)

type OrderBusiness struct {
	orderRepo order_domain.IOrderRepo
	packRepo  pack_domain.IPackRepo
}

func NewOrderBusiness(orderRepo order_domain.IOrderRepo, packRepo pack_domain.IPackRepo) order_domain.IOrderBusiness {
	orderBusinessOnce.Do(func() {
		orderBusinessInstance = &OrderBusiness{}
		orderBusinessInstance.orderRepo = orderRepo
		orderBusinessInstance.packRepo = packRepo
	})

	return orderBusinessInstance
}

func NewOrderBusiness_Mock(orderRepo order_domain.IOrderRepo, packRepo pack_domain.IPackRepo) order_domain.IOrderBusiness {
	orderBusinessMock := &OrderBusiness{}
	orderBusinessMock.orderRepo = orderRepo
	orderBusinessMock.packRepo = packRepo

	return orderBusinessMock
}

func (b *OrderBusiness) CalculatePackaging(availablePacks []pack_domain.Pack, size int32) []order_domain.OrderPack {
	orderPacks := []order_domain.OrderPack{}

	sort.Slice(availablePacks, func(i, j int) bool {
		return availablePacks[i].Size > availablePacks[j].Size
	})

	totalShiped := int32(0)
	remaining := size
	for i, availableSize := range availablePacks {
		currentPackCount := 0

		for remaining-availableSize.Size >= 0 {
			remaining -= availableSize.Size
			currentPackCount++
		}

		if i+1 == len(availablePacks) && remaining > 0 {
			currentPackCount++
		}

		totalShiped += int32(currentPackCount) * availableSize.Size

		orderPacks = append(orderPacks, order_domain.OrderPack{
			Amount: int32(currentPackCount),
			ID:     availableSize.ID,
			Size:   availableSize.Size,
		})
	}

	// Optimizing Package Amount
	for i := len(orderPacks) - 1; i >= 0; i-- {
		if i-1 < 0 {
			break
		}

		orderPack := orderPacks[i]
		prevOrderPack := orderPacks[i-1]

		if orderPack.Amount > 1 {
			current := (orderPack.Amount * orderPack.Size) + (prevOrderPack.Amount * prevOrderPack.Size)
			possible := (prevOrderPack.Amount + 1) * prevOrderPack.Size

			if possible <= current && (totalShiped-(current-possible)) >= size {
				totalShiped -= (current - possible)

				orderPacks[i-1].Amount++
				orderPacks = orderPacks[:i]
			} else {
				break
			}
		}
	}

	// Clear Unused Packages
	finalList := []order_domain.OrderPack{}
	for _, orderPack := range orderPacks {
		if orderPack.Amount > 0 {
			finalList = append(finalList, orderPack)
		}
	}

	return finalList
}

func (b *OrderBusiness) Cancel(orderID int32) error {
	return b.orderRepo.Cancel(orderID)
}

func (b *OrderBusiness) Create(size int32) (order_domain.Order, error) {
	// Size Data
	dbPacks, err := b.packRepo.GetAll()
	if err != nil {
		return order_domain.Order{}, fmt.Errorf("error getting the available package sizes: %s", err.Error())
	}

	if len(dbPacks) == 0 {
		return order_domain.Order{}, errors.New("cannot create any order because there are no Package sizes available, please create at least 1 first")
	}

	// Packaging Calculation
	orderPacks := b.CalculatePackaging(dbPacks, size)

	/// Save Order
	newOrder, err := b.orderRepo.Save(orderPacks, size)
	if err != nil {
		return order_domain.Order{}, fmt.Errorf("error saving the new Order on DB: %s", err.Error())
	}

	return newOrder, nil
}

func (b *OrderBusiness) GetAll() ([]order_domain.Order, error) {
	return b.orderRepo.GetAll()
}
