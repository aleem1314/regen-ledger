package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/regen-network/regen-ledger/v6/app"
	"github.com/regen-network/regen-ledger/v6/app/client/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()
	if err := cmd.Execute(rootCmd, app.EnvPrefix, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
