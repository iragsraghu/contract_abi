package helpers

import (
	"encoding/hex"
	"strconv"

	"ContractMethodAPI/mapping"

	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
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

// getting action name from user action
func GetMappingFields(chain_name string, protocol string, action string, contract_address string, rpcProviderURL string) (string, string, string, string) {
	return mapping.ABIMethodsMapping(chain_name, protocol, action, contract_address, rpcProviderURL)
}

func GetEncodeData(c *gin.Context, contract_address string, abi_data string, action_name string, amount float64, input_duration float64, lock_duration_exists bool, rpcProviderURL string, from_address string, chain_name string) {
	// only for gas estimation purpose
	client := ethrpc.New(rpcProviderURL)
	web3, err := web3.NewWeb3(rpcProviderURL)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error connecting to ethereum network",
		})
		return
	}
	contract, err := web3.Eth.NewContract(abi_data, contract_address)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error creating contract from web3",
		})
		return
	}
	bigIntAmount := web3.Utils.ToWei(float64(amount))           // convert amount to wei with 18 decimals
	bigIntDuration := web3.Utils.ToWei(float64(input_duration)) // convert duration to wei with 18 decimals
	var encoded_data []byte
	if lock_duration_exists {
		// if input_duration == 0 {
		// 	c.JSON(400, gin.H{
		// 		"error": "Duration Should be greater than 0",
		// 	})
		// 	return
		// }
		encoded_data, err = contract.EncodeABI(action_name, bigIntAmount, bigIntDuration)
	} else {
		encoded_data, err = contract.EncodeABI(action_name, bigIntAmount)
	}

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error encoding data for contract address " + contract_address + err.Error(),
		})
		return
	}
	// converting byte encoded data to hex string
	encodedString := hex.EncodeToString(encoded_data)
	token_address := common.HexToAddress(contract_address)
	wallet_address := common.HexToAddress(from_address)
	gas := int(2100000 * 5)
	gasLimit, err := client.EthEstimateGas(ethrpc.T{
		To:   contract_address,
		Data: "0x" + encodedString,
		From: from_address,
		Gas:  gas,
	})
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error getting gas limit " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"from":  wallet_address,
		"data":  encodedString,
		"to":    token_address,
		"chain": chain_name,
		"gas":   gasLimit,
		"value": 0,
	})
}
