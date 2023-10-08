package truora

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	truora "terraform-provider-truora/truora/client"
)

func parseCustomMessages(customMessages interface{}) []*truora.CustomFinalMessage {
	if customMessagesSlice, ok := customMessages.([]interface{}); ok {

		customFinalMessages := make([]*truora.CustomFinalMessage, len(customMessagesSlice))
		for i, customMessageMap := range customMessagesSlice {
			if customMessage, ok := customMessageMap.(map[string]interface{}); ok {
				customFinalMessages[i] = &truora.CustomFinalMessage{
					Message: getStringFromMap(customMessage, "message"),
					Status:  getStringFromMap(customMessage, "status"),
				}
			}
		}
		return customFinalMessages
	}
	return nil
}
func parseMessages(messagesMap map[string]interface{}) *truora.Messages {
	if messagesMap == nil {
		return nil
	}

	message := truora.Messages{
		FailureMessage:           getStringFromMap(messagesMap, "failure_message"),
		SuccessMessage:           getStringFromMap(messagesMap, "success_message"),
		PendingMessage:           getStringFromMap(messagesMap, "pending_message"),
		ExitMessage:              getStringFromMap(messagesMap, "exit_message"),
		WaitingForResultsMessage: getStringFromMap(messagesMap, "waiting_for_results_message"),
		CustomMessages:           parseCustomMessages(messagesMap["custom_messages"]),
	}
	return &message
}

func parseFlowConfig(data map[string]interface{}) *truora.IdentityFlowConfig {
	flowConfig := &truora.IdentityFlowConfig{
		EnableDesktopFlow:       getBoolFromMap(data, "enable_desktop_flow"),
		Lang:                    getStringFromMap(data, "lang"),
		ContinueFlowInNewDevice: getBoolFromMap(data, "continue_flow_in_new_device"),
		EnableFollowUp:          getBoolFromMap(data, "enable_follow_up"),
		FollowUpDelay:           getInt64FromMap(data, "follow_up_delay"),
	}

	flowConfig.FollowUpMessage = getStringFromMap(data, "follow_up_message")
	flowConfig.StartBusinessHours = getTimeFromMap(data, "start_business_hours")
	flowConfig.EndBusinessHours = getTimeFromMap(data, "end_business_hours")

	messagesMap := getMapFromSingleElementList(data, "messages")
	flowConfig.Messages = parseMessages(messagesMap)

	return flowConfig
}

func parseVerifications(data []interface{}) []*truora.IdentityVerification {
	verifications := make([]*truora.IdentityVerification, len(data))

	for i, verificationMap := range data {
		verification := parseVerification(verificationMap.(map[string]interface{}))
		verifications[i] = verification
	}

	return verifications
}

func parseVerification(verificationMap map[string]interface{}) *truora.IdentityVerification {
	verification := &truora.IdentityVerification{
		Name: getStringFromMap(verificationMap, "name"),
	}

	verification.Logic = parseStringArray(verificationMap, "logic")

	parseSteps(verificationMap, verification)

	return verification
}

func parseSteps(verificationMap map[string]interface{}, verification *truora.IdentityVerification) {
	if v, ok := verificationMap["steps"]; ok {
		stepsMapList := v.([]interface{})
		steps := make([]*truora.Step, len(stepsMapList))

		for j, stepMap := range stepsMapList {
			step := parseStep(stepMap.(map[string]interface{}))
			steps[j] = step
		}

		verification.Steps = steps
	}
}

func parseStep(stepMap map[string]interface{}) *truora.Step {
	step := truora.Step{
		Type:        getStringFromMap(stepMap, "type"),
		Title:       getStringFromMap(stepMap, "title"),
		Description: getStringFromMap(stepMap, "description"),
	}

	parseExpectedInputs(stepMap, &step)

	return &step
}

func parseExpectedInputs(stepMap map[string]interface{}, step *truora.Step) {
	if v, ok := stepMap["expected_inputs"]; ok {
		expectedInputs := v.([]interface{})
		step.ExpectedInputs = make([]*truora.Input, len(expectedInputs))

		for j, expectedInputMap := range expectedInputs {
			expectedInput := parseExpectedInput(expectedInputMap.(map[string]interface{}))
			step.ExpectedInputs[j] = expectedInput
		}
	}
}

func parseExpectedInput(expectedInputMap map[string]interface{}) *truora.Input {
	expectedInput := truora.Input{
		Type: getStringFromMap(expectedInputMap, "type"),
		Name: getStringFromMap(expectedInputMap, "name"),
	}

	parseResponseOptions(expectedInputMap, &expectedInput)

	return &expectedInput
}

func parseResponseOptions(expectedInputMap map[string]interface{}, expectedInput *truora.Input) {
	if v, ok := expectedInputMap["response_options"]; ok {
		responseOptions := v.([]interface{})
		expectedInput.ResponseOptions = make([]*truora.ResponseOption, len(responseOptions))

		for k, responseOptionMap := range responseOptions {
			responseOption := parseResponseOption(responseOptionMap.(map[string]interface{}))
			expectedInput.ResponseOptions[k] = responseOption
		}
	}
}

func parseResponseOption(responseOptionMap map[string]interface{}) *truora.ResponseOption {
	return &truora.ResponseOption{
		Value: responseOptionMap["value"].(string),
		Alias: responseOptionMap["alias"].(string),
	}
}

func flowFromResourceData(d *schema.ResourceData) *truora.IdentityProcessFlow {
	flow := truora.IdentityProcessFlow{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	if v, ok := d.GetOk("config.0"); ok {
		mapConfig := v.(map[string]interface{})
		flow.Config = parseFlowConfig(mapConfig)
	}

	if v, ok := d.GetOk("verification"); ok {
		verificationsData := v.([]interface{})
		verifications := parseVerifications(verificationsData)
		flow.IdentityVerifications = verifications
	}

	return &flow
}
