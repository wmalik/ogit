package service

import (
	"fmt"
	"math"
	"time"
)

type APIUsage struct {
	Name          string
	Authenticated bool
	User          string
	Limit         int
	Remaining     int
	ResetsAt      time.Time
}

func (a *APIUsage) String() string {
	auth := "[not authenticated]"
	if a.Authenticated {
		auth = fmt.Sprintf("[Authenticated as %s]", a.User)
	}

	if a.Remaining > 60 {
		return auth
	}

	return fmt.Sprintf("%s [%s API Usage (%d of %d) (resets in %d mins)]",
		auth,
		a.Name,
		(a.Limit - a.Remaining),
		a.Limit,
		int(math.Ceil(time.Until(a.ResetsAt).Minutes())),
	)
}
