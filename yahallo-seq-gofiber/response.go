package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sirupsen/logrus"
)

type ResponseData struct {
	Data      interface{}        `json:"data"`
	Status    StatusResponseData `json:"status"`
	TimeStamp time.Time          `json:"time_stamp"`
}

type StatusResponseData struct {
	Code    int    `json:"status_code"`
	Message string `json:"message"`
}

func GenerateTimeJakarta() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Errorf("fail to load location time| %s", err.Error())
	}

	return time.Now().In(loc)
}

func Response_Log(ctx *fiber.Ctx, log *logrus.Logger, code int, message string, data interface{}) error {
	RespData := ResponseData{
		Data: data,
		Status: StatusResponseData{
			Code:    code,
			Message: message,
		},
		TimeStamp: GenerateTimeJakarta(),
	}

	CreateLog(ctx, log, code, message, RespData)
	return ctx.Status(code).JSON(RespData)
}
