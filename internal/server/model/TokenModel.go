package model

import (
	"bytes"
	"fmt"
	"time"
)

const (
	AccessTokenName  = "access_token"
	RefreshTokenName = "refresh_token"
)

type TokenModel struct {
	Name    string
	Value   string
	ExpTime time.Time
}

// MarshalTokenToJsonStr  crates string type " {'name':'value'} "
func MarshalTokenToJsonStr(tokens ...TokenModel) string {
	var s bytes.Buffer
	fmt.Fprintf(&s, "{")

	for i := 0; i < len(tokens)-1; i++ {
		fmt.Fprintf(&s, "'%s':'%s'", tokens[i].Name, tokens[i].Value)

		if i < len(tokens)-2 {
			fmt.Fprintf(&s, ",")
		} else {
			fmt.Fprintf(&s, "}")
		}

	}
	return s.String()
}
