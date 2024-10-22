package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptEncryptionCost = 12
	minFirstNameLen      = 2
	minLastNameLen       = 2
	minPasswordLen       = 7
)

type CreateUserParam struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"Password"`
}

func (param *CreateUserParam) MapToUser() (*User, error) {
	encrPw, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcryptEncryptionCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         param.FirstName,
		LastName:          param.LastName,
		Email:             param.Email,
		EncryptedPassword: string(encrPw),
	}, nil
}

func (param *CreateUserParam) Validate() []string {
	var errors []string
	if len(param.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen))
	}
	if len(param.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen))
	}
	if len(param.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password length should be at least %d characters", minPasswordLen))
	}

	if !isValidEmail(param.Email) {
		errors = append(errors, "Email is invalid")
	}
	return errors

}

func isValidEmail(email string) bool {
	regx := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regx.MatchString(email)
}

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}
