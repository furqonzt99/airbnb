package helper

import "time"

func countNight(checkinDate, checkoutDate time.Time) int {
	const HOURS_IN_A_DAY = 24

	differentTimeInhours := int(checkoutDate.Sub(checkinDate).Hours())
	convertToDay := differentTimeInhours / HOURS_IN_A_DAY

	return convertToDay
}