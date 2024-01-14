package flow

import (
	"context"
	"fmt"

	"github.com/forge4flow/forge4flow-core/pkg/database"
	"github.com/pkg/errors"
)

type FlowEventRepository interface {
	Create(ctx context.Context, event Model) (int64, error)
	CreateEventAction(ctx context.Context, eventAction ActionModel) (int64, error)
	GetById(ctx context.Context, id int64) (Model, error)
	GetByType(ctx context.Context, eventType string) (Model, error)
	GetAllEvents(ctx context.Context) ([]Model, error)
	GetActionsForEvent(ctx context.Context, eventType string) ([]ActionModel, error)
	DeleteById(ctx context.Context, id int64) error
	DeleteByType(ctx context.Context, eventType string) error
	UpdateLastBlockHeightByType(ctx context.Context, eventType string, lastBlockHeight uint64) error
}

func NewEventRepository(db database.Database) (FlowEventRepository, error) {
	switch db.Type() {
	case database.TypeMySQL:
		mysql, ok := db.(*database.MySQL)
		if !ok {
			return nil, errors.New(fmt.Sprintf("invalid %s database config", database.TypeMySQL))
		}

		return NewMySQLRepository(mysql), nil
	//TODO: Finish Repositories for PostgresSQL and SQLite
	// case database.TypePostgres:
	// 	postgres, ok := db.(*database.Postgres)
	// 	if !ok {
	// 		return nil, errors.New(fmt.Sprintf("invalid %s database config", database.TypePostgres))
	// 	}

	// 	return NewPostgresRepository(postgres), nil
	// case database.TypeSQLite:
	// 	sqlite, ok := db.(*database.SQLite)
	// 	if !ok {
	// 		return nil, errors.New(fmt.Sprintf("invalid %s database config", database.TypeSQLite))
	// 	}

	// 	return NewSQLiteRepository(sqlite), nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported database type %s specified", db.Type()))
	}
}
