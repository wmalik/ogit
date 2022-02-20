package upstream

import (
	"log"
)

func logPaginationStatus(upstream string, owner string, numRepos, remainingPages int, remainingAPILimit string) {
	if owner == "" {
		owner = "current_user"
	}
	log.Printf("[%s:%s] fetched %d repositories, remaining pages %d, remaining API calls: %s",
		upstream,
		owner,
		numRepos,
		remainingPages,
		remainingAPILimit,
	)
}

func logAuthenticatedUser(upstream string, username string) {
	log.Printf("Authenticated with %s as %s", upstream, username)
}
