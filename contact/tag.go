package contact

import (
	"fmt"

	wechatwork "github.com/dfang/wechat-work-go"
	"github.com/pkg/errors"
)

type TagManagement struct {
	App *wechatwork.App
	// tags []*Tag
}

// @todo 关联 contact
func (contact *Contact) NewTagManagement() *TagManagement {
	contact.TagManageMent = &TagManagement{
		App: contact.App,
	}
	return contact.TagManageMent
}

type Tag struct {
	Tagname string `json:"tagname"`
	TagID   int    `json:"tagid"`
}

func NewTag(tagname string, tagid int) *Tag {
	return &Tag{}
}

var apiPath = "https://qyapi.weixin.qq.com/cgi-bin/tag"

const (
	createURL = "%s/create?access_token=%s"
	updateURL = "%s/update?access_token=%s"
	deleteURL = "%s/delete?access_token=%s&tagid=%d"
	getURL    = "%s/get?access_token=%s&tagid=%d"
	listURL   = "%s/list?access_token=%s"

	addTagUsersURL = "%s/addtagusers?access_token=%s"
	delTagUsersURL = "%s/deltagusers?access_token=%s"
)

// Create 创建标签
//
// 请求方式：POST（HTTPS）
// 请求地址：https://qyapi.weixin.qq.com/cgi-bin/tag/create?access_token=ACCESS_TOKEN
type RespCreateTag struct {
	RespCommon
	TagID int `json:"tagid"`
}

func (manager *TagManagement) Create(body Tag) (RespCreateTag, error) {
	accessToken := manager.App.GetAccessToken()
	uri := fmt.Sprintf(createURL, apiPath, accessToken)
	var result RespCreateTag

	err := manager.App.SimplePost(uri, body, &result)
	if err != nil {
		return RespCreateTag{}, err
	}

	// manager.tags = append(manager.tags, &body)

	return result, nil
}

// @todo: 注意，标签总数不能超过3000个。
// func (manager *TagManagement) CheckTagLimit(tag Tag) (*, error) {
// 	len(manager.tags)
// }

// Update 更新标签名字
// 请求方式：POST（HTTPS）
// 请求地址：https://qyapi.weixin.qq.com/cgi-bin/tag/update?access_token=ACCESS_TOKEN
type RespUpdateTag struct {
	RespCommon
}

func (manager *TagManagement) Update(body Tag) (RespUpdateTag, error) {
	accessToken := manager.App.GetAccessToken()
	uri := fmt.Sprintf(updateURL, apiPath, accessToken)
	var result RespUpdateTag

	err := manager.App.SimplePost(uri, body, &result)
	if err != nil {
		return RespUpdateTag{}, err
	}

	// 更新 本地 manager
	// for _, tag := range manager.tags {
	// 	if tag.TagID == body.TagID {
	// 		tag.Tagname = body.Tagname
	// 	}
	// }

	return result, nil
}

// Delete 删除标签
// 请求方式：GET（HTTPS）
// 请求地址：https://qyapi.weixin.qq.com/cgi-bin/tag/delete?access_token=ACCESS_TOKEN&tagid=TAGID
func (manager *TagManagement) Delete(tagid int) (RespCommon, error) {
	accessToken := manager.App.GetAccessToken()
	uri := fmt.Sprintf(deleteURL, apiPath, accessToken, tagid)

	var result RespCommon

	err := manager.App.SimpleGet(uri, &result)
	if err != nil {
		return RespCommon{}, err
	}

	// 更新 本地 manager
	// for index, tag := range manager.tags {
	// 	if tag.TagID == tagid {
	// 		manager.tags = append(manager.tags[:index], manager.tags[index+1:]...)
	// 	}
	// }
	return result, nil
}

// GetTagMembers 获取标签成员
// 请求方式：GET（HTTPS）
// 请求地址：https://qyapi.weixin.qq.com/cgi-bin/tag/get?access_token=ACCESS_TOKEN&tagid=TAGID
type RespGetTagUsers struct {
	RespCommon
	Tagname  string `json:"tagname"`
	UserList []struct {
		UserID string `json:"userid"`
		Name   string `json:"name"`
	} `json:"userlist"`
	PartyList []int `json:"partylist"`
}

func (manager *TagManagement) GetTagUsers(tagid int) (RespGetTagUsers, error) {
	accessToken := manager.App.GetAccessToken()
	uri := fmt.Sprintf(getURL, apiPath, accessToken, tagid)

	var result RespGetTagUsers

	err := manager.App.SimpleGet(uri, &result)
	if err != nil {
		return RespGetTagUsers{}, err
	}

	// 本地 manager

	return result, nil
}

// 获取标签列表
// 请求方式：GET（HTTPS）
// 请求地址：https://qyapi.weixin.qq.com/cgi-bin/tag/list?access_token=ACCESS_TOKEN
type RespGetTags struct {
	RespCommon
	TagList []Tag
}

func (manager *TagManagement) List() (RespGetTags, error) {
	accessToken := manager.App.GetAccessToken()
	uri := fmt.Sprintf(listURL, apiPath, accessToken)
	var result RespGetTags

	err := manager.App.SimpleGet(uri, &result)
	if err != nil {
		return RespGetTags{}, err
	}

	return result, nil
}

type RespTagUsers struct {
	RespCommon
	InValidList  string `json:"invalidlist"`
	InValidParty string `json:"invalidparty"`
}

type ReqAddTagUsers struct {
	TagID     int      `json:"tagid"`
	UserList  []string `json:"userlist"`
	PartyList []int    `json:"partylist"`
}

// 三种返回结果
// @todo check 注意，每个标签下部门、人员总数不能超过3万个。
// @todo 返回状态信息标准
func (manager *TagManagement) AddTagUsers(body ReqAddTagUsers) (interface{}, error) {
	accessToken := manager.App.GetAccessToken()
	uri := fmt.Sprintf(addTagUsersURL, apiPath, accessToken)
	var result RespTagUsers

	err := manager.App.SimplePost(uri, body, &result)
	// @note http 请求错误 状态码需要处理把?
	// @todo wechat 返回码以外的系统的错误返回处理 内部码?
	if err != nil {
		return RespTagUsers{}, err
	}

	// c)当包含userid、partylist全部非法时返回
	if result.ErrCode == 40070 {
		return result, errors.New(result.ErrMsg)
	}

	// 部分成功
	if result.ErrCode == 0 && result.ErrMsg == "ok" {
		// b)若部分userid、partylist非法，则返回
		if result.InValidList != "" || result.InValidParty != "" {
			return result, errors.New(result.ErrMsg)
		}
	}

	// a)正确时返回
	return result, nil
}

//
// 删除标签成员
//
// 请求方式：POST（HTTPS）
// 请求地址：https://qyapi.weixin.qq.com/cgi-bin/tag/deltagusers?access_token=ACCESS_TOKEN
type ReqDelTagUsers struct {
	TagID     int      `json:"tagid"`
	UserList  []string `json:"userlist"`
	PartyList []int    `json:"partylist"`
}

func (manager *TagManagement) DelTagUsers(body ReqDelTagUsers) (interface{}, error) {
	accessToken := manager.App.GetAccessToken()
	uri := fmt.Sprintf(delTagUsersURL, apiPath, accessToken)
	var result RespTagUsers

	err := manager.App.SimplePost(uri, body, &result)
	// @note http 请求错误 状态码需要处理把?
	// @todo wechat 返回码以外的系统的错误返回处理 内部码?
	if err != nil {
		return RespTagUsers{}, err
	}

	// c)当包含的userid、partylist全部非法时返回
	if result.ErrCode == 40031 {
		return result, errors.New(result.ErrMsg)
	}

	// 部分成功
	if result.ErrCode == 0 && result.ErrMsg == "deleted" {
		// b)若部分userid、partylist非法，则返回
		if result.InValidList != "" || result.InValidParty != "" {
			return result, errors.New(result.ErrMsg)
		}
		// a)正确时返回
		return result, nil
	}
	return result, nil
}
