package business

import (
	"fmt"
	"sync"

	order_domain "github.com/CamiloC-pvt/re-challenge/app/order/domain"
	pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"
)

var (
	packBusinessInstance *PackBusiness
	packBusinessOnce     sync.Once
)

type PackBusiness struct {
	orderRepo order_domain.IOrderRepo
	packRepo  pack_domain.IPackRepo
}

func NewPackBusiness(orderRepo order_domain.IOrderRepo, packRepo pack_domain.IPackRepo) pack_domain.IPackBusiness {
	packBusinessOnce.Do(func() {
		packBusinessInstance = &PackBusiness{}
		packBusinessInstance.orderRepo = orderRepo
		packBusinessInstance.packRepo = packRepo
	})

	return packBusinessInstance
}

func (b *PackBusiness) Create(size int32) (int32, error) {
	// Get Current Packs
	dbCurrentPacks, err := b.packRepo.GetAll()
	if err != nil {
		return -1, err
	}

	for _, pack := range dbCurrentPacks {
		if pack.Size == size {
			return -1, fmt.Errorf("error, there is already a pack with size '%d'", size)
		}
	}

	// Create
	newPackID, err := b.packRepo.Create(size)
	if err != nil {
		return -1, fmt.Errorf("error creating pack size on DB: %s", err.Error())
	}

	return newPackID, nil
}

func (b *PackBusiness) Delete(packID int32) error {
	// Get Package to Delete
	currentPack, err := b.packRepo.GetByID(packID)
	if err != nil {
		return err
	}

	// Get Current orders
	dbCurrentOrders, err := b.orderRepo.GetAll()
	if err != nil {
		return err
	}

	// Verify Orders
	for _, order := range dbCurrentOrders {
		for _, usedPackage := range order.Packs {
			if usedPackage.Size == currentPack.Size {
				return fmt.Errorf("error, cannot delete package because is used by the order '%d'", order.ID)
			}
		}
	}

	return b.packRepo.Delete(packID)
}

func (b *PackBusiness) GetAll() ([]pack_domain.Pack, error) {
	return b.packRepo.GetAll()
}
