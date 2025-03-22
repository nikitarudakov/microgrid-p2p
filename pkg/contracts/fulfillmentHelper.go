package contracts

import "time"

func parseServiceWindowEdge(edge string) (time.Time, error) {
	parsedTime, err := time.Parse("15:04:05", edge)
	if err != nil {
		return time.Time{}, err
	}

	// Get today's date
	now := time.Now()
	combined := time.Date(
		now.Year(), now.Month(), now.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
		0, now.Location(),
	)

	return combined, nil
}
