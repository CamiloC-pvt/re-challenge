package order

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	db "github.com/CamiloC-pvt/re-challenge/app/db"
	order_domain "github.com/CamiloC-pvt/re-challenge/app/order/domain"
)

var (
	orderPostgresRepoInstance *OrderPostgresRepo
	orderPostgresRepoOnce     sync.Once
)

type OrderPostgresRepo struct {
	postgresConnection *db.PostgresConnection
}

func NewOrderPostgresRepo(postgresConnection *db.PostgresConnection) order_domain.IOrderRepo {
	orderPostgresRepoOnce.Do(func() {
		orderPostgresRepoInstance = &OrderPostgresRepo{}
		orderPostgresRepoInstance.postgresConnection = postgresConnection
	})

	return orderPostgresRepoInstance
}

func (r *OrderPostgresRepo) Cancel(orderID int32) error {
	ctx := context.Background()

	_, err := r.postgresConnection.Connection.Exec(ctx, fmt.Sprintf("UPDATE \"order\" SET deleted = TRUE, modified = NOW() WHERE id = %d;", orderID))
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPostgresRepo) GetAll() ([]order_domain.Order, error) {
	ctx := context.Background()

	query := `SELECT o.created, o.deleted, o.id, o.modified, o.packs, o.size FROM "order" o WHERE o.deleted = FALSE;`

	rows, err := r.postgresConnection.Connection.Query(ctx, query)
	if err != nil {
		return []order_domain.Order{}, err
	}
	defer rows.Close()

	var orders []order_domain.Order = []order_domain.Order{}
	for rows.Next() {
		var o order_domain.Order
		err := rows.Scan(&o.Created, &o.Deleted, &o.ID, &o.Modified, &o.Packs, &o.Size)
		if err != nil {
			return []order_domain.Order{}, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}

func (r *OrderPostgresRepo) Save(packs []order_domain.OrderPack, size int32) (order_domain.Order, error) {
	ctx := context.Background()

	// Pack bytes
	packsBytes, err := json.Marshal(packs)
	if err != nil {
		return order_domain.Order{}, err
	}

	// Save
	var o order_domain.Order = order_domain.Order{}
	row := r.postgresConnection.Connection.QueryRow(ctx, "INSERT INTO \"order\" (size, packs) VALUES ($1, $2) RETURNING created, deleted, id, modified, packs, size;", size, packsBytes)
	err = row.Scan(&o.Created, &o.Deleted, &o.ID, &o.Modified, &o.Packs, &o.Size)
	if err != nil {
		return order_domain.Order{}, err
	}

	return o, nil
}
