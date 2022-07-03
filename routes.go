package main

import (
	"io/ioutil"
	"strings"

	"ContractMethodAPI/helpers"
	"ContractMethodAPI/inputs"

	"github.com/gin-gonic/gin"
)

func indexPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Running up.....",
	})
}

func contractSourceCode(c *gin.Context) {
	// loading yaml file
	// var protocol_data = config.LoadProtocol().Protocols.ProtocolData

	// get input data from user request
	inputData, currProtocol, requiredData, errors := inputs.GetInputData(c)
	if errors != nil {
		c.JSON(400, gin.H{
			"error": "Error while getting input data : " + strings.Join(errors, ", "),
		})
		return
	}

	// currProtocol, currAction := helpers.GetProtocolsData(protocol_data, inputData)
	if currProtocol.ContractAddress == "" {
		c.JSON(400, gin.H{
			"error": "ABI file is not found",
		})
		return
	}

	// get abi data from abi file name
	abi_data, err := ioutil.ReadFile("ABI/" + currProtocol.ContractAddress + ".abi")
	if err != nil {
		c.JSON(400, gin.H{
			"error": currProtocol.ContractAddress + " Error reading abi file",
		})
		return
	}

	// get encode data
	helpers.GetEncodeData(c, string(abi_data), inputData, currProtocol, requiredData)
}
