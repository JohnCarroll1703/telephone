package postgresql

import (
	"database/sql"
	"fmt"
	"git.tarlanpayments.kz/pkg/golog"
	gormlog "git.tarlanpayments.kz/pkg/golog/contrib/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

const (
	connectionsNum = 15
)

func NewDB(sources []string, replicas []string, logger golog.ContextLogger) (*sql.DB, *gorm.DB, error) {
	logger.Infow("Start connecting to DB")

	if len(sources) == 0 {
		return nil, nil, fmt.Errorf("Sources for connection are empty")
	}

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  sources[0],
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}),
		&gorm.Config{
			Logger: gormlog.NewGormLogAdapter(logger),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	if len(sources) > 1 || len(replicas) > 0 {
		if err = setSourcesAndReplicas(logger, db, sources, replicas); err != nil {
			return nil, nil, err
		}
	}
	conn, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	conn.SetConnMaxLifetime(time.Hour)
	conn.SetMaxIdleConns(connectionsNum)
	conn.SetMaxOpenConns(connectionsNum)

	logger.Infow("DB is connected")
	return conn, db, nil
}

func setSourcesAndReplicas(logger golog.ContextLogger, db *gorm.DB, sources, replicas []string) error {
	logger.Infow("Database: setup multi source/replica")

	sources = sources[1:]
	sourceConns := make([]gorm.Dialector, len(sources))
	replicaConns := make([]gorm.Dialector, len(replicas))

	for i, replicaDSN := range replicas {
		replicaConns[i] = postgres.Open(replicaDSN)
	}

	for i, sourceDSN := range sources {
		sourceConns[i] = postgres.Open(sourceDSN)
	}

	return db.Use(dbresolver.Register(
		dbresolver.Config{
			Sources:  sourceConns,
			Replicas: replicaConns,
			Policy:   dbresolver.RandomPolicy{},
		},
	))

}
