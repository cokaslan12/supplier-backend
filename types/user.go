package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      int = 12
	minFirstNameLen int = 2
	minLastNameLen  int = 2
	minPasswordLen      = 7
)

type CreateUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (params CreateUser) Validate() map[string]string {

	errors := map[string]string{}

	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName lenght should be at least %d characters", minFirstNameLen)
	}

	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName lenght should be at least %d characters", minLastNameLen)
	}

	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password lenght should be at least %d characters", minPasswordLen)
	}

	if !isEmailValid(params.Email) {
		errors["email"] = "please enter a valid email address"
	}

	return errors
}

func (params UpdateUser) ToBSON() bson.M {
	maps := bson.M{}
	if len(params.FirstName) > 0 {
		maps["firstname"] = params.FirstName
	}

	if len(params.LastName) > 0 {
		maps["lastname"] = params.LastName
	}

	return maps
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty`
	FirstName         string             `bjson:"firstName" json:"firstName"`
	LastName          string             `bjson:"lastName" json:"lastName"`
	Email             string             `bjson:"email" json:"email"`
	EncryptedPassword string             `bjson:"encryptedPassord" json:"-"`
}

func NewUserFromParams(params CreateUser) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
