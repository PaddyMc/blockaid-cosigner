package main

import (
	"log"
	"os"

	"github.com/PaddyMc/blockaid-cosigner/cmd/seed"
	"github.com/PaddyMc/blockaid-cosigner/pkg/config"

	"github.com/spf13/cobra"
)

// Test data for the seeds to run
const (
	GrpcConnectionTimeoutSeconds = 10
	TestKeyUser1                 = "aec234a57cb59b801a6b1cf2a5f84ff124161b08be28e9bb64383b835c933e40"
	TestKeyUser2                 = "dd7701a7e359a2969208a967d1516efed5f1ca9c3e3204c5e73413b86ada65e4"
	TestKeyUser3                 = "3e51a79a23000d9c4164520dbf8def01a4158b2cb36bd4655b51e1e327ed354c"

	// TestUser4 is not in the auth store
	TestKeyUser4         = "1717ec286e6a208a4c471c97ee3d32abc2b93ddbcb4b7e5895814b390660b57b"
	AccountAddressPrefix = "osmo"
	ChainID              = "smartaccount"
	addr                 = "164.92.247.225:9090"
)

var DefaultDenoms = map[string]string{
	"OsmoDenom":     "uosmo",
	"IonDenom":      "uion",
	"StakeDenom":    "stake",
	"AtomDenom":     "uatom",
	"DaIBCiDenom":   "ibc/0CD3A0285E1341859B5E86B6AB7682F023D03E97607CCC1DC95706411D866DF7",
	"OsmoIBCDenom":  "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518",
	"StakeIBCDenom": "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B7787",
	"UstIBCDenom":   "ibc/BE1BB42D4BE3C30D50B68D7C41DB4DFCE9678E8EF8C539F6E6A9345048894FCC",
	"LuncIBCDenom":  "ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0",
	"AtomIBCDenom":  "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
	"UsdcIBCDenom":  "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4",
}

const (
	appName = "blockaid-cosigner"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// NewRootCmd returns the root command for parser.
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "blockaid-cosigner is a tool to test integration with third party signers and the authenticator module",
	}

	config := config.SetUp(
		ChainID,
		addr,
		[]string{
			TestKeyUser1,
			TestKeyUser2,
			TestKeyUser3,
			TestKeyUser4,
		},
		DefaultDenoms,
	)

	rootCmd.AddCommand(
		seed.SeedCreateOneClickTradingAccount(config),
		seed.SeedSwapCmd(config),
		seed.SeedRemoveAllAuthenticators(config),
		seed.SeedCreateCosigner(config),
		seed.StartBankSendFlow(config),
	)

	return rootCmd
}
