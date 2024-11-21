package infinirewards

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"time"

	"infinirewards/logs"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// InvokeTransactionMaster invokes a transaction on the master account
//	@param		contractAddressStr:		The	address		of	the	contract
//	@param		functionSelectorStr:	The	selector	of	the	function
//	@param		calldata:				The	calldata	of	the	function
//	@return:	The transaction receipt and an error
func InvokeTransactionMaster(contractAddressStr string, functionSelectorStr string, calldata []*felt.Felt) (*rpc.TransactionReceiptWithBlockInfo, error) {
	contractAddress, err := utils.HexToFelt(contractAddressStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert contract address %s to felt: %w", contractAddressStr, err)
	}
	// Get the current nonce for the master account
	nonce, err := masterAccnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, masterAccnt.AccountAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Build the InvokeTx struct
	invokeTx := rpc.InvokeTxnV3{
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: masterAccnt.AccountAddress,
		Version:       rpc.TransactionV3,
		// Signature: ,
		ResourceBounds: rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x2710",
				MaxPricePerUnit: "0x174876E800",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		// MaxFee:        new(felt.Felt).SetUint64(1000000000000000),
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
		Nonce:                 nonce,
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
	}

	// Prepare the function call
	fnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(functionSelectorStr),
		Calldata:           calldata,
	}

	// Format the calldata
	invokeTx.Calldata, err = masterAccnt.FmtCalldata([]rpc.FunctionCall{fnCall})
	if err != nil {
		return nil, fmt.Errorf("failed to format calldata: %w", err)
	}

	// Sign the transaction
	err = SignInvokeTransaction(context.Background(), masterAccnt, &invokeTx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Execute the transaction
	resp, err := masterAccnt.AddInvokeTransaction(context.Background(), rpc.BroadcastInvokev3Txn{InvokeTxnV3: invokeTx})
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}

	// Wait for the transaction to be accepted
	txStatus, err := waitForTransaction(resp.TransactionHash, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	// Get the transaction details
	receipt, err := Client.TransactionReceipt(context.Background(), resp.TransactionHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction details: %w", err)
	}

	if txStatus.ExecutionStatus != rpc.TxnExecutionStatusSUCCEEDED {
		return nil, fmt.Errorf("transaction failed with status: %s, error: %s", txStatus.ExecutionStatus, receipt.TransactionReceipt.RevertReason)
	}
	fmt.Printf("Gas used: %.6e %s\n", float64(utils.FeltToBigInt(receipt.TransactionReceipt.ActualFee.Amount).Int64()), receipt.TransactionReceipt.ActualFee.Unit)

	return receipt, nil
}

// InvokeTransaction invokes a transaction on an account
//	@param		account:				The	account		of	the	user
//	@param		contractAddressStr:		The	address		of	the	contract
//	@param		functionSelectorStr:	The	selector	of	the	function
//	@param		calldata:				The	calldata	of	the	function
//	@return:	The transaction receipt and an error
func InvokeTransaction(account *account.Account, contractAddressStr string, functionSelectorStr string, calldata []*felt.Felt) (*rpc.TransactionReceiptWithBlockInfo, error) {
	contractAddress, err := utils.HexToFelt(contractAddressStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert contract address to felt: %w", err)
	}
	// Get the current nonce for the master account
	nonce, err := account.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, account.AccountAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Build the InvokeTx struct
	invokeTx := rpc.InvokeTxnV3{
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: account.AccountAddress,
		Version:       rpc.TransactionV3,
		ResourceBounds: rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x2710",
				MaxPricePerUnit: "0x174876E800",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		// MaxFee:        new(felt.Felt).SetUint64(1000000000000000),
		Nonce:                 nonce,
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
	}

	// Prepare the function call
	fnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(functionSelectorStr),
		Calldata:           calldata,
	}

	// Format the calldata
	invokeTx.Calldata, err = account.FmtCalldata([]rpc.FunctionCall{fnCall})
	if err != nil {
		return nil, fmt.Errorf("failed to format calldata: %w", err)
	}

	// Sign the transaction
	err = SignInvokeTransaction(context.Background(), account, &invokeTx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Execute the transaction
	resp, err := account.AddInvokeTransaction(context.Background(), rpc.BroadcastInvokev3Txn{InvokeTxnV3: invokeTx})
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}

	// Wait for the transaction to be accepted
	txStatus, err := waitForTransaction(resp.TransactionHash, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	// Get the transaction details
	receipt, err := Client.TransactionReceipt(context.Background(), resp.TransactionHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction details: %w", err)
	}

	if txStatus.ExecutionStatus != rpc.TxnExecutionStatusSUCCEEDED {
		return nil, fmt.Errorf("transaction failed with status: %s, error: %s", txStatus.ExecutionStatus, receipt.TransactionReceipt.RevertReason)
	}
	fmt.Printf("Gas used: %.6e %s\n", float64(utils.FeltToBigInt(receipt.TransactionReceipt.ActualFee.Amount).Int64()), receipt.TransactionReceipt.ActualFee.Unit)

	return receipt, nil
}

func SignInvokeTransaction(ctx context.Context, account *account.Account, invokeTx *rpc.InvokeTxnV3) error {
	txHash, err := account.TransactionHashInvoke(*invokeTx)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return err
	}
	invokeTx.Signature = signature
	return nil
}

// waitForTransaction waits for a transaction to be accepted
//	@param		txHash:		The	hash	of		the	transaction
//	@param		maxRetries:	The	maximum	number	of	retries
//	@return:	The transaction status and an error
func waitForTransaction(txHash *felt.Felt, maxRetries int) (*rpc.TxnStatusResp, error) {
	for i := 0; i < maxRetries; i++ {
		status, err := Client.GetTransactionStatus(context.Background(), txHash)
		if err != nil {
			logs.Logger.Error("failed to get transaction status",
				slog.String("handler", "waitForTransaction"),
				slog.String("tx_hash", txHash.String()),
				slog.String("error", err.Error()),
				slog.Int("attempt", i+1),
			)
			return nil, err
		}

		if status.FinalityStatus == rpc.TxnStatus_Accepted_On_L2 {
			logs.Logger.Debug("transaction confirmed",
				slog.String("handler", "waitForTransaction"),
				slog.String("tx_hash", txHash.String()),
				slog.Int("attempts", i+1),
			)
			return status, nil
		}

		time.Sleep(5 * time.Second)
	}

	logs.Logger.Error("transaction confirmation timeout",
		slog.String("handler", "waitForTransaction"),
		slog.String("tx_hash", txHash.String()),
		slog.Int("max_retries", maxRetries),
	)
	return nil, fmt.Errorf("transaction not confirmed after %d attempts", maxRetries)
}

// GetAccount gets an account
//	@param		provider:		The	provider
//	@param		privateKey:		The	private	key	of	the	account
//	@param		accountAddress:	The	address	of	the	account
//	@return:	The account and an error
func GetAccount(privateKey string, publicKey string, accountAddress string) (*account.Account, error) {
	keyStore := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert privKey to bigInt")
	}
	keyStore.Put(publicKey, privKeyBI)
	// Here we are converting the account address to felt
	accountAddressInFelt, err := HexToFelt(accountAddress)
	if err != nil {
		return nil, err
	}
	// Initialize the account
	result, err := account.NewAccount(Client, accountAddressInFelt, publicKey, keyStore, 2)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CallContract calls a contract
//	@param		ctx:					The	context
//	@param		contractAddress:		The	address		of	the	contract
//	@param		functionSelectorStr:	The	selector	of	the	function
//	@param		calldata:				The	calldata	of	the	function
//	@return:	The response and an error
func CallContract(ctx context.Context, contractAddress *felt.Felt, functionSelectorStr string, calldata []*felt.Felt) ([]*felt.Felt, error) {
	// Get balance from specified account address. Make read contract call with calldata
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(functionSelectorStr),
		Calldata:           calldata,
	}
	resp, rpcErr := Client.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		panic(rpcErr)
	}
	return resp, nil
}

func FundAccount(address string) (string, error) {
	// Getting the nonce from the account
	nonce, err := masterAccnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, masterAccnt.AccountAddress)
	if err != nil {
		return "", err
	}

	// Building the InvokeTx struct
	InvokeTx := rpc.BroadcastInvokev1Txn{
		InvokeTxnV1: rpc.InvokeTxnV1{
			MaxFee:        new(felt.Felt).SetUint64(100000000000000),
			Version:       rpc.TransactionV1,
			Nonce:         nonce,
			Type:          rpc.TransactionType_Invoke,
			SenderAddress: masterAccnt.AccountAddress,
		}}

	// Converting the STRK contractAddress from hex to felt
	contractAddress, err := HexToFelt("0x4718F5A0FC34CC1AF16A1CDEE98FFB20C31F5CD61D6AB07201858F4287C938D")
	if err != nil {
		return "", err
	}
	// Sending 1 STRK to the account
	amount, _ := HexToFelt("0xDE0B6B3A7640000")
	receipient, _ := HexToFelt(address)
	// Building the functionCall struct, where :
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,                              //contractAddress is the contract that we want to call
		EntryPointSelector: utils.GetSelectorFromNameFelt("transfer"),    //this is the function that we want to call
		Calldata:           []*felt.Felt{receipient, amount, &felt.Zero}, //the calldata necessary to call the function. Here we are passing the "amount" value for the "mint" function
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	InvokeTx.Calldata, err = masterAccnt.FmtCalldata([]rpc.FunctionCall{FnCall})
	if err != nil {
		return "", err
	}

	// Signing of the transaction that is done by the account
	err = masterAccnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV1)
	if err != nil {
		return "", err
	}

	// Estimate the transaction fee
	feeRes, err := masterAccnt.EstimateFee(context.Background(), []rpc.BroadcastTxn{InvokeTx}, []rpc.SimulationFlag{}, rpc.WithBlockTag("latest"))
	if err != nil {
		return "", PanicRPC(err)
	}
	estimatedFee := feeRes[0].OverallFee
	// If the estimated fee is higher than the current fee, let's override it and sign again
	if estimatedFee.Cmp(InvokeTx.MaxFee) == 1 {
		newFee, err := strconv.ParseUint(estimatedFee.String(), 0, 64)
		if err != nil {
			return "", err
		}
		InvokeTx.MaxFee = new(felt.Felt).SetUint64(newFee + newFee/5) // fee + 20% to be sure
		// Signing the transaction again
		err = masterAccnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV1)
		if err != nil {
			return "", err
		}
	}

	// After the signing we finally call the AddInvokeTransaction in order to invoke the contract function
	resp, err := masterAccnt.AddInvokeTransaction(context.Background(), InvokeTx)
	if err != nil {
		return "", PanicRPC(err)
	}

	_, err = waitForTransaction(resp.TransactionHash, 5)
	if err != nil {
		return "", err
	}

	// This returns us with the transaction hash and status
	// fmt.Printf("Transaction hash response: %v\n", resp.TransactionHash)
	// fmt.Printf("Transaction execution status: %s\n", txStatus.ExecutionStatus)
	// fmt.Printf("Transaction status: %s\n", txStatus.FinalityStatus)
	return resp.TransactionHash.String(), nil
}

// PadZerosInFelt pads zeros to the left of a hex felt value to make it 64 characters long.
func PadZerosInFelt(hexFelt *felt.Felt) string {
	length := 66
	hexStr := hexFelt.String()

	// Check if the hex value is already of the desired length
	if len(hexStr) >= length {
		return hexStr
	}

	// Extract the hex value without the "0x" prefix
	hexValue := hexStr[2:]
	// Pad zeros after the "0x" prefix
	paddedHexValue := fmt.Sprintf("%0*s", length-2, hexValue)
	// Add back the "0x" prefix to the padded hex value
	paddedHexStr := "0x" + paddedHexValue

	return paddedHexStr
}

// HexToFelt converts a hex string to a felt
//	@param		hexStr:	The	hex	string
//	@return:	The felt and an error
func HexToFelt(hexStr string) (*felt.Felt, error) {
	return utils.HexToFelt(hexStr)
}

// StrToFelt converts a string to a felt
//	@param		s:	The	string
//	@return:	The felt and an error
func StrToFelt(s string) (*felt.Felt, error) {
	return utils.HexToFelt(hex.EncodeToString([]byte(s)))
}

// PanicRPC panics on an RPC error
//	@param		err:	The	error
//	@return:	The error
func PanicRPC(err error) error {
	return fmt.Errorf("RPC error: %w", err)
}

func HashPhoneNumber(phoneNumber string) *felt.Felt {
	h := sha256.New()
	h.Write([]byte(phoneNumber))
	// phoneNumberHash := hex.EncodeToString(h.Sum(nil))
	return new(felt.Felt).SetBytes(h.Sum(nil))
}

func BigInt256ToFelt(bi *big.Int) []*felt.Felt {
	// Convert tokenId (u256) to two felt252 values
	low := utils.BigIntToFelt(new(big.Int).And(bi, new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))))
	high := utils.BigIntToFelt(new(big.Int).Rsh(bi, 128))
	return []*felt.Felt{low, high}
}

func FeltArrToBigInt256(felt [2]*felt.Felt) *big.Int {
	lowBits := utils.FeltToBigInt(felt[0])
	highBits := utils.FeltToBigInt(felt[1])

	// Combine low and high bits into a single big.Int
	balance := new(big.Int).Lsh(highBits, 128) // Shift high bits left by 128 bits
	return new(big.Int).Or(balance, lowBits)
}
