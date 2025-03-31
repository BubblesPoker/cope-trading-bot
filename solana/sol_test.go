package solana

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/program/token"
	"log"
	"testing"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/types"
)

func TestTransferSOL(t *testing.T) {
	// Connect to Solana Testnet
	c := client.NewClient("https://api.testnet.solana.com")

	// Sender's private key (Base58 format)
	senderPrivateKey := "YOUR_PRIVATE_KEY"

	// Decode sender's private key
	sender, err := types.AccountFromBase58(senderPrivateKey)
	if err != nil {
		log.Fatalf("failed to load sender account: %v", err)
	}

	//Ensure you have enough SOL for transaction fees.
	//
	//Find Associated Token Accounts (ATA):
	//
	//If the receiver does not have an ATA, you must create one.
	//
	//Use the following function:
	//createATAInstruction := token.CreateAssociatedTokenAccount(
	//    sender.PublicKey,
	//    common.PublicKeyFromString(receiver),
	//    common.PublicKeyFromString(tokenMint),
	//)

	// SPL Token details
	tokenMint := "SPL_TOKEN_MINT_ADDRESS"              // Token Mint Address
	senderATA := "SENDER_ASSOCIATED_TOKEN_ACCOUNT"     // Sender's ATA (Associated Token Account)
	receiverATA := "RECEIVER_ASSOCIATED_TOKEN_ACCOUNT" // Receiver's ATA

	// Amount to send (in smallest unit, e.g., 1000000 = 1 USDC for 6 decimal places)
	amount := uint64(1000000)

	// Get latest blockhash
	recentBlockhash, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get latest blockhash: %v", err)
	}

	// Create transfer instruction
	transferInstruction := token.TransferChecked(token.TransferCheckedParam{
		From:     common.PublicKeyFromString(senderATA),
		To:       common.PublicKeyFromString(receiverATA),
		Mint:     common.PublicKeyFromString(tokenMint),
		Auth:     sender.PublicKey,
		Amount:   amount,
		Decimals: 6, // Adjust decimals based on token type (e.g., 6 for USDC)
	})

	// Create transaction
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{sender},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        sender.PublicKey,
			RecentBlockhash: recentBlockhash.Blockhash,
			Instructions:    []types.Instruction{transferInstruction},
		}),
	})
	if err != nil {
		log.Fatalf("failed to create transaction: %v", err)
	}

	// Send transaction
	txSig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send transaction: %v", err)
	}

	fmt.Println("Transaction successful! Signature:", txSig)
}
func TestTransferSPL(t *testing.T) {
	// Connect to Solana Testnet
	c := client.NewClient("https://api.devnet.solana.com")

	// Sender's private key (Base58 format)
	senderPrivateKey := ""

	// Decode sender's private key

	sender, err := types.AccountFromBase58(senderPrivateKey)

	//types.AccountFromSeed("“)
	if err != nil {
		log.Fatalf("failed to load sender account: %v", err)
	}

	// Receiver's public key
	receiver := "72AeY74VVznEUjWCyzUTr9Go7sxx3VhvUhntF5id4kyo"

	// Amount to send (in lamports, 1 SOL = 1_000_000_000 lamports)
	//amount := uint64(1000000) // 0.001 SOL
	amount := uint64(100000000) // 0.1 SOL

	// Get latest blockhash
	recentBlockhash, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get latest blockhash: %v", err)
	}

	// Create transfer instruction
	transferInstruction := system.Transfer(system.TransferParam{
		From:   sender.PublicKey,
		To:     common.PublicKeyFromString(receiver),
		Amount: amount,
	})

	// Create transaction
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{sender},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        sender.PublicKey,
			RecentBlockhash: recentBlockhash.Blockhash,
			Instructions:    []types.Instruction{transferInstruction},
		}),
	})
	if err != nil {
		log.Fatalf("failed to create transaction: %v", err)
	}

	// Send transaction
	txSig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send transaction: %v", err)
	}

	fmt.Println("Transaction successful! Signature:", txSig)
}
