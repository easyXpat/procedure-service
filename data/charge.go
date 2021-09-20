package data

import (
	"gorm.io/gorm"
)

type Charge struct {
	gorm.Model
	// amount to be charged
	//
	// required: true
	Amount int64 `json:"amount" validate:"required" sql:"amount"`

	// email to send receipt to
	//
	// required: false
	// max length: 200
	ReceiptEmail string `json:"receipt_email" validate:"required" sql:"receipt_email"`

	// procedure for the charge
	//
	// required: true
	ProcedureID string `json:"procedure_id" sql:"procedure_id"`

}

func (c *Charge) TableName() string {
	return "charge"

}