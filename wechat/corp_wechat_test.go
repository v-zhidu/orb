package wechat

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/v-zhidu/orb/logging"
)

func init() {
	logging.SetLevel("debug")
}

func TestCorpWechat_NewCorpWechat(t *testing.T) {
	type args struct {
		corpID     string
		corpSecret string
		agentID    int
	}
	tests := []struct {
		name string
		args args
		want *CorpWechat
	}{
		{
			name: "Succeed",
			args: args{
				corpID:     "corpid",
				corpSecret: "corpSecret",
			},
			want: &CorpWechat{
				CorpID:     "corpid",
				CorpSecret: "corpSecret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCorpWechat(tt.args.corpID, tt.args.corpSecret)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCorpWechat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCorpWechat_GetAccessToken(t *testing.T) {
	type fields struct {
		CorpID     string
		CorpSecret string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *CorpWechatAccessTokenResponse
		wantErr bool
	}{
		{
			name: "succeed",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			want: &CorpWechatAccessTokenResponse{
				CorpWechatResponse: CorpWechatResponse{
					ErrorCode: 0,
					ErrorMsg:  "ok",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &CorpWechat{
				CorpID:     tt.fields.CorpID,
				CorpSecret: tt.fields.CorpSecret,
			}

			got, err := w.GetAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("CorpWechat.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ErrorCode != tt.want.ErrorCode {
				t.Errorf("GetAccessToken() = %d, want %d", got.ErrorCode, tt.want.ErrorCode)
			}
		})
	}
}

func TestCorpWechat_CreateChat(t *testing.T) {
	randomChatID := strconv.FormatInt(time.Now().Unix(), 10)
	type fields struct {
		CorpID     string
		CorpSecret string
	}
	type args struct {
		chat *CorpWechatChatInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CorpWechatCreateChatResponse
		wantErr bool
	}{
		{
			name: "succeed",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			args: args{
				chat: &CorpWechatChatInfo{
					Name:     randomChatID,
					ChatID:   randomChatID,
					Owner:    "DuZhiQiang",
					UserList: []string{"DuZhiQiang", "Hu"},
				},
			},
			want: &CorpWechatCreateChatResponse{
				ChatID: randomChatID,
				CorpWechatResponse: CorpWechatResponse{
					ErrorCode: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "existed chat",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			args: args{
				chat: &CorpWechatChatInfo{
					Name:     randomChatID,
					ChatID:   randomChatID,
					Owner:    "DuZhiQiang",
					UserList: []string{"DuZhiQiang", "Hu"},
				},
			},
			want: &CorpWechatCreateChatResponse{
				ChatID: "",
				CorpWechatResponse: CorpWechatResponse{
					ErrorCode: 86215,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &CorpWechat{
				CorpID:     tt.fields.CorpID,
				CorpSecret: tt.fields.CorpSecret,
			}
			accessToken, _ := w.GetAccessToken()
			got, err := w.CreateChat(accessToken.AccessToken, tt.args.chat)
			if (err != nil) != tt.wantErr {
				t.Errorf("CorpWechat.CreateChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ErrorCode != tt.want.ErrorCode {
				t.Errorf("CorpWechat.CreateChat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCorpWechat_sendChatMessage(t *testing.T) {
	randomChatID := strconv.FormatInt(time.Now().Unix(), 10)
	type fields struct {
		CorpID     string
		CorpSecret string
	}
	type args struct {
		message *CorpWechatChatMessageRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CorpWechatResponse
		wantErr bool
	}{
		{
			name: "text message succeed",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			args: args{
				message: &CorpWechatChatMessageRequest{
					ChatID:      randomChatID,
					MessageType: "text",
					Text: Text{
						Content: "text message",
					},
					Safe: 0,
				},
			},
			want: &CorpWechatResponse{
				ErrorCode: 0,
				ErrorMsg:  "ok",
			},
			wantErr: false,
		},
		{
			name: "textcard message succeed",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			args: args{
				message: &CorpWechatChatMessageRequest{
					ChatID:      randomChatID,
					MessageType: "textcard",
					TextCard: TextCard{
						Title:       "领奖通知",
						Description: "<div class=\"gray\">2016年9月26日</div> <div class=\"normal\"> 恭喜你抽中iPhone 7一台，领奖码:520258</div><div class=\"highlight\">请于2016年10月10日前联系行 政同事领取</div>",
						URL:         "https://zhidao.baidu.com/question/2073647112026042748.html",
						Btntxt:      "更多",
					},
					Safe: 0,
				},
			},
			want: &CorpWechatResponse{
				ErrorCode: 0,
				ErrorMsg:  "ok",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &CorpWechat{
				CorpID:     tt.fields.CorpID,
				CorpSecret: tt.fields.CorpSecret,
			}
			//获取accesstoken
			accessToken, err := w.GetAccessToken()
			if err != nil {
				t.Errorf("CorpWechat.SendChatMessage() get access token faild. %v", err)
				return
			}
			//创建测试群聊
			_, err = w.CreateChat(accessToken.AccessToken, &CorpWechatChatInfo{
				Name:     randomChatID,
				ChatID:   randomChatID,
				Owner:    "DuZhiQiang",
				UserList: []string{"DuZhiQiang", "Hu"},
			})
			if err != nil {
				t.Errorf("CorpWechat.SendChatMessage() create chat faild. %v", err)
				return
			}
			//发送测试消息
			got, err := w.SendChatMessage(accessToken.AccessToken, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("CorpWechat.SendChatMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CorpWechat.SendChatMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCorpWechat_GetChatInfo(t *testing.T) {
	randomChatID := strconv.FormatInt(time.Now().Unix(), 10)
	type fields struct {
		CorpID     string
		CorpSecret string
	}
	type args struct {
		accessToken string
		chatid      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CorpWechatChatInfo
		wantErr bool
	}{
		{
			name: "get existed chat",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			args: args{
				chatid: randomChatID,
			},
			want: &CorpWechatChatInfo{
				Name:     randomChatID,
				ChatID:   randomChatID,
				Owner:    "DuZhiQiang",
				UserList: []string{"DuZhiQiang", "hu"},
			},
			wantErr: false,
		},
		{
			name: "get not existed chat",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			args: args{
				chatid: "123456",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &CorpWechat{
				CorpID:     tt.fields.CorpID,
				CorpSecret: tt.fields.CorpSecret,
			}
			//获取accesstoken
			accessToken, err := w.GetAccessToken()
			if err != nil {
				t.Errorf("CorpWechat.SendChatMessage() get access token faild. %v", err)
				return
			}
			//创建测试群聊
			_, err = w.CreateChat(accessToken.AccessToken, &CorpWechatChatInfo{
				Name:     randomChatID,
				ChatID:   randomChatID,
				Owner:    "DuZhiQiang",
				UserList: []string{"DuZhiQiang", "hu"},
			})
			if err != nil {
				t.Errorf("CorpWechat.SendChatMessage() create chat faild. %v", err)
				return
			}
			got, err := w.GetChatInfo(accessToken.AccessToken, tt.args.chatid)
			if (err != nil) != tt.wantErr {
				t.Errorf("CorpWechat.GetChatInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CorpWechat.GetChatInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCorpWechat_EditChat(t *testing.T) {
	randomChatID := strconv.FormatInt(time.Now().Unix(), 10)
	oldChat := &CorpWechatChatInfo{
		ChatID:   randomChatID,
		Name:     randomChatID,
		Owner:    "DuZhiQiang",
		UserList: []string{"DuZhiQiang", "hu"},
	}
	type fields struct {
		CorpID     string
		CorpSecret string
	}
	type args struct {
		chat *CorpWechatChatInfo
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *CorpWechatChatInfo
		wantCode int
		wantErr  bool
	}{
		{
			name: "get existed chat",
			fields: fields{
				CorpID:     "wwfa024bce44f41ca1",
				CorpSecret: "MN0PYHNGMAp_705LBsZ3xyjsQdLXR_hvEIqrFOAifCY",
			},
			args: args{
				chat: &CorpWechatChatInfo{
					ChatID:   randomChatID,
					Owner:    "hu",
					UserList: []string{"DuZhiQiang", "hu"},
				},
			},
			want: &CorpWechatChatInfo{
				ChatID:   randomChatID,
				Owner:    "hu",
				UserList: []string{"DuZhiQiang", "hu"},
			},
			wantCode: 0,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &CorpWechat{
				CorpID:     tt.fields.CorpID,
				CorpSecret: tt.fields.CorpSecret,
			}
			//获取accesstoken
			accessToken, err := w.GetAccessToken()
			if err != nil {
				t.Errorf("CorpWechat.SendChatMessage() get access token faild. %v", err)
				return
			}
			//创建测试群聊
			_, err = w.CreateChat(accessToken.AccessToken, oldChat)
			result, _ := w.EditChat(accessToken.AccessToken, tt.args.chat)
			if result.ErrorCode != tt.wantCode {
				t.Errorf("CorpWechat.EditChat() code = %v, wantCode %v", result.ErrorCode, tt.wantCode)
				return
			}
			w.SendChatMessage(accessToken.AccessToken, &CorpWechatChatMessageRequest{
				ChatID:      randomChatID,
				MessageType: "text",
				Text: Text{
					Content: "text message",
				},
				Safe: 0,
			})
		})
	}
}
