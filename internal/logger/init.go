package logger

import (
	"github.com/phuslu/log"
)

func Init() {
	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Writer: &log.MultiEntryWriter{
			&log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
			&log.FileWriter{
				Filename:     "logs/latest.log",
				FileMode:     0600,
				MaxSize:      100 * 1024 * 1024,
				MaxBackups:   7,
				EnsureFolder: true,
				LocalTime:    true,
			},
		},
	}
}
