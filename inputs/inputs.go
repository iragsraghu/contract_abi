package inputs

import (
	"contract_abi/config"
	"contract_abi/helpers"
	"contract_abi/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetInputData(c *gin.Context) (models.InputFields, config.ProtocolData, [3]string, []string) {
	var inputData models.InputFields     // create object of input data
	var currProtocol config.ProtocolData // create object of current protocol
	var requiredData [3]string           // create array of request data
	var errorStrings []string            // create array of errorStrings

	// get from address from request
	from_address := c.PostForm("from_address")
	if from_address == "" {
		err1 := fmt.Errorf("from_address is required")
		errorStrings = append(errorStrings, err1.Error())
	}

	// get action name from user
	user_action := c.PostForm("action")
	if user_action == "" {
		err2 := fmt.Errorf("action is required")
		errorStrings = append(errorStrings, err2.Error())
	}

	// get chain name from user
	chain := c.PostForm("chain")
	if chain == "" {
		err3 := fmt.Errorf("chain is required")
		errorStrings = append(errorStrings, err3.Error())
	}

	// get protocol name from user
	protocol := c.PostForm("protocol")
	if protocol == "" {
		err4 := fmt.Errorf("protocol is required")
		errorStrings = append(errorStrings, err4.Error())
	}

	// getting current protocol data, required data and error strings
	currProtocol, requiredData = helpers.GetProtocolData(protocol, chain, user_action)
	// updating wallet address in current protocol data
	currProtocol.WalletAddress = from_address

	var input_amount float64 // create variable of input amount of type float64
	// checking if amount is present in request
	if requiredData[1] == "true" {
		// get amount from user
		amount := c.PostForm("amount")
		if amount == "" {
			err5 := fmt.Errorf("amount is required")
			errorStrings = append(errorStrings, err5.Error())
		}
		// get convert amount from string to int
		input_amount, _ = strconv.ParseFloat(amount, 64)
	}

	var input_duration float64 // create variable of input duration of type float64
	// checking if duration is present in request
	if requiredData[2] == "true" {
		lock_duration := c.PostForm("duration")
		if lock_duration == "" {
			err6 := fmt.Errorf("duration is required")
			errorStrings = append(errorStrings, err6.Error())
		}
		input_duration = helpers.ConvertStringToFloat(lock_duration)
	}

	// updating input data in object
	inputData = models.InputFields{
		FromAddress: from_address,
		Amount:      input_amount,
		Duration:    input_duration,
		Action:      user_action,
		Chain:       chain,
		Protocol:    protocol,
	}

	// checking if error strings is empty or not
	if errorStrings != nil {
		return inputData, currProtocol, requiredData, errorStrings
	} else {
		return inputData, currProtocol, requiredData, nil
	}
}
