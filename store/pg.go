package store

import (
	"context"
	"github.com/easyXpat/procedure-service/config"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
)

const (
	ProcedureTableDDL = `
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
	StepTableDDL = `
		create table if not exists step (
			id	Varchar(64) not null,
			name Varchar(255) not null,
			procedure_name	Varchar(64),
			city Varchar(32),
			description Varchar(1000),
			created_at  Timestamp not null,
			updated_at  Timestamp not null,
			Primary Key(id)
		);
	`

	ProcedureStepTableDDL = `
		create table if not exists procedure_step (
			procedure_id Varchar(64) REFERENCES procedure (id) ON UPDATE CASCADE ON DELETE CASCADE,
			step_id Varchar(64) REFERENCES step (id) ON UPDATE CASCADE ON DELETE CASCADE,
			sequence int NOT NULL,
			CONSTRAINT procedure_step_pk PRIMARY KEY (procedure_id, step_id)
		);
	`
)

func NewConnection(logger hclog.Logger, config *config.Configuration) (*pgx.Conn, error) {
	logger.Info("Connecting to postgres DB", "database_url", config.DatabaseURL)
	conn, err := pgx.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		logger.Error("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}