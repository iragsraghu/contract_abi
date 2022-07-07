package helpers

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"contract_abi/config"

	"contract_abi/models"

	"github.com/chenzhijie/go-web3"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
)

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
			"status_code": 400,
			"error":       "Error connecting to ethereum network",
		})
		return
	}
	// creating contract object
	contract, err := web3.Eth.NewContract(abi_data, currProtocol.ContractAddress)
	if err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"error":       "Error creating contract from web3",
		})
		return
	}

	// float to bigInt conversion
	bigIntAmount := web3.Utils.ToWei(float64(inputData.Amount)) // convert amount to wei with 18 decimals
	bigFloatAmount := web3.Utils.FromWei(bigIntAmount)
	bigIntDuration := web3.Utils.ToWei(float64(inputData.Duration)) // convert duration to wei with 18 decimals

	args := EncodeArguments(bigIntAmount, bigIntDuration, requiredData) // check arguments for particular action

	// encoding data for particular action
	encoded_data, err := contract.EncodeABI(requiredData[0], args...)
	if err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"error":       "Error encoding ABI " + err.Error(),
		})
		return
	}

	// converting byte encoded data to hex string
	encodedString := "0x" + hex.EncodeToString(encoded_data)
	gasLimit, err := calculateGasLimit(web3, currProtocol, encodedString, bigFloatAmount)
	if err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"error":       err.Error(),
		})
		return
	}

	// creating output data for the response
	result := models.OutputFields{
		EndcodedData: encodedString,
		From:         currProtocol.WalletAddress,
		To:           currProtocol.ContractAddress,
		Gas:          gasLimit,
	}
	c.JSON(200, gin.H{
		"status_code": 200,
		"data":        result, // contains encoded data, from address, to address, gas limit
	})
}

// get matched protocol and required fields
func GetProtocolData(protocol string, chain string, action string) (config.ProtocolData, [3]string) {
	var currProtocol config.ProtocolData // current protocol data
	var reqData [3]string                // required data for particular action

	// looping through protocols
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
func calculateGasLimit(web3 *web3.Web3, currProtocol config.ProtocolData, encodedString string, bigFloatAmount *big.Float) (int, error) {
	client := ethrpc.New(currProtocol.RPC)
	gas := int(2100000 * 5)

	// getting balance of the wallet
	balance, err := client.EthGetBalance(currProtocol.WalletAddress, "latest")
	if err != nil {
		return 0, fmt.Errorf("error getting balance of the wallet %v", err)
	}

	// convert to wei
	bigIntBalance := new(big.Int).SetUint64(balance.Uint64())
	bigFloatBalance := web3.Utils.FromWei(bigIntBalance)

	// check if balance is sufficient
	if bigFloatBalance.Cmp(bigFloatAmount) < 0 {
		return 0, fmt.Errorf("insufficient balance")
	}

	// getting gas price
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
