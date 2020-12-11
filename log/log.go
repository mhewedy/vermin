package log

import (
	"fmt"
	"github.com/mhewedy/vermin/db"
	"os"
	"path/filepath"
	"time"
)

func Debug(format string, a ...interface{}) {
	if _, ok := os.LookupEnv("VERMIN_DEBUG"); ok {
		Info(format, a)
	}
}

func Info(format string, a ...interface{}) {
	year, month, _ := time.Now().Date()
	logFilePath := filepath.Join(db.BaseDir, fmt.Sprintf("log_%d-%d.log", year, month))

	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString(fmt.Sprintf(format+"\n", a...))
}
