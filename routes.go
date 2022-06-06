package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"encoding/hex"

	"github.com/chenzhijie/go-web3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	api_key := os.Getenv("ETH_API_KEY") // get api key from .env file

	contract_address := c.PostForm("contract_address")
	if contract_address == "" {
		c.JSON(400, gin.H{
			"message": "Contract address is required",
			"error":   err,
		})
		return
	}

	user_action := c.PostForm("action")
	if user_action == "" {
		c.JSON(400, gin.H{
			"message": "Action is required",
			"error":   err,
		})
		return
	}

	user_amount := c.PostForm("amount")
	if user_amount == "" {
		c.JSON(400, gin.H{
			"message": "Amount is required",
			"error":   err,
		})
		return
	}

	fmt.Println(user_action, contract_address)

	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", contract_address, api_key)
	fmt.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error reading request body",
		})
		return
	}

	// convert string to json object and then to map of string to interface
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error unmarshalling json",
			"error":   err,
		})
		return
	}

	// get the data from the map
	abi_data := data["result"].(string)
	// abi_code, err := abi.JSON(strings.NewReader(abi_data))
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"message": "Error parsing abi",
	// 		"error":   err,
	// 	})
	// 	return
	// }

	// fmt.Println("abi_code", abi_code)

	checkMethodName(c, contract_address, user_amount, abi_data, user_action)

}

// check the method name
func checkMethodName(c *gin.Context, contract_address string, user_amount string, abi_data string, user_action string) {
	if user_action == "STAKE" {
		commonMethod(c, contract_address, user_amount, abi_data, "enter")
	} else if user_action == "UNSTAKE" {
		commonMethod(c, contract_address, user_amount, abi_data, "leave")
	} else {
		c.JSON(400, gin.H{
			"message": "Invalid action",
		})
	}
}

// open the dao function
func commonMethod(c *gin.Context, contract_address string, user_amount string, abi_data string, user_action string) {
	fmt.Println("user_action", user_action)
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
	encoded_data, err := contract.EncodeABI(user_action, bigIntAmount)
	fmt.Println("encoded_data", encoded_data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error encoding abi",
			"error":   err,
		})
		return
	}
	encodedString := hex.EncodeToString(encoded_data)
	c.JSON(200, gin.H{
		"message":      "Successfully Encoded ABI",
		"encoded data": encodedString,
	})
}
