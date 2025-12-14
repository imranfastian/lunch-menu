// Load environment variables, app settings
package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	DBSSLmode  string
	DBTimeZone string // <-- need to test
}

type CookieConfig struct {
	Domain            string
	Path              string
	Secure            bool
	HTTPOnly          bool
	SameSite          string
	MaxAge            int
	AccessCookieName  string
	AccessCookieAge   int
	RefreshCookieName string
	RefreshCookieAge  int
	CSRFCookieName    string
	CSRFCookieAge     int
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "db_user"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "lunch_menu"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		DBSSLmode:  getEnv("DB_SSLMODE", "disable"),
		DBTimeZone: getEnv("DB_TIMEZONE", "UTC"), // <-- Add this line
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func MustLoadConfig() {
	LoadConfig()
	if AppConfig.DBPassword == "" || AppConfig.JWTSecret == "" {
		log.Fatal("Critical environment variables missing: DB_PASSWORD or JWT_SECRET")
	}
}

func GetCookieConfig() CookieConfig {
	secure, _ := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	httpOnly, _ := strconv.ParseBool(os.Getenv("COOKIE_HTTPONLY"))
	maxAge, _ := strconv.Atoi(os.Getenv("COOKIE_AGE"))
	accessAge, _ := strconv.Atoi(os.Getenv("ACCESS_COOKIE_AGE"))
	refreshAge, _ := strconv.Atoi(os.Getenv("REFRESH_COOKIE_AGE"))
	csrAge, _ := strconv.Atoi(os.Getenv("CSRF_COOKIE_AGE"))
	return CookieConfig{
		Domain:            os.Getenv("COOKIE_DOMAIN"),
		Path:              os.Getenv("COOKIE_PATH"),
		Secure:            secure,
		HTTPOnly:          httpOnly,
		SameSite:          os.Getenv("COOKIE_SAMESITE"),
		MaxAge:            maxAge,
		AccessCookieName:  os.Getenv("ACCESS_COOKIE_NAME"),
		AccessCookieAge:   accessAge,
		RefreshCookieName: os.Getenv("REFRESH_COOKIE_NAME"),
		RefreshCookieAge:  refreshAge,
		CSRFCookieName:    os.Getenv("CSRF_COOKIE_NAME"),
		CSRFCookieAge:     csrAge,
	}
}
