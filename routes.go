package main

import (
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
	// get contract address from user
	contract_address := c.PostForm("contract_address")
	if contract_address == "" {
		c.JSON(400, gin.H{
			"error": "Contract address is required",
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

	// get amount from user
	amount := c.PostForm("amount")
	if amount == "" {
		c.JSON(400, gin.H{
			"error": "Amount is required",
		})
		return
	}
	// get convert amount from string to int
	input_amount := helpers.ConvertStringToInt(amount)
	if input_amount == 0 {
		c.JSON(400, gin.H{
			"error": "Invalid amount",
		})
		return
	}

	// get chain name from contract address
	abi_file_name := mapping.GetChainName(contract_address)
	if abi_file_name == "" {
		c.JSON(400, gin.H{
			"error": "ABI file name is required",
		})
		return
	}

	// get action name from user action
	action_name := helpers.GetActionName(abi_file_name, user_action)
	if action_name == "" {
		c.JSON(400, gin.H{
			"error": "Action Name not valid",
		})
		return
	}

	// get lock duration from user
	var input_duration int
	lock_duration_exists := mapping.GetLockDurationExist(abi_file_name, user_action)
	// lock_duration_exist_file := slices.Contains(lock_duration_exist_chains, abi_file_name)
	if lock_duration_exists {
		lock_duration := c.PostForm("lock_duration")
		if lock_duration == "" {
			c.JSON(400, gin.H{
				"Error": "Lock duration is required",
			})
			return
		}
		input_duration = helpers.ConvertStringToInt(lock_duration)
	}

	// get abi data from abi file name
	abi_data, err := ioutil.ReadFile("ABI/" + abi_file_name + ".abi")
	if err != nil {
		c.JSON(400, gin.H{
			"error": abi_file_name + "Error reading abi file",
		})
		return
	}

	// get encode data from abi data
	helpers.GetEncodeData(c, contract_address, string(abi_data), action_name, input_amount, input_duration, lock_duration_exists)
}
