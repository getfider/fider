package query

import "time"

type FetchRecentSupressions struct {
	StartTime time.Time

	//Output
	EmailAddresses []string
}
