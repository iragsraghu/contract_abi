package helpers

import (
	"encoding/hex"
	"os"
	"strconv"

	"ContractMethodAPI/mapping"

	"github.com/chenzhijie/go-web3"
	"github.com/gin-gonic/gin"
)

const rpcProviderURL = "https://mainnet.infura.io/v3/7ba7186d11d24eddbf53996feb6dbabf"

// get user amount from string
func ConvertStringToInt(amount string) int {
	amount_int, err := strconv.Atoi(amount)
	if err != nil {
		return 0
	}
	return amount_int
}

// getting chain name from contract address
func GetChainName(contract_address string) string {
	return os.Getenv(contract_address)
}

// getting action name from user action
func GetActionName(abi_file_name string, user_action string) string {
	return mapping.MappedMethods(abi_file_name, user_action)
}

func GetEncodeData(c *gin.Context, contract_address string, abi_data string, action_name string, amount int, input_duration int, lock_duration_file bool) {
	web3, err := web3.NewWeb3(rpcProviderURL)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error connecting to ethereum network",
		})
		return
	}
	contract, err := web3.Eth.NewContract(abi_data, contract_address)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error creating contract",
			"error":   err,
		})
		return
	}
	bigIntAmount := web3.Utils.ToWei(float64(amount)) // convert amount to wei with 18 decimals

	bigIntDuration := web3.Utils.ToWei(float64(input_duration)) // convert duration to wei with 18 decimals

	var encoded_data []byte
	if lock_duration_file {
		if input_duration == 0 {
			c.JSON(400, gin.H{
				"message": "Duration Should be greater than 0",
			})
			return
		}
		encoded_data, err = contract.EncodeABI(action_name, bigIntAmount, bigIntDuration)
	} else {
		encoded_data, err = contract.EncodeABI(action_name, bigIntAmount)
	}

	if err != nil {
		c.JSON(400, gin.H{
			"message": "aaaa Error encoding data",
		})
		return
	}

	// converting byte encoded data to hex string
	encodedString := hex.EncodeToString(encoded_data)
	c.JSON(200, gin.H{
		"message":      "Success",
		"encoded_data": encodedString,
	})
}
