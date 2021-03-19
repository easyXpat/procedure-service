package store

import (
	"context"
	"github.com/easyXpat/procedure-service/config"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
)

const (
	ProcedureTableQ = `
		create table if not exists procedure (
			id	Varchar(64) not null,
			name Varchar(255) not null,
			description Varchar(1000),
			city Varchar(32),
			created_at  Timestamp not null,
			updated_at  Timestamp not null,
			Primary Key(id)
		);
	`
)

func NewConnection(logger hclog.Logger, config *config.Configuration) (*pgx.Conn, error) {
	logger.Info("Connecting to postgres DB")
	conn, err := pgx.Connect(context.Background(), config.GetPGConnectionString())
	if err != nil {
		logger.Error("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}