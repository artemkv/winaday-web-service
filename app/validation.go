package app

import "time"

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
	if len(text) > 1000 {
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

func isWinPriorotyListValid(priorities []string) bool {
	if len(priorities) > 100 {
		return false
	}

	for _, p := range priorities {
		if p == "" || len(p) > 100 {
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

	return active <= 9 && total <= 200
}

func isPriorityTextValid(text string) bool {
	if len(text) > 100 {
		return false
	}

	return true
}

func isPriorityIdValid(id string) bool {
	return id != "" && len(id) < 100
}

func isPriorityColorValid(color int) bool {
	return color >= 0 && color < 100
}
