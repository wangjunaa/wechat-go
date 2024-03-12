package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"wechat/dao"
	Model "wechat/models"
)

var (
	friendPrefix    = "friendSet:"
	friendReqPrefix = "friendReq:"
)

func IsFriend(id1 string, id2 string) bool {
	list, err := GetFriendList(id1)
	if err != nil {
		log.Println("handler.user.IsFriend", err)
		return false
	}
	for _, user := range list {
		if user.ID == id2 {
			return true
		}
	}

	return false
}
func FriendReq(senderId string, receiveId string) error {
	m := &Model.Message{
		SenderId:   senderId,
		ReceiverId: receiveId,
		MsgType:    Model.MFriendReq,
	}
	//判断对方是否存在
	exist := IsUserExist(receiveId)
	if !exist {
		return errors.New("无法向不存在用户发送好友申请")
	}
	//判读是否为好友
	isFriend := IsFriend(senderId, receiveId)
	if isFriend {
		return errors.New("不可重复添加好友")
	}
	//将申请记录存入缓存
	err := dao.Rdb.SAdd(context.Background(), friendReqPrefix+receiveId, senderId).Err()
	if err != nil {
		return err
	}
	//若不是好友则发送申请
	err = SendMsg(m)
	return err
}
func FriendAgree(senderId string, receiveId string) error {
	m := &Model.Message{
		SenderId:   senderId,
		ReceiverId: receiveId,
		MsgType:    Model.MFriendAgree,
	}
	//判读是否为好友
	isFriend := IsFriend(senderId, receiveId)
	if isFriend {
		return errors.New("不可重复添加好友")
	}
	//判断是否发过申请
	exist, err := dao.Rdb.SIsMember(context.Background(), friendReqPrefix+senderId, receiveId).Result()
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("同意未请求申请")
	}
	//移除申请记录
	err = dao.Rdb.SRem(context.Background(), friendReqPrefix+senderId, receiveId).Err()
	if err != nil {
		return err
	}
	//添加好友并发送通知
	err = AddFriend(senderId, receiveId)
	if err != nil {
		return err
	}
	err = SendMsg(m)
	return err
}
func addFriendToRdb(id string, friendIDs ...string) error {
	for _, friendID := range friendIDs {
		err := dao.Rdb.SAdd(context.Background(), friendPrefix+id, friendID).Err()
		if err != nil {
			return err
		}
		err = dao.Rdb.SAdd(context.Background(), friendPrefix+friendID, id).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func addFriendToDB(id string, friendIDs ...string) error {
	for _, friendID := range friendIDs {
		friendShip := Model.FriendShip{
			UserId1: id,
			UserId2: friendID,
		}
		err := dao.DB.Create(friendShip).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func AddFriend(id1 string, id2 string) error {
	err := addFriendToDB(id1, id2)
	if err != nil {
		return err
	}
	err = addFriendToRdb(id1, id2)
	return err
}
func getFriendListFromDB(id string) ([]Model.ShowUser, error) {
	var friendShips []Model.FriendShip
	err := dao.DB.Preload("User1").Preload("User2").Where("user_id1=? or user_id2 = ?", id, id).Find(&friendShips).Error
	if err != nil {
		return nil, err
	}
	var friends []Model.ShowUser
	for _, friendShip := range friendShips {
		if friendShip.UserId1 == id {
			friends = append(friends, friendShip.User2.ToShowUser())
		} else {
			friends = append(friends, friendShip.User1.ToShowUser())
		}
	}
	return friends, nil
}
func getFriendListFromRdb(id string) ([]Model.ShowUser, error) {
	var friends []Model.ShowUser
	result, err := dao.Rdb.SMembers(context.Background(), friendPrefix+id).Result()
	if err != nil {
		return nil, err
	}
	for _, id := range result {
		user, err := FindUser(id)
		if err != nil {
			return nil, err
		}
		friends = append(friends, user.ToShowUser())
	}
	return friends, nil
}
func GetFriendList(id string) ([]Model.ShowUser, error) {
	friends, err := getFriendListFromRdb(id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if !errors.Is(err, redis.Nil) {
		return friends, nil
	}
	friends, err = getFriendListFromDB(id)
	if err != nil {
		return nil, err
	}
	var friendIDs []string
	for _, friend := range friends {
		friendIDs = append(friendIDs, friend.ID)
	}
	err = addFriendToRdb(id, friendIDs...)
	return friends, nil
}
func deleteFriendFromDB(id string, deletedId string) error {
	err := dao.DB.Model(Model.FriendShip{}).
		Where("user_id1 = ? AND user_id2 = ?", id, deletedId).
		Or("user_id1 = ? AND user_id2 = ?", deletedId, id).
		Delete(nil).Error
	return err

}
func deleteFriendFromRdb(id string, deletedId string) error {
	err := dao.Rdb.SRem(context.Background(), friendPrefix+id, deletedId).Err()
	if err != nil {
		return err
	}
	err = dao.Rdb.SRem(context.Background(), friendPrefix+deletedId, id).Err()
	return err
}
func DeleteFriend(id string, deletedId string) error {
	err := deleteFriendFromDB(id, deletedId)
	if err != nil {
		return err
	}
	err = deleteFriendFromRdb(id, deletedId)
	return nil
}
