# EasyXpat

## Procedure Service

The procedure service handles operations related to procedures offered by EasyXpat

### Running the procedure service

### Running local

`docker run -d -p 5432:5432 --env-file db.env postgres:alpine`
`go run .`


## Method Usage
### GET

#### Example


### DELETE

#### Example

### POST
A key/value pair for the "name", "description", and "city" keys is required in the JSON body of the POST request.
#### Example
```json
{
	"name": "Blue Card",
	"description": "Process to apply for a Blue Card",
	"city": "Frankfurt"
}
```

### PUT