package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type dateContainerData struct {
	Date string `uri:"dt" binding:"required"`
}

type winData struct {
	Text          string   `json:"text"`
	OverallResult int      `json:"overall"`
	Priorities    []string `json:"priorities"`
}

func handleGetWin(c *gin.Context, userId string, email string) {
	// get date from URL
	var dateContainer dateContainerData
	if err := c.ShouldBindUri(&dateContainer); err != nil {
		toBadRequest(c, err)
		return
	}

	win, err := getWin(userId, dateContainer.Date)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	if win == nil {
		win = &winData{
			Text:          "",
			OverallResult: 0,
			Priorities:    []string{},
		}
	}

	if win.Priorities == nil {
		win.Priorities = []string{}
	}

	// TODO: this is for testing, remove when no more useful
	//time.Sleep(300 * time.Millisecond)
	//toBadRequest(c, fmt.Errorf("Something went wrong returning win"))
	//return

	toSuccess(c, win)
}

func handlePostWin(c *gin.Context, userId string, email string) {
	// get date from URL
	var dateContainer dateContainerData
	if err := c.ShouldBindUri(&dateContainer); err != nil {
		toBadRequest(c, err)
		return
	}

	// get win data from the POST body
	var win winData
	if err := c.ShouldBindJSON(&win); err != nil {
		toBadRequest(c, err)
		return
	}

	// sanitize
	if !isDateValid(dateContainer.Date) {
		err := fmt.Errorf("invalid value '%s' for 'date'", dateContainer.Date)
		toBadRequest(c, err)
		return
	}
	if !isWinTextValid(win.Text) {
		err := fmt.Errorf("invalid value '%s' for 'text', should be less than 1000 characters long", win.Text)
		toBadRequest(c, err)
		return
	}
	if !isWinOverallResultValid(win.OverallResult) {
		err := fmt.Errorf("invalid value '%s' for 'overall', should be a number in [0:4] range", win.Text)
		toBadRequest(c, err)
		return
	}
	if !isWinPriorotyListValid(win.Priorities) {
		err := fmt.Errorf("invalid value '%s' for 'priorities', max 100 items allowed, non-empty and less than 100 characters long", win.Text)
		toBadRequest(c, err)
		return
	}

	err := updateWin(userId, dateContainer.Date, win)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	toSuccess(c, win)
}
