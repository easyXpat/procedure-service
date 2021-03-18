package store

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"os"
)

type PGClient struct {
	logger hclog.Logger
	Conn *pgx.Conn
}

const (
	procedureSchema = `
		create table if not exists procedure (
			id	Varchar(64) not null,
			name Varchar(255) not null,
			description Varchar(1000),
			city Varchar(32),
			Primary Key(id)
		);
	`
)

func NewPGClient(logger hclog.Logger) *PGClient {
	viper.SetDefault("DATABASE_URL", MakeConnString(logger))
	logger.Info("DATABASE_URL", viper.GetString("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		logger.Error("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer conn.Close(context.Background())
	return &PGClient{
		logger: logger,
		Conn: conn,
	}
}

func MakeConnString(logger hclog.Logger) string {
	logger.Debug("Creating postgres connection string")
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASS")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	database := viper.GetString("DB_NAME")
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?", username, password, host, port, database)
	logger.Debug("Postgres connection string is", connString)
	return connString
}

func (pg *PGClient) CreateProcedureDB() {
	pg.logger.Debug("Creating procedure relation")
	tag, err := pg.Conn.Exec(context.Background(), procedureSchema)
	if err != nil {
		pg.logger.Error("Error creating procedure relation")
	}
	pg.logger.Debug("procedure relation is present in db", tag.String())
}