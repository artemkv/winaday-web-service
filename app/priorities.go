package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type priorityListData struct {
	Items []priorityData `json:"items"`
}

// When adding properties here, do not forget to update encode/decode functions!
type priorityData struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	Color     int    `json:"color"`
	IsDeleted bool   `json:"deleted"`
}

func handleGetPriorities(c *gin.Context, userId string, email string) {

	// TODO: this is for testing, remove when no more useful
	//time.Sleep(300 * time.Millisecond)
	//toBadRequest(c, fmt.Errorf("Something went wrong returning priorities"))
	//return

	priorityList, err := getPriorities(userId)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	if priorityList == nil {
		priorityList = &priorityListData{
			Items: []priorityData{},
		}
	}

	toSuccess(c, priorityList)
}

func handlePostPriorities(c *gin.Context, userId string, email string) {
	// get win data from the POST body
	var priorities priorityListData
	if err := c.ShouldBindJSON(&priorities); err != nil {
		toBadRequest(c, err)
		return
	}

	// sanitize
	if !isPriorityListLengthValid(priorities) {
		err := fmt.Errorf("too many items in a priority list, max 9 active and 100 total allowed")
		toBadRequest(c, err)
		return
	}
	for _, p := range priorities.Items {
		if !isPriorityIdValid(p.Id) {
			err := fmt.Errorf("invalid id, should not be empty and less than 100 characters long")
			toBadRequest(c, err)
			return
		}
		if !isPriorityTextValid(p.Text) {
			err := fmt.Errorf("invalid value '%s' for 'text', should be less than 100 characters long", p.Text)
			toBadRequest(c, err)
			return
		}
		if !isPriorityColorValid(p.Color) {
			err := fmt.Errorf("invalid value '%d' for 'color', should be a number 0 <= x < 100", p.Color)
			toBadRequest(c, err)
			return
		}
	}

	updatedAt := generateTimestamp()
	err := updatePriorities(userId, priorities, updatedAt)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	toSuccess(c, priorities)
}
