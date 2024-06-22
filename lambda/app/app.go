package app

import (
	"go-aws/lambda/api"
	"go-aws/lambda/database"
)

type App struct {
    ApiHandler api.ApiHandler
}

func NewApp() App {

    db := database.NewDynamoDBClient()
    apiHandler := api.NewApiHandler(db)

    return App{
        ApiHandler: apiHandler,
    }

}
