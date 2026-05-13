package utils

import "time"

func GetExpectedDelivery() string {

	start := time.Now().AddDate(0, 0, 4)
	end := time.Now().AddDate(0, 0, 5)

	return start.Format("02 Jan") +
		" - " +
		end.Format("02 Jan")
}