package huawei

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

type Vers struct {
	Ver   string `json:"ver"`
	AppID string `json:"appId"`
}

//请求参数
type PushSendParam struct {
	DeviceToken string `json:"device_token"` //JSON数值字符串，单次最多只是100个。
	Payload     string `json:"payload"`      //描述投递消息的JSON结构体，描述PUSH消息的：类型、内容、显示、点击动作、报表统计和扩展信息
	NspCtx      string `json:"nsp_ctx"`
}

type PushSendResult struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"requestId"`
	Ext       string `json:"ext"`
}

//发送消息
//https://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush_agent.html?page=hmssdk_huaweipush_api_reference_agent_s2#4%20%E8%BF%94%E5%9B%9E%E5%80%BC
func PushSend(client_id string, client_secret string, param *PushSendParam) (*PushSendResult, error, *BizErr) {

	accessToken := GetToken(client_id, client_secret)
	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(param.NspCtx)

	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           "openpush.message.api.send",
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())), //服务请求时间戳，自GMT 时间 1970-1-1 0:0:0至今的秒数。如果传入的时间与服务器时间相差5分钟以上，服务器可能会拒绝请求
		"device_token_list": param.DeviceToken,
		"payload":           param.Payload,
		"expire_time":       time.Now().Format("2006-01-02T15:04"),
	}

	form_param := make(url.Values)
	form_param.Set("access_token", originParam["access_token"])
	form_param.Set("nsp_svc", originParam["nsp_svc"])
	form_param.Set("nsp_ts", originParam["nsp_ts"])
	form_param.Set("device_token_list", originParam["device_token_list"])
	form_param.Set("payload", originParam["payload"])

	// push
	res, err, bizErr := SendFormPost(reqUrl, form_param)
	if err != nil {
		return nil, err, nil
	}

	var pushSendResult = &PushSendResult{}
	err = json.Unmarshal(res, pushSendResult)
	if err != nil {
		return nil, err, nil
	}

	return pushSendResult, nil, bizErr
}
