package init_storage

import (
	"context"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage(ctx context.Context, storageType string, dsn string) (storage.Storage, error) {
	var db storage.Storage
	switch storageType {
	case "SQL":
		storageSQL := sqlstorage.New()
		pool, err := sqlstorage.Connect(ctx, dsn)
		if err != nil {
			return nil, err
		}
		// defer pool.Close()
		storageSQL.Pool = pool
		db = storageSQL
	default:
		storageMemory := memorystorage.New()
		db = storageMemory
	}
	return db, nil
}
