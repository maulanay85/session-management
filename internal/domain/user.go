package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID       string `gorm:"primaryKey;column:id" json:"id"`
	FullName string `gorm:"column:full_name" json:"fullName"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"-"`
}

type DummyUser map[string]User

func InitializeData() DummyUser {

	var dummyUsers = make(map[string]User)
	dummyUsers["test1@gmail.com"] = User{
		ID:       fmt.Sprintf("%d", time.Now().Unix()),
		FullName: "maulana",
		Email:    "test1@gmail.com",
		Password: "test123",
	}
	dummyUsers["test2@gmail.com"] = User{
		ID:       fmt.Sprintf("%d", time.Now().Unix()),
		FullName: "maulana 2",
		Email:    "test2@gmail.com",
		Password: "test456",
	}

	return dummyUsers
}
