package user

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)

type Model struct {
	ID 	int64 `json:"id"`
	Firstname	string `json:"firstname"`
	Username	string `json:"username"`
	password	string
	HashedPassword	[]byte `json:"-"`
	CreateAt	time.Time `json:"creat_at"`
}

func GeneratePassword(userPassword string) ([]byte,error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword),bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed,[]byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}