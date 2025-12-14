package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"lunch_menu/internal/database"
	"lunch_menu/internal/models"
	"lunch_menu/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateMenuItem godoc
// @Summary      Create a new menu item
// @Description  Admin only: create a new menu item
// @Tags         menu-items
// @Accept       json
// @Produce      json
// @Param        menu_item  body      models.MenuItem  true  "Menu Item Input"
// @Success      201  {object}  models.StandardResponse
// @Failure      400  {object}  models.StandardResponse
// @Failure      401  {object}  models.StandardResponse
// @Router       /menu-items [post]
// @Security     BearerAuth
func CreateMenuItem(c *gin.Context) {
	const method = "CreateMenuItem"
	var input models.MenuItem
	var status = http.StatusCreated
	var message = "Menu item created successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	if err := c.ShouldBindJSON(&input); err != nil {
		status = http.StatusBadRequest
		message = "Invalid menu item data"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else {
		created, err := database.CreateMenuItem(&input)
		if err != nil {
			status = http.StatusInternalServerError
			message = "Failed to create menu item"
			errResp = &models.ErrorResponse{
				Error:   "DATABASE_ERROR",
				Message: fmt.Sprintf("%s: %v", method, err.Error()),
			}
		} else {
			data = created
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// UpdateMenuItem godoc
// @Summary      Update a menu item
// @Description  Admin only: update an existing menu item
// @Tags         menu-items
// @Accept       json
// @Produce      json
// @Param        id         path      int           true  "Menu Item ID"
// @Param        menu_item  body      models.MenuItemUpdateInput  true  "Menu Item Input"
// @Success      200  {object}  models.StandardResponse
// @Failure      400  {object}  models.StandardResponse
// @Failure      401  {object}  models.StandardResponse
// @Failure      404  {object}  models.StandardResponse
// @Router       /menu-items/{id} [put]
// @Security     BearerAuth
func UpdateMenuItem(c *gin.Context) {
	const method = "UpdateMenuItem"
	var input models.MenuItemUpdateInput
	var status = http.StatusOK
	var message = "Menu item updated successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		status = http.StatusBadRequest
		message = "Invalid menu item ID"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
		utils.Respond(c, status, message, data, errResp)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		status = http.StatusBadRequest
		message = "Invalid menu item data"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
		utils.Respond(c, status, message, data, errResp)
		return
	}

	updated, err := database.PartialUpdateMenuItem(uint(id), &input)
	if err != nil {
		status = http.StatusInternalServerError
		message = "Failed to update menu item"
		errResp = &models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else {
		data = updated
	}

	utils.Respond(c, status, message, data, errResp)
}

// DeleteMenuItem godoc
// @Summary      Delete a menu item
// @Description  Admin only: delete a menu item
// @Tags         menu-items
// @Produce      json
// @Param        id   path      int  true  "Menu Item ID"
// @Success      200  {object}  models.StandardResponse
// @Failure      400  {object}  models.StandardResponse
// @Failure      401  {object}  models.StandardResponse
// @Failure      404  {object}  models.StandardResponse
// @Router       /menu-items/{id} [delete]
// @Security     BearerAuth
func DeleteMenuItem(c *gin.Context) {
	const method = "DeleteMenuItem"
	var status = http.StatusOK
	var message = "Menu item deleted successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		status = http.StatusBadRequest
		message = "Invalid menu item ID"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else if err := database.DeleteMenuItem(uint(id)); err != nil {
		status = http.StatusInternalServerError
		message = "Failed to delete menu item"
		errResp = &models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// GetMenuItems godoc
// @Summary      List menu items
// @Description  Get a paginated list of menu items
// @Tags         menu-items
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {object}  models.StandardResponse
// @Failure      500  {object}  models.StandardResponse
// @Router       /menu-items [get]
func GetMenuItems(c *gin.Context) {
	const method = "GetMenuItems"
	var status = http.StatusOK
	var message = "Menu items fetched successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	menuItems, total, err := database.GetMenuItems(0, limit, offset)
	if err != nil {
		status = http.StatusInternalServerError
		message = "Failed to retrieve menu items"
		errResp = &models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else {
		data = models.MenuItemsResponse{
			MenuItems: menuItems,
			Total:     total,
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// GetMenuItem godoc
// @Summary      Get a menu item
// @Description  Get a menu item by ID
// @Tags         menu-items
// @Produce      json
// @Param        id   path      int  true  "Menu Item ID"
// @Success      200  {object}  models.StandardResponse
// @Failure      400  {object}  models.StandardResponse
// @Failure      404  {object}  models.StandardResponse
// @Router       /menu-items/{id} [get]
func GetMenuItem(c *gin.Context) {
	const method = "GetMenuItem"
	var status = http.StatusOK
	var message = "Menu item fetched successfully"
	var data interface{}
	var errResp *models.ErrorResponse

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		status = http.StatusBadRequest
		message = "Invalid menu item ID"
		errResp = &models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: fmt.Sprintf("%s: %v", method, err.Error()),
		}
	} else {
		menuItem, err := database.GetMenuItemByID(uint(id))
		if err != nil {
			status = http.StatusNotFound
			message = "Menu item not found"
			errResp = &models.ErrorResponse{
				Error:   "NOT_FOUND",
				Message: fmt.Sprintf("%s: %v", method, err.Error()),
			}
		} else {
			data = menuItem
		}
	}

	utils.Respond(c, status, message, data, errResp)
}

// GetRestaurantMenu godoc
// @Summary      Get menu items for a restaurant
// @Description  Get a paginated list of menu items for a specific restaurant
// @Tags         menu-items
// @Produce      json
// @Param        id      path      int  true  "Restaurant ID"
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {object}  models.StandardResponse
// @Failure      400  {object}  models.StandardResponse
// @Failure      404  {object}  models.StandardResponse
// @Router       /restaurants/{id}/menu [get]
func GetRestaurantMenu(c *gin.Context) {
	const method = "GetRestaurantMenu"
	var status = http.StatusOK
	var message = "Menu items fetched successfully"
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
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		menuItems, total, err := database.GetMenuItems(uint(id), limit, offset)
		if err != nil {
			status = http.StatusInternalServerError
			message = "Failed to retrieve menu items"
			errResp = &models.ErrorResponse{
				Error:   "DATABASE_ERROR",
				Message: fmt.Sprintf("%s: %v", method, err.Error()),
			}
		} else {
			data = models.MenuItemsResponse{
				MenuItems: menuItems,
				Total:     total,
			}
		}
	}

	utils.Respond(c, status, message, data, errResp)
}
