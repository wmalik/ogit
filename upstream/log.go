package upstream

import "fmt"

func logPaginationStatus(upstream string, owner string, numRepos int, remainingAPILimit string) {
	if owner == "" {
		owner = "current_user"
	}
	fmt.Printf("    [%s/%s] fetched %d repositories, remaining api calls: %s\n",
		upstream,
		owner,
		numRepos,
		remainingAPILimit,
	)
}

func logAuthenticatedUser(upstream string, username string) {
	fmt.Printf("\nAuthenticated with %s as %s\n", upstream, username)
}
