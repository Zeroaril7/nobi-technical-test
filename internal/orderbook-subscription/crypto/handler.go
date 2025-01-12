package crypto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CryptoHandler interface {
	Add(c *gin.Context)
	FindAll(c *gin.Context)
	Delete(c *gin.Context)
}

type cryptoHandler struct {
	service CryptoService
}

func NewHandler(service CryptoService) CryptoHandler {
	return &cryptoHandler{
		service: service,
	}
}

func (h *cryptoHandler) Add(c *gin.Context) {
	request := new(CreateCrypto)
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse data", "details": err.Error()})
		return
	}

	payload := new(CryptoEntity)
	payload = request.ToCryptoEntity(payload)

	if err := h.service.Add(payload, c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully added crypto data"})
}

func (h *cryptoHandler) FindAll(c *gin.Context) {
	data, total, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved all crypto data",
		"data":    data,
		"total":   total,
	})
}

func (h *cryptoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	if err := h.service.Delete(id, c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted crypto data"})
}
