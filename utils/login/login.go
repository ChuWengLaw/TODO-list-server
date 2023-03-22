package login

import (
	"context"
    "fmt"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/github"
    "net/http"
)

var (
    githubClientID     = "Iv1.e296d38b65295a80"
    githubClientSecret = "b892a9b4dbe8d439d3a7f910d4555f1c6b7f6216"
    githubRedirectURL  = "http://localhost:8080/"
)

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
    conf := &oauth2.Config{
        ClientID:     githubClientID,
        ClientSecret: githubClientSecret,
        Scopes:       []string{"repo", "user"},
        Endpoint:     github.Endpoint,
    }

    // Redirect the user to the GitHub authentication page.
    url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
    fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

    // Wait for the user to authorize the app and receive the authorization code.
    var code string
    fmt.Print("Enter the authorization code: ")
    fmt.Scan(&code)

    // Exchange the authorization code for an access token.
    token, err := conf.Exchange(context.Background(), code)
    if err != nil {
        panic(err)
    }

    // Use the access token to make authenticated API calls.
    client := conf.Client(context.Background(), token)
    resp, err := client.Get("https://api.github.com/user")
    if err != nil {
        panic(err)
    }

    // Print the user's profile.
    defer resp.Body.Close()
    fmt.Println(resp.Status)
    fmt.Println(resp.Header)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	// Extract value of login method from request param
	method := r.URL.Query().Get("method")

	// If login via Github
	if method == "2" {
		handleGitHubLogin(w, r)
	}
}