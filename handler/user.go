package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
	"wechat/dao"
	"wechat/models"
	"wechat/utils/encryption"
	"wechat/utils/token"
)

var (
	userPrefix     = "user:"
	expirationTime = 7 * 24 * 60 * 60 * time.Second
	delayDelTime   = 3 * time.Second
)

func setUserNotExistToRdb(id string) error {
	err := dao.Rdb.Set(context.Background(), userPrefix+id, "", expirationTime).Err()
	return err
}

// 从缓存获取用户数据
func findUserFromRdb(id string) (user *models.UserBasic, err error) {
	result, err := dao.Rdb.Get(context.Background(), userPrefix+id).Result()
	if err != nil {
		return user, err
	}
	//判断数据库中是否存在值
	if result == "" {
		return nil, gorm.ErrRecordNotFound
	}
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		return user, err
	}
	//log.Println("findUserFromRdb:", user)
	return user, nil
}

// 从数据库获取用户数据
func findUserFromDB(id string) (user *models.UserBasic, err error) {
	user = &models.UserBasic{}
	if err = dao.DB.Where("id=?", id).First(user).Error; err != nil {
		return user, err
	}
	//log.Println("findUserFromDB:", user)
	return user, err
}

// FindUser 用户不存在，err返回gorm.ErrRecordNotFound
func FindUser(id string) (user *models.UserBasic, err error) {
	user, err = findUserFromRdb(id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	//在redis中查询到结果则直接返回
	if err == nil {
		return user, nil
	}
	//从数据库中获取
	user, err = findUserFromDB(id)
	if err != nil {
		//若数据库不存在信息，则将无信息记录进缓存
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := setUserNotExistToRdb(id); err != nil {
				return nil, err
			}
		}
		return nil, err
	}
	//将数据加入rdb
	err = addUserToRdb(user)
	return user, err
}

func addUserToRdb(user *models.UserBasic) error {
	marshal, err := json.Marshal(user)
	err = dao.Rdb.Set(context.Background(), userPrefix+user.ID, marshal, expirationTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func addUserToDB(user models.UserBasic) error {
	if err := dao.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func CreateUser(phone string, userName string, password string) (string, error) {
	id := phone
	u, err := FindUser(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if u != nil {
		return "", errors.New("用户名重复")
	}
	user := models.UserBasic{
		Phone:    phone,
		UserName: userName,
		ID:       id,
		Password: encryption.Encode(password),
	}
	if err := addUserToDB(user); err != nil {
		return "", err
	}
	//防止缓存还认为数据库无数据
	if err := deleteUserFromRdb(user); err != nil {
		return "", err
	}
	t, err := token.CreateToken(id)
	if err != nil {
		return "", err
	}
	return t, nil
}

func deleteUserFromRdb(user models.UserBasic) error {
	err := dao.Rdb.Del(context.Background(), userPrefix+user.ID).Err()
	return err
}

func deleteUserFromDB(user models.UserBasic) error {
	begin := dao.DB.Begin()
	//删除好友关系
	err := begin.Model(&models.FriendShip{}).
		Where("user_id1 =? or user_id2 =?", user.ID, user.ID).Delete(nil).Error
	if err != nil {
		begin.Rollback()
		return err
	}
	//删除群组关系
	err = begin.Table("user_groups").Where("user_basic_id=?", user.ID).Delete(nil).Error
	if err != nil {
		begin.Rollback()
		return err
	}
	//删除用户
	err = begin.Delete(user).Error
	if err != nil {
		begin.Rollback()
		return err
	}
	err = begin.Commit().Error
	if err != nil {
		begin.Rollback()
		return err
	}
	return err
}

func DeleteUser(id string, password string) error {
	user := models.UserBasic{ID: id}
	res := dao.DB.Model(user).First(&user)
	if res.Error != nil {
		return res.Error
	}

	if !encryption.IsEqualAfterEncode(password, user.Password) {
		return errors.New("密码错误")
	}

	err := deleteUserFromDB(user)
	if err != nil {
		return err
	}
	err = deleteUserFromRdb(user)
	return err
}

func updateUserFromDB(user models.UserBasic) error {
	err := dao.DB.Model(user).Updates(&user).Error
	return err
}

func UpdateUser(user models.UserBasic) error {
	//延迟双删
	if err := deleteUserFromRdb(user); err != nil {
		return err
	}
	if err := updateUserFromDB(user); err != nil {
		return err
	}
	go func() {
		time.Sleep(delayDelTime)
		_ = deleteUserFromRdb(user)
	}()
	return nil
}

func Login(id string, password string) (string, error) {

	var user models.UserBasic
	res := dao.DB.Where("id = ? and password = ?", id, encryption.Encode(password)).First(&user)
	if res.Error != nil {
		return "", errors.New("用户名或密码错误")
	}
	now := time.Now()
	tk, err := token.CreateToken(id)
	if err != nil {
		log.Println("handler.user_handle.Login:", err)
		return "", err
	}
	res = dao.DB.Model(&user).Updates(models.UserBasic{LoginTime: &now})
	if res.Error != nil {
		//log.Println("handler.user_handle.Login:", err)
		return "", err
	}
	return tk, nil
}

func IsUserExist(phone string) bool {
	user, _ := FindUser(phone)
	return user != nil
}
