package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	as "github.com/PaddyMc/blockaid-cosigner/pkg/authenticator"
	"github.com/PaddyMc/blockaid-cosigner/pkg/config"
	pm "github.com/PaddyMc/blockaid-cosigner/pkg/poolmanager"
)

func SeedCreateOneClickTradingAccount(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-one-click-trading-flow",
		Short: "this command goes through a series of tasks to test the one click trading flow",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[2]
			bob := seedConfig.Keys[3]
			cosigners := make(map[int][]cryptotypes.PrivKey)
			OsmoDenom := seedConfig.DenomMap["OsmoDenom"]
			AtomIBCDenom := seedConfig.DenomMap["AtomIBCDenom"]
			LuncIBCDenom := seedConfig.DenomMap["LuncIBCDenom"]
			osmoAtomClPool := uint64(1400)
			luncOsmoBalancerPool := uint64(561)
			selectedAuthenticator := []uint64{1}

			spendLimitContractAddress := "osmo1x57l2yanux5277kht6udkgdxkkuynv6ndm5836x5ll4hgwcxfhlstmnwp3"

			log.Printf("Starting spend limit authenticator flow")
			log.Printf("Adding spend limit authenticator")
			err := as.CreateOneClickTradingAccount(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				spendLimitContractAddress,
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
				bob,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				osmoAtomClPool,
				100,
			)
			if err != nil {
				return err
			}

			log.Printf("Starting swappping to Lunc, should fail")
			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				LuncIBCDenom,
				luncOsmoBalancerPool,
				1000,
			)
			if err != nil {
				// we expected this to fail
				log.Println("error", err.Error())
			}

			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				osmoAtomClPool,
				100,
			)
			if err != nil {
				return err
			}

			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				osmoAtomClPool,
				100,
			)
			if err != nil {
				return err
			}

			log.Printf("Removing spend limit authenticator")
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
