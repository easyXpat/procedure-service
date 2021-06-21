# EasyXpat

## Procedure Service

The procedure service handles operations related to procedures offered by EasyXpat

### Running the procedure service locally

*Create postgres DB instance*

`docker run -d -p 5432:5432 --env-file db.env postgres:alpine`

*Run service*

`go build .`
`go run main.go`

### API Documentation

https://easyxpat-procedure.herokuapp.com/docs

## Methods Usage Examples

### GET

Get all procedures 

`curl https://easyxpat-procedure.herokuapp.com/procedures`

Get specific procedure

`curl https://easyxpat-procedure.herokuapp.com/procedures/{id}`

### POST 

Create procedure

`curl -X POST -d {json} https://easyxpat-procedure.herokuapp.com/procedures`

```json
{
	"name": "Blue Card",
	"description": "Process to apply for a Blue Card",
	"city": "Frankfurt"
}
```

### PUT 

Update procedure 

`curl -X PUT -d {json} https://easyxpat-procedure.herokuapp.com/procedures/{id}`

```json
{
	"id": "09e1ae77-52b1-449b-801c-bb90357b7ed4"
	"name": "Blue Card",
	"description": "Process to apply for a Blue Card",
	"city": "Frankfurt"
}
```

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
