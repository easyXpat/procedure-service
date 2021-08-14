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
	InsertStepDML  = "INSERT INTO step (id, name, procedure_name, city, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	GetAllStepsDML = "SELECT * FROM step"
	GetStepDML     = "SELECT * FROM step WHERE id = $1"
	UpdateStepDML  = "UPDATE step SET name = $2, procedure_name = $3, city = $4, description = $5, updated_at = $6 where id = $1"
	DeleteStep  = "DELETE FROM step WHERE id = $1"
)

// StepPostgres is an implementation of the Steps DB interface
type StepPostgres struct {
	logger hclog.Logger
	db     *pgx.Conn
}

func NewStepPostgres(l hclog.Logger, db *pgx.Conn) *StepPostgres {
	return &StepPostgres{l, db}
}


// AddStep adds a step to the postgres database
func (spg *StepPostgres) AddStep(ctx context.Context, s *Step) error {
	s.ID = uuid.NewV4().String()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	spg.logger.Info("creating step", "id", s.ID, "name", s.Name)
	_, err := spg.db.Exec(ctx, InsertStepDML, s.ID, s.Name, s.ProcedureName, s.City, s.Description, s.CreatedAt, s.UpdatedAt)
	return err
}

// GetAllSteps gets all steps from postgres database
func (spg *StepPostgres) GetAllSteps(ctx context.Context) (Steps, error) {
	spg.logger.Info("fetching steps from db")
	var s Steps

	err := pgxscan.Select(ctx, spg.db, &s, GetAllStepsDML)
	if err != nil {
		spg.logger.Error("Error fetching steps from db", "error", err)
		return nil, err
	}

	return s, nil
}

// GetStep retrieves a step from the postgres database
func (spg *StepPostgres) GetStep(ctx context.Context, id string) (Steps, error) {
	spg.logger.Info("fetching step from db")
	var s Steps

	err := pgxscan.Select(ctx, spg.db, &s, GetStepDML, id)
	if err != nil {
		spg.logger.Error("Error fetching steps from db", "error", err)
		return nil, err
	}

	return s, nil
}

// UpdateStep step from DB using id
func (spg *StepPostgres) UpdateStep(ctx context.Context, s *Step) (Steps, error) {
	s.UpdatedAt = time.Now()

	_, err := spg.db.Exec(ctx, UpdateStepDML, s.ID, s.Name, s.ProcedureName, s.City, s.Description, s.UpdatedAt)
	if err != nil {
		spg.logger.Error("Error updating step in db", "error", err)
		return nil, err
	}
	updatedStep, err := spg.GetStep(ctx, s.ID)
	if err != nil {
		spg.logger.Error("Failed fetching updated record from db", "error", err)
	}
	return updatedStep, nil
}

// DeleteStep deletes a step from the DB using its id
func (spg *StepPostgres) DeleteStep(ctx context.Context, id string) error {

	spg.logger.Info("deleting step", "id", id)
	_, err := spg.db.Exec(ctx, DeleteStep, id)
	if err != nil {
		spg.logger.Error("Error deleting step in db", "error", err)
		return err
	}
	return nil
}
