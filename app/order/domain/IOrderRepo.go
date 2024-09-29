package order

type IOrderRepo interface {
	Cancel(orderID int32) error
	GetAll() ([]Order, error)
	Save(packs []OrderPack, size int32) (Order, error)
}
