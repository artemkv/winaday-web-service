package app

import (
	"time"
)

const (
	WIN_TEXT_MAX_LENGTH     = 1000
	WIN_PRIORITIES_MAX_SIZE = 100

	PRIORITY_ID_MAX_LENGTH   = 100
	PRIORITY_TEXT_MAX_LENGTH = 100

	PRIORITIES_MAX_TOTAL        = 200
	PRIORITIES_ACTIVE_MAX_TOTAL = 9

	WINS_INTERVAL_REQUESTED_MAX_DAYS     = 50
	WIN_DAYS_INTERVAL_REQUESTED_MAX_DAYS = 50
	STATS_INTERVAL_REQUESTED_MAX_DAYS    = 400
)

func isUserIdValid(userId string) bool {
	return userId != ""
}

func isEmailValid(email string) bool {
	// TODO: check email format
	return email != ""
}

func isDateValid(date string) bool {
	d, err := time.Parse("20060102", date)
	if err != nil {
		return false
	}

	if d.Year() < 1900 || d.Year() > 2100 {
		return false
	}

	return true
}

func isWinTextValid(text string) bool {
	if len(text) > WIN_TEXT_MAX_LENGTH {
		return false
	}

	return true
}

func isWinOverallResultValid(r int) bool {
	if r < 0 || r > 4 {
		return false
	}

	return true
}

func isWinPriorityListValid(priorities []string) bool {
	if len(priorities) > WIN_PRIORITIES_MAX_SIZE {
		return false
	}

	for _, p := range priorities {
		if p == "" || len(p) > PRIORITY_ID_MAX_LENGTH {
			return false
		}
	}

	return true
}

func isPriorityListLengthValid(priorities priorityListData) bool {
	active := 0
	total := 0

	for _, p := range priorities.Items {
		if !p.IsDeleted {
			active++
		}
		total++
	}

	return active <= PRIORITIES_ACTIVE_MAX_TOTAL && total <= PRIORITIES_MAX_TOTAL
}

func isPriorityTextValid(text string) bool {
	if len(text) > PRIORITY_TEXT_MAX_LENGTH {
		return false
	}

	return true
}

func isPriorityIdValid(id string) bool {
	return id != "" && len(id) < PRIORITY_ID_MAX_LENGTH
}

func isPriorityColorValid(color int) bool {
	return color >= 0 && color < 100
}

func isIntervalValid(from string, to string, maxLength float64) bool {
	start, err := time.Parse("20060102", from)
	if err != nil {
		return false
	}

	end, err := time.Parse("20060102", to)
	if err != nil {
		return false
	}

	if start.After(end) {
		return false
	}

	duration := end.Sub(start)
	if duration.Hours()/24 > maxLength {
		return false
	}

	return true
}
