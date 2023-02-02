package config_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/wmalik/ogit/internal/config"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfig(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}

func (s *ConfigTestSuite) TestReadConfigReturnsAnErrorIfTheConfigFileDoesNotExist() {
	filePath := "./testdata/missing-file.json"
	cfg, err := config.ReadConfig(filePath)
	s.Require().Nil(cfg)
	s.Error(err)
	s.Equal(config.ErrReadingFile, errors.Cause(err))
}

func (s *ConfigTestSuite) TestReadConfigReturnsAnErrorIfTheConfigFileIsNotValidJSON() {
	filePath := "./testdata/invalid.json"
	cfg, err := config.ReadConfig(filePath)
	s.Require().Nil(cfg)
	s.Error(err)
	s.Equal(config.ErrJSONFile, errors.Cause(err))
}

func (s *ConfigTestSuite) TestReadConfigReturnsAnErrorIfTheConfigFileHasInvalidData() {
	filePath := "./testdata/invalid-fields.json"
	cfg, err := config.ReadConfig(filePath)
	s.Require().Nil(cfg)
	s.Error(err)
	s.Equal(config.ErrJSONFile, errors.Cause(err))
}

func (s *ConfigTestSuite) TestReadConfigReturnsTheFullConfigData() {
	filePath := "./testdata/config.json"
	cfg, err := config.ReadConfig(filePath)
	s.Require().NoError(err)
	s.NotNil(cfg)
	s.NotNil(cfg.Colors.ClonedRepoFG)
	s.NotNil(cfg.Colors.DimmedColorFG)
	s.NotNil(cfg.Colors.SelectedColorFG)
	s.NotNil(cfg.Colors.SelectedColorBG)
	s.NotNil(cfg.Colors.TitleBarFG)
	s.NotNil(cfg.Colors.TitleBarBG)
	s.NotNil(cfg.Colors.StatusMessageFG)
	s.NotNil(cfg.Colors.StatusErrorFG)
}
