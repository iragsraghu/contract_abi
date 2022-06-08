package helpers

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/chenzhijie/go-web3"
	"github.com/gin-gonic/gin"
)

func ActionName(c *gin.Context, action string) string {
	action = strings.ToLower(action)
	switch action {
	case "stake":
		return "stake"
	case "unstake":
		return "leave"
	case "deposit":
		return "deposit"
	case "withdraw":
		return "withdraw"
	default:
		c.JSON(400, gin.H{
			"message": "aaaa Invalid action",
		})
		return "INVALID"
	}
}

func CommonMethod(c *gin.Context, contract_address string, user_amount string, abi_data string, action string, lock_duration string) {
	fmt.Println("user_action", action)
	var rpcProviderURL = "https://mainnet.infura.io/v3/7ba7186d11d24eddbf53996feb6dbabf"
	web3, err := web3.NewWeb3(rpcProviderURL)
	if err != nil {
		fmt.Println("web3 error", err)
	}

	contract, err := web3.Eth.NewContract(abi_data, contract_address)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error creating contract",
			"error":   err,
		})
		return
	}
	amount, err := strconv.Atoi(user_amount)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error converting amount to int",
			"error":   err,
		})
		return
	}
	bigIntAmount := web3.Utils.ToWei(float64(amount))
	if lock_duration != "" {
		duration, err := strconv.Atoi(lock_duration)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Error converting amount to int",
				"error":   err,
			})
			return
		}
		bigIntDuration := web3.Utils.ToWei(float64(duration))
		encoded_data, err := contract.EncodeABI(action, bigIntAmount, bigIntDuration)
		fmt.Println("encoded_data", encoded_data)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "a Error encoding abi",
				"error":   err,
			})
			return
		}
		encodedString := hex.EncodeToString(encoded_data)
		c.JSON(200, gin.H{
			"message":          "Successfully Encoded ABI",
			"encoded data":     encodedString,
			"contract_address": contract_address,
			"method_name":      action,
		})
	} else {
		encoded_data, err := contract.EncodeABI(action, bigIntAmount)
		fmt.Println("encoded_data", encoded_data)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "b Error encoding abi",
				"error":   err,
			})
			return
		}
		encodedString := hex.EncodeToString(encoded_data)
		c.JSON(200, gin.H{
			"message":          "Successfully Encoded ABI",
			"encoded data":     encodedString,
			"contract_address": contract_address,
			"method_name":      action,
		})
	}
}
