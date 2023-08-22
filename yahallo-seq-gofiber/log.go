package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nullseed/logruseq"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

const (
	projectName           = "yahallo seq and gofiber"
	FullTimeFormat string = "2006-01-02 15:04:05"
)

var (
	logger     *logrus.Logger
	loggerInit sync.Once
)

func LogrusGetLevel(logLevel string) logrus.Level {
	switch strings.ToLower(logLevel) {
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	}
	return logrus.InfoLevel
}

// GetLogger returns a singleton instance of the logger.
func NewLogger() *logrus.Logger {
	loggerInit.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&easy.Formatter{
			TimestampFormat: FullTimeFormat,
			LogFormat:       fmt.Sprintf("%s\n", `[%lvl%]: "%time%" %msg%`),
		})
		logger.SetLevel(LogrusGetLevel("debug"))
		logger.AddHook(logruseq.NewSeqHook("http://localhost:5341"))
	})

	return logger
}

// custome logging // tambahin variable log
func CreateLog(c *fiber.Ctx, log *logrus.Logger, code int, message string, respData ResponseData) {
	reqBody := string(c.Request().Body())
	fmt.Println("reqBody: ", reqBody)
	path := c.Request().URI().String()
	fmt.Println("path: ", path)
	fmt.Println("method: ", c.Method())
	fmt.Println("ini code: ", code)
	if code == fiber.StatusOK || code == fiber.StatusAccepted || code == fiber.StatusCreated {
		log.WithFields(logrus.Fields{
			"At":            time.Now(),
			"Method":        c.Method(),
			"Path":          path,
			"ParamRequest":  reqBody,
			"ParamResponse": respData,
			"Message":       message,
			"Duration":      time.Since(time.Now()),
			"Project":       projectName,
		}).Info(c.Method() + " " + path)
	} else if code == fiber.StatusBadRequest || code == fiber.StatusConflict || code == fiber.StatusUnauthorized || code == fiber.StatusNotFound {
		log.WithFields(logrus.Fields{
			"At":            time.Now(),
			"Method":        c.Method(),
			"Path":          path,
			"ParamRequest":  reqBody,
			"ParamResponse": respData,
			"Message":       message,
			"Duration":      time.Since(time.Now()),
			"Project":       projectName,
		}).Warn(c.Method() + " " + path)
	} else {
		log.WithFields(logrus.Fields{
			"At":            time.Now(),
			"Method":        c.Method(),
			"Path":          path,
			"ParamRequest":  reqBody,
			"ParamResponse": respData,
			"Message":       message,
			"Duration":      time.Since(time.Now()),
			"Project":       projectName,
		}).Error(c.Method() + " " + path)
	}
}
