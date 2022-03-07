package clear

import (
	"context"
	"log"
	"ogit/internal/gitconfig"
	"os"
	"path"
)

func HandleCommandDefault(ctx context.Context) error {
	gitConf, err := gitconfig.ReadGitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	dbPath := path.Join(gitConf.StoragePath(), "ogit.db")
	return os.Remove(dbPath)
}
