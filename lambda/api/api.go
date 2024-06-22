package api

import (
	"encoding/json"
	"fmt"
	"go-aws/lambda/database"
	"go-aws/lambda/types"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
    dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
    return ApiHandler{
        dbStore: dbStore,
    }
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
    var registerUser types.RegisterUser

    err := json.Unmarshal([]byte(request.Body), &registerUser)
    if err != nil {
        return events.APIGatewayProxyResponse{
            Body: "Invalid Request ",
            StatusCode: http.StatusNotFound,
        }, err
    }

    if registerUser.Username == "" || registerUser.Password == "" {
        return events.APIGatewayProxyResponse{
            Body: "Invalid Request - fields empty",
            StatusCode: http.StatusNotFound,
        }, err
    }

    userExists, err := api.dbStore.DoesUserExist(registerUser.Username)
    if err != nil {
        return events.APIGatewayProxyResponse{
            Body: "Internal server error",
            StatusCode: http.StatusInternalServerError,
        }, err
    }
    if userExists {
        return events.APIGatewayProxyResponse{
            Body: "User already exists",
            StatusCode: http.StatusConflict,
        }, nil
    }

    user, err := types.NewUser(registerUser)
    if err != nil {
        return events.APIGatewayProxyResponse{
            Body: "Internal server error",
            StatusCode: http.StatusInternalServerError,
        }, fmt.Errorf("could not create new user %w", err)
    }

    err = api.dbStore.InserUser(user)
    if err != nil {
        return events.APIGatewayProxyResponse{
            Body: "internal server error",
            StatusCode: http.StatusInternalServerError,
        }, nil
    }

    return events.APIGatewayProxyResponse{
        Body: "Successfully registered user",
        StatusCode: http.StatusOK,
    }, nil

}

func (api *ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid login credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	accessToken := types.CreateToken(user)
    successMsg := fmt.Sprintf("Access token: %s",string(accessToken))

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}
