package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/forge4flow/forge4flow-core/pkg/config"
	"github.com/forge4flow/forge4flow-core/pkg/database"
	"github.com/forge4flow/forge4flow-core/pkg/event"
	"github.com/forge4flow/forge4flow-core/pkg/flow"
	"github.com/forge4flow/forge4flow-core/pkg/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	MySQLDatastoreMigrationVersion     = 10
	MySQLEventstoreMigrationVersion    = 3
	PostgresDatastoreMigrationVersion  = 5
	PostgresEventstoreMigrationVersion = 4
	SQLiteDatastoreMigrationVersion    = 4
	SQLiteEventstoreMigrationVersion   = 3
)

type ServiceEnv struct {
	Datastore  database.Database
	Eventstore database.Database
}

func (env ServiceEnv) DB() database.Database {
	return env.Datastore
}

func (env ServiceEnv) EventDB() database.Database {
	return env.Eventstore
}

func (env *ServiceEnv) InitDB(cfg config.Config) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	if cfg.GetDatastore().MySQL.Hostname != "" {
		db := database.NewMySQL(*cfg.GetDatastore().MySQL)
		err := db.Connect(ctx)
		if err != nil {
			return err
		}

		if cfg.GetAutoMigrate() {
			err = db.Migrate(ctx, MySQLDatastoreMigrationVersion)
			if err != nil {
				return err
			}
		}

		env.Datastore = db
		return nil
	}

	if cfg.GetDatastore().Postgres.Hostname != "" {
		db := database.NewPostgres(*cfg.GetDatastore().Postgres)
		err := db.Connect(ctx)
		if err != nil {
			return err
		}

		if cfg.GetAutoMigrate() {
			err = db.Migrate(ctx, PostgresDatastoreMigrationVersion)
			if err != nil {
				return err
			}
		}

		env.Datastore = db
		return nil
	}

	if cfg.GetDatastore().SQLite.Database != "" {
		db := database.NewSQLite(*cfg.GetDatastore().SQLite)
		err := db.Connect(ctx)
		if err != nil {
			return err
		}

		if cfg.GetAutoMigrate() {
			err = db.Migrate(ctx, SQLiteDatastoreMigrationVersion)
			if err != nil {
				return err
			}
		}

		env.Datastore = db
		return nil
	}

	return errors.New("invalid database configuration provided")
}

func (env *ServiceEnv) InitEventDB(config config.Config) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	if config.GetEventstore().MySQL.Hostname != "" {
		db := database.NewMySQL(*config.GetEventstore().MySQL)
		err := db.Connect(ctx)
		if err != nil {
			return err
		}

		if config.GetAutoMigrate() {
			err = db.Migrate(ctx, MySQLEventstoreMigrationVersion)
			if err != nil {
				return err
			}
		}

		env.Eventstore = db
		return nil
	}

	if config.GetEventstore().Postgres.Hostname != "" {
		db := database.NewPostgres(*config.GetEventstore().Postgres)
		err := db.Connect(ctx)
		if err != nil {
			return err
		}

		if config.GetAutoMigrate() {
			err = db.Migrate(ctx, PostgresEventstoreMigrationVersion)
			if err != nil {
				return err
			}
		}

		env.Eventstore = db
		return nil
	}

	if config.GetEventstore().SQLite.Database != "" {
		db := database.NewSQLite(*config.GetEventstore().SQLite)
		err := db.Connect(ctx)
		if err != nil {
			return err
		}

		if config.GetAutoMigrate() {
			err = db.Migrate(ctx, SQLiteEventstoreMigrationVersion)
			if err != nil {
				return err
			}
		}

		env.Eventstore = db
		return nil
	}

	return errors.New("invalid database configuration provided")
}

func NewServiceEnv() ServiceEnv {
	return ServiceEnv{
		Datastore:  nil,
		Eventstore: nil,
	}
}

func main() {
	cfg := config.NewConfig()
	svcEnv := NewServiceEnv()
	err := svcEnv.InitDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize and connect to the configured datastore. Shutting down.")
	}

	err = svcEnv.InitEventDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize and connect to the configured eventstore. Shutting down.")
	}

	// Init event repo and service
	eventRepository, err := event.NewRepository(svcEnv.EventDB())
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize EventRepository")
	}
	eventSvc := event.NewService(svcEnv, eventRepository, cfg.Eventstore.SynchronizeEvents, nil)

	// Init the flow service
	flowSerice := flow.NewService(&svcEnv, cfg)

	svcs := []service.Service{
		eventSvc,
		flowSerice,
	}

	router, err := service.NewRouter(cfg, "", svcs, service.PassthroughAuthMiddleware, []service.Middleware{}, []service.Middleware{})
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize service router")
	}

	log.Info().Msgf("Listening on port %d", cfg.GetPort())
	shutdownErr := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetPort()), router)
	log.Fatal().Err(shutdownErr).Msg("")
}
