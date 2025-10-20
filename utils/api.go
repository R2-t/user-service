package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/medama-io/go-useragent"
	"santiagotorres.me/user-service/models"
)

// Extract device information from request
func ExtractDeviceInfo(c *gin.Context) models.DeviceInfo {
	uaParser := useragent.NewParser()

	userAgent := uaParser.Parse(c.GetHeader("User-Agent"))

	return models.DeviceInfo{
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		DeviceType: userAgent.Device().String(),
		Browser:    userAgent.Browser().String(),
		OS:         userAgent.OS().String(),
	}
}
