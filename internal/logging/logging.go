package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var (
	Logger zerolog.Logger
)

func init() {
	logConfig := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Kitchen,
	}

	logConfig.FormatLevel = func(i interface{}) string {
		if i == "info" {
			return "\x1b[38;5;120mINF\x1b[0m \x1b[38;5;239m>\x1b[0m"
		} else if i == "debug" {
			return "\x1b[38;5;221mDBG\x1b[0m \x1b[38;5;239m>\x1b[0m"
		} else if i == "warn" {
			return "\x1b[38;5;203mWRN\x1b[0m \x1b[38;5;239m>\x1b[0m"
		} else if i == "error" {
			return "\033[1m\x1b[38;5;203mWRN\x1b[0m\033[0m \x1b[38;5;239m>\x1b[0m"
		} else if i == "fatal" {
			return "\033[1m\x1b[38;5;209mFTL\x1b[0m\033[0m \x1b[38;5;239m>\x1b[0m"
		} else {
			return i.(string)
		}
	}

	logConfig.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\u001B[38;5;239m%v=\u001B[0m", i)
	}

	logConfig.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("\u001B[38;5;239m%v=\u001B[0m", i)
	}

	logConfig.FormatErrFieldValue = func(i interface{}) string {
		return i.(string)
	}

	log.Logger = log.Output(logConfig)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	Logger = log.Logger
}
