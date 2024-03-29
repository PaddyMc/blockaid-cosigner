# Blockaid Cosigner Flow for Authenticators

This repository contains a Go client used to interact with the Osmosis authenticator module, designed to facilitate testing and integration with third-party signers.

For detailed information, please refer to the [Osmosis Smart Accounts documentation](https://github.com/osmosis-labs/osmosis/tree/feat/smart-accounts).

## Overview of the Blockaid Cosigner Tool

The `blockaid-cosigner` tool is designed to test integration with third-party signers and the authenticator module in the Osmosis blockchain environment. It provides a range of commands to simulate various transaction flows involving smart accounts and cosigners.

### Available Commands

The tool offers the following commands:

- `start-one-click-trading-flow`: Tests the one-click trading flow.
- `start-swap-with-signature-authenticator-flow`: Creates a SignatureVerificationAuthenticator and executes a swap in a pool.
- `start-remove-all-authenticators-flow`: Removes all authenticators associated with an account.
- `start-cosigner-flow`: Creates a cosigner key and performs transactions.
- `start-bank-send-flow`: Executes bank sends to multiple accounts.
- `help`: Provides help and information about the commands.

To use the tool, run commands using the following syntax:

```bash
go run cmd/main.go <command>
```

### Configuration and Defaults

The tool is configured with sensible defaults for ease of testing:
```
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

```

## Blockaid Flow

The Blockaid flow in this client is a simplified version of what is expected in production implementations.

### Steps in the Flow

1. **Cosigner Creation**: Alice and Bob create a cosigner for Alice's account.
2. **Token Swap**: Execute a token swap in pool 1400 (osmo/atom), signed with the cosigner.
3. **Authenticator Removal**: Remove the cosigner authenticator from the account.

The cosigner authenticator is constructed using a `PartitionedAllOfAuthenticator`. This authenticator splits the transaction signatures and routes them to the appropriate authenticator.

### Key Areas of the Implementation

- **Cosigning Keys**: [Keys for cosigning](https://github.com/PaddyMc/blockaid-cosigner/blob/main/cmd/seed/cosigner_flow.go#L27-L30)
- **Authenticator Creation**: [Creating the complex authenticator](https://github.com/PaddyMc/blockaid-cosigner/blob/main/pkg/authenticator/cosigner.go#L52-L114)
- **Transaction Signing**: [Signing the signature data with both keys](https://github.com/PaddyMc/blockaid-cosigner/blob/main/pkg/key/sign_transaction.go#L129-L159)

### Running the Blockaid Flow

Execute the cosigner flow with the following command:

```bash
go run cmd/main.go start-cosigner-flow
```

#### Step 1: Create a Cosigner
```
2024/03/02 19:57:20 cosigner_flow.go:33: Starting cosigner authenticator flow
2024/03/02 19:57:20 cosigner_flow.go:34: Adding cosigner authenticator
2024/03/02 19:57:20 cosigner.go:49: Querying authenticators for account osmo1mveeh0ruel03usw3k4agxf68l5dmltyemgdlsk
2024/03/02 19:57:20 cosigner.go:50: Number of authenticators: 0
2024/03/02 19:57:20 cosigner.go:107: Adding authenticator for account osmo1mveeh0ruel03usw3k4agxf68l5dmltyemgdlsk first authenticator
2024/03/02 19:57:20 sign_and_broadcast_msg.go:27: Signing and broadcasting message flow
2024/03/02 19:57:20 sign_and_broadcast_msg.go:69: Broadcasting...
2024/03/02 19:57:20 sign_and_broadcast_msg.go:80: Transaction Hash: FF6B5AD7657B4B60A8AFBB2AE734178EC94AD4D48B8C57F1B1E94A9C340BCBB8
2024/03/02 19:57:20 sign_and_broadcast_msg.go:82: Transaction failed reason: []
2024/03/02 19:57:26 sign_and_broadcast_msg.go:87: Verifing...
2024/03/02 19:57:26 sign_and_broadcast_msg.go:98: Transaction Success...
2024/03/02 19:57:26 sign_and_broadcast_msg.go:103: Gas Used: 106734
2024/03/02 19:57:26 cosigner.go:128: Number of authenticators post: 1
2024/03/02 19:57:26 cosigner.go:132: Added authenticator
2024/03/02 19:57:26 cosigner.go:135: Add authenticator completed.
```

#### Step 2: Swap Tokens Using the Cosigner
```
2024/03/02 19:57:26 cosigner_flow.go:46: Starting swap flow
2024/03/02 19:57:26 swap.go:38: Starting token swap...
2024/03/02 19:57:26 swap.go:72: Initial balance of ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2: 90
2024/03/02 19:57:26 sign_and_broadcast_msg.go:27: Signing and broadcasting message flow
2024/03/02 19:57:26 sign_and_broadcast_msg.go:69: Broadcasting...
2024/03/02 19:57:26 sign_and_broadcast_msg.go:80: Transaction Hash: C79D01E055F1EC343AA0C8E9AD376CB37068AA663E09E06A3D0335B73122CEDA
2024/03/02 19:57:26 sign_and_broadcast_msg.go:82: Transaction failed reason: []
2024/03/02 19:57:32 sign_and_broadcast_msg.go:87: Verifing...
2024/03/02 19:57:32 sign_and_broadcast_msg.go:98: Transaction Success...
2024/03/02 19:57:32 sign_and_broadcast_msg.go:103: Gas Used: 194130
2024/03/02 19:57:32 swap.go:110: Post-swap balance of ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2: 105
2024/03/02 19:57:32 swap.go:114: Balance of toToken has increased after the swap.
2024/03/02 19:57:32 swap.go:119: Token swap completed.
```

#### Step 3: Remove the Authenticator
```
2024/03/02 21:00:48 cosigner_flow.go:67: Removing cosigner authenticator
2024/03/02 21:00:48 remove_authenticator.go:43: Querying authenticators for account osmo1mveeh0ruel03usw3k4agxf68l5dmltyemgdlsk
2024/03/02 21:00:48 remove_authenticator.go:44: Number of authenticators: 1
2024/03/02 21:00:48 remove_authenticator.go:51: Removing authenticator for account osmo1mveeh0ruel03usw3k4agxf68l5dmltyemgdlsk
2024/03/02 21:00:48 sign_and_broadcast_msg.go:27: Signing and broadcasting message flow
2024/03/02 21:00:48 sign_and_broadcast_msg.go:69: Broadcasting...
2024/03/02 21:00:48 sign_and_broadcast_msg.go:80: Transaction Hash: F86D6FBAEF3A5C1C44C1812408E05B08B6F6CA08DF0AC78AA28542D9373FFC4B
2024/03/02 21:00:48 sign_and_broadcast_msg.go:82: Transaction failed reason: []
2024/03/02 21:00:54 sign_and_broadcast_msg.go:87: Verifing...
2024/03/02 21:00:54 sign_and_broadcast_msg.go:98: Transaction Success...
2024/03/02 21:00:54 sign_and_broadcast_msg.go:103: Gas Used: 64883
2024/03/02 21:00:54 remove_authenticator.go:72: Number of authenticators post: 0
2024/03/02 21:00:54 remove_authenticator.go:76: Removed authenticator
2024/03/02 21:00:54 remove_authenticator.go:79: Remove authenticator completed.
```

#### Helpful Commands

Query Authenticators For Account:
```
osmosisd query authenticator authenticators osmo1mveeh0ruel03usw3k4agxf68l5dmltyemgdlsk --node http://164.92.247.225:26657 --output json | jq
```

Query Bank Balance For Account: 
```
osmosisd query bank balances osmo1mveeh0ruel03usw3k4agxf68l5dmltyemgdlsk --node http://164.92.247.225:26657 --output json | jq
```
