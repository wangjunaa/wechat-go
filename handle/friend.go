package handle

import (
	"demo/config"
	Model "demo/models"
	"errors"
	"log"
)

func IsFriend(id1 string, id2 string) bool {
	list, err := GetFriendList(id1)
	if err != nil {
		log.Println("handle.user.IsFriend", err)
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
	//判读是否为好友
	isFriend := IsFriend(senderId, receiveId)
	if isFriend {
		return errors.New("不可重复添加好友")
	}
	//若不是好友则发送申请
	err := SendMsg(m)
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
	//若不是好友则发送通知并添加好友
	err := AddFriend(senderId, receiveId)
	if err != nil {
		return err
	}
	err = SendMsg(m)
	return err
}
func AddFriend(id1 string, id2 string) error {
	sFriendMux.Lock()
	defer sFriendMux.UnLock()
	friendShip := Model.FriendShip{
		UserId1: id1,
		UserId2: id2,
	}
	err := config.DB.Create(friendShip).Error
	return err
}
func GetFriendList(id string) ([]Model.ShowUser, error) {
	var friendShips []Model.FriendShip
	err := config.DB.Preload("User1").Preload("User2").Where("user_id1=? or user_id2 = ?", id, id).Find(&friendShips).Error
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
func DeleteFriend(id string, deletedId string) error {
	sFriendMux.Lock()
	defer sFriendMux.UnLock()
	var friendship Model.FriendShip
	if err := config.DB.
		Where("user_id1 = ? AND user_id2 = ?", id, deletedId).
		Or("user_id1 = ? AND user_id2 = ?", deletedId, id).
		First(&friendship).Error; err != nil {
		return err
	}

	// 删除好友关系记录
	if err := config.DB.Delete(&friendship).Error; err != nil {
		return err
	}
	return nil
}
