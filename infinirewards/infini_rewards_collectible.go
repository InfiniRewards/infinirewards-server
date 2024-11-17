package infinirewards

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/utils"
)

func MintCollectible(account *account.Account, collectibleAddress string, to string, tokenId *big.Int, amount *big.Int) (string, error) {
	toFelt, err := utils.HexToFelt(to)
	if err != nil {
		return "", fmt.Errorf("failed to convert 'to' address to felt: %w", err)
	}
	amountFelt := BigInt256ToFelt(amount)
	tokenIdFelt := BigInt256ToFelt(tokenId)
	calldata := []*felt.Felt{toFelt, tokenIdFelt[0], tokenIdFelt[1], amountFelt[0], amountFelt[1], &felt.Zero}

	resp, err := InvokeTransaction(account, collectibleAddress, "mint", calldata)
	if err != nil {
		return "", fmt.Errorf("failed to mint collectible: %w", err)
	}

	return resp.TransactionHash.String(), nil
}

func BalanceOf(ctx context.Context, account *account.Account, collectibleAddress string, tokenId *big.Int) (*big.Int, error) {
	collectibleAddressFelt, err := utils.HexToFelt(collectibleAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert collectible address to felt: %w", err)
	}
	tokenIdFelt := BigInt256ToFelt(tokenId)
	calldata := []*felt.Felt{account.AccountAddress, tokenIdFelt[0], tokenIdFelt[1]}
	resp, err := CallContract(ctx, collectibleAddressFelt, "balanceOf", calldata)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	balance := FeltArrToBigInt256([2]*felt.Felt{resp[0], resp[1]})
	return balance, nil
}

func URI(ctx context.Context, collectibleAddress string, tokenId *big.Int) (string, error) {
	collectibleAddressFelt, err := utils.HexToFelt(collectibleAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert collectible address to felt: %w", err)
	}
	tokenIdFelt := BigInt256ToFelt(tokenId)
	calldata := []*felt.Felt{tokenIdFelt[0], tokenIdFelt[1]}

	resp, err := CallContract(ctx, collectibleAddressFelt, "uri", calldata)
	if err != nil {
		return "", fmt.Errorf("failed to get token URI: %w", err)
	}

	uri, err := utils.ByteArrFeltToString(resp)
	if err != nil {
		return "", fmt.Errorf("failed to convert felt array to string: %w", err)
	}

	return uri, nil
}

func SetTokenData(account *account.Account, collectibleAddress string, tokenId *big.Int, pointsContract string, price *big.Int, expiry uint64, description string) (string, error) {

	tokenIdFelt := BigInt256ToFelt(tokenId)
	pointsContractFelt, err := utils.HexToFelt(pointsContract)
	if err != nil {
		return "", fmt.Errorf("failed to convert points contract address to felt: %w", err)
	}
	priceFelt := BigInt256ToFelt(price)
	expiryFelt := utils.Uint64ToFelt(expiry)

	descriptionFelt, err := utils.StringToByteArrFelt(description)
	if err != nil {
		return "", fmt.Errorf("failed to convert description to felt: %w", err)
	}
	calldata := append([]*felt.Felt{tokenIdFelt[0], tokenIdFelt[1], pointsContractFelt, priceFelt[0], priceFelt[1], expiryFelt}, descriptionFelt...)

	resp, err := InvokeTransaction(account, collectibleAddress, "set_token_data", calldata)
	if err != nil {
		return "", fmt.Errorf("failed to set token data: %w", err)
	}

	return resp.TransactionHash.String(), nil
}

func GetTokenData(ctx context.Context, collectibleAddress string, tokenId *big.Int) (string, *big.Int, uint64, string, error) {
	collectibleAddressFelt, err := utils.HexToFelt(collectibleAddress)
	if err != nil {
		return "", nil, 0, "", fmt.Errorf("failed to convert collectible address to felt: %w", err)
	}
	tokenIdFelt := BigInt256ToFelt(tokenId)
	calldata := []*felt.Felt{tokenIdFelt[0], tokenIdFelt[1]}

	resp, err := CallContract(ctx, collectibleAddressFelt, "get_token_data", calldata)
	if err != nil {
		return "", nil, 0, "", fmt.Errorf("failed to get token data: %w", err)
	}

	pointsContract := PadZerosInFelt(resp[0])
	price := FeltArrToBigInt256([2]*felt.Felt{resp[1], resp[2]})
	expiry := utils.FeltToBigInt(resp[3]).Uint64()
	description, err := utils.ByteArrFeltToString(resp[4:])
	if err != nil {
		return "", nil, 0, "", fmt.Errorf("failed to convert felt array to string: %w", err)
	}

	return pointsContract, price, expiry, description, nil
}

func Redeem(ctx context.Context, account *account.Account, collectibleAddress string, user string, tokenId *big.Int, amount *big.Int) (string, error) {
	userFelt, err := utils.HexToFelt(user)
	if err != nil {
		return "", fmt.Errorf("failed to convert user address to felt: %w", err)
	}
	tokenIdFelt := BigInt256ToFelt(tokenId)
	amountFelt := BigInt256ToFelt(amount)
	calldata := []*felt.Felt{userFelt, tokenIdFelt[0], tokenIdFelt[1], amountFelt[0], amountFelt[1]}

	resp, err := InvokeTransaction(account, collectibleAddress, "redeem", calldata)
	if err != nil {
		return "", fmt.Errorf("failed to redeem: %w", err)
	}

	return resp.TransactionHash.String(), nil
}

func GetDetails(ctx context.Context, collectibleAddress string) (string, string, []*big.Int, []*big.Int, []uint64, []string, error) {
	collectibleAddressFelt, err := utils.HexToFelt(collectibleAddress)
	if err != nil {
		return "", "", nil, nil, nil, nil, fmt.Errorf("failed to convert collectible address to felt: %w", err)
	}
	resp, err := CallContract(ctx, collectibleAddressFelt, "get_details", []*felt.Felt{})
	if err != nil {
		return "", "", nil, nil, nil, nil, fmt.Errorf("failed to get details: %w", err)
	}

	// Parse the response
	count := utils.FeltToBigInt(resp[0]).Int64()
	description, err := utils.ByteArrFeltToString(resp[0 : 3+count])
	if err != nil {
		return "", "", nil, nil, nil, nil, fmt.Errorf("failed to convert description to string: %w", err)
	}
	currentIndex := 3 + count
	pointsContract := PadZerosInFelt(resp[currentIndex])
	currentIndex += 1

	// Parse token IDs
	tokenIDsLen := utils.FeltToBigInt(resp[currentIndex]).Int64()
	currentIndex += 1
	tokenIDs := make([]*big.Int, tokenIDsLen)
	for i := int64(0); i < tokenIDsLen; i += 2 {
		tokenIDs[i] = FeltArrToBigInt256([2]*felt.Felt{resp[currentIndex+i], resp[currentIndex+i+1]})
	}
	currentIndex += tokenIDsLen * 2

	// Parse token prices
	tokenPricesLen := utils.FeltToBigInt(resp[currentIndex]).Int64()
	currentIndex += 1
	tokenPrices := make([]*big.Int, tokenPricesLen)
	for i := int64(0); i < tokenPricesLen; i += 2 {
		tokenPrices[i] = FeltArrToBigInt256([2]*felt.Felt{resp[currentIndex+i], resp[currentIndex+i+1]})
	}
	currentIndex += tokenPricesLen * 2

	// Parse token expiries
	tokenExpiriesLen := utils.FeltToBigInt(resp[currentIndex]).Int64()
	currentIndex += 1
	tokenExpiries := make([]uint64, tokenExpiriesLen)
	for i := int64(0); i < tokenExpiriesLen; i++ {
		tokenExpiries[i] = utils.FeltToBigInt(resp[currentIndex+i]).Uint64()
	}
	currentIndex += tokenExpiriesLen

	// Parse token descriptions
	tokenDescriptionsLen := utils.FeltToBigInt(resp[currentIndex]).Int64()
	currentIndex += 1
	tokenDescriptions := make([]string, tokenDescriptionsLen)
	for i := int64(0); i < tokenDescriptionsLen; i++ {
		descriptionLen := utils.FeltToBigInt(resp[currentIndex]).Int64()
		desc, err := utils.ByteArrFeltToString(resp[currentIndex : currentIndex+3+descriptionLen])
		if err != nil {
			return "", "", nil, nil, nil, nil, fmt.Errorf("failed to convert token description to string: %w", err)
		}
		currentIndex += 3 + descriptionLen
		tokenDescriptions[i] = desc
	}

	return description, pointsContract, tokenIDs, tokenPrices, tokenExpiries, tokenDescriptions, nil
}

func IsValid(ctx context.Context, collectibleAddress string, tokenId *big.Int) (bool, error) {
	collectibleAddressFelt, err := utils.HexToFelt(collectibleAddress)
	if err != nil {
		return false, fmt.Errorf("failed to convert collectible address to felt: %w", err)
	}
	tokenIdFelt := BigInt256ToFelt(tokenId)
	calldata := []*felt.Felt{tokenIdFelt[0], tokenIdFelt[1]}

	resp, err := CallContract(ctx, collectibleAddressFelt, "is_valid", calldata)
	if err != nil {
		return false, fmt.Errorf("failed to check validity: %w", err)
	}

	return utils.FeltToBigInt(resp[0]).Uint64() == 1, nil
}

func Purchase(ctx context.Context, account *account.Account, collectibleAddress string, user string, tokenId *big.Int, amount *big.Int) (string, error) {
	userFelt, err := utils.HexToFelt(user)
	if err != nil {
		return "", fmt.Errorf("failed to convert user address to felt: %w", err)
	}
	tokenIdFelt := BigInt256ToFelt(tokenId)
	amountFelt := BigInt256ToFelt(amount)
	calldata := []*felt.Felt{userFelt, tokenIdFelt[0], tokenIdFelt[1], amountFelt[0], amountFelt[1]}

	resp, err := InvokeTransaction(account, collectibleAddress, "purchase", calldata)
	if err != nil {
		return "", fmt.Errorf("failed to purchase: %w", err)
	}

	return resp.TransactionHash.String(), nil
}
