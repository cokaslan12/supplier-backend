package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"supplier-backend/db/fixtures"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	tDb := setup()
	defer tDb.tearDown(t)
	insertedUser := fixtures.AddUser(tDb.store, "muzaffer", "cokaslan", false)

	//MARK: INIT FIBER
	app := fiber.New()

	//MARK: HANDLER INITIALIZATION
	authHandler := NewAuthHandler(tDb.store)

	//MARK: AUTH API
	app.Post("/auth", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "muzaffer@cokaslan.com",
		Password: "muzaffer_cokaslan",
	}

	paramsJson, jsonErr := json.Marshal(params)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	reg := httptest.NewRequest("POST", "/auth", bytes.NewReader(paramsJson))
	reg.Header.Add("Content-Type", "application/json")

	resp, regErr := app.Test(reg)
	if regErr != nil {
		t.Fatal(regErr)
	}

	if resp.StatusCode != 200 {
		t.Fatal("status code is not ok", resp.StatusCode)
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

	data, exist := m["data"].(map[string]any)
	if !exist {
		t.Error("User not exist")
	}

	dataJson, dataErr := json.Marshal(data)
	if dataErr != nil {
		t.Fatal(dataErr)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(dataJson, &authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected the JWT TOKEN to be present in the auth response")
	}

	//Set the encryted password to an empty string, because we do not return that in any json response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatal("expected user to be present in the inserted user")
	}

}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tDb := setup()
	defer tDb.tearDown(t)
	fixtures.AddUser(tDb.store, "muzaffer", "cokaslan", false)

	//MARK: INIT FIBER
	app := fiber.New()

	//MARK: HANDLER INITIALIZATION
	authHandler := NewAuthHandler(tDb.store)

	//MARK: AUTH API
	app.Post("/auth", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "muzaffer@cokaslan.com",
		Password: "wrong_password",
	}

	paramsJson, jsonErr := json.Marshal(params)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	reg := httptest.NewRequest("POST", "/auth", bytes.NewReader(paramsJson))
	reg.Header.Add("Content-Type", "application/json")

	resp, regErr := app.Test(reg)
	if regErr != nil {
		t.Fatal(regErr)
	}

	if resp.StatusCode != 400 {
		t.Fatal("status code is not bad request", resp.StatusCode)
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

	data, exist := m["error"].(string)
	if !exist {
		t.Error("error not exist")
	}

	if data != "invalid credentials" {
		t.Fatalf("expected response error to be <invalid credentials> but go %s", data)
	}

}
