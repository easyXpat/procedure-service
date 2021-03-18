package data 

// Procedure defines the structure for an API procedure 
// swagger:model
type Procedure struct {
	// unique id for the procedure
	//
	// required: true
	// min: 1
	// max length: 64
	ID int `json:"id" validate:"required"`

	// name for the procedure
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// description for the procedure
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// city for the procedure
	//
	// required: false
	// max length: 32
	City string `json:"city"`
}

type Procedures []*Procedure

//type ProcedureDB struct {
//
//}