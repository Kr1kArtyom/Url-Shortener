package pkg //зачем, если есть postgres/url-shortener.go

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DataBase struct {
	db *pgxpool.Pool
}

func NewDataBase(dbCfg string) (*DataBase, error) {
	config, err := pgxpool.ParseConfig(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to establish database connection: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf(": %w", err) //переделать
	}

	return &DataBase{db: db}, nil
}

func (d *DataBase) GetDB() *pgxpool.Pool {
	return d.db
}

func (d *DataBase) Close() {
	d.db.Close()
}
