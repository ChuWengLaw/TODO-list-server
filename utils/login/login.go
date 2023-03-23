package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Global variable for used in main
var IsAuth bool = false

// Constants holder for GitHub OAuth2
var (
	githubClientID     = "Iv1.e296d38b65295a80"
	githubClientSecret = "b892a9b4dbe8d439d3a7f910d4555f1c6b7f6216"
	githubRedirectURL  = "http://localhost:8080/Login"
)

// Create an OAuth2.Config struct for GitHub
var conf = &oauth2.Config{
	ClientID:     githubClientID,
	ClientSecret: githubClientSecret,
	Scopes:       []string{"repo", "user"},
	Endpoint:     github.Endpoint,
	RedirectURL:  githubRedirectURL,
}

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the GitHub authentication page.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Fprintf(w, "Visit the URL for obtaining the GitHub's OAuth code: %v\n", url)
}

func authorizeUser(w http.ResponseWriter, code string) {
	// Exchange the authorization code for an access token.
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP client with the access token.
	client := conf.Client(context.Background(), token)

	// Make a request to the GitHub API to get the authenticated user's details.
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Check if the request was successful.
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("API request failed with status code %v", resp.StatusCode))
	}

	// Parse the response body to get the authenticated user's details.
	defer resp.Body.Close()
	var user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		panic(err)
	}

	// Print the authenticated user's details.
	IsAuth = true
	fmt.Fprintf(w, "You have successfully logged in as: %s (%s)\n", user.Name, user.Login)
	fmt.Fprintf(w, "You are authorized to use http://localhost:8080/{path}\n")
	fmt.Fprintf(w, "Below are 4 of the available paths:-\n")
	fmt.Fprintf(w, "1. List\n")
	fmt.Fprintf(w, "2. Add\n")
	fmt.Fprintf(w, "3. Mark-complete\n")
	fmt.Fprintf(w, "4. Delete\n")
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	// Extract value of login method from request param
	method := r.URL.Query().Get("method")
	code := r.URL.Query().Get("code")

	if code != "" {
		// Wait for the user to authorize the app and receive the authorization code.
		authorizeUser(w, code)
	} else {
		// If login via GitHub
		if method == "3" {
			handleGitHubLogin(w, r)
		} else {
			fmt.Fprintf(w, "Please select a valid authentication method: 1. Google, 2. Facebook, 3. GitHub\n")
		}
	}
}

func AuthMsg(w http.ResponseWriter) {
	fmt.Fprintf(w, "Please login to gain access to the server.\n")
	fmt.Fprintf(w, "You can use http://localhost:8080/Login?method={auth_type} to get the authentication tokens that can be passed in via Authorization header or as part of the POST body.\n")
	fmt.Fprintf(w, "Below is the list of available auth_type:-\n")
	fmt.Fprintf(w, "1. Google\n")
	fmt.Fprintf(w, "2. Facebook\n")
	fmt.Fprintf(w, "3. GitHub\n")
}
