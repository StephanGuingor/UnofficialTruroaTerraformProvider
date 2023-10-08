package truora

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	truora "terraform-provider-truora/truora/client"
)

func dataSourceFlow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowRead,
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lang": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable_desktop_flow": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"identity_verifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"verification_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						},
						"steps": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"step_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expected_inputs": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"response_options": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"alias": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*truora.TruoraClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	flowID := d.Get("flow_id").(string)

	flow, err := c.GetFlow(flowID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", flow.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("creation_date", flow.CreationDate.Format("2006-01-02T15:04:05.000Z")); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("update_date", flow.UpdateDate.Format("2006-01-02T15:04:05.000Z")); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("version", flow.Version); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("type", flow.Type); err != nil {
		return diag.FromErr(err)
	}

	if flow.Config != nil {
		mapConfig := mapFlowConfig(flow.Config)

		if err = d.Set("config", mapConfig); err != nil {
			return diag.FromErr(err)
		}
	}

	if flow.IdentityVerifications != nil {
		mapVerification := mapFlowIdentityVerifications(flow.IdentityVerifications)

		if err = d.Set("identity_verifications", mapVerification); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(flow.FlowID)

	return diags
}

func mapFlowConfig(config *truora.IdentityFlowConfig) []interface{} {
	mapConfig := make(map[string]interface{})
	mapConfig["lang"] = config.Lang
	mapConfig["enable_desktop_flow"] = config.EnableDesktopFlow

	return []interface{}{mapConfig}
}

func mapExpectedInput(expectedInput *truora.Input) map[string]interface{} {
	mapExpectedInput := make(map[string]interface{})
	mapExpectedInput["type"] = expectedInput.Type
	mapExpectedInput["name"] = expectedInput.Name

	if expectedInput.ResponseOptions != nil {
		mapResponseOptions := make([]interface{}, len(expectedInput.ResponseOptions))

		for i, responseOption := range expectedInput.ResponseOptions {
			mapResponseOption := make(map[string]interface{})
			mapResponseOption["value"] = responseOption.Value
			mapResponseOption["alias"] = responseOption.Alias
			mapResponseOptions[i] = mapResponseOption
		}

		mapExpectedInput["response_options"] = mapResponseOptions
	}

	return mapExpectedInput
}

func mapStep(step *truora.Step) map[string]interface{} {
	mapStep := make(map[string]interface{})
	mapStep["step_id"] = step.StepID
	mapStep["type"] = step.Type

	if step.ExpectedInputs != nil {
		mapExpectedInputs := make([]interface{}, len(step.ExpectedInputs))

		for i, expectedInput := range step.ExpectedInputs {
			mappedInput := mapExpectedInput(expectedInput)

			mapExpectedInputs[i] = mappedInput
		}

		mapStep["expected_inputs"] = mapExpectedInputs
	}

	return mapStep
}

func mapVerification(verification *truora.IdentityVerification) map[string]interface{} {
	mapVerification := make(map[string]interface{})

	mapVerification["verification_id"] = verification.VerificationID
	mapVerification["name"] = verification.Name

	if verification.Steps != nil {
		mapSteps := make([]interface{}, len(verification.Steps))

		for i, step := range verification.Steps {
			mapStep := mapStep(step)

			mapSteps[i] = mapStep
		}

		mapVerification["steps"] = mapSteps
	}

	return mapVerification
}

func mapFlowIdentityVerifications(identityVerifications []*truora.IdentityVerification) []interface{} {
	mapIdentityVerifications := make([]interface{}, len(identityVerifications))

	for i, identityVerification := range identityVerifications {

		mapIdentityVerification := mapVerification(identityVerification)

		mapIdentityVerifications[i] = mapIdentityVerification
	}

	return mapIdentityVerifications
}
