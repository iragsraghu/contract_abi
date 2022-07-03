package inputs

import (
	"ContractMethodAPI/config"
	"ContractMethodAPI/helpers"
	"ContractMethodAPI/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetInputData(c *gin.Context) (models.InputFields, config.ProtocolData, [3]string, []string) {
	var inputData models.InputFields     // create object of input data
	var currProtocol config.ProtocolData // create object of current protocol
	var requiredData [3]string           // create array of request data
	var errstrings []string

	// get from address from request
	from_address := c.PostForm("from_address")
	if from_address == "" {
		err1 := fmt.Errorf("from_address is required")
		errstrings = append(errstrings, err1.Error())
	}

	// get action name from user
	user_action := c.PostForm("action")
	if user_action == "" {
		err2 := fmt.Errorf("action is required")
		errstrings = append(errstrings, err2.Error())
	}

	// get chain name from user
	chain := c.PostForm("chain")
	if chain == "" {
		err3 := fmt.Errorf("chain is required")
		errstrings = append(errstrings, err3.Error())
	}

	// get protocol name from user
	protocol := c.PostForm("protocol")
	if protocol == "" {
		err4 := fmt.Errorf("protocol is required")
		errstrings = append(errstrings, err4.Error())
	}

	currProtocol, requiredData = helpers.GetProtocolData(protocol, chain, user_action)
	currProtocol.WalletAddress = from_address

	var input_amount float64
	if requiredData[1] == "true" {
		// get amount from user
		amount := c.PostForm("amount")
		if amount == "" {
			err5 := fmt.Errorf("amount is required")
			errstrings = append(errstrings, err5.Error())
		}
		// get convert amount from string to int
		input_amount = helpers.ConvertStringToFloat(amount)
	}

	// get lock duration from user
	var input_duration float64
	if requiredData[2] == "true" {
		lock_duration := c.PostForm("duration")
		if lock_duration == "" {
			err6 := fmt.Errorf("duration is required")
			errstrings = append(errstrings, err6.Error())
		}
		input_duration = helpers.ConvertStringToFloat(lock_duration)
	}

	inputData = models.InputFields{
		FromAddress: from_address,
		Amount:      input_amount,
		Duration:    input_duration,
		Action:      user_action,
		Chain:       chain,
		Protocol:    protocol,
	}

	if errstrings != nil {
		return inputData, currProtocol, requiredData, errstrings
	} else {
		return inputData, currProtocol, requiredData, nil
	}
}
