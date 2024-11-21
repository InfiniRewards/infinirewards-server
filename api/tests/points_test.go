package tests

import (
	"bytes"
	"encoding/json"
	"infinirewards/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointsManagement(t *testing.T) {
	router := setupTest(t)
	testMerchant := createTestMerchantWithAuth(t, router)

	t.Run("Points Contract Operations", func(t *testing.T) {
		// Create points contract using merchant token
		createReq := models.CreatePointsContractRequest{
			Name:        "Test Points",
			Symbol:      "TST",
			Description: "Test points contract",
			Decimals:    "18",
		}
		reqBody, _ := json.Marshal(createReq)

		req := httptest.NewRequest("POST", "/points-contracts", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createResp models.CreatePointsContractResponse
		err := json.Unmarshal(w.Body.Bytes(), &createResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, createResp.Address)

		// Create a test recipient user
		testRecipient := createTestUserWithAuth(t, router)

		// Mint points to the recipient's address
		mintReq := models.MintPointsRequest{
			PointsContract: createResp.Address,
			Recipient:      testRecipient.User.AccountAddress, // Use actual user address
			Amount:         "100",
		}
		reqBody, _ = json.Marshal(mintReq)

		req = httptest.NewRequest("POST", "/points/mint", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testMerchant.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Create another test user for transfer target
		testTransferTarget := createTestUserWithAuth(t, router)

		// Transfer points
		transferReq := models.TransferPointsRequest{
			PointsContract: createResp.Address,
			To:             testTransferTarget.User.AccountAddress, // Use recipient's address
			Amount:         "50",
		}
		reqBody, _ = json.Marshal(transferReq)

		req = httptest.NewRequest("POST", "/points/transfer", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testRecipient.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Burn points (using recipient's remaining balance)
		burnReq := models.BurnPointsRequest{
			PointsContract: createResp.Address,
			Amount:         "25",
		}
		reqBody, _ = json.Marshal(burnReq)

		req = httptest.NewRequest("POST", "/points/burn", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testRecipient.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
