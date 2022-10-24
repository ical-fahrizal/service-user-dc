package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	conf "kafka-user-dc/config"

	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"
	// "github.com/satori/go.uuid"
)

// User struct
// type User struct {
// 	gorm.Model
// 	Username string `gorm:"unique_index;not null" json:"username"`
// 	Email    string `gorm:"unique_index;not null" json:"email"`
// 	Password string `gorm:"not null" json:"password"`
// 	Names    string `json:"names"`
// }

type User struct {
	Id               string `json:"id" default:"0"`
	Status           string `json:"status" default:"1"`
	Locked           string `json:"locked" default:"0"`
	FailedLoginCount string `json:"failedLoginCount" default:"0"`
	Fullname         string `json:"fullname"`
	Email            string `json:"email"`
	Hp               string `json:"hp"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	CompanyId        string `json:"companyId" default:"0"`
	RoleId           string `json:"roleId" default:"0"`
	OfficeId         string `json:"officeId" default:"0"`
	DepartementId    string `json:"departementId" default:"0"`
	// Total            uint   `json:"-"`
	// MsgError         string `json:"-"`
}

type MessageKafka struct {
	Topic string `json:"topic"`
	Proc  string `json:"proc"`
	Data  string `json:"data"`
	From  string `json:"-"`
	Exp   int64  `json:"exp"`
}

// data diperbaikin
func (user *User) Validated() *User {

	generateId := uuid.Must(uuid.NewV4())
	user.Id = fmt.Sprintf("%v", generateId)

	if user.Status == "" {
		user.Status = "1"
	}
	if user.Locked == "" {
		user.Locked = "0"
	}
	if user.FailedLoginCount == "" {
		user.FailedLoginCount = "0"
	}
	if user.CompanyId == "" {
		user.CompanyId = "0"
	}
	if user.RoleId == "" {
		user.RoleId = "0"
	}
	if user.OfficeId == "" {
		user.OfficeId = "0"
	}
	if user.DepartementId == "" {
		user.DepartementId = "0"
	}

	return user
}

func (user *User) UpdateData(dataRequest *User) *User {
	if dataRequest.Locked != "" {
		user.Locked = dataRequest.Locked
	}
	if dataRequest.Fullname != "" {
		user.Fullname = dataRequest.Fullname
	}
	if dataRequest.Email != "" {
		user.Email = dataRequest.Email
	}
	if dataRequest.Hp != "" {
		user.Hp = dataRequest.Hp
	}
	if dataRequest.Password != "" {
		user.Password = dataRequest.Password
	}
	if dataRequest.CompanyId != "" {
		user.CompanyId = dataRequest.CompanyId
	}
	if dataRequest.RoleId != "" {
		user.RoleId = dataRequest.RoleId
	}
	if dataRequest.OfficeId != "" {
		user.OfficeId = dataRequest.OfficeId
	}
	if dataRequest.DepartementId != "" {
		user.DepartementId = dataRequest.DepartementId
	}

	if dataRequest.Status != "" {
		user.Status = dataRequest.Status
	}

	return user
}

func GetUser(username string) (*User, error) {
	user := &User{}
	keyRedis := "user:profile:" + username
	userRedis, err := redis.Bytes(conf.GetReJson().JSONGet(keyRedis, "."))
	if err != nil {
		log.Printf("err : %v", err)
		return user, nil
	}

	err = json.Unmarshal(userRedis, &user)
	if err != nil {
		log.Printf("err : %v", err)
		return user, err
	}

	return user, nil
}

// CreateUser new user
func CreateUser(user *User) (*User, error) {
	//check user exists
	checkUser, err := GetUser(user.Username)
	if err != nil {
		return user, nil
	}

	//user exists
	if checkUser.Username != "" {
		log.Printf("data user is exists")
		return user, nil
	}

	//user not exists
	resp := user.Validated()

	userData := &User{}
	dataJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error parse to JSON : %v\n", err)
		return userData, nil
	}
	log.Printf("dataJson : %v", string(dataJson))

	keyRedis := "user:profile:" + resp.Username
	log.Printf("keyRedis : %v", keyRedis)
	res, err := conf.GetReJson().JSONSet(keyRedis, ".", resp)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return userData, err
	}

	var err1 error
	if res.(string) == "OK" {

	} else {
		return userData, err1
	}

	return resp, nil

}

func UpdateUser(user *User) (*User, error) {
	//check user exists
	checkUser, err := GetUser(user.Username)
	if err != nil {
		return user, nil
	}

	//user exists
	if checkUser.Username == "" {
		log.Printf("data user not is exists")
		return user, nil
	}

	log.Printf("checkUser : %v", checkUser)

	newData := checkUser.UpdateData(user)

	userData := &User{}
	keyRedis := "user:profile:" + newData.Username

	res, err := conf.GetReJson().JSONSet(keyRedis, ".", newData)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return userData, err
	}

	var err1 error
	if res.(string) == "OK" {

	} else {
		return userData, err1
	}

	newData.Password = ""
	return newData, nil
}

func DeleteUser(user *User) (*User, error) {
	//check user exists
	checkUser, err := GetUser(user.Username)
	if err != nil {
		return user, nil
	}

	//user exists
	if checkUser.Username == "" {
		log.Printf("data user is not exists")
		return user, nil
	}

	userData := &User{}
	keyRedis := "user:profile:" + user.Username

	res, err := conf.GetReJson().JSONSet(keyRedis, "status", "0")
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return userData, err
	}

	var err1 error
	if res.(string) == "OK" {

	} else {
		return userData, err1
	}

	return user, nil
}

func SearchUser(user *User) ([]User, error) {

	var userData []User
	docs, err := dataFind(user)
	if err != nil {
		return userData, err
	}

	userData = docs.Values
	keyRedis := "user:profile:" + user.Username
	res, err := conf.GetReJson().JSONSet(keyRedis, ".", user)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return userData, err
	}

	var err1 error
	if res.(string) == "OK" {

	} else {
		return userData, err1
	}

	return userData, nil
}

func dataFind(user *User) (*Document, error) {

	docs := &Document{}
	docsError := &Document{}

	ctx := context.Background()
	redisUser := conf.GetRedis().Do(ctx, "FT.SEARCH", "user-json-idx", "@status:1").Val()
	res, err := json.Marshal(redisUser)
	if err != nil {
		log.Printf("Error : %s\n", err)
		return docsError, err
	}
	// log.Printf("hasilres : %v", string(res))

	var dat []interface{}
	if err := json.Unmarshal(res, &dat); err != nil {
		log.Println(err)
		return docsError, nil
	}

	total := 0
	if total, err = conf.Int(dat[0], nil); err != nil {
		log.Printf("err : %v", err)
		return docsError, err
	}

	docs.Total = total
	total = total * 2
	skip := 2

	for i := skip; i <= total; i += skip {
		// log.Printf("dat[%v] : %v", i, dat[i])
		// log.Printf("dat[%v] : %v", i-1, dat[i-1])
		strKey := fmt.Sprintf("%v", dat[i-1])
		// log.Printf("strKey : %v", dat[i-1])

		strData := fmt.Sprintf("%v", dat[i])
		strData = strings.ReplaceAll(strData, "[$", "")
		strData = strings.ReplaceAll(strData, "]", "")

		log.Printf("strData : %v", strData)
		strData = strings.Trim(strData, "")
		data := make(map[string]interface{})
		data[strKey] = strData
		docs.Data = append(docs.Data, data)

		var valuesData User
		if err := json.Unmarshal([]byte(strData), &valuesData); err != nil {
			log.Println(err)
			return docsError, err
		}
		docs.Values = append(docs.Values, valuesData)

	}

	// // log.Printf("hasil : %v", docs)
	// docsByte, err := json.Marshal(docs)
	// if err != nil {
	// 	log.Printf("Error : %s\n", err)
	// 	return docsError, err
	// }
	// log.Printf("hasilres : %v", string(docsByte))

	return docs, nil
}

type Document struct {
	Total  int
	Data   []map[string]interface{}
	Values []User
}
