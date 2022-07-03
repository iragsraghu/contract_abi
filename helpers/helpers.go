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
func GetEncodeData(c *gin.Context, abi_data string, inputData models.InputFields, currProtocol config.ProtocolData, requiredData [3]string) {

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
	bigIntAmount := web3.Utils.ToWei(float64(inputData.Amount))     // convert amount to wei with 18 decimals
	bigIntDuration := web3.Utils.ToWei(float64(inputData.Duration)) // convert duration to wei with 18 decimals

	args := EncodeArguments(bigIntAmount, bigIntDuration, requiredData) // check arguments for particular action

	// encoding data for particular action
	encoded_data, err := contract.EncodeABI(requiredData[0], args...)
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
		"encoded_data": encodedString,                // encoded data
		"from":         currProtocol.WalletAddress,   // from address
		"to":           currProtocol.ContractAddress, // protocol data
		"gas":          gasLimit,                     // gas limit
		"value":        0,
	})
}

func GetProtocolData(protocol string, chain string, action string) (config.ProtocolData, [3]string) {
	var currProtocol config.ProtocolData // current protocol data
	var reqData [3]string
	// var currAction string // current action
	for _, currData := range config.LoadProtocol().Protocols.ProtocolData {
		if currData.ProtocolName == protocol && currData.ChainName == chain {
			if action == "stake" {
				reqData = currData.Stake.RequiredArray
			} else if action == "unstake" {
				reqData = currData.Unstake.RequiredArray
			}
			currProtocol = currData
		}

	}
	return currProtocol, reqData
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
func EncodeArguments(input_amount *big.Int, input_duration *big.Int, requiredData [3]string) []interface{} {
	var args []interface{}
	if requiredData[2] == "true" {
		args = append(args, input_amount, input_duration)
	} else if requiredData[1] == "false" {
		args = nil
	} else {
		args = append(args, input_amount)
	}
	return args
}
