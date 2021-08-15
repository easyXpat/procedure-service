package data

import (
	"context"
	"fmt"
	"strings"
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
	GetProcedureStepsDML     = "SELECT * FROM step WHERE procedure_id = $1"
	DeleteProcedureQ  = "DELETE FROM procedure WHERE id = $1"
	InsertProcedureStepsDML = "INSERT INTO procedure_step (procedure_id, step_id, sequence) VALUES"
	GetStepsWithProcedureID = "SELECT * FROM step WHERE id IN (SELECT step_id FROM procedure_step WHERE procedure_id = $1)"
	DeleteProcedureStepsRelation = "DELETE FROM procedure_step WHERE procedure_id = $1"
)

// ProcedurePostgres is an implementation of the Procedure DB interface
type ProcedurePostgres struct {
	logger hclog.Logger
	db     *pgx.Conn
}

// StepSequence is an implementation of the procedure step mapping result from postgres
type StepSequence struct {
	StepID string `json:"step_id" sql:"step_id"`
	Sequence int `json:"sequence" sql:"sequence"`
}

// NewProcedurePostgres Creates a client for interacting with the procedure relation in postgres DB
func NewProcedurePostgres(l hclog.Logger, db *pgx.Conn) *ProcedurePostgres {
	return &ProcedurePostgres{l, db}
}

// AddProcedure adds a procedure and stores it in the DB
func (ppg *ProcedurePostgres) AddProcedure(ctx context.Context, p *Procedure) error {
	p.ID = uuid.NewV4().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	ppg.logger.Info("creating procedure", "id", p.ID, "name", p.Name)
	_, err := ppg.db.Exec(ctx, InsertProcedureDML, p.ID, p.Name, p.Description, p.City, p.CreatedAt, p.UpdatedAt)
	return err
}

// GetAllProcedures gets all procedures from DB
func (ppg *ProcedurePostgres) GetAllProcedures(ctx context.Context) (Procedures, error) {
	ppg.logger.Info("fetching procedures from db")
	var p Procedures

	//rows, err := ppg.db.Query(context.Background(), GetAllProceduresDML)
	err := pgxscan.Select(ctx, ppg.db, &p, GetAllProceduresDML)
	if err != nil {
		ppg.logger.Error("Error fetching procedures from db", "error", err)
		return nil, err
	}

	for k, v := range p {
		var s Steps
		err = pgxscan.Select(ctx, ppg.db, &s, GetStepsWithProcedureID, v.ID)
		if err != nil {
			ppg.logger.Error("Error fetching procedure", "error", err)
			return nil, err
		}
		p[k].Steps = s
	}

	return p, nil
}

// UpdateProcedure procedure from DB using id
func (ppg *ProcedurePostgres) UpdateProcedure(ctx context.Context, p *Procedure) (*Procedure, error) {
	p.UpdatedAt = time.Now()

	_, err := ppg.db.Exec(ctx, UpdateProcedureDML, p.Name, p.Description, p.City, p.UpdatedAt, p.ID)
	if err != nil {
		ppg.logger.Error("Error updating procedure in db", "error", err)
		return nil, err
	}
	if p.StepsMapping != nil {
		_, err = ppg.db.Exec(ctx, DeleteProcedureStepsRelation, p.ID)
		if err != nil {
			ppg.logger.Error("Error deleting procedure step mapping", "error", err)
			return nil, err
		}
		err = ppg.MapProcedureSteps(context.Background(), p.ID, p.StepsMapping)
		if err != nil {
			ppg.logger.Error("Error mapping procedure steps", "error", err)
			return nil, err
		}
	}
	updatedProcedure, err := ppg.GetProcedure(ctx, p.ID)
	if err != nil {
		ppg.logger.Error("Failed fetching updated record from db", "error", err)
	}
	return updatedProcedure, nil
}

// GetProcedure fetch a procedure from DB using id
func (ppg *ProcedurePostgres) GetProcedure(ctx context.Context, id string) (*Procedure, error) {
	ppg.logger.Debug("querying for procedure", "id", id)
	var p Procedure
	err := pgxscan.Get(ctx, ppg.db, &p, GetProcedureDML, id)
	if err != nil {
		ppg.logger.Error("Error fetching procedure", "error", err)
		return nil, err
	}
	var s Steps
	err = pgxscan.Select(ctx, ppg.db, &s, GetStepsWithProcedureID, id)
	if err != nil {
		ppg.logger.Error("Error fetching procedure", "error", err)
		return nil, err
	}
	p.Steps = s
	return &p, nil
}

func (ppg *ProcedurePostgres) GetProcedureSteps(ctx context.Context, id string) (Steps, error) {
	ppg.logger.Info("fetching steps from db for procedure", "procedure", id)
	var s Steps

	err := pgxscan.Select(ctx, ppg.db, &s, GetProcedureStepsDML, id)
	if err != nil {
		ppg.logger.Error("Error fetching steps from db", "error", err)
		return nil, err
	}

	return s, nil
}

// DeleteProcedure deletes a procedure from the DB using its id
func (ppg *ProcedurePostgres) DeleteProcedure(ctx context.Context, id string) error {
	ppg.logger.Info("deleting procedure", "id", id)
	_, err := ppg.db.Exec(ctx, DeleteProcedureQ, id)
	if err != nil {
		ppg.logger.Error("Error deleting procedure in db", "error", err)
		return err
	}
	_, err = ppg.db.Exec(ctx, DeleteProcedureStepsRelation, id)
	if err != nil {
		ppg.logger.Error("Error deleting procedure step mapping", "error", err)
		return err
	}
	return nil
}

// MapProcedureSteps maps steps to procedures using relationship table
func (ppg *ProcedurePostgres) MapProcedureSteps(ctx context.Context, id string, m map[int]string) error {
	ppg.logger.Info("mapping steps for", "procedure", id)
	var lines []string
	for k, v := range m {
		line := fmt.Sprintf("('%s', '%s', %d)", id, v, k)
		lines = append(lines, line)
	}
	values := strings.Join(lines, ",")
	sql := fmt.Sprintf("%s %s", InsertProcedureStepsDML, values)
	_, err := ppg.db.Exec(ctx, sql)
	if err != nil {
		ppg.logger.Error("Error mapping procedure steps", "error", err)
		return err
	}
	return nil
}
