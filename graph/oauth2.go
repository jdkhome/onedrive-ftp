package graph

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	authCodeURL     = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	authTokenURL    = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	authRedirectURL = "https://login.live.com/oauth20_desktop.srf"
	authClientID    = "3470c3fa-bc10-45ab-a0a9-2d30836485d1"
	authFile        = "auth_tokens.json"
)

// Auth represents a set of oauth2 authentication tokens
type Auth struct {
	ExpiresIn    int64  `json:"expires_in"` // only used for parsing
	ExpiresAt    int64  `json:"expires_at"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ToFile writes auth tokens to a file
func (a Auth) ToFile(file string) error {
	byteData, _ := json.Marshal(a)
	return ioutil.WriteFile(file, byteData, 0600)
}

// FromFile populates an auth struct from a file
func (a *Auth) FromFile(file string) error {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, a)
}

// Refresh auth tokens if expired.
func (a *Auth) Refresh() {
	if a.ExpiresAt <= time.Now().Unix() {
		log.Info("Auth tokens expired, attempting renewal.")
		oldTime := a.ExpiresAt

		postData := strings.NewReader("client_id=" + authClientID +
			"&redirect_uri=" + authRedirectURL +
			"&refresh_token=" + a.RefreshToken +
			"&grant_type=refresh_token")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		req, err := http.NewRequest("POST", authTokenURL, postData)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err :=client.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("Could not POST to renew tokens, exiting.")
		}

		//resp, err := http.Post(authTokenURL,
		//	"application/x-www-form-urlencoded",
		//	postData)
		//if err != nil {
		//	log.WithFields(log.Fields{
		//		"err": err,
		//	}).Fatal("Could not POST to renew tokens, exiting.")
		//}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &a)
		if a.ExpiresAt == oldTime {
			a.ExpiresAt = time.Now().Unix() + a.ExpiresIn
		}
		if a.AccessToken == "" || a.RefreshToken == "" {
			os.Remove(authFile)
			log.Fatalf("Failed to renew access tokens. Response from server:\n%s\n", string(body))
		}
		a.ToFile(authFile)
	}
}

// Get the appropriate authentication URL for the Graph OAuth2 challenge.
func getAuthURL() string {
	return authCodeURL +
		"?client_id=" + authClientID +
		"&scope=" + url.PathEscape("files.readwrite.all offline_access") +
		"&response_type=code" +
		"&redirect_uri=" + authRedirectURL
}

// Exchange an auth code for a set of access tokens
func getAuthTokens(authCode string) Auth {
	postData := strings.NewReader(
		"client_id=" + authClientID +
			"&redirect_uri=" + authRedirectURL +
			"&code=" + authCode +
			"&grant_type=authorization_code")


	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", authTokenURL, postData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err :=client.Do(req)
	//resp, err := http.Post(authTokenURL,
	//	"application/x-www-form-urlencoded",
	//	postData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var auth Auth
	json.Unmarshal(body, &auth)
	if auth.ExpiresAt == 0 {
		auth.ExpiresAt = time.Now().Unix() + auth.ExpiresIn
	}
	if auth.AccessToken == "" || auth.RefreshToken == "" {
		log.Fatalf("Failed to retrieve access tokens. Response from server:\n%s\n", string(body))
	}
	return auth
}

// Authenticate performs first-time authentication to Graph
func Authenticate() *Auth {
	var auth Auth
	_, err := os.Stat(authFile)
	if os.IsNotExist(err) {
		// no tokens found, gotta start oauth flow from beginning
		code := getAuthCode()
		auth = getAuthTokens(code)
		auth.ToFile(authFile)
	} else {
		// we already have tokens, no need to force a refresh
		auth.FromFile(authFile)
		auth.Refresh()
	}
	return &auth
}
