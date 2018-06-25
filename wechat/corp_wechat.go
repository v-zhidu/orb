package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/v-zhidu/orb/http"
	"github.com/v-zhidu/orb/logging"
	"github.com/v-zhidu/orb/set"
)

const (
	//CorpWechatAccessTokenURL 获取Access Token
	CorpWechatAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	//CorpWechatSendTextMessageURL 发送文本信息
	CorpWechatSendTextMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	//CorpWechatSendChatMessageURL 发送群聊消息
	CorpWechatSendChatMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send?access_token=%s"
	//CorpWechatCreateChatURL 创建群组
	CorpWechatCreateChatURL = "https://qyapi.weixin.qq.com/cgi-bin/appchat/create?access_token=%s"
	//CorpWechatEditChatURL 修改群组
	CorpWechatEditChatURL = "https://qyapi.weixin.qq.com/cgi-bin/appchat/update?access_token=%s"
	//CorpWehcatChatInfoURL 获取群组信息
	CorpWehcatChatInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/appchat/get"
)

//CorpWechat ...
type CorpWechat struct {
	CorpID     string
	CorpSecret string
}

//NewCorpWechat ...
func NewCorpWechat(corpID string, corpSecret string) *CorpWechat {
	return &CorpWechat{
		CorpID:     corpID,
		CorpSecret: corpSecret,
	}
}

//GetAccessToken 获取企业微信Access Token
func (w *CorpWechat) GetAccessToken() (*CorpWechatAccessTokenResponse, error) {
	data, err := http.Get(CorpWechatAccessTokenURL, map[string]string{
		"corpid":     w.CorpID,
		"corpsecret": w.CorpSecret,
	})

	if err != nil {
		logging.WithError("get wechat access token failed", err)
		return nil, err
	}

	var response CorpWechatAccessTokenResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		logging.WithError("get wechat access token failed", err)
		return nil, err
	}

	return &response, nil
}

//CreateChat returns the result that create a wechat group.
//errcode=0 创建成功
//errcode=86215 chatid或chatname已经存在
func (w *CorpWechat) CreateChat(accessToken string,
	chat *CorpWechatChatInfo) (*CorpWechatCreateChatResponse, error) {
	logging.Info("create chat request", logging.Fields{
		"chatName": chat.Name,
		"chatId":   chat.ChatID,
		"owner":    chat.Owner,
		"userList": chat.UserList,
	})

	body, err := json.Marshal(chat)
	if err != nil {
		logging.WithError("json marshal error", err)
		return nil, err
	}
	data, err := http.PostJSON(fmt.Sprintf(CorpWechatCreateChatURL, accessToken), body, nil)
	if err != nil {
		logging.WithError("create chat failed", err)
		return nil, err
	}

	var response CorpWechatCreateChatResponse
	if err := json.Unmarshal(data, &response); err != nil {
		logging.WithError("json unmarshal error", err)
		return nil, err
	}
	logging.Info("create chat response", logging.Fields{
		"errcode": response.ErrorCode,
		"errmsg":  response.ErrorMsg,
		"chatid":  response.ChatID,
	})

	return &response, nil
}

//EditChat 修改群聊会话信息
func (w *CorpWechat) EditChat(accessToken string, chat *CorpWechatChatInfo) (*CorpWechatResponse, error) {
	logging.Info("edit wechat request", logging.Fields{
		"chatName": chat.Name,
		"chatId":   chat.ChatID,
		"owner":    chat.Owner,
		"userList": chat.UserList,
	})
	//获取会话消息
	chatInfo, err := w.GetChatInfo(accessToken, chat.ChatID)
	if err != nil {
		return nil, err
	}
	if chatInfo == nil {
		return &CorpWechatResponse{
			ErrorCode: 40050,
			ErrorMsg:  "chat not existed",
		}, nil
	}
	//检查群组信息变更
	oldUserSet := set.New(chatInfo.UserList...)
	newUserSet := set.New(chat.UserList...)
	delUserList := oldUserSet.Minus([]*set.StringSet{newUserSet}...).StringSlice()
	addUserList := newUserSet.Minus([]*set.StringSet{oldUserSet}...).StringSlice()

	body := map[string]interface{}{
		"chatid":        chat.ChatID,
		"chat":          chat.Name,
		"owner":         chat.Owner,
		"del_user_list": delUserList,
		"add_user_list": addUserList,
	}

	data, err := http.PostMap(fmt.Sprintf(CorpWechatEditChatURL, accessToken), body, nil)
	if err != nil {
		logging.WithError("edit chat information failed", err)
		return nil, err
	}

	var response CorpWechatResponse
	if err := json.Unmarshal(data, &response); err != nil {
		logging.WithError("json unmarshal error", err)
		return nil, err
	}
	logging.Info("edit chat information response", logging.Fields{
		"chatid":  chat.ChatID,
		"errcode": response.ErrorCode,
		"errmsg":  response.ErrorMsg,
	})

	return &response, nil
}

//GetChatInfo 获取群聊会话信息
func (w *CorpWechat) GetChatInfo(accessToken string, chatid string) (*CorpWechatChatInfo, error) {
	data, err := http.Get(CorpWehcatChatInfoURL, map[string]string{
		"access_token": accessToken,
		"chatid":       chatid,
	})
	if err != nil {
		logging.WithError("get chat information failed", err)
		return nil, err
	}

	var response CorpWechatChatInfoResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		logging.WithError("json unmarshal error", err)
		return nil, err
	}
	if response.CorpWechatResponse.ErrorCode == 0 {
		return &response.CorpWechatChatInfo, nil
	}

	return nil, nil
}

//SendChatMessage 发送群聊消息
func (w *CorpWechat) SendChatMessage(accessToken string,
	message *CorpWechatChatMessageRequest) (*CorpWechatResponse, error) {
	logging.Info("send chat message", logging.Fields{
		"chatid":      message.ChatID,
		"messageType": message.MessageType,
	})

	body, err := json.Marshal(message)
	if err != nil {
		logging.WithError("json marshal error", err)
		return nil, err
	}
	data, err := http.PostJSON(fmt.Sprintf(CorpWechatSendChatMessageURL, accessToken), body, nil)
	if err != nil {
		logging.WithError("send chat message failed", err)
		return nil, err
	}

	var response CorpWechatResponse
	if err := json.Unmarshal(data, &response); err != nil {
		logging.WithError("json unmarshal error", err)
		return nil, err
	}
	logging.Info("send chat message response", logging.Fields{
		"chatid":      message.ChatID,
		"messageType": message.MessageType,
		"response":    response.ErrorCode,
	})

	return &response, nil
}

// //SendTextMessage ...
// func (w *CorpWechat) SendTextMessage(message *CorpWechatMessageRequest) (*CorpWechatMessageResponse, error) {
// 	message.MessagegType = "text"
// 	return w.sendMessage(message)
// }

// //SendTextCardMessage ...
// func (w *CorpWechat) SendTextCardMessage(message *CorpWechatMessageRequest) (*CorpWechatMessageResponse, error) {
// 	message.MessagegType = "textcard"
// 	return w.sendMessage(message)
// }

// //sendMessage
// func (w *CorpWechat) sendMessage(message *CorpWechatMessageRequest) (*CorpWechatMessageResponse, error) {
// 	accessToken, err := w.GetAccessToken()
// 	if err != nil {
// 		return nil, errors.New("get wehcat access token failed")
// 	}

// 	data, err := HttpPostStruct(fmt.Sprintf(CorpWechatSendTextMessageURL, accessToken.AccessToken), message, nil)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"error": err,
// 		}).Error("send text message failed")
// 		return nil, errors.New("send text message failed")
// 	}

// 	var response CorpWechatMessageResponse
// 	err = json.Unmarshal(*data, &response)
// 	if err != nil || response.ErrorCode != 0 {
// 		log.WithFields(log.Fields{
// 			"errcode": response.ErrorCode,
// 			"errmsg":  response.ErrorMsg,
// 		}).Error("send text message failed")
// 	}

// 	return &response, nil
// }

//CorpWechatResponse ...
type CorpWechatResponse struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

//CorpWechatAccessTokenResponse ...
type CorpWechatAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	CorpWechatResponse
}

//CorpWechatChatInfo ...
type CorpWechatChatInfo struct {
	Name     string   `json:"name"`
	Owner    string   `json:"owner"`
	UserList []string `json:"userlist"`
	ChatID   string   `json:"chatid"`
}

//CorpWechatChatInfoResponse ...
type CorpWechatChatInfoResponse struct {
	CorpWechatChatInfo `json:"chat_info"`
	CorpWechatResponse
}

//CorpWechatCreateChatResponse ...
type CorpWechatCreateChatResponse struct {
	ChatID string `json:"chatid"`
	CorpWechatResponse
}

//CorpWechatChatMessageRequest ...
type CorpWechatChatMessageRequest struct {
	ChatID      string   `json:"chatid"`
	MessageType string   `json:"msgtype"`
	Safe        int      `json:"safe"`
	Text        Text     `json:"text"`
	TextCard    TextCard `json:"textcard"`
}

//CorpWechatMessageResponse ...
type CorpWechatMessageResponse struct {
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
	CorpWechatResponse
}

//CorpWechatMessageRequest ...
type CorpWechatMessageRequest struct {
	ToUser       string   `json:"touser"`
	ToParty      string   `json:"toparty"`
	ToTag        string   `json:"totag"`
	MessagegType string   `json:"msgtype"`
	AgentID      int      `json:"agentid"`
	Safe         int      `json:"safe"`
	Text         Text     `json:"text"`
	TextCard     TextCard `json:"textcard"`
}

//Text Message
type Text struct {
	Content string `json:"content"`
}

//TextCard Message
type TextCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}
