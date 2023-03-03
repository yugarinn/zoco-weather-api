package internal

import (
	"time"
)

// This function will go boom boom if the flatValues len is not
// 168. TODO: avoid that; maybe precheck and instantiate the map
// accordingly?
func CastHourlyValuesToWeek(flatValues []int) map[string][]int {
	dayCounter := 0
	weeklyValues := make(map[string][]int)

	for j := 0; j < 168; j += 24 {
		date := time.Now().AddDate(0, 0, dayCounter).Format("2006-01-02")

		weeklyValues[date] = flatValues[dayCounter * 24:dayCounter * 24 + 24]
		dayCounter++
    }

	return weeklyValues
}
