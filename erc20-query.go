package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"

	"github.com/wade-liwei/go-web3-examples/erc20"
)

var rpcAddrStr, accountAddrStr, contractAddr string

func init() {
	RootCmd.AddCommand(queryBalanceCmd)
	queryBalanceCmd.PersistentFlags().StringVar(&rpcAddrStr, "rpc", "http://127.0.0.1:8545", "eth node addr")
	queryBalanceCmd.PersistentFlags().StringVar(&accountAddrStr, "addr", "0x4cb75889e2918954a63853af1a2ba2a5bc7c5f2d", "eth address")
	queryBalanceCmd.PersistentFlags().StringVar(&contractAddr, "ctct", "0xad4d74d4ad72e92523b45a1cf64c6f027ce37ff9", "eth contract addr")
}

var queryBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "query addr balance",
	Run: func(cmd *cobra.Command, args []string) {

		erc20Entity, err := NewErc20(rpcAddrStr, common.HexToAddress(contractAddr))
		if err != nil {
			fmt.Println(err)
			return
		}

		balance, err := erc20Entity.BalanceOf(nil, common.HexToAddress(accountAddrStr))
		//  balacne, err :=a.EthForeignBridge.Balances(opts, arg0)
		if err != nil {
			fmt.Println(err)
			return
		}

		tokenName, err := erc20Entity.Name(nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		total, err := erc20Entity.TotalSupply(nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		decimal, err := erc20Entity.Decimals(nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		//
		// client, err := rpc.Dial(rpcAddrStr)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// a.EthAuthor = app.SetPrvKey(prvkeyHex)

		fmt.Printf("balance: %d decimal: %d total:%d account addr: %s tokenName: %s contract addr: %s \n",
			balance, decimal, total, accountAddrStr, tokenName, contractAddr)
	},
}

//query contract  e.g. totalSupply
//and write without permission  e.g. transfer
func NewErc20(url string, contractAddr common.Address) (*erc20.Erc20, error) {
	//	a.logger.Debugf("new eth cli param: \nurl:\t\t%s\ncontAddr:\t%s \n", url, contractAddr.String())
	conn, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	erc20, err := erc20.NewErc20(contractAddr, conn)
	if err != nil {
		return nil, err
	}
	return erc20, nil
}
