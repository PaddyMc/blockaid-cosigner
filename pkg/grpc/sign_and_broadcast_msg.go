package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	key "github.com/PaddyMc/blockaid-cosigner/pkg/key"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v24/app/params"
)

func SignAndBroadcastAuthenticatorMsgMultiSigners(
	senderPrivKeys []cryptotypes.PrivKey,
	signerPrivKeys []cryptotypes.PrivKey,
	cosignerPrivKeys map[int][]cryptotypes.PrivKey,
	encCfg params.EncodingConfig,
	ac authtypes.QueryClient,
	txClient txtypes.ServiceClient,
	chainID string,
	msgs []sdk.Msg,
	selectedAuthenticators []uint64,
) error {
	log.Println("Signing and broadcasting message flow")

	var accNums []uint64
	var accSeqs []uint64

	for _, privKey := range senderPrivKeys {
		// Generate the account address from the private key
		addr := sdk.AccAddress(privKey.PubKey().Address()).String()

		// Get the account information
		res, err := ac.Account(
			context.Background(),
			&authtypes.QueryAccountRequest{Address: addr},
		)
		if err != nil {
			return err
		}

		var acc authtypes.AccountI
		if err := encCfg.InterfaceRegistry.UnpackAny(res.Account, &acc); err != nil {
			return err
		}

		accNums = append(accNums, acc.GetAccountNumber())
		accSeqs = append(accSeqs, acc.GetSequence())
	}

	// Sign the message
	txBytes, _ := key.SignAuthenticatorMsg(
		encCfg.TxConfig,
		msgs,
		sdk.Coins{sdk.NewInt64Coin("uosmo", 7000)},
		1700000,
		chainID,
		accNums,
		accSeqs,
		senderPrivKeys,
		signerPrivKeys,
		cosignerPrivKeys,
		selectedAuthenticators,
	)

	log.Println("Simulate...")
	respSim, err := txClient.Simulate(
		context.Background(),
		&txtypes.SimulateRequest{
			TxBytes: txBytes,
		},
	)
	fmt.Println(respSim)
	if err != nil {
		panic(err)
	}

	log.Println("Broadcasting...")
	resp, err := txClient.BroadcastTx(
		context.Background(),
		&txtypes.BroadcastTxRequest{
			Mode:    txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return err
	}
	log.Println("Transaction Hash:", resp.TxResponse.TxHash)
	if resp.TxResponse.RawLog != "" {
		log.Println("Transaction failed reason:", resp.TxResponse.RawLog)
	}

	time.Sleep(6 * time.Second)

	log.Println("Verifing...")
	tx, err := txClient.GetTx(
		context.Background(),
		&txtypes.GetTxRequest{
			Hash: resp.TxResponse.TxHash,
		},
	)
	if err != nil {
		return err
	} else {
		if tx.TxResponse.Code == 0 {
			log.Println("Transaction Success...")
		} else {
			log.Println(tx.TxResponse)
		}
	}
	log.Println("Gas Used:", tx.TxResponse.GasUsed)

	return nil
}
