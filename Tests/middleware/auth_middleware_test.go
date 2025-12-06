package middleware_test

import (
	infrastructure "example/task_manager/Infrastructure"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.Use(infrastructure.AuthMiddleware())
	r.GET("/test-auth", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test case 1: No Authorization header
	req, _ := http.NewRequest(http.MethodGet, "/test-auth", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")

	// Test case 2: Invalid token format
	req, _ = http.NewRequest(http.MethodGet, "/test-auth", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_format")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")

	// Test case 3: Valid token
	username := "testuser"
	isAdmin := false
	token, _ := infrastructure.GenerateJWT(username, isAdmin)

	req, _ = http.NewRequest(http.MethodGet, "/test-auth", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Helper to create a gin context with user claims
	createContextWithClaims := func(username string, isAdmin bool) *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("username", username)
		c.Set("is_admin", isAdmin)
		return c
	}

	r := gin.Default()
	r.Use(infrastructure.AdminMiddleware())
	r.GET("/test-admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test case 1: Admin user (This part of the test was problematic in original, simplified for chained test)
	adminUsername := "adminuser"
	adminIsAdmin := true
	adminToken, _ := infrastructure.GenerateJWT(adminUsername, adminIsAdmin)

	reqAdmin, _ := http.NewRequest(http.MethodGet, "/test-admin", nil)
	reqAdmin.Header.Set("Authorization", "Bearer "+adminToken)

	// Manually set claims in context for middleware test, as AuthMiddleware is not chained here
	cAdmin := createContextWithClaims(adminUsername, adminIsAdmin)
	_ = cAdmin // Ignore cAdmin as it's not used directly after simplified logic
	// We will rely on the chained middleware tests below for actual functionality validation
	routerWithAuthAdmin := gin.Default()
	routerWithAuthAdmin.Use(infrastructure.AuthMiddleware())
	routerWithAuthAdmin.Use(infrastructure.AdminMiddleware())
	routerWithAuthAdmin.GET("/admin-protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome, Admin!"})
	})

	// Test case 1.1: Admin user with full middleware chain
	adminUser := "superadmin"
	adminTokenFull, _ := infrastructure.GenerateJWT(adminUser, true)
	reqAdminFull, _ := http.NewRequest(http.MethodGet, "/admin-protected", nil)
	reqAdminFull.Header.Set("Authorization", "Bearer "+adminTokenFull)
	wAdminFull := httptest.NewRecorder()
	routerWithAuthAdmin.ServeHTTP(wAdminFull, reqAdminFull)
	assert.Equal(t, http.StatusOK, wAdminFull.Code)
	assert.Contains(t, wAdminFull.Body.String(), "Welcome, Admin!")

	// Test case 1.2: Non-admin user with full middleware chain
	nonAdminUser := "regularuser"
	nonAdminTokenFull, _ := infrastructure.GenerateJWT(nonAdminUser, false)
	reqNonAdminFull, _ := http.NewRequest(http.MethodGet, "/admin-protected", nil)
	reqNonAdminFull.Header.Set("Authorization", "Bearer "+nonAdminTokenFull)
	wNonAdminFull := httptest.NewRecorder()
	routerWithAuthAdmin.ServeHTTP(wNonAdminFull, reqNonAdminFull)
	assert.Equal(t, http.StatusForbidden, wNonAdminFull.Code)
	assert.Contains(t, wNonAdminFull.Body.String(), "Admin access required")

	// Test case 1.3: No token (AuthMiddleware will catch this first)
	reqNoToken, _ := http.NewRequest(http.MethodGet, "/admin-protected", nil)
	wNoToken := httptest.NewRecorder()
	routerWithAuthAdmin.ServeHTTP(wNoToken, reqNoToken)
	assert.Equal(t, http.StatusUnauthorized, wNoToken.Code)
	assert.Contains(t, wNoToken.Body.String(), "Authorization header required")

	// Test case 1.4: AdminMiddleware directly with missing 'is_admin' in context
	rAdminOnly := gin.Default()
	rAdminOnly.Use(infrastructure.AdminMiddleware())
	rAdminOnly.GET("/admin-direct", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Should not reach here"})
	})
	reqDirectAdmin, _ := http.NewRequest(http.MethodGet, "/admin-direct", nil)
	wDirectAdmin := httptest.NewRecorder()
	rAdminOnly.ServeHTTP(wDirectAdmin, reqDirectAdmin)
	assert.Equal(t, http.StatusForbidden, wDirectAdmin.Code)
	assert.Contains(t, wDirectAdmin.Body.String(), "Admin access required")
}
