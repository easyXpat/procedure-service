package data

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	InsertProcedureQ = "insert into procedure (id, name, description, city, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)"
	GetAllProceduresQ = "SELECT * FROM procedure"
)

// ProcedurePG is an implementation of the Procedure DB interface
type ProcedurePG struct {
	logger      hclog.Logger
	db		*pgx.Conn
}

// NewProcedurePG Creates a client for interacting with the procedure relation in postgres DB
func NewProcedurePG(l hclog.Logger, db *pgx.Conn) *ProcedurePG {
	return &ProcedurePG{l, db}
}

// AddProcedure adds a procedure and stores it in the DB
func (ppg *ProcedurePG) AddProcedure(ctx context.Context, p *Procedure) error {
	p.ID = uuid.NewV4().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	ppg.logger.Info("creating procedure", hclog.Fmt("%#v", p))
	_, err := ppg.db.Exec(ctx, InsertProcedureQ, p.ID, p.Name, p.Description, p.City, p.CreatedAt, p.UpdatedAt)
	return err
}

func (ppg *ProcedurePG) GetAllProcedures(ctx context.Context) (Procedures, error) {
	ppg.logger.Info("fetching procedures from db")
	var p Procedures

	//rows, err := ppg.db.Query(context.Background(), GetAllProceduresQ)
	err := pgxscan.Select(context.Background(), ppg.db, &p, GetAllProceduresQ)
	if err != nil {
		ppg.logger.Error("Error fetching procedures from db", "error", err)
		return nil, err
	}

	return p, nil
}