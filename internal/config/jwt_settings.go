package config

import "fmt"

// JwtSettings contains the settings for the JWT token.
type JwtSettings struct {
	SecretKey Secret   `json:"secret"`
	Expire    Duration `json:"expire"`
}

func (js JwtSettings) String() string {
	return fmt.Sprintf("{SecretKey: %s, Expire: %v}", js.SecretKey, js.Expire)
}
