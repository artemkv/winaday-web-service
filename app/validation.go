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
	if len(text) == 0 || len(text) > 1000 {
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
