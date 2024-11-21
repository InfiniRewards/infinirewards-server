package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerchantManagement(t *testing.T) {
	router := setupTest(t)
	testMerchant := createTestMerchantWithAuth(t, router)

	t.Run("Get Merchant Details", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/merchant", nil)
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected status code to be 200: %v", w.Body.String())
	})
}
