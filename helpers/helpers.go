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
func GetEncodeData(c *gin.Context, abi_data string, amount float64, input_duration float64, objectData models.RequiredFields) {

	// connecting to the blockchain from the given rpc provider
	web3, err := web3.NewWeb3(objectData.RPC)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error connecting to ethereum network",
		})
		return
	}
	// creating contract object
	contract, err := web3.Eth.NewContract(abi_data, objectData.ContractAddress)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error creating contract from web3",
		})
		return
	}
	bigIntAmount := web3.Utils.ToWei(float64(amount))                   // convert amount to wei with 18 decimals
	bigIntDuration := web3.Utils.ToWei(float64(input_duration))         // convert duration to wei with 18 decimals
	args := RequiredArguments(bigIntAmount, bigIntDuration, objectData) // check arguments for particular action

	c.JSON(200, gin.H{
		"arguments": args,
	})
	// encoding data for particular action
	encoded_data, err := contract.EncodeABI(objectData.Action, args...)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error encoding ABI" + err.Error(),
			"data":  args,
		})
		return
	}

	// converting byte encoded data to hex string
	encodedString := "0x" + hex.EncodeToString(encoded_data)
	gasLimit, err := calculateGasLimit(objectData, encodedString)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error getting gas limit " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"encoded_data": encodedString, // encoded data
		"Data":         objectData,    // protocol data
		"gas":          gasLimit,      // gas limit
		"value":        0,
	})
}

var currentProtocolData config.ProtocolsData

func GetProtocolsData(protocols_data []config.ProtocolsData, protocol string, chain string, user_action string) models.RequiredFields {
	var objectData models.RequiredFields
	for _, protocol_data := range protocols_data {
		if protocol_data.ProtocolName == protocol && protocol_data.ChainName == chain {
			objectData.ContractAddress, objectData.RPC, objectData.Protocol, objectData.Chain = protocol_data.ContractAddress, protocol_data.RPC, protocol_data.ProtocolName, protocol_data.ChainName
			if user_action == "stake" {
				objectData.Action = protocol_data.Stake.Action
			} else if user_action == "unstake" {
				objectData.Action = protocol_data.Unstake.Action
			}
			currentProtocolData = protocol_data
		}
	}
	return objectData
}

// calculate gas limit
func calculateGasLimit(objectData models.RequiredFields, encodedString string) (int, error) {
	// creating transaction object clinet from given rpc provider
	client := ethrpc.New(objectData.RPC)
	gas := int(2100000 * 5)
	gasLimit, err := client.EthEstimateGas(ethrpc.T{
		To:   objectData.ContractAddress,
		Data: encodedString,
		From: objectData.WalletAddress,
		Gas:  gas,
	})
	if err != nil {
		return 0, err
	} else {
		return gasLimit, nil
	}
}

// check arguments for particular action
func RequiredArguments(input_amount *big.Int, input_duration *big.Int, obj models.RequiredFields) []interface{} {
	var args []interface{}
	currProto := currentProtocolData // assign current protocol data to a variable
	if obj.Action == currProto.Stake.Action && currProto.Stake.AmountRequired == "true" && currProto.Stake.DurationRequired == "true" {
		args = append(args, input_amount, input_duration)
	} else if obj.Action == currProto.Stake.Action && currProto.Stake.AmountRequired == "true" && currProto.Stake.DurationRequired == "false" {
		args = append(args, input_amount)
	} else if obj.Action == currProto.Unstake.Action && currProto.Unstake.AmountRequired == "true" && currProto.Unstake.DurationRequired == "false" {
		args = append(args, input_amount)
	} else if obj.Action == currProto.Unstake.Action && currProto.Unstake.AmountRequired == "false" && currProto.Unstake.DurationRequired == "false" {
		args = nil
	}
	return args
}
