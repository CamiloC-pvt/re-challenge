package business

import (
	"fmt"
	"sync"

	pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"
)

var (
	packBusinessInstance *PackBusiness
	packBusinessOnce     sync.Once
)

type PackBusiness struct {
	packRepo pack_domain.IPackRepo
}

func NewPackBusiness(packRepo pack_domain.IPackRepo) pack_domain.IPackBusiness {
	packBusinessOnce.Do(func() {
		packBusinessInstance = &PackBusiness{}
		packBusinessInstance.packRepo = packRepo
	})

	return packBusinessInstance
}

func NewPackBusiness_Mock(packRepo pack_domain.IPackRepo) pack_domain.IPackBusiness {
	packBusinessMock := &PackBusiness{}
	packBusinessMock.packRepo = packRepo

	return packBusinessMock
}

func (b *PackBusiness) Create(size int32) (int32, error) {
	newPackID, err := b.packRepo.Create(size)
	if err != nil {
		return -1, fmt.Errorf("error creating pack size on DB: %s", err.Error())
	}

	return newPackID, nil
}

func (b *PackBusiness) Delete(packID int32) error {
	return b.packRepo.Delete(packID)
}

func (b *PackBusiness) GetAll() ([]pack_domain.Pack, error) {
	return b.packRepo.GetAll()
}
