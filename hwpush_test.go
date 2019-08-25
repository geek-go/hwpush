package huawei

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
)

var Cfg = &HwpushConfig{
	AppId:     "",
	AppSecret: "",
	Package:   "",
}

//测试单推
func TestHwPush_SendByCid(t *testing.T) {
	hwpush, err := NewHuawei(Cfg)
	if err != nil {
		t.Error(err)
		os.Exit(1)
	}

	cid := "xxx"
	payLoad := Payload{"这是测试title", "这是测试内容", "1", ""}
	err = hwpush.SendByCid(cid, &payLoad)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("ok")
	}
}

//测试群推
func TestHwPush_SendByCids(t *testing.T) {
	hwpush, err := NewHuawei(Cfg)
	if err != nil {
		t.Error(err)
		os.Exit(1)
	}

	cids := []string{"xxx"}
	payLoad := Payload{"这是测试title", "这是测试内容", "1", ""}
	err = hwpush.SendByCids(cids, &payLoad)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("ok")
	}
}

type HwpushConfig struct {
	AppId     string `toml:"app_id"`
	AppSecret string `toml:"app_secret"`
	Package   string `toml:"package"`
}

//消息payload，根据业务自定义
type Payload struct {
	PushTitle    string `json:"push_title"`
	PushBody     string `json:"push_body"`
	IsShowNotify string `json:"is_show_notify"`
	Ext          string `json:"ext"`
}

type HuaweiPush struct {
	Config *HwpushConfig
}

//存储每个cid失败的原因
type PushSendResultMsg struct {
	Success       int64    `json:"success"`
	Failure       int64    `json:"failure"`
	IllegalTokens []string `json:"illegal_tokens"`
}

//获取实例
func NewHuawei(config *HwpushConfig) (*HuaweiPush, error) {

	if config.AppId == "" || config.Package == "" || config.AppSecret == "" {
		return nil, errors.New("请检查配置")
	}

	hwpush := &HuaweiPush{
		Config: config,
	}

	return hwpush, nil
}

//根据用户cid推送
func (h *HuaweiPush) SendByCid(cid string, payload *Payload) error {
	cids := []string{cid}
	return h.SendByCids(cids, payload)
}

//组装消息体
func NewMessage(appPkgName string, payload *Payload) *Message {

	msgType := 1                     //默认透传
	if payload.IsShowNotify == "1" { //通知栏
		msgType = 3
	}

	payload_str, _ := json.Marshal(payload)

	return &Message{
		Hps: Hps{
			Msg: Msg{
				Type: msgType, //1, 透传异步消息; 3, 系统通知栏异步消息。注意：2和4以后为保留后续扩展使用
				Body: Body{
					Content: payload.PushBody,
					Title:   payload.PushTitle,
				},
				Action: Action{
					Type: 3, //1, 自定义行为; 2, 打开URL; 3, 打开App;
					Param: Param{
						AppPkgName: appPkgName,
					},
				},
			},
			Ext: Ext{
				Customize: []string{string(payload_str)},
			},
		},
	}
}

//根据用户cids批量推送
func (h *HuaweiPush) SendByCids(cids []string, payload *Payload) error {

	//nspCtx
	vers := &Vers{
		Ver:   "1",
		AppID: h.Config.AppId,
	}
	nspCtx, _ := json.Marshal(vers)

	//hps
	hps, err := json.Marshal(NewMessage(h.Config.Package, payload))
	if err != nil {
		return err
	}

	deviceToken, err := json.Marshal(cids)
	if err != nil {
		return err
	}

	pushSendParam := &PushSendParam{
		DeviceToken: string(deviceToken),
		NspCtx:      string(nspCtx),
		Payload:     string(hps),
	}

	pushSendResult, err, bizErr := PushSend(h.Config.AppId, h.Config.AppSecret, pushSendParam)
	if err != nil {
		return err
	}

	fmt.Println(pushSendResult, bizErr)

	return nil
}
