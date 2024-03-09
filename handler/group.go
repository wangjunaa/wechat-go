package handler

import (
	"demo/dao"
	Model "demo/models"
	"encoding/json"
	"errors"
	"time"
)

func CheckOwner(ownerId string, groupId string) error {
	var group Model.GroupBasic
	err := dao.DB.First(&group, "id=?", groupId).Error
	if err != nil {
		return err
	}
	if group.OwnerId != ownerId {
		return errors.New("非群主操作")
	}
	return nil

}
func CreateGroup(ownerId string, membersId []string) (Model.GroupBasic, error) {
	sGroupMux.Lock()
	defer sGroupMux.UnLock()
	name := ""
	group := Model.GroupBasic{
		Name:    "",
		OwnerId: ownerId,
		Icon:    nil,
		Type:    0,
	}
	for i, id := range membersId {
		user, err := GetUser(id)
		if err != nil {
			return group, err
		}
		//初始化群名
		if i < 3 {
			name += user.UserName
			if i < 2 {
				name += "、"
			}
		}
		group.Members = append(group.Members, *user)
	}
	err := dao.DB.Create(&group).Error
	return group, err
}
func DeleteGroup(groupId string) error {
	sGroupMux.Lock()
	defer sGroupMux.UnLock()
	err := dao.DB.Table("user_groups").Where("group_basic_id=?", groupId).Delete(nil).Error
	if err != nil {
		return err
	}
	err = dao.DB.Model(&Model.GroupBasic{}).Delete("id=?", groupId).Error
	if err != nil {
		return err
	}
	return nil
}
func RemoveGroupMember(groupId string, deletedId string) (*Model.GroupBasic, error) {
	//log.Println(groupId, deletedId)

	group, err := GetGroupById(groupId)
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

	sGroupMux.Lock()
	defer sGroupMux.UnLock()
	// 更新数据库中的群组信息
	if err := dao.DB.Table("user_groups").
		Where("group_basic_id=? and user_basic_id=?", groupId, deletedId).
		Delete(nil).Error; err != nil {
		return nil, err
	}
	return group, nil
}

//	func AddToGroup(gid string, uid string) error {
//		u := Model.UserBasic{ID: uid}
//		g := Model.GroupBasic{}
//		err := dao.DB.First(&g, "id=?", gid).Error
//		if err != nil {
//			return err
//		}
//
//		g.Members = append(g.Members, u)
//		sGroupMux.Lock()
//		defer sGroupMux.UnLock()
//		err = dao.DB.Model(g).Update("Members", g.Members).Error
//		return err
//	}
func AddToGroup(groupId string, invitedMembers []string) (*Model.GroupBasic, error) {
	sGroupMux.Lock()
	defer sGroupMux.UnLock()
	group, err := GetGroupById(groupId)
	if err != nil {
		return nil, err
	}
	for _, id := range invitedMembers {
		user, err := GetUser(id)
		if err != nil {
			return nil, err
		}
		group.Members = append(group.Members, *user)
	}
	err = dao.DB.Model(&group).Update("Members", group.Members).Error
	return group, err
}
func GetGroupByName(name string) (*Model.GroupBasic, error) {
	g := &Model.GroupBasic{}
	err := dao.DB.Preload("Members").First(&g, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return g, err
}

func GetGroupById(gid string) (*Model.GroupBasic, error) {
	g := &Model.GroupBasic{}
	err := dao.DB.Preload("Members").First(&g, gid).Error
	if err != nil {
		return nil, err
	}
	return g, err
}
func getGroupOwner(gid string) (string, error) {
	g := Model.GroupBasic{}
	err := dao.DB.Select("id", "owner_id").First(&g, "id", gid).Error
	if err != nil {
		return "", err
	}
	return g.OwnerId, nil
}
func EnterGroupReq(gid string, uid string) error {
	t := time.Now()
	user, err := GetUser(uid)
	if err != nil {
		return err
	}
	//marshal, err := json.Marshal(user)
	//if err != nil {
	//	return err
	//}
	owner, err := getGroupOwner(gid)
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
	var uList []string
	uList = append(uList, uid)
	_, err := AddToGroup(gid, uList)
	if err != nil {
		return err
	}
	user, err := GetUser(uid)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}
	m := &Model.Message{
		SenderId:   gid,
		ReceiverId: uid,
		Content:    marshal,
		MsgType:    Model.MGroupAgree,
	}
	err = SendMsg(m)
	return err
}
