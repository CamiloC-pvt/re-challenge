package pack

import (
	"context"
	"fmt"
	"sync"

	db "github.com/CamiloC-pvt/re-challenge/app/db"
	pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"
)

var (
	packPostgresRepoInstance *PackPostgresRepo
	packPostgresRepoOnce     sync.Once
)

type PackPostgresRepo struct {
	postgresConnection *db.PostgresConnection
}

func NewPackPostgresRepo(postgresConnection *db.PostgresConnection) pack_domain.IPackRepo {
	packPostgresRepoOnce.Do(func() {
		packPostgresRepoInstance = &PackPostgresRepo{}
		packPostgresRepoInstance.postgresConnection = postgresConnection
	})

	return packPostgresRepoInstance
}

func (b *PackPostgresRepo) Create(size int32) (int32, error) {
	ctx := context.Background()

	var newPackID int32 = -1
	row := b.postgresConnection.Connection.QueryRow(ctx, "INSERT INTO pack_size (size) VALUES ($1) RETURNING id;", size)
	err := row.Scan(&newPackID)
	if err != nil {
		return -1, err
	}

	return newPackID, nil
}

func (b *PackPostgresRepo) Delete(packID int32) error {
	ctx := context.Background()

	_, err := b.postgresConnection.Connection.Exec(ctx, fmt.Sprintf("DELETE FROM pack_size WHERE id = %d;", packID))
	if err != nil {
		return err
	}

	return nil
}

func (b *PackPostgresRepo) GetAll() ([]pack_domain.Pack, error) {
	ctx := context.Background()

	query := "SELECT p.id, p.size, p.created FROM pack_size p;"

	rows, err := b.postgresConnection.Connection.Query(ctx, query)
	if err != nil {
		return []pack_domain.Pack{}, err
	}
	defer rows.Close()

	var packs []pack_domain.Pack = []pack_domain.Pack{}
	for rows.Next() {
		var p pack_domain.Pack
		err := rows.Scan(&p.ID, &p.Size, &p.Created)
		if err != nil {
			return []pack_domain.Pack{}, err
		}

		packs = append(packs, p)
	}

	return packs, nil
}
