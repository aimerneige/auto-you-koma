package handler

import (
	"net/http"

	"github.com/aimerneige/auto-you-koma/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type ReqRegister struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req ReqRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "email": user.Email})
}

type ReqLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Passcode string `json:"passcode"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req ReqLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var token string
	var err error

	if req.Passcode != "" {
		token, _, err = h.svc.LoginWith2FA(c.Request.Context(), req.Email, req.Password, req.Passcode)
	} else {
		token, _, err = h.svc.Login(c.Request.Context(), req.Email, req.Password)
	}

	if err != nil {
		if err.Error() == "2fa_required" {
			c.JSON(http.StatusForbidden, gin.H{"error": "2fa_required"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Setup2FA(c *gin.Context) {
	userID := c.GetString("user_id")
	url, _, err := h.svc.Setup2FA(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to setup 2fa"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

type ReqVerify2FA struct {
	Passcode string `json:"passcode" binding:"required"`
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	userID := c.GetString("user_id")
	var req ReqVerify2FA
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.VerifyAndEnable2FA(c.Request.Context(), userID, req.Passcode); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid passcode"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	email := c.GetString("email")
	c.JSON(http.StatusOK, gin.H{"user_id": userID, "email": email})
}
