package template

import (
	"fmt"
	"path/filepath"
)

// List of template types.
const (
	Layout        = "email:layout"
	WelcomeEmail  = "email:welcome"
	ReminderEmail = "email:reminder"
)

func PathByType(taskType string) string {
	return map[string]string{
		Layout:       filepath.Join(cache.Dir, "layout.html"),
		WelcomeEmail: filepath.Join(cache.Dir, "welcome.html"),
		// ReminderEmail: "",
	}[taskType]
}

func FillLayout(content string) string {
	layout := string(cache.Get(PathByType(Layout)))
	return fmt.Sprintf(layout, content)
}
