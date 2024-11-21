package infinirewards

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
)

// GetPhoneNumber gets the phone number hash of a user
//	@param		account:	The	account	of	the	user
//	@return:	The phone number hash and an error
func GetPhoneNumber(ctx context.Context, account *account.Account) (string, error) {
	calldata := []*felt.Felt{}

	resp, err := CallContract(ctx, account.AccountAddress, "get_phone_number", calldata)
	if err != nil {
		return "", fmt.Errorf("failed to get user account: %v", err)
	}

	return PadZerosInFelt(resp[0]), nil
}
