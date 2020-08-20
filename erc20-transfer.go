package main

import (
	"fmt"
	"strconv"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/spf13/cobra"
)

var withdrawBalance int64
var toAddr string

func init() {
	queryBalanceCmd.AddCommand(transferCmd)
	transferCmd.PersistentFlags().Int64Var(&withdrawBalance, "balance", 1, "transfer blance to home chain")
	transferCmd.PersistentFlags().StringVar(&toAddr, "toAddr", "1111111111111111111111111111111111111111", "home chain addr")
	//transferCmd.PersistentFlags().StringVar(&withdrawCfgFile, "batch", "", "batch items file ")
}

// SetPrvKey creates keyed-transactor with specified private key.
func SetPrvKey(prvkeyHex string) *bind.TransactOpts {
	keyBytes := common.FromHex(prvkeyHex)
	key := crypto.ToECDSAUnsafe(keyBytes)
	return bind.NewKeyedTransactor(key)
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer balance to other addr",
	Run: func(cmd *cobra.Command, args []string) {

		erc20Entity, err := NewErc20(rpcAddrStr, common.HexToAddress(contractAddr))
		if err != nil {
			fmt.Println(err)
			return
		}

		author := SetPrvKey(prvkeyHex)

		client, err := rpc.Dial(rpcAddrStr)
		if err != nil {
			fmt.Println(err)
			return
		}

		blockNum, err := GetEthBlockNumByRpc(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		tokenName, err := erc20Entity.Name(nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		balance, err := erc20Entity.BalanceOf(nil, author.From)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("before withdraw:  balance: %d  my eth addr: %s contract addr: %s blockNum: %d \n",
			balance, author.From, contractAddr, blockNum)

		pendingNonce := 3

		txHashs := []*types.Transaction{}

		for i := 0; i < 10; i++ {
			author.Nonce = big.NewInt(int64(pendingNonce) + int64(i))
			author.GasLimit = 800000
			author.GasPrice = big.NewInt(6)

			txHash, err := erc20Entity.Transfer(author, common.HexToAddress("0x026c94c25066b473d92568ca83950877cb358164"), big.NewInt(100))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("transfer home via relay idx:%d txHash: %s  tokenName: %s ipcDestAddr:%s  withdraw value: %d \n", i, txHash.Hash().String(), tokenName, "0x026c94c25066b473d92568ca83950877cb358164", withdrawBalance)
			txHashs = append(txHashs, txHash)
		}

		for k, v := range txHashs {
			_, _ = k, v
			// ctx := context.Background()
			// receipt, err := bind.WaitMined(ctx, client, v)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// fmt.Printf("idx: %d  txHash %s   receipt.status:%d \n", k, v.Hash().String(), receipt.Status)
		}

		fmt.Printf("before transfer:  balance: %d account addr: %s contract addr: %s blockNum: %d \n",
			balance, accountAddrStr, contractAddr, blockNum)

		// if len( []byte(toAddr)) == 35 || len( []byte(toAddr)) == 36 {
		// } else {
		// 	fmt.Printf("expect ipc addr len  35 or 36 but actual is %v \n", len([]byte(toAddr)))
		// 	return
		// }

		// if !common.IsHexAddress(toAddr) {
		// 	fmt.Println("to addr  err")
		// 	return
		// }
		//
		// fmt.Printf("to addr:  %s \n", common.HexToAddress(toAddr).Hex())
		//
		// txHash, err := a.EthForeignBridge.Transfer(author, common.HexToAddress(toAddr), big.NewInt(withdrawBalance))
		//
		// //a.EthForeignBridge.Transfer(opts, _to, _value)
		// //txHash,err := a.EthForeignBridge.TransferHomeViaRelay(author, []byte(toAddr), big.NewInt(withdrawBalance))
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		//
		// fmt.Printf("transfer  txHash: %s  tokenName: %s total:%d  withdraw value: %d \n", txHash.Hash().String(), tokenName, total, withdrawBalance)
		//
		// conn, err := ethclient.Dial(rpcAddrStr)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		//
		// ctx := context.Background()
		//
		// receipt, err := bind.WaitMined(ctx, conn, txHash)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		//
		// receiptAsJson, err := receipt.MarshalJSON()
		//
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		//
		// fmt.Printf("receipt.Status: %d \n\n receipt: %v", receipt.Status, string(receiptAsJson))

	},
}

func GetEthBlockNumByRpc(r *rpc.Client) (uint64, error) {

	var blockNumHex string
	err := r.Call(&blockNumHex, "eth_blockNumber")
	if err != nil {
		return 0, err
	}

	blockNumInt64, err := strconv.ParseInt(blockNumHex[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return uint64(blockNumInt64), nil
}
