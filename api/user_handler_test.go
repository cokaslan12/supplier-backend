package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"supplier-backend/types"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	tDb := setup()
	defer tDb.tearDown(t)

	//MARK: INIT FIBER
	app := fiber.New()

	//MARK: HANDLER INITIALIZATION
	userHandler := NewUserHandler(tDb.store)

	//MARK: USERS API
	app.Post("/user", userHandler.HandlePostUser)
	params := types.CreateUser{
		FirstName: "muzaffer",
		LastName:  "çokaslan",
		Email:     "muzaffer@cokaslan.com",
		Password:  "muzaffer_cokaslan",
	}

	paramsJson, jsonErr := json.Marshal(params)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	reg := httptest.NewRequest("POST", "/user", bytes.NewReader(paramsJson))
	reg.Header.Add("Content-Type", "application/json")

	resp, regErr := app.Test(reg)
	if regErr != nil {
		t.Error(regErr)
	}

	if resp.StatusCode != 200 {
		t.Error("status code is not ok", resp.StatusCode)
	}

	bodyBytes, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		t.Error(ioErr)
	}

	var m map[string]any
	unMarshallErr := json.Unmarshal(bodyBytes, &m)
	if unMarshallErr != nil {
		t.Error(unMarshallErr)
	}

	userData, exist := m["data"]
	if !exist {
		t.Error("User not exist")
	}

	userJson, err := json.Marshal(userData)
	if err != nil {
		t.Error("Error marshalling map to JSON:", err)
	}

	var user types.User

	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the json response")
	}

	if err := json.Unmarshal(userJson, &user); err != nil {
		t.Error("Error unmarshalling JSON to struct:", err)
	}

	if user.FirstName != params.FirstName {

		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}

}
