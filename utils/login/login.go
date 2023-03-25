package login

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	d "server/utils/db_setting"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	gg_oauth "google.golang.org/api/oauth2/v2"
)

// Global variables
var method string = ""
var User_id int

/******************** Obtain Access Token with Authorisation Codes ********************/
// Constants holder for Google OAuth2
var (
	googleClientID     = "144575449170-njifilpn4vuu5ujmst66qctf5g1uufg5.apps.googleusercontent.com"
	googleClientSecret = "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"
	googleRedirectURL  = "http://localhost:8080/Login"
)

// Create an OAuth2.Config struct for Facebook
var gg_conf = &oauth2.Config{
	ClientID:     googleClientID,
	ClientSecret: googleClientSecret,
	Endpoint:     google.Endpoint,
	RedirectURL:  googleRedirectURL,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Google authentication page.
	url := gg_conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Fprintf(w, "Visit the URL for obtaining the Google's OAuth code: %v\n", url)
}

// Constants holder for Facebook OAuth2
var (
	facebookClientID     = "532885048762005"
	facebookClientSecret = "5f8d5c626ca964a41e432419928aabe9"
	facebookRedirectURL  = "http://localhost:8080/Login"
)

// Create an OAuth2.Config struct for Facebook
var fb_conf = &oauth2.Config{
	ClientID:     facebookClientID,
	ClientSecret: facebookClientSecret,
	Endpoint:     facebook.Endpoint,
	RedirectURL:  facebookRedirectURL,
	Scopes:       []string{"email"},
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Facebook authentication page.
	url := fb_conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Fprintf(w, "Visit the URL for obtaining the Facebook's OAuth code: %v\n", url)
}

// Constants holder for GitHub OAuth2
var (
	githubClientID     = "Iv1.e296d38b65295a80"
	githubClientSecret = "b892a9b4dbe8d439d3a7f910d4555f1c6b7f6216"
	githubRedirectURL  = "http://localhost:8080/Login"
)

// Create an OAuth2.Config struct for GitHub
var git_conf = &oauth2.Config{
	ClientID:     githubClientID,
	ClientSecret: githubClientSecret,
	Scopes:       []string{"repo", "user"},
	Endpoint:     github.Endpoint,
	RedirectURL:  githubRedirectURL,
}

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the GitHub authentication page.
	url := git_conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Fprintf(w, "Visit the URL for obtaining the GitHub's OAuth code: %v\n", url)
}

func googleAuth(w http.ResponseWriter, conf *oauth2.Config, code string) bool {
	// Exchange authorization code for access token
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintf(w, "Invalid access token.\n")
		AuthMsg(w)
		return false
	}
	// Verify token with Google
	client := conf.Client(context.Background(), token)
	oauth2Service, err := gg_oauth.New(client)
	if err != nil {
		fmt.Fprintf(w, "Verification with Google failed.\nPlease login with other methods.\n")
		return false
	}
	tokenInfo, err := oauth2Service.Tokeninfo().AccessToken(token.AccessToken).Do()
	if err != nil {
		fmt.Fprintf(w, "An error occured when calling Tokeninfo API.\n")
		return false
	}
	// Print the authenticated user's details.
	fmt.Fprintf(w, "You have successfully logged in as: %v\n", tokenInfo.Email)
	fmt.Fprintf(w, "Token: %s\n", token.AccessToken)
	welcomeMsg(w)
	InsertUser(w, strings.Split(tokenInfo.Email, "@")[0])
	return true
}

func facebookAuth(w http.ResponseWriter, conf *oauth2.Config, code string) bool {
	// Build the access token request
	tokenURL := "https://graph.facebook.com/v12.0/oauth/access_token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", facebookClientID)
	data.Set("client_secret", facebookClientSecret)
	data.Set("redirect_uri", facebookRedirectURL)
	// Make the access token request
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		fmt.Fprintf(w, "Invalid access token.\n")
		AuthMsg(w)
		return false
	}
	defer resp.Body.Close()
	// Parse the access token response
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		fmt.Fprintf(w, "An error occured when parsing the access token response.\n")
		return false
	}
	// Verify the access token by making a request to the Facebook Graph API
	graphAPIURL := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s", tokenResp.AccessToken)
	resp, err = http.Get(graphAPIURL)
	if err != nil {
		fmt.Fprintf(w, "Verification with Facebook failed.\n")
		return false
	}
	defer resp.Body.Close()
	// Parse the response from the Facebook Graph API
	var graphResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&graphResp)
	if err != nil {
		fmt.Fprintf(w, "An error occured when parsing the response from the Facebook Graph API.\n")
		return false
	}
	// Check if the response contains an error message
	if errorMessage, ok := graphResp["error"].(map[string]interface{}); ok {
		fmt.Printf("Error: %s\n", errorMessage["message"])
		return false
	}
	// Print the authenticated user's details.
	fmt.Fprintf(w, "You have successfully logged in as: %s\n", graphResp["name"])
	fmt.Fprintf(w, "Token: %s\n", tokenResp.AccessToken)
	welcomeMsg(w)
	InsertUser(w, fmt.Sprint(graphResp["name"]))
	return true
}

func githubAuth(w http.ResponseWriter, conf *oauth2.Config, code string) bool {
	// Exchange the authorization code for an access token.
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintf(w, "Invalid access token.\n")
		AuthMsg(w)
		return false
	}
	// Create a new HTTP client with the access token.
	client := conf.Client(context.Background(), token)
	// Make a request to the API to get the authenticated user's details.
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Fprintf(w, "Invalid client request.\n")
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Invalid client request.\n")
		return false
	}
	// Check if the request was successful.
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "API request failed with status code %v", resp.StatusCode)
		return false
	}
	// Parse the response body to get the authenticated user's details.
	defer resp.Body.Close()
	var user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Fprintf(w, "An error occured when decoding the data.\n")
		return false
	}
	// Print the authenticated user's details.
	fmt.Fprintf(w, "You have successfully logged in as: (%s)\n", user.Login)
	fmt.Fprintf(w, "Token: %s\n", token.AccessToken)
	welcomeMsg(w)
	InsertUser(w, user.Login)
	return true
}

/******************** Verify Token form Header ********************/
func verifyAccessToken(method int, token string) bool {
	// Create a new HTTP request to the debug_token endpoint
	var url string
	if method == 1 {
		url = fmt.Sprintf("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=%s", token)
	} else if method == 2 {
		url = fmt.Sprintf("https://graph.facebook.com/me?access_token=%s", token)
	} else if method == 3 {
		url = "https://api.github.com/user"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	if method == 3 {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	// Send the HTTP request and check the response status code
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API request failed with status code %d", resp.StatusCode)
		return false
	}

	// The access token is valid
	return true
}

func CheckCurlHeader(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("token")
	var IsAuth bool = false
	if token != "" {
		if verifyAccessToken(1, token) || verifyAccessToken(2, token) || verifyAccessToken(3, token) {
			IsAuth = true
		} else {
			fmt.Fprintf(w, "Invalid token, please login to gain access. %s\n", token)
			IsAuth = false
		}
	}
	return IsAuth
}

// helper function to redirect respective auths
func authorizeUser(w http.ResponseWriter, code string) {
	if method == "1" {
		googleAuth(w, gg_conf, code)
	} else if method == "2" {
		facebookAuth(w, fb_conf, code)
	} else if method == "3" {
		githubAuth(w, git_conf, code)
	} else {
		fmt.Fprintf(w, "Please select a valid authentication method: 1. Google, 2. Facebook, 3. GitHub\n")
	}
}

// main login function
func SignIn(w http.ResponseWriter, r *http.Request) {
	// Extract value of login method from request param
	if r.URL.Query().Get("method") != "" {
		method = r.URL.Query().Get("method")
	}
	code := r.URL.Query().Get("code")

	if code != "" {
		// Wait for the user to authorize the app and receive the authorization code.
		authorizeUser(w, code)
	} else {
		// If login via GitHub
		if method == "1" {
			handleGoogleLogin(w, r)
		} else if method == "2" {
			handleFacebookLogin(w, r)
		} else if method == "3" {
			handleGitHubLogin(w, r)
		} else {
			fmt.Fprintf(w, "Please select a valid authentication method: 1. Google, 2. Facebook, 3. GitHub\n")
		}
	}
}

/******************** Messages ********************/
func AuthMsg(w http.ResponseWriter) {
	fmt.Fprintf(w, "Please login to gain access to the server.\n")
	fmt.Fprintf(w, "You can use http://localhost:8080/Login?method={auth_type} to get the authentication tokens that can be passed in via Authorization header or as part of the POST body.\n")
	fmt.Fprintf(w, "Below is the list of available auth_type:-\n")
	fmt.Fprintf(w, "1. Google\n")
	fmt.Fprintf(w, "2. Facebook\n")
	fmt.Fprintf(w, "3. GitHub\n")
}

func welcomeMsg(w http.ResponseWriter) {
	fmt.Fprintf(w, "You are authorized to use http://localhost:8080/{path}\n")
	fmt.Fprintf(w, "Below are 4 of the available paths:-\n")
	fmt.Fprintf(w, "1. List\n")
	fmt.Fprintf(w, "2. Add\n")
	fmt.Fprintf(w, "3. Mark-complete\n")
	fmt.Fprintf(w, "4. Delete\n")
}

/******************** Register New User ********************/
func InsertUser(w http.ResponseWriter, username string) {
	// Connect to db
	db_settings := fmt.Sprintf("%s:%s@%s/%s", d.DbSettings()["user"], d.DbSettings()["pw"], d.DbSettings()["conn"], d.DbSettings()["schema"])
	db, err := sql.Open("mysql", db_settings)
	if err != nil {
		fmt.Fprintf(w, "An error occured when connecting to database.")
		return
	}

	defer db.Close()

	// Query data to check if the user exist then only we proceed to add user
	slt_stmt, err := db.Prepare("SELECT COUNT(id) AS length FROM users WHERE user_name = ?")
	if err != nil {
		fmt.Fprintf(w, "An error occured when preparing statement.")
		return
	}

	defer slt_stmt.Close()

	rows, err := slt_stmt.Query(username)
	if err != nil {
		fmt.Fprintf(w, "An error occured when executing statement.")
		return
	}

	defer rows.Close()

	// Insert user
	var length int
	for rows.Next() {
		err := rows.Scan(&length)
		if err != nil {
			fmt.Fprintf(w, "An error occured when scanning row.")
		}
		if length <= 0 {
			// Prepare statement
			insert_stmt, err := db.Prepare("INSERT INTO users(user_name, create_time) VALUES(?, ?)")
			if err != nil {
				fmt.Fprintf(w, "An error occured when preparing statement.")
			}

			// Execute statement
			_, err = insert_stmt.Exec(username, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				fmt.Fprintf(w, "An error occured when executing statement.")
			}

			defer insert_stmt.Close()
		}
	}

	// Query for user_id
	find_user_id, err := db.Prepare("SELECT * FROM users WHERE user_name = ?")
	if err != nil {
		fmt.Fprintf(w, "An error occured when preparing statement.")
		return
	}

	defer find_user_id.Close()

	user_ids, err := find_user_id.Query(username)
	if err != nil {
		fmt.Fprintf(w, "An error occured when executing statement.")
		return
	}

	// Extract data
	var id int
	var user_name string
	var create_time string
	for user_ids.Next() {
		err := user_ids.Scan(&id, &user_name, &create_time)
		if err != nil {
			fmt.Fprintf(w, "An error occured when scanning row.")
		}
		// Update global user_id
		User_id = id
		fmt.Fprintf(w, "Your user id is %d", id)
	}
}
