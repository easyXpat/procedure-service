package data

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	InsertProcedureDML  = "INSERT INTO procedure (id, name, description, city, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	GetAllProceduresDML = "SELECT * FROM procedure"
	GetProcedureDML     = "SELECT * FROM procedure WHERE id = $1"
	UpdateProcedureDML  = "UPDATE procedure SET name = $1, description = $2, city = $3, updated_at = $4 where id = $5"
	InsertStepDML  = "INSERT INTO step (id, procedure_id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	GetAllStepsDML = "SELECT * FROM step"
	GetProcedureStepsDML     = "SELECT * FROM step WHERE procedure_id = $1"
	GetStepDML     = "SELECT * FROM step WHERE id = $1"
	DeleteProcedureQ  = "DELETE FROM procedure WHERE id = $1"
)

// ProcedurePG is an implementation of the Procedure DB interface
type ProcedurePG struct {
	logger hclog.Logger
	db     *pgx.Conn
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
	_, err := ppg.db.Exec(ctx, InsertProcedureDML, p.ID, p.Name, p.Description, p.City, p.CreatedAt, p.UpdatedAt)
	return err
}

// GetAllProcedures gets all procedures from DB
func (ppg *ProcedurePG) GetAllProcedures(ctx context.Context) (Procedures, error) {
	ppg.logger.Info("fetching procedures from db")
	var p Procedures

	//rows, err := ppg.db.Query(context.Background(), GetAllProceduresDML)
	err := pgxscan.Select(ctx, ppg.db, &p, GetAllProceduresDML)
	if err != nil {
		ppg.logger.Error("Error fetching procedures from db", "error", err)
		return nil, err
	}

	return p, nil
}

// UpdateProcedure procedure from DB using id
func (ppg *ProcedurePG) UpdateProcedure(ctx context.Context, p *Procedure) (*Procedure, error) {
	p.UpdatedAt = time.Now()

	_, err := ppg.db.Exec(ctx, UpdateProcedureDML, p.Name, p.Description, p.City, p.UpdatedAt, p.ID)
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
	err := pgxscan.Get(ctx, ppg.db, &p, GetProcedureDML, id)
	if err != nil {
		ppg.logger.Error("Error fetching procedure", "error", err)
		return nil, err
	}
	return &p, nil

}

// AddStep adds a procedure and stores it in the DB
func (ppg *ProcedurePG) AddStep(ctx context.Context, p *Step) error {
	p.ID = uuid.NewV4().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	ppg.logger.Info("creating step", "id", p.ID, "name", p.Name)
	_, err := ppg.db.Exec(ctx, InsertStepDML, p.ID, p.ProcedureID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt)
	return err
}

// GetAllSteps gets all steps from DB
func (ppg *ProcedurePG) GetAllSteps(ctx context.Context) (Steps, error) {
	ppg.logger.Info("fetching steps from db")
	var s Steps

	err := pgxscan.Select(ctx, ppg.db, &s, GetAllStepsDML)
	if err != nil {
		ppg.logger.Error("Error fetching steps from db", "error", err)
		return nil, err
	}

	return s, nil
}

func (ppg *ProcedurePG) GetProcedureSteps(ctx context.Context, id string) (Steps, error) {
	ppg.logger.Info("fetching steps from db for procedure", "procedure", id)
	var s Steps

	err := pgxscan.Select(ctx, ppg.db, &s, GetProcedureStepsDML, id)
	if err != nil {
		ppg.logger.Error("Error fetching steps from db", "error", err)
		return nil, err
	}

	return s, nil
}

func (ppg *ProcedurePG) GetStep(ctx context.Context, id string) (Steps, error) {
	ppg.logger.Info("fetching step from db")
	var s Steps

	err := pgxscan.Select(ctx, ppg.db, &s, GetStepDML, id)
	if err != nil {
		ppg.logger.Error("Error fetching steps from db", "error", err)
		return nil, err
	}

	return s, nil
}

// DeleteProcedure deletes a procedure from the DB using its id
func (ppg *ProcedurePG) DeleteProcedure(ctx context.Context, id string) error {

	ppg.logger.Info("deleting procedure", "id", id)
	_, err := ppg.db.Exec(ctx, DeleteProcedureQ, id)
	if err != nil {
		ppg.logger.Error("Error deleting procedure in db", "error", err)
		return err
	}
	return nil
}
