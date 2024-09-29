package pack

type IPackRepo interface {
	Create(size int32) (int32, error)
	Delete(packID int32) error
	GetAll() ([]Pack, error)
	GetByID(id int32) (Pack, error)
}
