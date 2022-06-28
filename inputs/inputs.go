package inputs

import (
	"ContractMethodAPI/helpers"
	"ContractMethodAPI/mapping"
	"ContractMethodAPI/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetInputData(c *gin.Context) (models.InputFields, error) {
	var inputData models.InputFields // create object of input data
	var err error                    // error variable to store any error

	// get from address from request
	from_address := c.PostForm("from_address")
	if from_address == "" {
		err = fmt.Errorf("from address is required")
	}

	// get action name from user
	user_action := c.PostForm("action")
	if user_action == "" {
		err = fmt.Errorf("action is required")
	}

	// get chain name from user
	chain := c.PostForm("chain")
	if chain == "" {
		err = fmt.Errorf("chain is required")
	}

	// get protocol name from user
	protocol := c.PostForm("protocol")
	if protocol == "" {
		err = fmt.Errorf("protocol is required")
	}
	inputData.Protocol = protocol

	// get amount from user
	amount := c.PostForm("amount")
	if amount == "" {
		err = fmt.Errorf("amount is required")
	}

	// get convert amount from string to int
	input_amount := helpers.ConvertStringToFloat(amount)
	fmt.Println("input_amount: ", input_amount)
	if input_amount == 0 {
		err = fmt.Errorf("invalid amount")
	}

	// get lock duration from user
	var input_duration float64
	lock_duration_exists := mapping.GetLockDurationExist(protocol, user_action)
	// lock_duration_exist_file := slices.Contains(lock_duration_exist_chains, abi_file_name)
	if lock_duration_exists {
		lock_duration := c.PostForm("duration")
		if lock_duration == "" {
			err = fmt.Errorf("duration is required")
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

	if err != nil {
		return inputData, err
	} else {
		return inputData, nil
	}
}
