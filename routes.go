package main

import (
	"fmt"
	"io/ioutil"

	"ContractMethodAPI/helpers"
	"ContractMethodAPI/mapping"

	"github.com/gin-gonic/gin"
)

func indexPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Running up.....",
	})
}

// how to interact with smart contract on ethereum network using go with ethereum api
func contractSourceCode(c *gin.Context) {
	// get from address from request
	from_address := c.PostForm("from_address")
	if from_address == "" {
		c.JSON(400, gin.H{
			"error": "From address is required",
		})
		return
	}

	// get action name from user
	user_action := c.PostForm("action")
	if user_action == "" {
		c.JSON(400, gin.H{
			"error": "Action is required",
		})
		return
	}

	// get chain name from user
	chain := c.PostForm("chain")
	if chain == "" {
		c.JSON(400, gin.H{
			"error": "Chain name is required",
		})
		return
	}

	// get protocol name from user
	protocol := c.PostForm("protocol")
	if protocol == "" {
		c.JSON(400, gin.H{
			"error": "Protocol is required",
		})
		return
	}

	// get amount from user
	amount := c.PostForm("amount")
	if amount == "" {
		c.JSON(400, gin.H{
			"error": "Amount is required",
		})
		return
	}
	// get convert amount from string to int
	input_amount := helpers.ConvertStringToFloat(amount)
	fmt.Println("input_amount: ", input_amount)
	if input_amount == 0 {
		c.JSON(400, gin.H{
			"error": "Invalid amount",
		})
		return
	}
	// get action name from user action
	action_name, contract_address, rpcProviderURL, chain_name := helpers.GetMappingFields(chain, protocol, user_action, "contract_address", "rpcProviderURL")

	// get lock duration from user
	var input_duration float64
	lock_duration_exists := mapping.GetLockDurationExist(protocol, user_action)
	// lock_duration_exist_file := slices.Contains(lock_duration_exist_chains, abi_file_name)
	if lock_duration_exists {
		lock_duration := c.PostForm("duration")
		if lock_duration == "" {
			c.JSON(400, gin.H{
				"Error": "Lock duration is required",
			})
			return
		}
		input_duration = helpers.ConvertStringToFloat(lock_duration)
		fmt.Println("input_duration", input_duration)
	}

	if contract_address == "" {
		c.JSON(400, gin.H{
			"error": "ABI file is not found",
		})
		return
	}

	// get abi data from abi file name
	abi_data, err := ioutil.ReadFile("ABI/" + contract_address + ".abi")
	if err != nil {
		c.JSON(400, gin.H{
			"error": contract_address + " Error reading abi file",
		})
		return
	}

	// get encode data from abi data
	helpers.GetEncodeData(c, contract_address, string(abi_data), action_name, input_amount, input_duration, lock_duration_exists, rpcProviderURL, from_address, chain_name, protocol)
}
