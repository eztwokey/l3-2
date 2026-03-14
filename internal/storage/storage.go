package storage

import (
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
)

type Storage struct {
	db  *dbpg.DB
	rdb *redis.Client
}

func New(db *dbpg.DB, rdb *redis.Client) *Storage {
	return &Storage{db: db, rdb: rdb}
}
