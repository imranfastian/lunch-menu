package utils

import (
	"lunch_menu/internal/config"
	"lunch_menu/internal/models"
	"os"

	"github.com/gin-gonic/gin"
)

// Respond sends a standardized JSON response.
func Respond(c *gin.Context, status int, message string, data interface{}, err *models.ErrorResponse) {
	c.JSON(status, models.StandardResponse{
		Message: message,
		Data:    data,
		Error:   err,
	})
}

// SetAuthCookies sets access and refresh token cookies using config from .env
func SetAuthCookies(c *gin.Context, accessToken, refreshToken, csrfToken string) {
	cookieCfg := config.GetCookieConfig()

	accessAge := cookieCfg.AccessCookieAge
	if accessAge == 0 {
		accessAge = cookieCfg.MaxAge
	}
	refreshAge := cookieCfg.RefreshCookieAge
	if refreshAge == 0 {
		refreshAge = cookieCfg.MaxAge
	}
	csrfAge := cookieCfg.CSRFCookieAge
	if csrfAge == 0 {
		csrfAge = cookieCfg.MaxAge
	}

	// Set Access Token Cookie
	c.SetCookie(
		cookieCfg.AccessCookieName,
		accessToken,
		accessAge,
		cookieCfg.Path,
		cookieCfg.Domain,
		cookieCfg.Secure,
		cookieCfg.HTTPOnly,
	)

	// Set Refresh Token Cookie
	c.SetCookie(
		cookieCfg.RefreshCookieName,
		refreshToken,
		refreshAge,
		cookieCfg.Path,
		cookieCfg.Domain,
		cookieCfg.Secure,
		cookieCfg.HTTPOnly,
	)

	// Set CSRF Token Cookie
	c.SetCookie(
		cookieCfg.CSRFCookieName,
		csrfToken,
		csrfAge,
		cookieCfg.Path,
		cookieCfg.Domain,
		cookieCfg.Secure,
		false, // CSRF cookie should be readable by JS (not HttpOnly)
	)
	// c.Header("X-CSRF-Token", csrfToken) // X-CSRF-Token	JS clients (fetch/axios, etc.)	JavaScript frontend
}

// ExpireAuthCookies expires access and refresh token cookies.
func ExpireAuthCookies(c *gin.Context) {
	cookieCfg := config.GetCookieConfig()
	c.SetCookie(os.Getenv("ACCESS_COOKIE_NAME"), "", -1, cookieCfg.Path, cookieCfg.Domain, cookieCfg.Secure, cookieCfg.HTTPOnly)
	c.SetCookie(os.Getenv("REFRESH_COOKIE_NAME"), "", -1, cookieCfg.Path, cookieCfg.Domain, cookieCfg.Secure, cookieCfg.HTTPOnly)
	c.SetCookie(os.Getenv("CSRF_COOKIE_NAME"), "", -1, cookieCfg.Path, cookieCfg.Domain, cookieCfg.Secure, false) // Not HttpOnly
}
