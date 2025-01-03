package audit

import (
	"fmt"
	"io"
	"os"
	"time"

	"secretary/alpha/utils"
)

type AuditEntry struct {
	Timestamp string `json:"timestamp"`
	User      string `json:"user"`
	Action    string `json:"action"`
}

func createAuditFile() (io.Writer, error) {
	directory := "data/audit/"
	utils.MakeDir(directory)
	currentDate := time.Now().Format("1999-01-28")
	filename := fmt.Sprintf("audit_%s", currentDate)
	file, err := os.OpenFile(directory+filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Audit(message string) error {
	file, err := createAuditFile()
	if err != nil {
		utils.Logger("fatal", err.Error())
		return err
	}
	defer file.(*os.File).Close()

	_, err = file.Write([]byte("[" + utils.CurrentTime() + " - " + utils.UUID() + "]" + " | " + message + "\n"))
	// TODO write to DB as well
	if err != nil {
		utils.Logger("fatal", err.Error())
		return err
	}
	return nil
}
