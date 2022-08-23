package template_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}

type TemplateTestSuite struct {
	suite.Suite

	Cache   *template.TemplateCache
	BaseDir string
}

func (s *TemplateTestSuite) SetupSuite() {
	logger.Setup()

	s.BaseDir = "templates_test"
	cache, err := template.NewTemplateCache(&s.BaseDir)
	s.Nil(err)
	s.Cache = cache
}
