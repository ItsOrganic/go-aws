package main

import (
	"fmt"
	"go-aws/lambda/app"
	"net/http"
    "go-aws/lambda/middleware"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
    Username string `json:"username"`
}

func HandleRequest(event MyEvent) (string, error) {
    if event.Username == "" {
        return "", fmt.Errorf("username cannot be empty")
    }
    return fmt.Sprintf("Successfully called by - %s", event.Username), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
    return events.APIGatewayProxyResponse{
        Body: "This is a protected path",
        StatusCode: http.StatusOK,
    }, nil
}

func main() {
    myApp := app.NewApp()
    //Start the lambda function using one ApiHandler RegisterUserHandler
    //lambda.Start(myApp.ApiHandler.RegisterUserHandler)

    lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        switch request.Path {
        case "/register":
            return myApp.ApiHandler.RegisterUserHandler(request)
        case "/login":
            return myApp.ApiHandler.LoginUser(request)
        case "/protected":
            return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)
        default:
            return events.APIGatewayProxyResponse{
                Body: "Not found ",
                StatusCode: http.StatusNotFound,
            }, nil
        }
    })



}

