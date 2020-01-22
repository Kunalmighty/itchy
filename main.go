package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)


var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type book struct {
    NAME   string `json:"name"`
    Description  string `json:"description"`
    Size string `json:"size"`
    Price string `json:"price"`
}

//routing incoming requests
func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    switch req.HTTPMethod {
    case "GET":
        return show(req)
    case "POST":
        return create(req)
    default:
        return clientError(http.StatusMethodNotAllowed)
    }
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    // Get the 'name' query string parameter from the request and validate it.
    name := req.QueryStringParameters["name"]

    // Fetch the scratcher record from the database based on the name value.
    bk, err := getItem(name)
    if err != nil {
        return serverError(err)
    }
    if bk == nil {
        return clientError(http.StatusNotFound)
    }

    // The APIGatewayProxyResponse.Body field needs to be a string, so marshal the scratcher record into JSON.
    js, err := json.Marshal(bk)
    if err != nil {
        return serverError(err)
    }

    // Return a response with a 200 OK status and the JSON scratcher record as the body.
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

//Handling POST operation, including error checks
func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
    if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }
    bk := new(book)
    err := json.Unmarshal([]byte(req.Body), bk)
    if err != nil {
        return clientError(http.StatusUnprocessableEntity)
    }
    if bk.Description == "" || bk.Size == "" || bk.Price == "" {
        return clientError(http.StatusBadRequest)
    }
    err = putItem(bk)
    if err != nil {
        return serverError(err)
    }
    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Headers:    map[string]string{"Location": fmt.Sprintf("/scratchers?name=%s", bk.NAME)},
    }, nil
}

// A helper for handling errors. This logs any error to os.Stderr and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body: http.StatusText(http.StatusInternalServerError),
    }, nil
}

// Similarly adding a helper for sending responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body: http.StatusText(status),
    }, nil
}

func main() {
    lambda.Start(router)
}