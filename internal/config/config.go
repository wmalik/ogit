package config

import (
	"encoding/json"
	"os/user"
	"path"

	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
	"github.com/wmalik/ogit/internal/utils"
)

var (
	ErrReadingFile = errors.New("the provided file could not be read")
	ErrJSONFile    = errors.New("the provided file could not be parsed")
)

type Config struct {
	Colors struct {
		ClonedRepoFG    *lipgloss.AdaptiveColor `json:"cloned_repo_fg"`
		DimmedColorFG   *lipgloss.AdaptiveColor `json:"dimmed_color_fg"`
		SelectedColorFG *lipgloss.AdaptiveColor `json:"selected_color_fg"`
		SelectedColorBG *lipgloss.AdaptiveColor `json:"selected_color_bg"`
		TitleBarFG      *lipgloss.AdaptiveColor `json:"title_bar_fg"`
		TitleBarBG      *lipgloss.AdaptiveColor `json:"title_bar_bg"`
		StatusMessageFG *lipgloss.AdaptiveColor `json:"status_message_fg"`
		StatusErrorFG   *lipgloss.AdaptiveColor `json:"status_error_fg"`
	} `json:"colors"`
}

func ReadConfig(path string) (*Config, error) {
	var cfg Config
	configFile, err := utils.ReadFile(path)
	if err != nil {
		return nil, utils.ErrorWithCause(ErrReadingFile, err.Error())
	}

	if err = json.Unmarshal(configFile, &cfg); err != nil {
		return nil, utils.ErrorWithCause(ErrJSONFile, err.Error())
	}

	return &cfg, nil
}

func GetConfigPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(user.HomeDir, ".config", "ogit", "config.json"), nil
}
