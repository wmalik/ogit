package browser

type updateStatusMsg string

type updateBottomStatusBarMsg string

type openURLMsg string

type cloneRepoMsg struct {
	repo repoItem
}
