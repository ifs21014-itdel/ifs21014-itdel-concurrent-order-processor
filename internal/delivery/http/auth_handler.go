package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	uc "github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
)

type AuthHandler struct {
	uc *uc.AuthUsecase
}

func NewAuthHandler(rg *gin.RouterGroup, uc *uc.AuthUsecase) {
	h := &AuthHandler{uc: uc}
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
	rg.POST("/totp/setup/:id", h.SetupTOTP)   // id: user id (for demo)
	rg.POST("/totp/verify/:id", h.VerifyTOTP) // verify and enable
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.uc.Register(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": gin.H{"id": user.ID, "email": user.Email}})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TOTP     string `json:"totp"` // optional: required if user enabled
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, user, err := h.uc.Login(req.Email, req.Password, req.TOTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email, "totp_enabled": user.TOTPEnabled}})
}

func (h *AuthHandler) SetupTOTP(c *gin.Context) {
	idStr := c.Param("id")
	id64, _ := strconv.ParseUint(idStr, 10, 32)
	issuer := os.Getenv("APP_NAME")
	secret, uri, err := h.uc.GenerateTOTPForUser(uint(id64), issuer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// return secret and otpauth uri (frontend will generate QR code)
	c.JSON(http.StatusOK, gin.H{"secret": secret, "otpauth_uri": uri})
}

type verifyReq struct {
	Code string `json:"code" binding:"required"`
}

func (h *AuthHandler) VerifyTOTP(c *gin.Context) {
	idStr := c.Param("id")
	id64, _ := strconv.ParseUint(idStr, 10, 32)
	var req verifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := h.uc.VerifyAndEnableTOTP(uint(id64), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"enabled": true})
}
