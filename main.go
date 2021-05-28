package main

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
	"github.com/pin-yu/datalab-name-registration/backend"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var fullchain = filepath.Join(BasePath(), "certs/fullchain.pem")
var privkey = filepath.Join(BasePath(), "certs/privkey.pem")
var serviceAccountSecret = filepath.Join(BasePath(), "credential/google-service.json")

var ErrUnauthorized = errors.New("please login")

func authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("email") != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/register")
		}

		c.Next()
	}
}

func registerAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("email") != nil {
			c.Next()
		}

		c.Redirect(http.StatusSeeOther, "/login")
	}
}

func main() {
	r := gin.Default()

	// store := cookie.NewStore([]byte(backend.LoadSecret()))
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(backend.LoadSecret()))

	// store-level sessions
	store.Options(sessions.Options{
		MaxAge:   3600 * 24, // 24 hours
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	r.Use(sessions.Sessions("authentication-session", store))

	r.Static("/img", "./frontend/img")

	root := r.Group("/")
	root.Use(authentication())
	root.GET("/", func(c *gin.Context) {
	})

	register := r.Group("/register")
	register.Use(registerAuthentication())
	register.Static("/", filepath.Join(BasePath(), "frontend/register"))

	ctx := context.Background()
	srv, _ := sheets.NewService(ctx, option.WithCredentialsFile(serviceAccountSecret), option.WithScopes(sheets.SpreadsheetsScope))
	pCtrl := backend.PublicController{SheetService: srv}
	register.POST("/come", pCtrl.RegisterCome)
	register.POST("/leave", pCtrl.RegisterLeave)

	// separate a group because login page doesn't have to authenticate
	login := r.Group("/login")
	login.Static("/", filepath.Join(BasePath(), "frontend/login"))

	logout := r.Group("/logout")
	logout.GET("/", backend.GoogleOauthLogout)

	oauth := r.Group("/oauth")
	oauth.GET("/", backend.GoogleOauthLogin)

	callback := r.Group("/callback")
	callback.GET("/", backend.GoogleOauthCallBack)

	r.RunTLS(":5003", fullchain, privkey)
}

func BasePath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
