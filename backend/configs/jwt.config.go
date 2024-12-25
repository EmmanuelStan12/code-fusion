package configs

type JwtConfig struct {
	SecretKey  string `json:"secretKey"`
	ExpInHours int    `json:"expInHours"`
	Issuer     string `json:"issuer"`
	Audience   string `json:"audience"`
}
