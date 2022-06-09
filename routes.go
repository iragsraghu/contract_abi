package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"ContractMethodAPI/helpers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

func indexPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Running up.....",
	})
}

// how to interact with smart contract on ethereum network using go with ethereum api
func contractSourceCode(c *gin.Context) {
	err := godotenv.Load() // load .env file
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	// get contract address from user
	contract_address := c.PostForm("contract_address")
	if contract_address == "" {
		c.JSON(400, gin.H{
			"message": "Contract address is required",
		})
		return
	}

	// get action name from user
	user_action := c.PostForm("action")
	if user_action == "" {
		c.JSON(400, gin.H{
			"message": "Action is required",
		})
		return
	}

	// get amount from user
	amount := c.PostForm("amount")
	if amount == "" {
		c.JSON(400, gin.H{
			"message": "Amount is required",
		})
		return
	}
	// get convert amount from string to int
	input_amount := helpers.ConvertStringToInt(amount)
	if input_amount == 0 {
		c.JSON(400, gin.H{
			"message": "Invalid amount",
		})
		return
	}

	// get chain name from contract address
	abi_file_name := helpers.GetChainName(contract_address)
	if abi_file_name == "" {
		c.JSON(400, gin.H{
			"message": "ABI file name is required",
		})
		return
	}

	// get action name from user action
	action_name := helpers.GetActionName(abi_file_name, user_action)
	if action_name == "" {
		c.JSON(400, gin.H{
			"message": "Action Name not valid",
		})
		return
	}

	// get lock duration from user
	var input_duration int
	file_names := os.Getenv("lockDurations")                                  // get files for lock duration from .env file
	lock_duration_files := strings.Split(file_names, ",")                     // split files by comma
	lock_duration_file := slices.Contains(lock_duration_files, abi_file_name) // check if lock file name is in lock duration files
	if lock_duration_file {
		lock_duration := c.PostForm("lock_duration")
		if lock_duration == "" {
			c.JSON(400, gin.H{
				"message": "Lock duration is required",
			})
			return
		}
		input_duration = helpers.ConvertStringToInt(lock_duration)
	}

	// get abi data from abi file name
	abi_data, err := ioutil.ReadFile("ABI/" + abi_file_name + ".abi")
	if err != nil {
		c.JSON(400, gin.H{
			"message": abi_file_name + "Error reading abi file",
		})
		return
	}

	// get encode data from abi data
	helpers.GetEncodeData(c, contract_address, string(abi_data), action_name, input_amount, input_duration, lock_duration_file)
}
