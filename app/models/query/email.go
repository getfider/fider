package query

import "time"

type FetchRecentSuppressions struct {
	StartTime time.Time

	//Output
	EmailAddresses []string
}
