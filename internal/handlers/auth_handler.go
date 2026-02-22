package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/domain/ports/service"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/sidji-omnichannel/internal/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService service.AuthService
	cfg         *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// Register handles user registration
// @Summary      Register a new user
// @Description  Register a new user and organization
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.RegisterInput true "Register Input"
// @Success      201  {object}  models.AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(&input)
	if err != nil {
		if err == services.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		log.Printf("[ERROR] Register failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user login
// @Summary      Login user
// @Description  Login with email and password to get access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.LoginInput true "Login Input"
// @Success      200  {object}  models.AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&input)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		log.Printf("[ERROR] Login failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Me returns the current user
// @Summary      Get current user
// @Description  Get the currently logged in user profile
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserPublic
// @Failure      401  {object}  map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDValue.(uuid.UUID)
	authResponse, err := h.authService.GetMe(userID)
	if err != nil {
		log.Printf("[ERROR] GetMe failed for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

func (h *AuthHandler) getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  h.cfg.Google.RedirectURI,
		ClientID:     h.cfg.Google.ClientID,
		ClientSecret: h.cfg.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleLogin redirects to Google consent page
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state := uuid.New().String()
	// Optionally we could store state in cookie or session
	
	url := h.getGoogleOauthConfig().AuthCodeURL(state)
	// Return the URL for frontend to redirect
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	var request struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code in request body"})
		return
	}

	token, err := h.getGoogleOauthConfig().Exchange(context.Background(), request.Code)
	if err != nil {
		log.Printf("[ERROR] Google token exchange failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange authorization code"})
		return
	}

	client := h.getGoogleOauthConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("[ERROR] Failed to get user info from Google: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo models.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("[ERROR] Failed to parse google user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	authResponse, err := h.authService.GoogleLogin(&userInfo)
	if err != nil {
		log.Printf("[ERROR] Google login failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}
