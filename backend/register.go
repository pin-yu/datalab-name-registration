package backend

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"google.golang.org/api/sheets/v4"
)

var status = "status"
var inLab = "in_lab"
var notInLab = "not_in_lab"
var spreadsheetId = "1bIXTsOsL0UCYYyknoYQ3n3Faw2snmQtQiOP7uoIMWU4"

type PublicController struct {
	SheetService *sheets.Service
}

func ConvertRowValue(data1D []string) *sheets.ValueRange {
	s1D := make([]interface{}, len(data1D))
	for i, v := range data1D {
		s1D[i] = v
	}
	s2D := [][]interface{}{}
	s2D = append(s2D, s1D)
	rb := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         s2D,
	}
	return rb
}

func (ctx *PublicController) RegisterCome(c *gin.Context) {
	session := sessions.Default(c)

	currentTime := time.Now()

	name := session.Get("name")
	email := session.Get("email")
	currentStatus := fmt.Sprintf("%v", session.Get(status))

	if currentStatus != inLab {
		session.Set(status, inLab)
		session.Save()

		c.String(http.StatusCreated, fmt.Sprintf("%s come to the lab at %s", name, currentTime.Format(time.UnixDate)))

		registerType := "come"
		data1D := []string{fmt.Sprintf("%s", name), fmt.Sprintf("%d", currentTime.Unix()), registerType, fmt.Sprintf("%s", email)}
		RowValue := ConvertRowValue(data1D)
		ctx.SheetService.Spreadsheets.Values.Append(spreadsheetId, "DatalabService", RowValue).ValueInputOption("USER_ENTERED").Context(context.Background()).Do()

	} else {
		c.String(http.StatusAlreadyReported, fmt.Sprintf("%s have registered today", name))
	}
}

func (ctx *PublicController) RegisterLeave(c *gin.Context) {
	session := sessions.Default(c)

	currentTime := time.Now()

	name := session.Get("name")
	email := session.Get("email")
	currentStatus := fmt.Sprintf("%v", session.Get(status))

	if currentStatus != inLab {
		c.String(http.StatusBadRequest, fmt.Sprintf("%s has no register record today", name))
	} else {
		session.Set(status, notInLab)
		session.Save()
		c.String(http.StatusCreated, fmt.Sprintf("%s leave the lab at %s", name, currentTime.Format(time.UnixDate)), fmt.Sprintf("%s", email))

		registerType := "leave"
		data1D := []string{fmt.Sprintf("%s", name), fmt.Sprintf("%d", currentTime.Unix()), registerType}
		RowValue := ConvertRowValue(data1D)
		ctx.SheetService.Spreadsheets.Values.Append(spreadsheetId, "DatalabService", RowValue).ValueInputOption("USER_ENTERED").Context(context.Background()).Do()
	}
}
