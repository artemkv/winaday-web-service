package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type winListShortData struct {
	Items []winOnDayShortData `json:"items"`
}

type winOnDayShortData struct {
	Date string       `json:"date"`
	Win  winShortData `json:"win"`
}

type winShortData struct {
	OverallResult int      `json:"overall"`
	Priorities    []string `json:"priorities"`
}

func handleGetWinStats(c *gin.Context, userId string, email string) {
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
		STATS_INTERVAL_REQUESTED_MAX_DAYS) {
		err := fmt.Errorf("invalid value for the interval '%s' - '%s', 'from' should be earlier than 'to', max %d days allowed",
			dateIntervalContainer.From,
			dateIntervalContainer.To,
			STATS_INTERVAL_REQUESTED_MAX_DAYS)
		toBadRequest(c, err)
		return
	}

	winDays, err := getWinDayStats(userId, dateIntervalContainer.From, dateIntervalContainer.To)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	winDayList := winListShortData{
		Items: winDays,
	}

	// TODO: this is for testing, remove when no more useful
	// time.Sleep(300 * time.Millisecond)
	/*toBadRequest(c, fmt.Errorf("Something went wrong returning stats"))
	return*/

	toSuccess(c, winDayList)
}
