package huawei

type Message struct {
	Hps Hps `json:"hps"`
}

//华为Push消息总结构体
type Hps struct {
	Msg Msg `json:"msg"`
	Ext Ext `json:"ext,omitempty"`
}

//PUSH消息定义。包括：消息类型type、消息内容body、消息动作action
type Msg struct {
	Type   int    `json:"type"` //1, 透传异步消息; 3, 系统通知栏异步消息。注意：2和4以后为保留后续扩展使用
	Body   Body   `json:"body"`
	Action Action `json:"action,omitempty"`
}

//消息内容。注意：对于透传类的消息可以是字符串，不必是JSON Object。
type Body struct {
	Content string `json:"content"` //消息内容体
	Title   string `json:"title"`   //消息标题
}

//消息点击动作
type Action struct {
	Type  int   `json:"type,omitempty"` //1 自定义行为：行为由参数intent定义;2 打开URL：URL地址由参数url定义;3 打开APP：默认值，打开App的首页。注意：富媒体消息开放API不支持。
	Param Param `json:"param,omitempty"`
}

//关于消息点击动作的参数
type Param struct {
	Intent     string `json:"intent,omitempty"` //Action的type为1的时候表示自定义行为。
	Url        string `json:"url,omitempty"`    //Action的type为2的时候表示打开URL地址
	AppPkgName string `json:"appPkgName"`       //需要拉起的应用包名，必须和注册推送的包名一致。
}

//扩展信息，含BI消息统计，特定展示风格，消息折叠
type Ext struct {
	BadgeAddNum string   `json:"badgeAddNum,omitempty"` //设置应用角标数值，取值范围1-99。
	BadgeClass  string   `json:"badgeClass,omitempty"`  //桌面图标对应的应用入口Activity类。
	BiTag       string   `json:"biTag,omitempty"`       //设置消息标签，如果带了这个标签，会在回执中推送给CP用于检测某种类型消息的到达率和状态。
	Customize   []string `json:"customize,omitempty"`   //用于触发onEvent点击事件，扩展样例：[{"season":"Spring"},{"weather":"raining"}] 。说明：这个字段类型必须是JSON Array，里面是key-value的一组扩展信息。
}
