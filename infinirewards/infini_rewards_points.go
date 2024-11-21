package infinirewards

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/utils"
)

// MintPoints mints points
//	@param		account:		The	account	of	the		merchant
//	@param		pointsContract:	The	address	of	the		points	contract
//	@param		recipient:		The	address	of	the		recipient
//	@param		amount:			The	amount	of	points	to	mint
//	@return:	The transaction hash and an error
func MintPoints(account *account.Account, pointsContract string, recipient string, amount *big.Int) (string, error) {
	recipientFelt, err := utils.HexToFelt(recipient)
	if err != nil {
		return "", fmt.Errorf("failed to convert recipient address to felt: %w", err)
	}
	amountFelt := BigInt256ToFelt(amount)

	resp, err := InvokeTransaction(account, pointsContract, "mint", []*felt.Felt{recipientFelt, amountFelt[0], amountFelt[1]})
	if err != nil {
		return "", fmt.Errorf("failed to mint points: %w", err)
	}

	return resp.TransactionHash.String(), nil
}

// BurnPoints burns points
//	@param		account:		The	account	of	the		merchant
//	@param		pointsContract:	The	address	of	the		points	contract
//	@param		amount:			The	amount	of	points	to		burn
//	@return:	The transaction hash and an error
func BurnPoints(account *account.Account, pointsContract string, amount *big.Int) (string, error) {
	amountFelt := BigInt256ToFelt(amount)

	resp, err := InvokeTransaction(account, pointsContract, "burn", []*felt.Felt{amountFelt[0], amountFelt[1]})
	if err != nil {
		return "", fmt.Errorf("failed to burn points: %w", err)
	}

	return resp.TransactionHash.String(), nil
}

// GetBalance gets the balance of a points contract
//	@param		account:		The	account	of	the	merchant
//	@param		pointsContract:	The	address	of	the	points	contract
//	@return:	The balance and an error
func GetBalance(ctx context.Context, account *account.Account, pointsContract string) (*big.Int, error) {
	contractAddress, err := utils.HexToFelt(pointsContract)
	if err != nil {
		return nil, fmt.Errorf("failed to convert collectible contract address: %w", err)
	}

	resp, err := CallContract(ctx, contractAddress, "balanceOf", []*felt.Felt{account.AccountAddress})
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	balance := FeltArrToBigInt256([2]*felt.Felt{resp[0], resp[1]})
	return balance, nil
}

// TransferPoints transfers points
//	@param		account:		The	account	of	the		merchant
//	@param		pointsContract:	The	address	of	the		points	contract
//	@param		from:			The	address	of	the		sender
//	@param		to:				The	address	of	the		recipient
//	@param		amount:			The	amount	of	points	to	transfer
//	@return:	The transaction hash and an error
func TransferPoints(ctx context.Context, account *account.Account, pointsContract string, to string, amount *big.Int) (string, error) {
	toFelt, err := utils.HexToFelt(to)
	if err != nil {
		return "", fmt.Errorf("failed to convert 'to' address to felt: %w", err)
	}
	amountFelt := BigInt256ToFelt(amount)
	calldata := []*felt.Felt{toFelt, amountFelt[0], amountFelt[1]}

	resp, err := InvokeTransaction(account, pointsContract, "transfer", calldata)
	if err != nil {
		return "", fmt.Errorf("failed to transfer points: %w", err)
	}

	return resp.TransactionHash.String(), nil
}

// GetPointsContractDetails gets the details of a points contract
//	@param		pointsContract:	The	address	of	the	points	contract
//	@return:	The name, symbol, description, decimals, total supply, and an error
func GetPointsContractDetails(ctx context.Context, pointsContract string) (string, string, string, uint64, uint64, error) {
	contractAddress, err := utils.HexToFelt(pointsContract)
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to convert points contract address: %w", err)
	}

	resp, err := CallContract(ctx, contractAddress, "get_details", []*felt.Felt{})
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to get details: %w", err)
	}

	count := utils.FeltToBigInt(resp[0]).Int64()
	name, err := utils.ByteArrFeltToString(resp[0 : 3+count])
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to convert name to string: %w", err)
	}
	currentIndex := 3 + count

	count = utils.FeltToBigInt(resp[currentIndex]).Int64()
	symbol, err := utils.ByteArrFeltToString(resp[currentIndex : currentIndex+3+count])
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to convert symbol to string: %w", err)
	}
	currentIndex += (count + 3)

	count = utils.FeltToBigInt(resp[currentIndex]).Int64()
	description, err := utils.ByteArrFeltToString(resp[currentIndex : currentIndex+3+count])
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to convert description to string: %w", err)
	}
	currentIndex += (count + 3)

	decimals := utils.FeltToBigInt(resp[currentIndex]).Uint64()
	currentIndex += 1
	totalSupply := utils.FeltToBigInt(resp[currentIndex]).Uint64()

	return name, symbol, description, decimals, totalSupply, nil
}
