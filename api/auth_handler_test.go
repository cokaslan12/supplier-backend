package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"supplier-backend/db"
	"supplier-backend/types"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	params := types.CreateUser{
		Email:     "cokaslanmuzaffer@gmail.com",
		FirstName: "Muzaffer",
		LastName:  "Çokaslan",
		Password:  "123456",
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Fatal(err)
	}

	_, InsertedErr := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(InsertedErr)
	}

	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tDb := setup()
	defer tDb.tearDown(t)
	insertedUser := insertTestUser(t, tDb.store.UserStore)

	//MARK: INIT FIBER
	app := fiber.New()

	//MARK: HANDLER INITIALIZATION
	authHandler := NewAuthHandler(tDb.store)

	//MARK: AUTH API
	app.Post("/auth", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "cokaslanmuzaffer@gmail.com",
		Password: "123456",
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
	insertedUser := insertTestUser(t, tDb.store.UserStore)

	//MARK: INIT FIBER
	app := fiber.New()

	//MARK: HANDLER INITIALIZATION
	authHandler := NewAuthHandler(tDb.store)

	//MARK: AUTH API
	app.Post("/auth", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "cokaslanmuzaffer@gmail.com",
		Password: "wrongpassword",
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
