package internal

import (
	"os"

	"secretary/alpha/pkg/utils"
)

// PrintBanner prints the application banner
func PrintBanner() {
	content, err := os.ReadFile("banner.txt")
	if err != nil {
		utils.Error("Error reading the banner.txt file: " + err.Error())
		return
	}
	utils.Info("\n" + string(content))
}
