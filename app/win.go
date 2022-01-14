package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type dateContainerData struct {
	Date string `uri:"dt" binding:"required"`
}

type dateIntervalContainerData struct {
	From string `uri:"from" binding:"required"`
	To   string `uri:"to" binding:"required"`
}

type winListData struct {
	Items []winOnDayData `json:"items"`
}

type winOnDayData struct {
	Date string  `json:"date"`
	Win  winData `json:"win"`
}

type winData struct {
	Text          string   `json:"text"`
	OverallResult int      `json:"overall"`
	Priorities    []string `json:"priorities"`
}

type winDayListData struct {
	Items []string `json:"items"`
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
	/*time.Sleep(300 * time.Millisecond)
	toBadRequest(c, fmt.Errorf("Something went wrong returning win"))
	return*/

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
		err := fmt.Errorf("invalid value '%s' for 'text', should be less than %d characters long",
			win.Text,
			WIN_TEXT_MAX_LENGTH)
		toBadRequest(c, err)
		return
	}
	if !isWinOverallResultValid(win.OverallResult) {
		err := fmt.Errorf("invalid value '%s' for 'overall', should be a number in [0:4] range",
			win.Text)
		toBadRequest(c, err)
		return
	}
	if !isWinPriorityListValid(win.Priorities) {
		err := fmt.Errorf("invalid value '%s' for 'priorities', max %d items allowed, non-empty and less than %d characters long",
			win.Text,
			WIN_PRIORITIES_MAX_SIZE,
			PRIORITY_ID_MAX_LENGTH)
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

func handleGetWins(c *gin.Context, userId string, email string) {
	// get date from URL
	var dateIntervalContainer dateIntervalContainerData
	if err := c.ShouldBindUri(&dateIntervalContainer); err != nil {
		toBadRequest(c, err)
		return
	}

	// sanitize
	if !isDateValid(dateIntervalContainer.From) {
		err := fmt.Errorf("invalid value '%s' for 'from'", dateIntervalContainer.From)
		toBadRequest(c, err)
		return
	}
	if !isDateValid(dateIntervalContainer.To) {
		err := fmt.Errorf("invalid value '%s' for 'to'", dateIntervalContainer.To)
		toBadRequest(c, err)
		return
	}
	if !isIntervalValid(dateIntervalContainer.From,
		dateIntervalContainer.To,
		WINS_INTERVAL_REQUESTED_MAX_DAYS) {
		err := fmt.Errorf(
			"invalid value for the interval '%s' - '%s', 'from' should be earlier than 'to', max %d days allowed",
			dateIntervalContainer.From,
			dateIntervalContainer.To,
			WINS_INTERVAL_REQUESTED_MAX_DAYS)
		toBadRequest(c, err)
		return
	}

	wins, err := getWins(userId, dateIntervalContainer.From, dateIntervalContainer.To)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	for _, winOnDay := range wins {
		if winOnDay.Win.Priorities == nil {
			winOnDay.Win.Priorities = []string{}
		}
	}

	winList := winListData{
		Items: wins,
	}

	// TODO: this is for testing, remove when no more useful
	/*time.Sleep(300 * time.Millisecond)
	toBadRequest(c, fmt.Errorf("Something went wrong returning win list"))
	return*/

	toSuccess(c, winList)
}

func handleGetWinDays(c *gin.Context, userId string, email string) {
	// get date from URL
	var dateIntervalContainer dateIntervalContainerData
	if err := c.ShouldBindUri(&dateIntervalContainer); err != nil {
		toBadRequest(c, err)
		return
	}

	// sanitize
	if !isDateValid(dateIntervalContainer.From) {
		err := fmt.Errorf("invalid value '%s' for 'from'", dateIntervalContainer.From)
		toBadRequest(c, err)
		return
	}
	if !isDateValid(dateIntervalContainer.To) {
		err := fmt.Errorf("invalid value '%s' for 'to'", dateIntervalContainer.To)
		toBadRequest(c, err)
		return
	}
	if !isIntervalValid(
		dateIntervalContainer.From,
		dateIntervalContainer.To,
		WIN_DAYS_INTERVAL_REQUESTED_MAX_DAYS) {
		err := fmt.Errorf("invalid value for the interval '%s' - '%s', 'from' should be earlier than 'to', max %d days allowed",
			dateIntervalContainer.From,
			dateIntervalContainer.To,
			WIN_DAYS_INTERVAL_REQUESTED_MAX_DAYS)
		toBadRequest(c, err)
		return
	}

	winDays, err := getWinDays(userId, dateIntervalContainer.From, dateIntervalContainer.To)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	winDayList := winDayListData{
		Items: winDays,
	}

	// TODO: this is for testing, remove when no more useful
	/*time.Sleep(300 * time.Millisecond)
	toBadRequest(c, fmt.Errorf("Something went wrong returning win days list"))
	return*/

	//time.Sleep(2000 * time.Millisecond)
	toSuccess(c, winDayList)
}
