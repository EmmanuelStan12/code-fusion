package configs

type JwtConfig struct {
	SecretKey []byte
	Exp       int
	Issuer    string
	Audience  string
}
