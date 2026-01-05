package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GetGoogleAuthURL(state string) string {
	return GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func ExchangeGoogleCode(code string) (*oauth2.Token, error) {
	return GoogleOAuthConfig.Exchange(context.Background(), code)
}

func GetGoogleUserInfo(token *oauth2.Token) (*models.GoogleUserInfo, error) {
	client := GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo models.GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func VerifyGoogleIDToken(idToken string) (*models.GoogleUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v1/tokeninfo?id_token=%s", idToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenInfo struct {
		Email    string `json:"email"`
		UserID   string `json:"user_id"`
		Audience string `json:"audience"`
	}

	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return nil, err
	}

	if tokenInfo.Audience != GoogleOAuthConfig.ClientID {
		return nil, fmt.Errorf("token audience mismatch")
	}

	userResp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", idToken))
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()

	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo models.GoogleUserInfo
	err = json.Unmarshal(userBody, &userInfo)
	if err != nil {
		return nil, err
	}

	userInfo.ID = tokenInfo.UserID
	return &userInfo, nil
}
