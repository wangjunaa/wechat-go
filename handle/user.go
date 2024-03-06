package handle

import (
	"demo/config"
	Model "demo/models"
	"demo/tools/encryption"
	"demo/tools/token"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// UserListToShow 将userBasic列表转为showUser列表
func UserListToShow(users []Model.UserBasic) []Model.ShowUser {
	var su []Model.ShowUser
	for _, user := range users {
		su = append(su, user.ToShowUser())
	}
	return su
}

func getUserFromRdb(id string) (user *Model.UserBasic, err error) {
	result, err := config.Rdb.HGet(config.BgCtx, hUserKey, id).Result()
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		return user, err
	}
	//log.Println("getUserFromRdb:", user)
	return user, nil
}
func getUserFromDB(id string) (user *Model.UserBasic, err error) {
	if err = config.DB.First(&user, "id=?", id).Error; err != nil {
		return user, err
	}
	//log.Println("getUserFromDB:", user)
	return user, err
}

// GetUser 当发生错误或未找到记录 user返回nil
func GetUser(id string) (user *Model.UserBasic, err error) {
	user, err = getUserFromRdb(id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return user, err
	}
	//在redis中查询到结果则直接返回
	if err == nil {
		return user, nil
	}
	//从数据库中获取
	user, err = getUserFromDB(id)
	if err != nil {
		return nil, err
	}
	//将数据加入rdb
	err = addUserToRdb(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func addUserToRdb(user *Model.UserBasic) error {
	rUserMux.Lock()
	defer rUserMux.UnLock()
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = config.Rdb.HSet(config.BgCtx, hUserKey, user.ID, bytes).Err()
	if err != nil {
		return err
	}

	return nil
}
func addUserToDB(user Model.UserBasic) error {
	sUserMux.Lock()
	defer sUserMux.UnLock()
	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
func CreateUser(json Model.RegisterJson) error {
	u, err := GetUser(json.ID)
	if err != nil && u != nil {
		return err
	}
	if u != nil {
		return errors.New("用户名重复")
	}

	var user Model.UserBasic
	password := encryption.Encode(json.Password1)
	user.Phone = json.ID
	user.UserName = json.UserName
	user.ID = json.ID
	user.Password = password
	if err := addUserToDB(user); err != nil {
		return err
	}
	return nil
}

func deleteUserFromRdb(user Model.UserBasic) error {
	rUserMux.Lock()
	defer rUserMux.UnLock()
	err := config.Rdb.HDel(config.BgCtx, hUserKey, user.ID).Err()
	return err
}
func deleteUserFromDB(user Model.UserBasic) error {
	sUserMux.Lock()
	defer sUserMux.UnLock()
	err := config.DB.Delete(user).Error
	return err
}
func DeleteUser(id string, password string) error {
	user := Model.UserBasic{}
	res := config.DB.Where("id = ?", id).First(&user)

	if res.Error != nil {
		return res.Error
	}
	if !encryption.IsEqualAfterEncode(password, user.Password) {
		return errors.New("密码错误")
	}

	//延时双删
	err := deleteUserFromRdb(user)
	if err != nil {
		return err
	}
	err = deleteUserFromDB(user)
	if err != nil {
		return err
	}
	time.Sleep(3)
	err = deleteUserFromRdb(user)
	return err
}

func updateUserFromDB(user Model.UserBasic) error {
	sUserMux.Lock()
	defer sUserMux.UnLock()
	err := config.DB.Updates(user).Error
	return err
}
func UpdateUser(user Model.UserBasic) error {
	//延迟双删

	if err := deleteUserFromRdb(user); err != nil {
		return err
	}
	if err := updateUserFromDB(user); err != nil {
		return err
	}
	time.Sleep(3)
	if err := deleteUserFromRdb(user); err != nil {
		return err
	}
	return nil
}

func Login(id string, password string) (string, error) {

	var user Model.UserBasic
	res := config.DB.Where("id = ? and password = ?", id, encryption.Encode(password)).First(&user)
	if res.Error != nil {
		return "", errors.New("用户名或密码错误")
	}
	now := time.Now()
	tk, err := token.CreateToken(id)
	if err != nil {
		log.Println("handle.user_handle.Login:", err)
		return "", err
	}
	res = config.DB.Model(&user).Updates(Model.UserBasic{LoginTime: &now})
	if res.Error != nil {
		log.Println("handle.user_handle.Login:", err)
		return "", err
	}
	return tk, nil
}
