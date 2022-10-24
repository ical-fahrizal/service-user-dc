package controllers

import (
	conf "kafka-user-dc/config"
	"log"
)

type User struct {
	// Id               uint   `json:"id"`
	// Status           bool   `json:"status"`
	Locked           bool   `json:"locked" default:"false"`
	FailedLoginCount int    `json:"failedLoginCount" default:"0"`
	Fullname         string `json:"fullname"`
	Email            string `json:"email"`
	Hp               string `json:"hp"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	// FcmToken         string `json:"fcmToken"`
	// ApiKey           string `json:"apiKey"`
	// CompanyId        int    `json:"companyId"`
	// ParentUserId     int    `json:"parentUserId"`
	// RootId           int    `json:"rootId"`
	// OfficeId         int    `json:"officeId"`
	// DepartementId    int    `json:"departementId"`
	// LastLoginTime    string `json:"lastLoginTime"`
	// Token            string `json:"token"`
	// Total            uint   `json:"-"`
	// MsgError         string `json:"-"`
}

// CreateUser new user
func CreateUser(user *User) (*User, error) {
	userData := &User{}
	keyRedis := "user:profile:" + user.Username
	res, err := conf.GetReJson().JSONSet(keyRedis, ".", user)
	if err != nil {
		log.Printf("Failed to JSONSet")
		return userData, err
	}

	var err1 error
	if res.(string) == "OK" {

	} else {
		return userData, err1
	}

	return user, nil

}
