package huawei

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	NSP_STATUS_6   = 6
	NSP_STATUS_102 = 102
	NSP_STATUS_105 = 105
	NSP_STATUS_111 = 111
	NSP_STATUS_112 = 112
	NSP_STATUS_113 = 113
	NSP_STATUS_114 = 114
	NSP_STATUS_199 = 199
	NSP_STATUS_403 = 403
)

var NSP_STATUS_MSG = map[int]string{
	NSP_STATUS_6:   "session过期",
	NSP_STATUS_102: "无效的SESSION_KEY",
	NSP_STATUS_105: "参数错误",
	NSP_STATUS_111: "系统、服务处理忙",
	NSP_STATUS_112: "找不到对应服务",
	NSP_STATUS_113: "请求服务失败",
	NSP_STATUS_114: "服务不可达、无路由",
	NSP_STATUS_199: "未知错误",
	NSP_STATUS_403: "无权限",
}

//对方业务异常
type BizErr struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}

func NewBizErr(errno int, errmsg string) *BizErr {
	return &BizErr{
		Errno:  errno,
		Errmsg: errmsg,
	}
}

//post请求
func SendFormPost(url string, data url.Values) ([]byte, error, *BizErr) {
	body := ioutil.NopCloser(strings.NewReader(data.Encode()))
	response, err := http.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {

		return []byte(""), err, nil
	}
	defer response.Body.Close()

	//HTTP协议错误码处理
	if response.StatusCode == 500 {
		return []byte(""), nil, NewBizErr(response.StatusCode, "推送服务器系统错误")
	}

	if response.StatusCode == 503 {
		return []byte(""), nil, NewBizErr(response.StatusCode, "流量控制错误，发送流量过高")
	}

	//系统级错误码(通过扩展HTTP协议头 NSP_STATUS)
	NSP_STATUS := response.Header.Get("NSP_STATUS")
	NSP_STATUS_INT, _ := strconv.Atoi(NSP_STATUS)
	if NSP_STATUS_INT != 0 {
		return []byte(""), nil, NewBizErr(NSP_STATUS_INT, NSP_STATUS_MSG[NSP_STATUS_INT])
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err, nil
	}
	return result, err, nil
}
