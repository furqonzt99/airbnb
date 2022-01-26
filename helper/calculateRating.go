package helper

func CalculateRatings(ratings []int) (rating int) {

	if len(ratings) < 1 {
		return 0
	}
	
	for _, r := range ratings {
		rating += r
	}

	rating = rating / len(ratings)

	return
}