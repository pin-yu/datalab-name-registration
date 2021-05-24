package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var status = "status"
var inLab = "in_lab"
var notInLab = "not_in_lab"

func RegisterCome(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		MaxAge:   3600 * 16, // 16 hours
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	currentTime := time.Now()

	name := session.Get("name")
	currentStatus := fmt.Sprintf("%v", session.Get(status))

	if currentStatus != inLab {
		session.Set(status, inLab)
		session.Options(sessions.Options{
			MaxAge: 3600 * 16, // 16hrs
		})
		session.Save()

		c.String(http.StatusCreated, fmt.Sprintf("%s come to the lab at %s", name, currentTime.Format(time.UnixDate)))
	} else {
		c.String(http.StatusAlreadyReported, fmt.Sprintf("%s have registered today", name))
	}
}

func RegisterLeave(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		MaxAge:   3600 * 8, // 16 hours
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	currentTime := time.Now()

	name := session.Get("name")
	currentStatus := fmt.Sprintf("%v", session.Get(status))

	if currentStatus != inLab {
		c.String(http.StatusBadRequest, fmt.Sprintf("%s has no register record today", name))
	} else {
		session.Set(status, notInLab)
		session.Options(sessions.Options{
			MaxAge: 3600 * 8, // 16hrs
		})
		session.Save()
		c.String(http.StatusCreated, fmt.Sprintf("%s leave the lab at %s", name, currentTime.Format(time.UnixDate)))
	}
}
