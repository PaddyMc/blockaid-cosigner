package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	as "github.com/PaddyMc/blockaid-cosigner/pkg/authenticator"
	"github.com/PaddyMc/blockaid-cosigner/pkg/config"
	pm "github.com/PaddyMc/blockaid-cosigner/pkg/poolmanager"
)

func SeedCreateCosigner(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-cosigner-flow",
		Short: "creates a cosigner key and does transactions",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig
			selectedAuthenticator := []uint64{1}

			alice := seedConfig.Keys[1]
			blockaid := seedConfig.Keys[2]

			// This is where we add the blockaid public key!
			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, blockaid}

			OsmoDenom := seedConfig.DenomMap["OsmoDenom"]
			AtomIBCDenom := seedConfig.DenomMap["AtomIBCDenom"]
			osmoAtomClPool := uint64(1400)

			log.Printf("Starting cosigner authenticator flow")
			log.Printf("Adding cosigner authenticator")
			err := as.CreateCosignerAccount(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				blockaid,
			)
			if err != nil {
				return err
			}

			log.Printf("Starting swap flow")
			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				blockaid,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				osmoAtomClPool,
				100,
			)
			if err != nil {
				log.Println("error", err.Error())
			}

			log.Printf("Removing cosigner authenticator")
			err = as.RemoveLatestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				alice,
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
