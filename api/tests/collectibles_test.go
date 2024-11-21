package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"infinirewards/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectiblesManagement(t *testing.T) {
	router := setupTest(t)
	testUser := createTestUserWithAuth(t, router)
	testMerchant := createTestMerchantWithAuth(t, router)

	t.Run("Collectible Contract Operations", func(t *testing.T) {
		// Create collectible using merchant token
		createReq := models.CreateCollectibleRequest{
			Name:        "Test Collectible",
			Description: "Test collectible contract",
		}
		reqBody, _ := json.Marshal(createReq)

		req := httptest.NewRequest("POST", "/collectibles", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createResp models.CreateCollectibleResponse
		err := json.Unmarshal(w.Body.Bytes(), &createResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, createResp.Address)

		t.Logf("Created collectible with address: %s", createResp.Address)

		// Get merchant's points contracts first
		req = httptest.NewRequest("GET", "/merchant/points-contracts", nil)
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var pointsContractsResp models.GetPointsContractsResponse
		err = json.Unmarshal(w.Body.Bytes(), &pointsContractsResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, pointsContractsResp.Contracts)

		// Set token data using the first points contract
		tokenDataReq := models.SetTokenDataRequest{
			PointsContract: pointsContractsResp.Contracts[0].Address,
			Price:          "100",
			Expiry:         1735689600,
			Description:    "Test token",
		}
		reqBody, _ = json.Marshal(tokenDataReq)

		req = httptest.NewRequest("PUT", fmt.Sprintf("/collectibles/%s/token-data/1", createResp.Address), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Set token data failed: %s", w.Body.String())

		// Get token data
		req = httptest.NewRequest("GET", fmt.Sprintf("/collectibles/%s/token-data/1", createResp.Address), nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var tokenData models.GetTokenDataResponse
		err = json.Unmarshal(w.Body.Bytes(), &tokenData)
		assert.NoError(t, err)
		assert.Equal(t, tokenDataReq.Description, tokenData.Description)

		// Mint collectible to the user's address
		mintReq := models.MintCollectibleRequest{
			CollectibleAddress: createResp.Address,
			To:                 testUser.User.AccountAddress,
			TokenId:            "1",
			Amount:             "1",
		}
		reqBody, _ = json.Marshal(mintReq)

		req = httptest.NewRequest("POST", "/collectibles/mint", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Get balance
		req = httptest.NewRequest("GET", fmt.Sprintf("/collectibles/%s/balance/1", createResp.Address), nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var balanceResp models.GetCollectibleBalanceResponse
		err = json.Unmarshal(w.Body.Bytes(), &balanceResp)
		assert.NoError(t, err)
		assert.Equal(t, "1", balanceResp.Balance)

		// Get collectible URI
		req = httptest.NewRequest("GET", fmt.Sprintf("/collectibles/%s/uri/1", createResp.Address), nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var uriResp models.GetCollectibleURIResponse
		err = json.Unmarshal(w.Body.Bytes(), &uriResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, uriResp.URI)

		// Check validity
		req = httptest.NewRequest("GET", fmt.Sprintf("/collectibles/%s/valid/1", createResp.Address), nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var validResp models.IsCollectibleValidResponse
		err = json.Unmarshal(w.Body.Bytes(), &validResp)
		assert.NoError(t, err)
		assert.True(t, validResp.IsValid)

		// Redeem collectible
		redeemReq := models.RedeemCollectibleRequest{
			User:    testUser.User.AccountAddress,
			TokenId: "1",
			Amount:  "1",
		}
		reqBody, _ = json.Marshal(redeemReq)

		req = httptest.NewRequest("POST", fmt.Sprintf("/collectibles/%s/redeem", createResp.Address), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
