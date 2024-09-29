package db

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
)

const (
	ENV_POSTGRES_DB         string = "POSTGRES_DB"
	ENV_POSTGRES_PASSWORD   string = "POSTGRES_PASSWORD"
	ENV_POSTGRES_PORT       string = "POSTGRES_PORT"
	ENV_POSTGRES_URL        string = "POSTGRES_URL"
	ENV_POSTGRES_USER       string = "POSTGRES_USER"
	POSTGRES_CONNECTION_URI string = "postgres://POSTGRES_USER:POSTGRES_PASSWORD@POSTGRES_URL:POSTGRES_PORT/POSTGRES_DB"
)

type PostgresConnection struct {
	Connection *pgx.Conn
}

var (
	postgresConnection *PostgresConnection
	once               sync.Once
)

func NewPostgresConnection() *PostgresConnection {
	once.Do(func() {
		postgresConnection = &PostgresConnection{}
	})

	return postgresConnection
}

func (p *PostgresConnection) Connect() {
	connectionUrl := POSTGRES_CONNECTION_URI

	postgresDb := os.Getenv(ENV_POSTGRES_DB)
	postgresPass := os.Getenv(ENV_POSTGRES_PASSWORD)
	postgresPort := os.Getenv(ENV_POSTGRES_PORT)
	postgresUrl := os.Getenv(ENV_POSTGRES_URL)
	postgresUsername := os.Getenv(ENV_POSTGRES_USER)

	connectionUrl = strings.ReplaceAll(connectionUrl, "POSTGRES_DB", postgresDb)
	connectionUrl = strings.ReplaceAll(connectionUrl, "POSTGRES_PASSWORD", postgresPass)
	connectionUrl = strings.ReplaceAll(connectionUrl, "POSTGRES_PORT", postgresPort)
	connectionUrl = strings.ReplaceAll(connectionUrl, "POSTGRES_URL", postgresUrl)
	connectionUrl = strings.ReplaceAll(connectionUrl, "POSTGRES_USER", postgresUsername)

	log.Printf("Connecting to Postgres: %s", connectionUrl)

	conn, err := pgx.Connect(context.Background(), connectionUrl)
	if err != nil {
		log.Panicf("error connecting to Postgres DB: %s", err.Error())
		os.Exit(1)
	}

	p.Connection = conn

	log.Printf("Connected to Postgres")
}
