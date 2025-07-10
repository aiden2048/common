package sysConfig

// OneAll登录验证配置
type OneAllConfig struct {
	Subdomain   string `json:"subdomain"`
	PublicKey   string `json:"public_key"`
	PrivateKey  string `json:"private_key"`
	APIEndpoint string `json:"api_endpoint"`
	CallBackUrl string `json:"call_back_url"`
}

type AllowMeConfig struct {
	Url       string `json:"url"`
	Apikey    string `json:"apikey"`
	ClientKey string `json:"client_key"`
}

type BigDataCorpApiConfig struct {
	Url         string `json:"url"`
	TokenId     string `json:"tokenId"`
	AccessToken string `json:"accessToken"`
}
