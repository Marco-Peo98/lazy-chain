package misc

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func LookupEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}

func LogEvent(log string) {
	if path := os.Getenv("DEBUG") ; path != "" {
		f, err := tea.LogToFile(path, log)
		if err != nil {
			fmt.Println("fatal: ", err)		// change to log library instead
		}
		defer f.Close()
	} 
}