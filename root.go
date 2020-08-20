package main

import (
	// "fmt"
	// "os"
	//"github.com/spf13/viper"
	// homedir "github.com/mitchellh/go-homedir"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().StringVar(&prvkeyHex, "prvkeyHex", "", "please provide your eth private key.")
}

var prvkeyHex string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "eth-tool",
	Short: "test web3",
	Long:  `batch eth txs test`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
