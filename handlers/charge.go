package handlers

import (
	"github.com/easyXpat/procedure-service/data"
	"github.com/easyXpat/procedure-service/vendor/github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-hclog"
	"gorm.io/gorm"
)

// ChargeKey is used as a key for storing the Charge object in context at middleware
type ChargeKey struct{}

// Charge wraps instances needed to perform operations on step object
type Charge struct {
	logger          hclog.Logger
	db                *gorm.DB
	validator         *data.Validation
}

// NewCharge creates a new charge handler
func NewCharge(l hclog.Logger, db *gorm.DB, v *data.Validation) *Charge {
	return &Charge{
		l,
		db,
		v,
	}
}

func (c *Charge) SavePayment(charge *data.Charge) (err error) {
	if err = c.db.Create(charge).Error; err != nil {
		return err
	}
	return nil
}