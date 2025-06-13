package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/models/migrated_models"
	"nmteasy_backend/utils"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []migrated_models.Product

	if err := models.DB.Where("is_active = ?", true).Find(&products).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

func ActivateSkin(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "Authentication required")
		return
	}

	var request struct {
		InventoryID string `json:"inventory_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	inventoryID, err := uuid.Parse(request.InventoryID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	// Start transaction
	tx := models.DB.Begin()

	// First, deactivate all other skins for this user
	if err := tx.Model(&migrated_models.UserInventory{}).
		Where("user_id = ?", user.ID).
		Update("is_active", false).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to deactivate other skins")
		return
	}

	// Then activate the selected one
	if err := tx.Model(&migrated_models.UserInventory{}).
		Where("id = ? AND user_id = ?", inventoryID, user.ID).
		Update("is_active", true).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to activate skin")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Transaction failed")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Skin activated"})
}

func PurchaseProduct(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "Authentication required")
		return
	}

	var request struct {
		ProductID string `json:"product_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	productID, err := uuid.Parse(request.ProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Get product with lock to prevent concurrent purchases
	var product migrated_models.Product
	if err := models.DB.First(&product, "id = ? AND is_active = ?", productID, true).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not available")
		return
	}

	// Check user balance
	if user.Points < product.Price {
		utils.RespondWithError(w, http.StatusBadRequest, "Insufficient points")
		return
	}

	var existingInventory migrated_models.UserInventory
	if err := models.DB.Where("user_id = ? AND product_id = ?", user.ID, product.ID).First(&existingInventory).Error; err == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "You already own this item")
		return
	}

	// Deduct points
	if err := models.DB.Model(&user).Update("points", user.Points-product.Price).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to deduct points")
		return
	}

	// Create purchase record
	purchase := migrated_models.Purchase{
		ID:        uuid.New(),
		UserID:    user.ID,
		ProductID: product.ID,
		PricePaid: product.Price,
	}

	if err := models.DB.Create(&purchase).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to record purchase")
		return
	}

	// Add to inventory
	inventory := migrated_models.UserInventory{
		ID:        uuid.New(),
		UserID:    user.ID,
		ProductID: product.ID,
		IsActive:  false, // Not active by default
	}

	if err := models.DB.Create(&inventory).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add to inventory")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Purchase successful",
		"balance": user.Points - product.Price,
	})
}

func GetUserInventory(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "Authentication required")
		return
	}

	var inventory []migrated_models.UserInventory
	if err := models.DB.Preload("Product").Where("user_id = ?", user.ID).Find(&inventory).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch inventory")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, inventory)
}
