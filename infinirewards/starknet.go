package infinirewards

import (
	"context"
	"fmt"
	"infinirewards/logs"
	"log/slog"
	"math/big"
	"os"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

var Client *rpc.Provider
var masterKs *account.MemKeystore
var masterAccnt *account.Account

// New variables for InfiniRewards contracts
var InfiniRewardsFactoryAddress string

func ConnectStarknet() error {

	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	var rpcProviderUrl string
	var masterPrivKey string
	var masterAccntAddress string

	network := os.Getenv("NETWORK")
	logs.Logger.Info("connecting to starknet",
		slog.String("network", network),
	)
	switch network {
	case "devnet":
		rpcProviderUrl = os.Getenv("RPC_PROVIDER_URL_DEVNET")
		masterPrivKey = os.Getenv("MASTER_PRIVATE_KEY_DEVNET")
		masterAccntAddress = os.Getenv("MASTER_ACCOUNT_ADDRESS_DEVNET")
		InfiniRewardsFactoryAddress = os.Getenv("INFINI_REWARDS_FACTORY_ADDRESS_DEVNET")
	case "sepolia":
		rpcProviderUrl = os.Getenv("RPC_PROVIDER_URL_SEPOLIA")
		masterPrivKey = os.Getenv("MASTER_PRIVATE_KEY_SEPOLIA")
		masterAccntAddress = os.Getenv("MASTER_ACCOUNT_ADDRESS_SEPOLIA")
		InfiniRewardsFactoryAddress = os.Getenv("INFINI_REWARDS_FACTORY_ADDRESS_SEPOLIA")
	case "mainnet":
		rpcProviderUrl = os.Getenv("RPC_PROVIDER_URL_MAINNET")
		masterPrivKey = os.Getenv("MASTER_PRIVATE_KEY_MAINNET")
		masterAccntAddress = os.Getenv("MASTER_ACCOUNT_ADDRESS_MAINNET")
		InfiniRewardsFactoryAddress = os.Getenv("INFINI_REWARDS_FACTORY_ADDRESS_MAINNET")
	default:
		return fmt.Errorf("invalid network configuration: %s", network)
	}

	Client, err = rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		return fmt.Errorf("failed to create RPC provider for URL %s: %w", rpcProviderUrl, err)
	}

	masterKs = account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(masterPrivKey, 0)
	if !ok {
		return fmt.Errorf("failed to convert private key to big.Int (key length: %d)", len(masterPrivKey))
	}

	masterKs.Put(masterAccntAddress, privKeyBI)
	accountAddressInFelt, err := HexToFelt(masterAccntAddress)
	if err != nil {
		return fmt.Errorf("failed to convert account address %s to felt: %w", masterAccntAddress, err)
	}

	_, err = Client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	masterAccnt, err = account.NewAccount(Client, accountAddressInFelt, masterAccntAddress, masterKs, 2)
	if err != nil {
		return fmt.Errorf("failed to create master account for address %s: %w", masterAccntAddress, err)
	}

	return nil
}

// func DeployAccount(phoneNumber string) (string, string, string, error) {
// 	// Initialise the client.

// 	// Get random keys for test purposes
// 	_, pub, privKey := account.GetRandomKeys()

// 	nonce, err := masterAccnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, masterAccnt.AccountAddress)
// 	if err != nil {
// 		fmt.Println("Failed to get the nonce")
// 		return "", "", "", err
// 	}

// 	// Build the InvokeTx struct
// 	InvokeTx := rpc.BroadcastInvokev1Txn{
// 		InvokeTxnV1: rpc.InvokeTxnV1{
// 			MaxFee:        new(felt.Felt).SetUint64(100000000000000),
// 			Version:       rpc.TransactionV1,
// 			Nonce:         nonce,
// 			Type:          rpc.TransactionType_Invoke,
// 			SenderAddress: masterAccnt.AccountAddress,
// 		}}

// 	// Convert the contractAddress from hex to felt
// 	contractAddress, err := HexToFelt(UDCAddress)
// 	if err != nil {
// 		fmt.Println("Failed to convert the contract address")
// 		return "", "", "", err
// 	}

// 	classHash, err := HexToFelt(accountContractClassHash)
// 	if err != nil {
// 		return "", "", "", err
// 	}

// 	randomInt := rand.Uint64()
// 	salt := new(felt.Felt).SetUint64(randomInt) // to prevent address clashes

// 	unique, err := HexToFelt("0x0") // see https://docs.starknet.io/architecture-and-concepts/accounts/universal-deployer/#deployment_types
// 	if err != nil {
// 		return "", "", "", err
// 	}

// 	hasher := sha256.New()
// 	hasher.Write([]byte(phoneNumber))
// 	hashedPhoneNumber := hasher.Sum(nil)
// 	phoneNumberFelt, err := HexToFelt(hex.EncodeToString(hashedPhoneNumber))
// 	if err != nil {
// 		return "", "", "", err
// 	}

// 	calldata := []*felt.Felt{
// 		pub,
// 		phoneNumberFelt,
// 	}

// 	length := int64(len(calldata))
// 	calldataLen, err := HexToFelt(strconv.FormatInt(length, 16))
// 	if err != nil {
// 		return "", "", "", err
// 	}

// 	udcCalldata := append([]*felt.Felt{classHash, salt, unique, calldataLen}, calldata...)

// 	deployedContractAddress := contracts.PrecomputeAddress(unique, salt, classHash, calldata)

// 	// Build the functionCall struct, where :
// 	FnCall := rpc.FunctionCall{
// 		ContractAddress:    contractAddress,                                     //contractAddress is the contract that we want to call
// 		EntryPointSelector: utils.GetSelectorFromNameFelt(deployContractMethod), //this is the function that we want to call
// 		Calldata:           udcCalldata,                                         //change this function content to your use case
// 	}

// 	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
// 	InvokeTx.Calldata, err = masterAccnt.FmtCalldata([]rpc.FunctionCall{FnCall})
// 	if err != nil {
// 		fmt.Println("Failed to format the calldata")
// 		return "", "", "", err
// 	}

// 	// Sign the transaction
// 	err = masterAccnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV1)
// 	if err != nil {
// 		fmt.Println("Failed to sign the transaction")
// 		return "", "", "", err
// 	}

// 	// Estimate the transaction fee
// 	feeRes, err := masterAccnt.EstimateFee(context.Background(), []rpc.BroadcastTxn{InvokeTx}, []rpc.SimulationFlag{}, rpc.WithBlockTag("latest"))
// 	if err != nil {
// 		fmt.Println("Failed to estimate the transaction fee")
// 		return "", "", "", PanicRPC(err)
// 	}
// 	estimatedFee := feeRes[0].OverallFee
// 	// If the estimated fee is higher than the current fee, let's override it and sign again
// 	if estimatedFee.Cmp(InvokeTx.MaxFee) == 1 {
// 		newFee, err := strconv.ParseUint(estimatedFee.String(), 0, 64)
// 		if err != nil {
// 			fmt.Println("Failed to parse the estimated fee")
// 			return "", "", "", err
// 		}
// 		InvokeTx.MaxFee = new(felt.Felt).SetUint64(newFee + newFee/5) // fee + 20% to be sure
// 		// Signing the transaction again
// 		err = masterAccnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV1)
// 		if err != nil {
// 			fmt.Println("Failed to sign the transaction again")
// 			return "", "", "", err
// 		}
// 	}

// 	// After the signing we finally call the AddInvokeTransaction in order to invoke the contract function
// 	resp, err := masterAccnt.AddInvokeTransaction(context.Background(), InvokeTx)
// 	if err != nil {
// 		fmt.Println("Failed to add the invoke transaction")
// 		return "", "", "", PanicRPC(err)
// 	}

// 	fmt.Println("Waiting for the transaction status...")
// 	time.Sleep(time.Second * 3) // Waiting 3 seconds

// 	//Getting the transaction status
// 	txStatus, err := Client.GetTransactionStatus(context.Background(), resp.TransactionHash)
// 	if err != nil {
// 		fmt.Println("Failed to get the transaction status")
// 		return "", "", "", PanicRPC(err)
// 	}

// 	// This returns us with the transaction hash and status
// 	fmt.Printf("Transaction hash response: %v\n", resp.TransactionHash)
// 	fmt.Printf("Transaction execution status: %s\n", txStatus.ExecutionStatus)
// 	fmt.Printf("Transaction status: %s\n", txStatus.FinalityStatus)
// 	return privKey.String(), pub.String(), PadZerosInFelt(deployedContractAddress), nil
// }
