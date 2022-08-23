package template_test

import (
	"fmt"
	"strings"

	"github.com/xn3cr0nx/email-service/internal/template"
)

func (s *TemplateTestSuite) TestPathByType() {
	path := template.PathByType("email:layout")
	s.Equal(path, fmt.Sprintf("%s/layout.html", s.BaseDir))
}

func (s *TemplateTestSuite) TestFillLayout() {
	path := template.PathByType("email:welcome")
	html := string(s.Cache.Get(path))

	name := "Test"
	URL := "test.org"

	filledHtml := fmt.Sprintf(html, name, URL, URL)
	s.False(strings.Contains(filledHtml, "%s"))

	final := template.FillLayout(filledHtml)
	s.False(strings.Contains(final, "%s"))
}
