package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    //"regexp"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

//var nameRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type book struct {
    NAME   string `json:"name"`
    Description  string `json:"description"`
    Size string `json:"size"`
    Price string `json:"price"`
}

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
    name := req.QueryStringParameters["name"]
    //if !nameRegexp.MatchString(name) {
        //return clientError(http.StatusBadRequest)
    //}
    bk, err := getItem(name)
    if err != nil {
        return serverError(err)
    }
    if bk == nil {
        return clientError(http.StatusNotFound)
    }
    js, err := json.Marshal(bk)
    if err != nil {
        return serverError(err)
    }
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
    if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }
    bk := new(book)
    err := json.Unmarshal([]byte(req.Body), bk)
    if err != nil {
        return clientError(http.StatusUnprocessableEntity)
    }
    //if !nameRegexp.MatchString(bk.NAME) {
      //  return clientError(http.StatusBadRequest)
    //}
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

func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body: http.StatusText(http.StatusInternalServerError),
    }, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body: http.StatusText(status),
    }, nil
}

func main() {
    lambda.Start(router)
}