package huawei

import (
	"encoding/json"
	"net/url"
)

type TokenResultStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

//获取auth_token
func GetToken(client_id string, client_secret string) string {
	reqUrl := TOKEN_URL
	param := make(url.Values)
	param["grant_type"] = []string{"client_credentials"}
	param["client_id"] = []string{client_id}
	param["client_secret"] = []string{client_secret}
	res, err, _ := SendFormPost(reqUrl, param)

	if nil != err {
		return ""
	}
	var tokenResult = &TokenResultStruct{}
	err = json.Unmarshal(res, tokenResult)
	if err != nil {
		return ""
	}
	return tokenResult.AccessToken
}
