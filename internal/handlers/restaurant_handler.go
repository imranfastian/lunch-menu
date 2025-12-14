package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"lunch_menu/internal/database"
	"lunch_menu/internal/models"
	"lunch_menu/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateRestaurant godoc
// @Summary      Create a new restaurant
// @Description  Admin only: create a new restaurant
// @Tags         restaurants
// @Accept       json
// @Produce      json
// @Param        restaurant  body      models.RestaurantInput  true  "Restaurant Input"
// @Success      201  {object}  models.RestaurantResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /restaurants [post]
// @Security     BearerAuth
func CreateRestaurant(c *gin.Context) {
	const method = "CreateRestaurant"
	var input models.RestaurantInput
	var status = http.StatusCreated
	var message = "Restaurant created successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	if err := c.ShouldBindJSON(&input); err != nil {
		status = http.StatusBadRequest
		message = "Invalid input"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else if ok, invalid := input.Validate(); !ok {
		status = http.StatusBadRequest
		message = "Invalid input fields"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_FIELDS",
			Message: fmt.Sprintf("%s: Invalid fields: %s", method, strings.Join(invalid, ", ")),
		}
	} else {
		restaurant := models.Restaurant{
			Name:        input.Name,
			Description: input.Description,
			Address:     input.Address,
			Coordinate:  input.Coordinate,
			Homepage:    input.Homepage,
			Region:      input.Region,
			Phone:       input.Phone,
			Email:       input.Email,
		}
		created, err := database.CreateRestaurant(&restaurant)
		if err != nil {
			status = http.StatusInternalServerError
			message = "Failed to create restaurant"
			errResp = &models.ErrorResponse{
				Error:   "DATABASE_ERROR",
				Message: fmt.Sprintf("%s: %v", method, err.Error()),
			}
		} else {
			data = models.RestaurantResponse{Restaurant: *created}
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// UpdateRestaurant godoc
// @Summary      Update a restaurant
// @Description  Admin only: update an existing restaurant
// @Tags         restaurants
// @Accept       json
// @Produce      json
// @Param        id          path      int                   true  "Restaurant ID"
// @Param        restaurant  body      models.RestaurantUpdateInput  true  "Restaurant Update Input"
// @Success      200  {object}  models.RestaurantResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /restaurants/{id} [put]
// @Security     BearerAuth
func UpdateRestaurant(c *gin.Context) {
	const method = "UpdateRestaurant"
	var input models.RestaurantUpdateInput
	var status = http.StatusOK
	var message = "Restaurant updated successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		status = http.StatusBadRequest
		message = "Invalid restaurant ID"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
		utils.Respond(c, status, message, data, errResp)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		status = http.StatusBadRequest
		message = "Invalid input"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
		utils.Respond(c, status, message, data, errResp)
		return
	}

	updated, err := database.PartialUpdateRestaurant(uint(id), &input)
	if err != nil {
		status = http.StatusInternalServerError
		message = "Failed to update restaurant"
		errResp = &models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else {
		data = models.RestaurantResponse{Restaurant: *updated}
	}

	utils.Respond(c, status, message, data, errResp)
}

// DeleteRestaurant godoc
// @Summary      Delete a restaurant
// @Description  Admin only: delete a restaurant
// @Tags         restaurants
// @Produce      json
// @Param        id   path      int  true  "Restaurant ID"
// @Success      200  {object}  models.StandardResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /restaurants/{id} [delete]
// @Security     BearerAuth
func DeleteRestaurant(c *gin.Context) {
	const method = "DeleteRestaurant"
	var status = http.StatusOK
	var message = "Restaurant deleted successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		status = http.StatusBadRequest
		message = "Invalid restaurant ID"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else if err := database.DeleteRestaurant(uint(id)); err != nil {
		status = http.StatusInternalServerError
		message = "Failed to delete restaurant"
		errResp = &models.ErrorResponse{
			Error:   "DELETE_FAILED",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// GetRestaurants godoc
// @Summary      List restaurants
// @Description  Get a paginated list of restaurants
// @Tags         restaurants
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {object}  models.RestaurantsResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /restaurants [get]
func GetRestaurants(c *gin.Context) {
	const method = "GetRestaurants"
	var status = http.StatusOK
	var message = "Restaurants fetched successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	restaurants, total, err := database.GetRestaurants(limit, offset)
	if err != nil {
		status = http.StatusInternalServerError
		message = "Failed to retrieve restaurants"
		errResp = &models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else {
		data = models.RestaurantsResponse{
			Restaurants: restaurants,
			Total:       total,
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// GetRestaurant godoc
// @Summary      Get a restaurant
// @Description  Get a restaurant by ID
// @Tags         restaurants
// @Produce      json
// @Param        id   path      int  true  "Restaurant ID"
// @Success      200  {object}  models.RestaurantResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /restaurants/{id} [get]
func GetRestaurant(c *gin.Context) {
	const method = "GetRestaurant"
	var status = http.StatusOK
	var message = "Restaurant fetched successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		status = http.StatusBadRequest
		message = "Invalid restaurant ID"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else {
		restaurant, err := database.GetRestaurantByID(uint(id))
		if err != nil {
			status = http.StatusNotFound
			message = "Restaurant not found"
			errResp = &models.ErrorResponse{
				Error:   "NOT_FOUND",
				Message: fmt.Sprintf("%s: %v", method, err.Error()),
			}
		} else {
			data = models.RestaurantResponse{
				Restaurant: *restaurant,
			}
		}
	}

	utils.Respond(c, status, message, data, errResp)
}
