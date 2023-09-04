package thirdparty

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/teamhanko/hanko/backend/config"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	ZaloAPIBase          = "https://oauth.zaloapp.com"
	ZaloAuthEndpoint     = ZaloAPIBase + "/v4/permission"
	ZaloTokenEndpoint    = ZaloAPIBase + "/v4/access_token"
	ZaloUserInfoEndpoint = "https://graph.zalo.me/v2.0/me"
)

var DefaultZaloScopes = []string{
	"name",
	"picture",
	"id",
}

type zaloProvider struct {
	*oauth2.Config
}

func NewZaloProvider(config config.ThirdPartyProvider, redirectURL string) (OAuthProvider, error) {
	if !config.Enabled {
		return nil, errors.New("zalo provider requested but disabled")
	}

	return &zaloProvider{
		Config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.Secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  ZaloAuthEndpoint,
				TokenURL: ZaloTokenEndpoint,
			},
			RedirectURL: redirectURL,
			Scopes:      DefaultZaloScopes,
		},
	}, nil
}

/*
https://developers.zalo.me/docs/social-api/tham-khao/user-access-token-v4
Web: https://oauth.zaloapp.com/v4/permission?app_id=<APP_ID>&redirect_uri=<CALLBACK_URL>&state=<STATE>
*/
func (a zaloProvider) AuthCodeURL(state string, args ...oauth2.AuthCodeOption) string {
	clientId := a.Config.ClientID
	//clientSecret := a.Config.ClientSecret
	redirectUrl := a.Config.RedirectURL

	base, _ := url.Parse(ZaloAuthEndpoint)

	// Query params
	params := url.Values{}
	params.Add("state", state)
	params.Add("app_id", clientId)
	params.Add("redirect_uri", redirectUrl)
	base.RawQuery = params.Encode()

	return base.String()
}

func (a zaloProvider) GetOAuthToken(code string) (*oauth2.Token, error) {
	return a.exchangeCodeForToken(code)
}

func (a zaloProvider) GetUserData(token *oauth2.Token) (*UserData, error) {

	// Create an HTTP client.
	client := &http.Client{}

	// Create a GET request.
	req, err := http.NewRequest("GET", ZaloUserInfoEndpoint+"?fields=id,name,picture,gender,birthday,phone,number", nil)
	if err != nil {
		return nil, err
	}

	// Set the access token as a header.
	req.Header.Set("access_token", token.AccessToken)

	// Perform the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code and handle it as needed.
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	// Read the response.
	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	responseBody := buffer.Bytes()

	// Create an instance of AccessTokenResponse
	var response zaloUserResponse

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	//parsedAccess, _, err := new(jwt.Parser).ParseUnverified(token.AccessToken, jwt.MapClaims{})
	//if err != nil {
	//	return nil, err
	//}
	//
	//claims, ok := parsedAccess.Claims.(jwt.MapClaims)
	//if !ok {
	//	return nil, errors.New("could not extract token claims")
	//}

	fakeEmail := response.ID + "@zalo.me"

	userData := &UserData{
		Emails: []Email{{
			Email:    fakeEmail,
			Verified: true,
			Primary:  true,
		}},
		Metadata: &Claims{
			Issuer:  "https://oauth.zaloapp.com",
			Subject: response.ID,
			//Aud:               claims["aud"].(string),
			//Iat:               claims["iat"].(float64),
			Name:              response.Name,
			PreferredUsername: response.Name,
			Picture:           response.Picture.Data.URL,
			Gender:            response.Gender,
			Birthdate:         response.Birthday,
			Email:             fakeEmail,
			EmailVerified:     true,
		},
	}

	return userData, nil
}

func (a zaloProvider) Name() string {
	return "zalo"
}

func (a zaloProvider) exchangeCodeForToken(code string) (*oauth2.Token, error) {
	// Define the API endpoint and your secret key.
	apiURL := ZaloTokenEndpoint

	// Prepare the form data.
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("app_id", a.Config.ClientID)
	formData.Set("grant_type", "authorization_code")
	//formData.Set("code_verifier", "your_code_verifier") We're trying to avoid using this because
	// we're acquiring the code on the backend hence this is not required

	// Create an HTTP client.
	client := &http.Client{}

	// Create a request.
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("secret_key", a.Config.ClientSecret)

	// Perform the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response.
	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	responseBody := buffer.Bytes()

	// Create an instance of AccessTokenResponse
	var response zaloAccessTokenResponse

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	expiryParsed, err := strconv.Atoi(response.ExpiresIn)
	if err != nil {
		return nil, err
	}

	expiry := expiryParsed - 100

	result := oauth2.Token{
		AccessToken:  response.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: response.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(expiry) * time.Second),
	}

	return &result, nil
}

type zaloAccessTokenResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             string `json:"expires_in"`
}

type zaloUserResponse struct {
	Birthday    string `json:"birthday"`
	Gender      string `json:"gender"`
	IsSensitive bool   `json:"is_sensitive"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Error       int    `json:"error"`
	Message     string `json:"message"`
	Picture     struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}
