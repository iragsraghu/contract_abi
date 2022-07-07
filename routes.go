package main

import (
	"io/ioutil"
	"strings"

	"contract_abi/helpers"
	"contract_abi/inputs"

	"github.com/gin-gonic/gin"
)

// root path
func indexPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Running up.....",
	})
}

// for getting encoded data from source code
func contractSourceCode(c *gin.Context) {
	// get input data from user request
	inputData, currProtocol, requiredData, errors := inputs.GetInputData(c)
	if errors != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"error":       "Error while getting input data : " + strings.Join(errors, ", "),
		})
		return
	}

	// currProtocol, currAction := helpers.GetProtocolsData(protocol_data, inputData)
	if currProtocol.ContractAddress == "" {
		c.JSON(400, gin.H{
			"status_code": 400,
			"error":       "ABI file is not found",
		})
		return
	}

	// get abi data from abi file name
	abi_data, err := ioutil.ReadFile("ABI/" + currProtocol.ContractAddress + ".abi")
	if err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"error":       currProtocol.ContractAddress + " Error reading abi file",
		})
		return
	}

	// get encode data
	helpers.GetEncodeData(c, string(abi_data), inputData, currProtocol, requiredData)
}
