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
	InsertProcedureQ = "INSERT INTO procedure (id, name, description, city, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	GetAllProceduresQ = "SELECT * FROM procedure"
	GetProcedureQ = "SELECT * FROM procedure WHERE id = $1"
	UpdateProcedureQ = "UPDATE procedure SET name = $1, description = $2, city = $3, updated_at = $4 where id = $5"
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

	ppg.logger.Info("creating procedure", "id", p.ID, "name", p.Name)
	_, err := ppg.db.Exec(ctx, InsertProcedureQ, p.ID, p.Name, p.Description, p.City, p.CreatedAt, p.UpdatedAt)
	return err
}

// GetAllProcedures gets all procedures from DB
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

// UpdateProcedure procedure from DB using id
func (ppg *ProcedurePG) UpdateProcedure(ctx context.Context, p *Procedure) (*Procedure, error) {
	p.UpdatedAt = time.Now()

	_, err := ppg.db.Exec(ctx, UpdateProcedureQ, p.Name, p.Description, p.City, p.UpdatedAt, p.ID)
	if err != nil {
		ppg.logger.Error("Error updating procedure in db", "error", err)
		return nil, err
	}
	updatedProcedure, err := ppg.GetProcedure(ctx, p.ID)
	if err != nil {
		ppg.logger.Error("Failed fetching updated record from db", "error", err)
	}
	return updatedProcedure, nil
}

// GetProcedure fetch a procedure from DB using id
func (ppg *ProcedurePG) GetProcedure(ctx context.Context, id string) (*Procedure, error) {
	ppg.logger.Debug("querying for procedure", "id", id)
	var p Procedure
	err := pgxscan.Get(context.Background(), ppg.db, &p, GetProcedureQ, id)
	if err != nil {
		ppg.logger.Error("Error fetching procedure", "error", err)
		return nil, err
	}
	return &p, nil

}