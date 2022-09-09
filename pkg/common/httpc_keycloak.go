package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	KeycloakServer     = "KEYCLOAK_SERVER_URL"
	KeycloakRealm      = "KEYCLOAK_REALM"
	KeycloakResource   = "KEYCLOAK_CLIENT_ID"
	KeycloakCredential = "KEYCLOAK_CLIENT_CREDENTIAL"
)

var token *AccessToken
var tokenExpireAt time.Time

type AccessToken struct {
	Token            string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

func KeycloakAuthenticator(request *http.Request) error {
	wrapToken := ""
	if token != nil && !tokenExpireAt.IsZero() && time.Now().Before(tokenExpireAt) {
		wrapToken = token.Token
	} else {
		accessToken, err := auth()
		if err != nil {
			return err
		}
		if accessToken != nil {
			token = accessToken
			wrapToken = token.Token
		}
	}
	var authHeaders map[string]string
	if !IsEmptyString(wrapToken) {
		authHeaders = map[string]string{
			"Authorization": "Bearer " + wrapToken,
		}
	}
	if len(authHeaders) <= 0 {
		return fmt.Errorf("no authorized token found! Please use login command first to get login token first")
	} else {
		for key, val := range authHeaders {
			request.Header.Add(key, val)
		}
	}
	return nil
}

func auth() (*AccessToken, error) {
	kcServer := os.Getenv(KeycloakServer)
	kcRealm := os.Getenv(KeycloakRealm)
	kcRes := os.Getenv(KeycloakResource)
	kcSecret := os.Getenv(KeycloakCredential)
	if IsEmptyString(kcServer) || IsEmptyString(kcRealm) || IsEmptyString(kcRes) || IsEmptyString(kcSecret) {
		fmt.Printf("Missing needed Keycloak configurations, no authentication needed.")
		return nil, nil
	}
	kcAuthURL := fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", kcServer, kcRealm)
	client := &http.Client{}
	values := url.Values{
		"grant_type": {"client_credentials"},
	}
	req, err := http.NewRequest("POST", kcAuthURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(kcRes, kcSecret)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("login failed: %s", resp.Status)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var accessToken = &AccessToken{}
	err = json.Unmarshal(content, accessToken)
	if err != nil {
		panic(err)
	}

	return accessToken, nil
}
