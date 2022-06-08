package main

import (
	"ContractMethodAPI/helpers"
	"fmt"
	"io/ioutil"
	"os"

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
	contract_address := c.PostForm("contract_address")
	if contract_address == "" {
		c.JSON(400, gin.H{
			"message": "Contract address is required",
		})
		return
	}

	user_action := c.PostForm("action")
	if user_action == "" {
		c.JSON(400, gin.H{
			"message": "Action is required",
		})
		return
	}

	user_amount := c.PostForm("amount")
	if user_amount == "" {
		c.JSON(400, gin.H{
			"message": "Amount is required",
		})
		return
	}

	if contract_address == os.Getenv("OPEN_DAO") {
		action := helpers.ActionName(c, user_action) // to get method name
		openDaoArray := []string{"enter", "leave"}
		user_method := slices.Contains(openDaoArray, action) // to check if action is valid
		abi_data, err := ioutil.ReadFile("ABI/OPEN_DAO.abi")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "openDao Error reading abi",
				"error":   err,
			})
			return
		}
		if user_method {
			helpers.CommonMethod(c, contract_address, user_amount, string(abi_data), action, "")
		} else {
			c.JSON(400, gin.H{
				"message": action + " method is not present in the abi",
			})
		}
	} else if contract_address == os.Getenv("PAN_CAKE") {
		action := helpers.ActionName(c, user_action)
		panCakeArray := []string{"deposit", "withdraw"}
		lock_duration := c.PostForm("lock_duration")
		user_method := slices.Contains(panCakeArray, action)
		if lock_duration == "" && action == "deposit" {
			c.JSON(400, gin.H{
				"message": "Lock duration is required for deposit",
			})
			return
		}
		abi_data, err := ioutil.ReadFile("ABI/PAN_CAKE.abi")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "panCake Error reading abi",
				"error":   err,
			})
			return
		}
		if user_method {
			helpers.CommonMethod(c, contract_address, user_amount, string(abi_data), action, lock_duration)
		} else {
			c.JSON(400, gin.H{
				"message": action + " method is not present in the abi",
			})
		}
	} else if contract_address == os.Getenv("BEEFY") {
		action := helpers.ActionName(c, user_action) // to get method name
		beefyArray := []string{"stake", "withdraw"}
		user_method := slices.Contains(beefyArray, action) // to check if action is valid
		abi_data, err := ioutil.ReadFile("ABI/BEEFY.abi")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Matic Error reading abi",
				"error":   err,
			})
			return
		}
		if user_method {
			helpers.CommonMethod(c, contract_address, user_amount, string(abi_data), action, "")
		} else {
			c.JSON(400, gin.H{
				"message": action + " method is not present in the abi",
			})
		}
	}

}
