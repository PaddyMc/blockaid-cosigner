package seed

import (
	"log"

	"github.com/spf13/cobra"

	as "github.com/PaddyMc/blockaid-cosigner/pkg/authenticator"
	"github.com/PaddyMc/blockaid-cosigner/pkg/config"
)

func SeedRemoveAllAuthenticators(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-remove-all-authenticators-flow",
		Short: "this command removes all the authenticators for an account",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[2]

			log.Printf("Starting remove all authenticators flow")

			log.Printf("Removing authenticator")
			err := as.RemoveLatestAuthenticator(
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
