package infinirewards

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/utils"
)

// CreateUser creates a user account
//	@param		publicKey:		The	public	key		of	the	user
//	@param		phoneNumber:	The	phone	number	of	the	user
//	@return:	The transaction hash, the address of the user, and an error
func CreateUser(publicKey string, phoneNumber string) (string, string, error) {
	// Convert publicKey and phoneNumberHash to felt
	publicKeyFelt, err := utils.HexToFelt(publicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert public key to felt: %w", err)
	}

	phoneNumberHash := sha256.Sum256([]byte(phoneNumber))
	phoneNumberHashFelt := new(felt.Felt).SetBytes(phoneNumberHash[:])

	receipt, err := InvokeTransactionMaster(InfiniRewardsFactoryAddress, "create_user", []*felt.Felt{publicKeyFelt, phoneNumberHashFelt})
	if err != nil {
		return "", "", fmt.Errorf("failed to create user: %w", err)
	}

	deployedAccountAddress := PadZerosInFelt(receipt.TransactionReceipt.Events[0].FromAddress)

	return receipt.TransactionHash.String(), deployedAccountAddress, nil
}

// CreateMerchant creates a merchant account
//	@param		publicKey:		The	public		key		of	the		merchant
//	@param		phoneNumber:	The	phone		number	of	the		merchant
//	@param		name:			The	name		of		the	initial	points	contract
//	@param		symbol:			The	symbol		of		the	initial	points	contract
//	@param		decimals:		The	decimals	of		the	initial	points	contract
//	@return:	The transaction hash, the address of the merchant, the address of the points, and an error
func CreateMerchant(publicKey string, phoneNumber string, name string, symbol string, decimals uint64) (string, string, string, error) {
	phoneNumberHashFelt := HashPhoneNumber(phoneNumber)
	// Convert publicKey and phoneNumberHash to felt
	publicKeyFelt, err := utils.HexToFelt(publicKey)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to convert public key to felt: %w", err)
	}

	nameFelt, err := utils.StringToByteArrFelt(name)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to convert name to felt: %w", err)
	}
	symbolFelt, err := utils.StringToByteArrFelt(symbol)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to convert symbol to felt: %w", err)
	}
	decimalsFelt := utils.Uint64ToFelt(decimals)

	calldata := []*felt.Felt{publicKeyFelt, phoneNumberHashFelt}
	calldata = append(calldata, nameFelt...)
	calldata = append(calldata, symbolFelt...)
	calldata = append(calldata, decimalsFelt)

	receipt, err := InvokeTransactionMaster(InfiniRewardsFactoryAddress, "create_merchant_contract", calldata)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create merchant: %w", err)
	}

	deployedAccountAddress := PadZerosInFelt(receipt.TransactionReceipt.Events[0].FromAddress)
	deployedPointsAddress := PadZerosInFelt(receipt.TransactionReceipt.Events[1].FromAddress)

	return receipt.TransactionHash.String(), deployedAccountAddress, deployedPointsAddress, nil
}

// CreateInfiniRewardsCollectible creates a collectible contract
//	@param		account:		The	account		of	the	merchant
//	@param		name:			The	name		of	the	collectible
//	@param		description:	The	description	of	the	collectible
//	@return:	The transaction hash, the address of the collectible, and an error
func CreateInfiniRewardsCollectible(account *account.Account, name string, description string) (string, string, error) {
	calldata := []*felt.Felt{}
	nameFelt, err := utils.StringToByteArrFelt(name)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert name to felt: %w", err)
	}
	calldata = append(calldata, nameFelt...)
	descriptionFelt, err := utils.StringToByteArrFelt(description)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert description to felt: %w", err)
	}

	calldata = append(calldata, descriptionFelt...)

	receipt, err := InvokeTransaction(account, InfiniRewardsFactoryAddress, "create_collectible_contract", calldata)
	if err != nil {
		return "", "", fmt.Errorf("failed to create collectible: %w", err)
	}

	return receipt.TransactionHash.String(), PadZerosInFelt(receipt.TransactionReceipt.Events[0].FromAddress), nil
}

// func GetInfiniRewardsCollectibleAddress(ctx context.Context, account *account.Account, salt *big.Int) (string, error) {
// 	calldata := []*felt.Felt{utils.BigIntToFelt(salt)}

// 	factoryAddressFelt, err := utils.HexToFelt(InfiniRewardsFactoryAddress)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to convert factory address to felt: %w", err)
// 	}

// 	resp, err := CallContract(ctx, factoryAddressFelt, "get_infini_rewards_collectible_address", calldata)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get InfiniRewardsCollectible address: %v", err)
// 	}

// 	return PadZerosInFelt(resp[0]), nil
// }

// CreateAdditionalPointsContract creates an additional points contract
//	@param		account:		The	account		of	the	merchant
//	@param		name:			The	name		of	the	points	contract
//	@param		symbol:			The	symbol		of	the	points	contract
//	@param		description:	The	description	of	the	points	contract
//	@param		decimals:		The	decimals	of	the	points	contract
//	@return:	The transaction hash, the address of the points, and an error
func CreateAdditionalPointsContract(ctx context.Context, account *account.Account, name, symbol, description string, decimals *big.Int) (string, string, error) {
	calldata := []*felt.Felt{}
	nameFelt, err := utils.StringToByteArrFelt(name)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert name to felt: %w", err)
	}
	symbolFelt, err := utils.StringToByteArrFelt(symbol)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert symbol to felt: %w", err)
	}
	descriptionFelt, err := utils.StringToByteArrFelt(description)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert description to felt: %w", err)
	}

	calldata = append(calldata, nameFelt...)
	calldata = append(calldata, symbolFelt...)
	calldata = append(calldata, descriptionFelt...)
	calldata = append(calldata, utils.BigIntToFelt(decimals))

	receipt, err := InvokeTransaction(account, InfiniRewardsFactoryAddress, "create_points_contract", calldata)
	if err != nil {
		return "", "", fmt.Errorf("failed to create points contract: %w", err)
	}

	return receipt.TransactionHash.String(), PadZerosInFelt(receipt.TransactionReceipt.Events[0].FromAddress), nil
}

// // GetUserAccount gets a user account
// // @param phoneNumberHash: The phone number hash of the user
// // @return: The address of the user and an error
// func GetUserAccount(ctx context.Context, phoneNumberHash string) (string, error) {
// 	phoneNumberHashFelt, err := utils.HexToFelt(phoneNumberHash)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to convert phone number hash to felt: %w", err)
// 	}
// 	factoryAddressFelt, err := utils.HexToFelt(InfiniRewardsFactoryAddress)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to convert factory address to felt: %w", err)
// 	}

// 	calldata := []*felt.Felt{phoneNumberHashFelt}

// 	resp, err := CallContract(ctx, factoryAddressFelt, "get_user_account", calldata)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get user account: %v", err)
// 	}

// 	return PadZerosInFelt(resp[0]), nil
// }

// // GetMerchantAccount gets a merchant account
// // @param phoneNumberHash: The phone number hash of the merchant
// // @return: The address of the merchant and an error
// func GetMerchantAccount(ctx context.Context, phoneNumberHash string) (string, error) {
// 	phoneNumberHashFelt, err := utils.HexToFelt(phoneNumberHash)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to convert phone number hash to felt: %w", err)
// 	}
// 	factoryAddressFelt, err := utils.HexToFelt(InfiniRewardsFactoryAddress)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to convert factory address to felt: %w", err)
// 	}

// 	calldata := []*felt.Felt{phoneNumberHashFelt}

// 	resp, err := CallContract(ctx, factoryAddressFelt, "get_merchant_account", calldata)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get merchant account: %v", err)
// 	}

// 	return PadZerosInFelt(resp[0]), nil
// }
