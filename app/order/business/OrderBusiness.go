package business

import (
	"errors"
	"fmt"
	"math"
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

	// Optimizing Items Amount
	for i := 0; i < len(orderPacks); i++ {
		checkEnd := false
		for !checkEnd && totalShiped > size {
			orderPack := orderPacks[i]

			if orderPack.Amount == 0 || i+1 == len(orderPacks) {
				checkEnd = true
				break
			}

			for j := i + 1; j < len(orderPacks); j++ {
				nextOrderPack := orderPacks[j]

				maxFound := false
				maxToReplace := int32(orderPack.Amount)
				for !maxFound && maxToReplace > 0 {
					maxNextPackPerCurrent := math.Floor(float64(orderPack.Size*maxToReplace) / float64(nextOrderPack.Size))

					current := (orderPack.Amount * orderPack.Size) + (nextOrderPack.Amount * nextOrderPack.Size)
					possible := ((orderPack.Amount - maxToReplace) * orderPack.Size) + ((nextOrderPack.Amount + int32(maxNextPackPerCurrent)) * nextOrderPack.Size)

					possibleTotal := totalShiped - current + possible

					if current > possible && possibleTotal >= size {
						orderPacks[i].Amount -= maxToReplace
						orderPacks[j].Amount += int32(maxNextPackPerCurrent)

						totalShiped = possibleTotal

						maxFound = true
						break
					}

					maxToReplace--
				}

				if maxFound {
					break
				}

				if j+1 == len(orderPacks) {
					checkEnd = true
				}
			}
		}
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

			possibleTotal := totalShiped - (current - possible)
			if possible <= current && possibleTotal >= size {
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
