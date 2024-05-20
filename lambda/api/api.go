package api

import (
	"fmt"
	"go-aws/lambda/database"
	"go-aws/lambda/types"
)

type ApiHandler struct {
    dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
    return ApiHandler{
        dbStore: dbStore,
    }
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
    if event.Username == "" || event.Password == "" {
        return fmt.Errorf("request has empty parameters")
    }

    userExists, err := api.dbStore.DoesUserExist(event.Username)
    if err != nil {
        return fmt.Errorf("there an error checking if user exists %w", err)
    }
    if userExists {
        return fmt.Errorf("a user with that username exists")
    }

    err = api.dbStore.InserUser(event)
    if err != nil {
        return fmt.Errorf("Error registering the user %w", err)
    }

    return nil

}
