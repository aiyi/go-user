package token

type Token struct {
	SlotId     int    `json:"slot_id"`
	Value      string `json:"value"`
	UserId     int64  `json:"user_id"`
	BindType   int64  `json:"bind_type"`
	ExpAccess  int64  `json:"exp_access"`
	ExpRefresh int64  `json:"exp_refresh"`
}

type TokenStorage struct {
	TokenList []Token
}
