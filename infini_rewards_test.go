package main

import (
	"context"
	"fmt"
	"infinirewards/infinirewards"
	"infinirewards/logs"
	"math/big"
	"testing"
	"time"

	"crypto/rand"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfiniRewards(t *testing.T) {
	t.Logf("Starting InfiniRewards Test")
	logs.InitHandler("")
	ctx := context.Background()
	err := infinirewards.ConnectStarknet()
	if err != nil {
		t.Logf("Error connecting to Starknet: %v", err)
	}
	require.NoError(t, err)

	// Test Merchant and User Flow
	t.Run("MerchantAndUserFlow", func(t *testing.T) {
		// Generate random keypair for merchant
		_, merchantPubKey, merchantPrivKey := account.GetRandomKeys()
		merchantPrivKeyStr := merchantPrivKey.String()

		// 1. Create merchant
		randomInt, _ := rand.Int(rand.Reader, big.NewInt(10000000000))
		merchantPhoneNumber := fmt.Sprintf("+1%010d", randomInt)
		name := "Test Merchant"
		symbol := "TM"
		decimals := uint64(18)

		t.Logf("Creating merchant with phone number: %s", merchantPhoneNumber)

		txHash, merchantAddress, pointsAddress, err := infinirewards.CreateMerchant(merchantPubKey.String(), merchantPhoneNumber, name, symbol, decimals)
		if err != nil {
			t.Logf("Error creating merchant: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		require.NotEmpty(t, merchantAddress)
		require.NotEmpty(t, pointsAddress)

		t.Logf("Created merchant with address: %s, points address: %s, tx hash: %s", merchantAddress, pointsAddress, txHash)

		// time.Sleep(30 * time.Second)

		// Fund the merchant account
		txHash, err = infinirewards.FundAccount(merchantAddress)
		if err != nil {
			t.Logf("Error funding merchant account: %v", err)
		}
		require.NoError(t, err)

		t.Logf("Funded merchant account with address: %s, tx hash: %s", merchantAddress, txHash)
		// time.Sleep(30 * time.Second)

		// Create an account instance for the merchant
		merchantAccount, err := infinirewards.GetAccount(merchantPrivKeyStr, merchantAddress)
		if err != nil {
			t.Logf("Error getting merchant account: %v", err)
		}
		require.NoError(t, err)

		// 2. Create collectible
		collectibleName := "Test Collectible"
		collectibleDescription := "Test description"

		txHash, collectibleAddress, err := infinirewards.CreateInfiniRewardsCollectible(merchantAccount, collectibleName, collectibleDescription)
		if err != nil {
			t.Logf("Error creating collectible: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		require.NotEmpty(t, collectibleAddress)

		t.Logf("Created collectible with address: %s, tx hash: %s", collectibleAddress, txHash)

		// 3. Get/Set Collectible Data
		tokenId := big.NewInt(1)
		price := big.NewInt(100)
		expiry := uint64(time.Now().Add(365 * 24 * time.Hour).Unix())
		txHash, err = infinirewards.SetTokenData(merchantAccount, collectibleAddress, tokenId, pointsAddress, price, expiry, "Updated description")
		if err != nil {
			t.Logf("Error setting token data: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)

		t.Logf("Set token data with address: %s, tx hash: %s", collectibleAddress, txHash)

		pointsContract, tokenPrice, tokenExpiry, tokenDescription, err := infinirewards.GetTokenData(ctx, collectibleAddress, tokenId)
		if err != nil {
			t.Logf("Error getting token data: %v", err)
		}
		require.NoError(t, err)
		assert.Equal(t, pointsAddress, pointsContract)
		assert.Equal(t, price, tokenPrice)
		assert.Equal(t, expiry, tokenExpiry)
		assert.Equal(t, "Updated description", tokenDescription)

		t.Logf("Got token data: points contract: %s, price: %s, expiry: %d, description: %s", pointsContract, tokenPrice, tokenExpiry, tokenDescription)

		// Create two user accounts
		createAndFundUser := func(phoneNumber string) (*account.Account, string) {
			_, userPubKey, userPrivKey := account.GetRandomKeys()
			userPrivKeyStr := userPrivKey.String()

			txHash, userAddress, err := infinirewards.CreateUser(userPubKey.String(), phoneNumber)
			if err != nil {
				t.Logf("Error creating user: %v", err)
			}
			require.NoError(t, err)
			require.NotEmpty(t, txHash)
			require.NotEmpty(t, userAddress)

			t.Logf("Created user with address: %s, tx hash: %s", userAddress, txHash)

			txHash, err = infinirewards.FundAccount(userAddress)
			if err != nil {
				t.Logf("Error funding user account: %v", err)
			}
			require.NoError(t, err)

			t.Logf("Funded user account with address: %s, tx hash: %s", userAddress, txHash)

			userAccount, err := infinirewards.GetAccount(userPrivKeyStr, userAddress)
			if err != nil {
				t.Logf("Error getting user account: %v", err)
			}
			require.NoError(t, err)

			return userAccount, userAddress
		}

		randomInt, _ = rand.Int(rand.Reader, big.NewInt(10000000000))
		user1PhoneNumber := fmt.Sprintf("+1%010d", randomInt)
		user1Account, user1Address := createAndFundUser(user1PhoneNumber)
		randomInt, _ = rand.Int(rand.Reader, big.NewInt(10000000000))
		user2PhoneNumber := fmt.Sprintf("+1%010d", randomInt)
		user2Account, user2Address := createAndFundUser(user2PhoneNumber)

		// 4. Mint Collectible to users
		amount := big.NewInt(5)
		txHash, err = infinirewards.MintCollectible(merchantAccount, collectibleAddress, user1Address, tokenId, amount)
		if err != nil {
			t.Logf("Error minting collectible to user1: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		t.Logf("Minted collectible to user1 with address: %s, tx hash: %s", user1Address, txHash)

		txHash, err = infinirewards.MintCollectible(merchantAccount, collectibleAddress, user2Address, tokenId, amount)
		if err != nil {
			t.Logf("Error minting collectible to user2: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		t.Logf("Minted collectible to user2 with address: %s, tx hash: %s", user2Address, txHash)

		// 5. Redeem Collectible for user1
		txHash, err = infinirewards.Redeem(ctx, user1Account, collectibleAddress, user1Address, tokenId, big.NewInt(1))
		if err != nil {
			t.Logf("Error redeeming collectible for user1: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		t.Logf("Redeemed collectible for user1 with address: %s, tx hash: %s", user1Address, txHash)

		// 6. GetBalance for Collectible
		balance1, err := infinirewards.BalanceOf(ctx, user1Account, collectibleAddress, tokenId)
		if err != nil {
			t.Logf("Error getting balance for user1: %v", err)
		}
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(4), balance1) // 5 minted - 1 redeemed
		t.Logf("Balance for user1: %s", balance1.String())

		balance2, err := infinirewards.BalanceOf(ctx, user2Account, collectibleAddress, tokenId)
		if err != nil {
			t.Logf("Error getting balance for user2: %v", err)
		}
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(5), balance2) // 5 minted
		t.Logf("Balance for user2: %s", balance2.String())
		// 7. Mint points to users
		pointsAmount := big.NewInt(1000)
		txHash, err = infinirewards.MintPoints(merchantAccount, pointsAddress, user1Address, pointsAmount)
		if err != nil {
			t.Logf("Error minting points to user1: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		t.Logf("Minted points to user1 with address: %s, tx hash: %s", user1Address, txHash)
		txHash, err = infinirewards.MintPoints(merchantAccount, pointsAddress, user2Address, pointsAmount)
		if err != nil {
			t.Logf("Error minting points to user2: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		t.Logf("Minted points to user2 with address: %s, tx hash: %s", user2Address, txHash)

		// 8. Burn points for user2
		burnAmount := big.NewInt(500)
		txHash, err = infinirewards.BurnPoints(user2Account, pointsAddress, burnAmount)
		if err != nil {
			t.Logf("Error burning points for user2: %v", err)
		}
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
		t.Logf("Burned points for user2 with address: %s, tx hash: %s", user2Address, txHash)

		// 9. Get balance for Points
		pointsBalance1, err := infinirewards.GetBalance(ctx, user1Account, pointsAddress)
		if err != nil {
			t.Logf("Error getting points balance for user1: %v", err)
		}
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(1000), pointsBalance1)
		t.Logf("Points balance for user1: %s", pointsBalance1.String())

		pointsBalance2, err := infinirewards.GetBalance(ctx, user2Account, pointsAddress)
		if err != nil {
			t.Logf("Error getting points balance for user2: %v", err)
		}
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(500), pointsBalance2)
		t.Logf("Points balance for user2: %s", pointsBalance2.String())
		// Get all balances for users
		checkUserBalances := func(userAccount *account.Account, t *testing.T) {
			collectibleContracts, err := infinirewards.GetCollectibleContracts(ctx, merchantAccount)
			if err != nil {
				t.Logf("Error getting collectible contracts: %v", err)
			}
			require.NoError(t, err)
			pointsContracts, err := infinirewards.GetPointsContracts(ctx, merchantAccount)
			if err != nil {
				t.Logf("Error getting points contracts: %v", err)
			}
			require.NoError(t, err)

			for _, contract := range collectibleContracts {
				t.Logf("Getting details for collectible contract: %s", contract)
				description, _, tokenIDs, _, _, _, err := infinirewards.GetDetails(ctx, contract)
				if err != nil {
					t.Logf("Error getting details for collectible contract %s: %v", contract, err)
				}
				require.NoError(t, err)
				t.Logf("Collectible Contract: %s, Description: %s", contract, description)
				for _, tokenID := range tokenIDs {
					balance, err := infinirewards.BalanceOf(ctx, userAccount, contract, tokenID)
					if err != nil {
						t.Logf("Error getting balance for token ID %s: %v", tokenID.String(), err)
					}
					require.NoError(t, err)
					t.Logf("Token ID: %s, Balance: %s", tokenID.String(), balance.String())
				}
			}

			for _, contract := range pointsContracts {
				balance, err := infinirewards.GetBalance(ctx, userAccount, contract)
				if err != nil {
					t.Logf("Error getting balance for points contract %s: %v", contract, err)
				}
				require.NoError(t, err)
				t.Logf("Points Contract: %s, Balance: %s", contract, balance.String())
			}
		}

		t.Log("User 1 Balances:")
		checkUserBalances(user1Account, t)

		t.Log("User 2 Balances:")
		checkUserBalances(user2Account, t)
	})
}
