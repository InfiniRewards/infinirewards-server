package infinirewards

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/utils"
)

// GetPointsContracts gets the points contracts of a merchant
//	@param		account:	The	account	of	the	merchant
//	@return:	The addresses of the points contracts and an error
func GetPointsContracts(ctx context.Context, account *account.Account) ([]string, error) {
	calldata := []*felt.Felt{}

	resp, err := CallContract(ctx, account.AccountAddress, "get_points_contracts", calldata)
	if err != nil {
		return nil, fmt.Errorf("failed to get user account: %v", err)
	}

	contractsLen := utils.FeltToBigInt(resp[0]).Int64()

	contracts := make([]string, contractsLen)

	for i := 0; i < int(contractsLen); i++ {
		contracts[i] = PadZerosInFelt(resp[i+1])
	}

	return contracts, nil
}

// GetCollectibleContracts gets the collectible contracts of a merchant
//	@param		account:	The	account	of	the	merchant
//	@return:	The addresses of the collectible contracts and an error
func GetCollectibleContracts(ctx context.Context, account *account.Account) ([]string, error) {
	calldata := []*felt.Felt{}

	resp, err := CallContract(ctx, account.AccountAddress, "get_collectible_contracts", calldata)
	if err != nil {
		return nil, fmt.Errorf("failed to get user account: %v", err)
	}

	contractsLen := utils.FeltToBigInt(resp[0]).Int64()

	contracts := make([]string, contractsLen)

	for i := 0; i < int(contractsLen); i++ {
		contracts[i] = PadZerosInFelt(resp[i+1])
	}

	return contracts, nil
}
