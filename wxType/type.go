package wxType

type BasicToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type Otoken struct {
	Access_token  string
	Expires_in    int
	Refresh_token string
	Openid        string
	Scope         string
	Errcode       int
}

type UserInfo struct {
	Openid     string
	Nickname   string
	Sex        int
	Language   string
	City       string
	Province   string
	Country    string
	Headimgurl string
}

type Ticket struct {
	Errcode    int
	Errmsg     string
	Ticket     string
	Expires_in int
}
