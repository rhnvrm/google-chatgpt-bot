package gchat

import (
	"context"
	"log"

	oidc "github.com/coreos/go-oidc/v3/oidc"
)

const (
	jwtURL     string = "https://www.googleapis.com/service_accounts/v1/jwk/"
	chatIssuer string = "chat@system.gserviceaccount.com"
)

func VerifyJWT(audience, token string) bool {
	ctx := context.Background()
	keySet := oidc.NewRemoteKeySet(ctx, jwtURL+chatIssuer)

	configLocal := &oidc.Config{
		SkipClientIDCheck: false,
		ClientID:          audience,
	}
	newVerifier := oidc.NewVerifier(chatIssuer, keySet, configLocal)
	test, err := newVerifier.Verify(ctx, token)
	if err != nil {
		log.Println("Audience doesnt match")
		return false
	}
	if len(test.Audience) == 1 {
		log.Println("Audience matches")
	}

	return true
}
