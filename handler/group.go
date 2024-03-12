package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"time"
	"wechat/dao"
	Model "wechat/models"
)

var (
	groupDataPreFix = "group:"
	groupReqPreFix  = "groupReq:"
)

func setGroupNotExist(gid string) error {
	err := dao.Rdb.Set(context.Background(), groupDataPreFix+gid, "", expirationTime).Err()
	return err
}

func getGroupFromDB(gid string) (*Model.GroupBasic, error) {
	group := &Model.GroupBasic{}
	err := dao.DB.Preload("Members").First(&group, "id=?", gid).Error
	if err != nil {
		return nil, err
	}
	return group, nil
}

func getGroupFromRdb(gid string) (*Model.GroupBasic, error) {
	result, err := dao.Rdb.Get(context.Background(), groupDataPreFix+gid).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, gorm.ErrRecordNotFound
	}
	group := &Model.GroupBasic{}
	err = json.Unmarshal([]byte(result), group)
	if err != nil {
		return nil, err
	}
	return group, nil

}

func addGroupToRdb(group *Model.GroupBasic) error {
	data, err := json.Marshal(group)
	if err != nil {
		return err
	}
	err = dao.Rdb.Set(context.Background(), groupDataPreFix+strconv.Itoa(int(group.ID)), data, expirationTime).Err()
	return err
}

func GetGroup(gid string) (*Model.GroupBasic, error) {
	group, err := getGroupFromRdb(gid)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	//redis查询到值
	if !errors.Is(err, redis.Nil) {
		return group, nil
	}

	group, err = getGroupFromDB(gid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := setGroupNotExist(gid); err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	err = addGroupToRdb(group)
	return group, err
}

func deleteGroupFromDB(gid string) error {
	err := dao.DB.Table("user_groups").Where("group_basic_id=?", gid).Delete(nil).Error
	if err != nil {
		return err
	}
	err = dao.DB.Model(Model.GroupBasic{}).Where("id=?", gid).Delete(nil).Error
	return err
}

func deleteGroupFromRdb(gid string) error {
	err := dao.Rdb.Del(context.Background(), groupDataPreFix+gid).Err()
	return err
}

func InitGroup(ownerId string, membersIDs []string) (*Model.GroupBasic, error) {
	group := &Model.GroupBasic{
		Name:    "",
		OwnerId: ownerId,
		Icon:    nil,
		Type:    0,
	}
	for i, id := range membersIDs {
		user, err := FindUser(id)
		if err != nil {
			return nil, err
		}
		//初始化群名
		if i < 3 {
			group.Name += user.UserName
			if i < 2 {
				group.Name += "、"
			}
		}
		group.Members = append(group.Members, *user)
	}
	return group, nil
}

func CreateGroup(ownerId string, membersIDs []string) (*Model.GroupBasic, error) {
	group, err := InitGroup(ownerId, membersIDs)
	err = dao.DB.Create(group).Error
	if err != nil {
		return nil, err
	}
	err = deleteGroupFromRdb(strconv.Itoa(int(group.ID)))
	return group, err
}

func getGroupOwner(gid string) (string, error) {
	group, err := GetGroup(gid)
	if err != nil {
		return "", err
	}
	return group.OwnerId, nil
}

// CheckOwner 若非群主则报错
func CheckOwner(ownerId string, groupId string) error {
	group, err := GetGroup(groupId)
	if err != nil {
		return err
	}
	if group.OwnerId != ownerId {
		return errors.New("非群主操作")
	}
	return nil
}

func DeleteGroup(gid string) error {
	if err := deleteGroupFromDB(gid); err != nil {
		return err
	}
	err := deleteGroupFromRdb(gid)
	return err
}

func RemoveGroupMember(groupId string, deletedId string) (*Model.GroupBasic, error) {
	group, err := GetGroup(groupId)
	if err != nil {
		return nil, err
	}
	// 找到要删除的成员并移除
	found := false
	for i, member := range group.Members {
		if member.ID == deletedId {
			group.Members = append(group.Members[:i], group.Members[i+1:]...)
			found = true
			break
		}
	}
	if found == false {
		return nil, errors.New("群中无此成员")
	}

	// 更新数据库中的群组信息
	if err := dao.DB.Table("user_groups").
		Where("group_basic_id=? and user_basic_id=?", groupId, deletedId).
		Delete(nil).Error; err != nil {
		return nil, err
	}
	//err = dao.DB.Model(&group).Update("Members", group.Members).Error
	err = deleteGroupFromRdb(groupId)
	return group, err
}

func AddToGroup(groupId string, invitedMembers []string) (*Model.GroupBasic, error) {
	group, err := GetGroup(groupId)
	if err != nil {
		return nil, err
	}
	for _, id := range invitedMembers {
		user, err := FindUser(id)
		if err != nil {
			return nil, err
		}
		group.Members = append(group.Members, *user)
	}
	err = dao.DB.Model(&group).Update("Members", group.Members).Error
	if err != nil {
		return nil, err
	}
	err = deleteGroupFromRdb(groupId)
	return group, err
}

func EnterGroupReq(gid string, uid string) error {
	t := time.Now()
	user, err := FindUser(uid)
	if err != nil {
		return err
	}
	owner, err := getGroupOwner(gid)
	if err != nil {
		return err
	}
	//将申请记录存入缓存
	err = dao.Rdb.SAdd(context.Background(), groupReqPreFix+gid, uid).Err()
	if err != nil {
		return err
	}

	m := &Model.Message{
		CreatedAt:  &t,
		SenderId:   uid,
		ReceiverId: owner,
		Content:    user,
		MsgType:    Model.MGroupReq,
	}
	err = SendMsg(m)
	return err
}

func EnterGroupAgree(gid string, uid string) error {
	//判断是否发过申请
	exist, err := dao.Rdb.SIsMember(context.Background(), groupReqPreFix+gid, uid).Result()
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("同意未请求申请")
	}

	owner, err := getGroupOwner(gid)
	if err != nil {
		return err
	}
	err = dao.Rdb.SRem(context.Background(), groupReqPreFix+gid, uid).Err()
	if err != nil {
		return err
	}

	_, err = AddToGroup(gid, []string{uid})
	if err != nil {
		return err
	}
	m := &Model.Message{
		SenderId:   gid,
		ReceiverId: uid,
		Content:    owner,
		MsgType:    Model.MGroupAgree,
	}
	err = SendMsg(m)
	return err
}
