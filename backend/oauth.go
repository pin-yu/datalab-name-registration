package backend

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	api "google.golang.org/api/oauth2/v2"
)

var ErrState = errors.New("state error")
var googleConfig = googleOauthConfig()

func googleOauthConfig() *oauth2.Config {
	credential := LoadGoogleCredential()
	conf := &oauth2.Config{
		ClientID:     credential.ClientID,
		ClientSecret: credential.ClientSecret,
		RedirectURL:  "https://shwu16.cs.nthu.edu.tw:5003/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return conf
}

func generateState() string {
	return uuid.New().String()
}

func GoogleOauthLogin(c *gin.Context) {
	state := generateState()

	redirectURL := googleConfig.AuthCodeURL(state)

	// prevent CSRF
	session := sessions.Default(c)
	session.Set("state", state)
	err := session.Save()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, redirectURL)
}

func GoogleOauthCallBack(c *gin.Context) {
	session := sessions.Default(c)

	// check whether the state is consistent. Prevent CSRF
	state := session.Get("state")
	if state != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, ErrState)
	}

	// after receiving the code from clients, send the code to the google api server to request a token
	// we are able to get users' information via the token.
	code := c.Query("code")
	token, err := googleConfig.Exchange(c, code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// get userInfo via the token
	client := googleConfig.Client(context.TODO(), token)
	svc, err := api.New(client) // don't care this deprecated now!!!
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userInfo, err := svc.Userinfo.V2.Me.Get().Do()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Set("name", userInfo.Name)
	session.Set("email", userInfo.Email)
	session.Options(sessions.Options{Path: "/", MaxAge: 3600 * 24 * 30}) // one month
	err = session.Save()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// redirect to the root
	redirectUrl, _ := url.Parse("/register/")
	c.Redirect(http.StatusSeeOther, redirectUrl.Path)
}

func GoogleOauthLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("email", "") // this will mark the session as "written"
	session.Set("name", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1}) // this sets the cookie with a MaxAge of 0
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
