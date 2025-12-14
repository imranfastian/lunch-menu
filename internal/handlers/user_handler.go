package handlers

import (
	"net/http"
	"strings"
	"time"

	"lunch_menu/internal/database"
	"lunch_menu/internal/models"
	"lunch_menu/internal/utils"

	"github.com/gin-gonic/gin"
)

// Helper to generate and set CSRF token
func setCSRFToken(c *gin.Context) string {
	csrfToken := utils.GenerateRandomToken()
	c.SetCookie("csrf_token", csrfToken, 3600, "/", "localhost", false, true)
	return csrfToken
}

// UserRegister godoc
// @Summary      Register a new admin user
// @Description  Register a new admin user (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserInput  true  "User Input"
// @Success      201  {object}  models.SafeUser
// @Failure      400  {object}  models.ErrorResponse
// @Router       /user/register [post]
func UserRegister(c *gin.Context) {
	const method = "AdminRegister"
	var input models.UserInput
	var status = http.StatusCreated
	var message = "Admin registered successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	if err := c.ShouldBindJSON(&input); err != nil {
		status = http.StatusBadRequest
		message = "Invalid User Input data"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: method + ": " + err.Error(),
		}
	} else if ok, invalid := input.Validate(); !ok {
		status = http.StatusBadRequest
		message = "Invalid input data"
		errResp = &models.ErrorResponse{
			Error:   "VALIDATION_FAILED",
			Message: method + ": " + strings.Join(invalid, ", "),
		}
	} else {
		// Always set role to admin for this endpoint
		user := &models.User{
			Username:     input.Username,
			PasswordHash: input.PasswordHash,
			Email:        input.Email,
			Role:         input.Role,
			IsActive:     true,
		}
		// Hash password before saving
		hashed, err := utils.HashPassword(user.PasswordHash)
		if err != nil {
			status = http.StatusInternalServerError
			message = "Failed to hash password"
			errResp = &models.ErrorResponse{
				Error:   "HASH_ERROR",
				Message: method + ": " + err.Error(),
			}
		} else {
			user.PasswordHash = hashed
			user, err := database.CreateUser(user)
			if err != nil {
				status = http.StatusBadRequest
				message = "Registration failed"
				errResp = &models.ErrorResponse{
					Error:   "REGISTRATION_FAILED",
					Message: method + ": " + err.Error(),
				}
			} else {
				csrfToken := setCSRFToken(c)
				data = gin.H{
					"user":       user.ToSafeUser(),
					"csrf_token": csrfToken,
				}
			}
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// UserLogin godoc
// @Summary      Login as admin
// @Description  Login as admin user and receive JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserLoginInput  true  "Login Input"
// @Success      200  {object}  models.SafeUser
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /user/login [post]
func UserLogin(c *gin.Context) {
	const method = "UserLogin"
	var input models.UserLoginInput
	var status = http.StatusOK
	var message = "Login successful"
	var data interface{}
	var errResp *models.ErrorResponse

	if err := c.ShouldBindJSON(&input); err != nil {
		status = http.StatusBadRequest
		message = "Invalid login data"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: method + ": " + err.Error(),
		}
	} else {
		user, err := database.GetUserByUsername(input.Username)
		if err != nil || !user.IsActive || user.Role != "admin" {
			status = http.StatusUnauthorized
			message = "Admin not found or not active or not authorized"
			errResp = &models.ErrorResponse{
				Error:   "AUTH_FAILED",
				Message: method + ": Admin not found or not active or not authorized",
			}
		} else if err := utils.CheckPasswordHash(input.PasswordHash, user.PasswordHash); err != nil {
			status = http.StatusUnauthorized
			message = "Invalid credentials"
			errResp = &models.ErrorResponse{
				Error:   "AUTH_FAILED",
				Message: method + ": Invalid credentials",
			}
		} else {
			// Generate access and refresh tokens
			accessToken, err := utils.GenerateJWT(user)
			if err != nil {
				status = http.StatusInternalServerError
				message = "Failed to generate access token"
				errResp = &models.ErrorResponse{
					Error:   "TOKEN_ERROR",
					Message: method + ": " + err.Error(),
				}
			} else {
				refreshToken := utils.GenerateRandomToken()
				refreshTokenAge := 7 * 24 * 3600 // 7 days
				expiresAt := time.Now().Add(time.Duration(refreshTokenAge) * time.Second)
				err = database.SaveRefreshToken(user.ID, refreshToken, c.Request.UserAgent(), c.ClientIP(), expiresAt)
				if err != nil {
					status = http.StatusInternalServerError
					message = "Failed to save refresh token"
					errResp = &models.ErrorResponse{
						Error:   "DB_ERROR",
						Message: method + ": " + err.Error(),
					}
				} else {
					csrfToken := setCSRFToken(c)
					utils.SetAuthCookies(c, accessToken, refreshToken, csrfToken)
					safeUser := user.ToSafeUser()
					data = gin.H{
						"user":  safeUser,
						"token": accessToken,
					}
				}
			}
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// UserLogout godoc
// @Summary      Logout user
// @Description  Logout the current user and expire tokens
// @Tags         users
// @Produce      json
// @Success      200  {object}  models.StandardResponse
// @Router       /user/logout [post]
func UserLogout(c *gin.Context) {
	// Get the access token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(token)
		if err == nil && claims["exp"] != nil {
			exp := int64(claims["exp"].(float64))
			expiresAt := time.Unix(exp, 0)
			_ = database.BlacklistToken(token, expiresAt)
		}
		// Also revoke refresh token
		refreshToken, err := c.Cookie("refresh_token")
		if err == nil && claims["user_id"] != nil {
			userID := uint(claims["user_id"].(float64))
			_ = database.DeleteRefreshToken(userID, refreshToken)
		}
	}
	utils.ExpireAuthCookies(c)
	utils.Respond(c, http.StatusOK, "Logout successful", nil, nil)
}
