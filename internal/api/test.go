package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"{{cookiecutter.app_name}}/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPet(c *gin.Context) {
	var pet models.Pet

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "wrong id")
	}

	err = models.DB.First(&pet, id).Error
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("failed to get pet with id \"%d\": %s", id, err))
	}

	c.JSON(http.StatusOK, pet)
}

func CreatePet(c *gin.Context) {
	var pet models.Pet

	err := json.NewDecoder(c.Request.Body).Decode(&pet)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("wrong json data: %s", err))
	}

	err = models.DB.Create(&pet).Error
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to create pet: %s", err))
	}

	c.String(http.StatusCreated, strconv.FormatUint(uint64(pet.ID), 10))
}

func UpdatePet(c *gin.Context) {
	var pet models.Pet
	jsonMap := make(map[string]interface{})

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "wrong id")
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pet)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("wrong json data: %s", err))
	}

	err = models.DB.Update(&pet).Error
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to update pet: %s", err))
	}

	c.String(http.StatusOK, fmt.Sprintf("pet with id \"%d\" updated successfully", id))
}

func DeletePet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "wrong id")
	}

	err = models.DB.Delete(&models.Pet{}, id).Error
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to delete pet with id \"%d\": %s", id, err))
	}

	c.String(http.StatusOK, fmt.Sprintf("pet with id \"%d\" deleted successfully", id))
}

func GetPets(c *gin.Context) {
	var pets []models.Pet

	err := models.DB.Find(&pets).Error
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to get all pets: %s", err))
	}

	c.JSONP(http.StatusOK, pets)
}
