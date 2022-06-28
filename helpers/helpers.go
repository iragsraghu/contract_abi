package helpers

import (
	"encoding/hex"
	"math/big"
	"strconv"

	"ContractMethodAPI/config"

	"ContractMethodAPI/models"

	"github.com/chenzhijie/go-web3"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
)

// const rpcProviderURL = "https://mainnet.infura.io/v3/7ba7186d11d24eddbf53996feb6dbabf"

// get user amount from string
func ConvertStringToFloat(input_data string) float64 {
	// input_data_int, err := strconv.Atoi(input_data)
	input_data_float, err := strconv.ParseFloat(input_data, 64)
	if err != nil {
		return 0
	}
	return input_data_float
}

// func GetEncodeData(c *gin.Context, contract_address string, abi_data string, action_name string, amount float64, input_duration float64, lock_duration_exists bool, rpcProviderURL string, from_address string, chain_name string, protocol string) {
func GetEncodeData(c *gin.Context, abi_data string, inputData models.InputFields, currProtocol config.ProtocolData, currAction string) {

	// connecting to the blockchain from the given rpc provider
	web3, err := web3.NewWeb3(currProtocol.RPC)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error connecting to ethereum network",
		})
		return
	}
	// creating contract object
	contract, err := web3.Eth.NewContract(abi_data, currProtocol.ContractAddress)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error creating contract from web3",
		})
		return
	}
	bigIntAmount := web3.Utils.ToWei(float64(inputData.Amount))                       // convert amount to wei with 18 decimals
	bigIntDuration := web3.Utils.ToWei(float64(inputData.Duration))                   // convert duration to wei with 18 decimals
	args := RequiredArguments(bigIntAmount, bigIntDuration, currProtocol, currAction) // check arguments for particular action

	// encoding data for particular action
	encoded_data, err := contract.EncodeABI(currAction, args...)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error encoding ABI " + err.Error(),
			"data":  args,
		})
		return
	}

	// converting byte encoded data to hex string
	encodedString := "0x" + hex.EncodeToString(encoded_data)
	gasLimit, err := calculateGasLimit(currProtocol, encodedString)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error getting gas limit " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"encoded_data": encodedString, // encoded data
		"Data":         currProtocol,  // protocol data
		"gas":          gasLimit,      // gas limit
		"value":        0,
	})
}

func GetProtocolsData(protocol_data []config.ProtocolData, inputData models.InputFields) (config.ProtocolData, string) {
	var currProtocol config.ProtocolData // current protocol data
	var currAction string                // current action
	for _, currData := range protocol_data {
		if currData.ProtocolName == inputData.Protocol && currData.ChainName == inputData.Chain {
			currData.WalletAddress = inputData.FromAddress
			if inputData.Action == "stake" {
				currAction = currData.Stake.Action
			} else if inputData.Action == "unstake" {
				currAction = currData.Unstake.Action
			}
			currProtocol = currData
		}
	}
	return currProtocol, currAction
}

// calculate gas limit
func calculateGasLimit(currProtocol config.ProtocolData, encodedString string) (int, error) {
	// creating transaction object clinet from given rpc provider
	client := ethrpc.New(currProtocol.RPC)
	gas := int(2100000 * 5)
	gasLimit, err := client.EthEstimateGas(ethrpc.T{
		To:   currProtocol.ContractAddress,
		Data: encodedString,
		From: currProtocol.WalletAddress,
		Gas:  gas,
	})
	if err != nil {
		return 0, err
	} else {
		return gasLimit, nil
	}
}

// check arguments for particular action
func RequiredArguments(input_amount *big.Int, input_duration *big.Int, currProtocol config.ProtocolData, currAction string) []interface{} {
	var args []interface{}
	if currAction == currProtocol.Stake.Action && currProtocol.Stake.AmountRequired == "true" && currProtocol.Stake.DurationRequired == "true" {
		args = append(args, input_amount, input_duration)
	} else if currAction == currProtocol.Stake.Action && currProtocol.Stake.AmountRequired == "true" && currProtocol.Stake.DurationRequired == "false" {
		args = append(args, input_amount)
	} else if currAction == currProtocol.Unstake.Action && currProtocol.Unstake.AmountRequired == "true" && currProtocol.Unstake.DurationRequired == "false" {
		args = append(args, input_amount)
	} else if currAction == currProtocol.Unstake.Action && currProtocol.Unstake.AmountRequired == "false" && currProtocol.Unstake.DurationRequired == "false" {
		args = nil
	}
	return args
}
