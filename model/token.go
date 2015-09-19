package model

type Token struct {
	UserId     int64  `json:"userid"`
	AuthType   int64  `json:"auth_type"`
	Value      string `json:"value"`
	ExpAccess  int64  `json:"exp_access"`
	ExpRefresh int64  `json:"exp_refresh"`
}
