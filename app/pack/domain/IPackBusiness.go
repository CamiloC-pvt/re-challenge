package pack

type IPackBusiness interface {
	Create(size int32) (int32, error)
	Delete(packID int32) error
	GetAll() ([]Pack, error)
}
