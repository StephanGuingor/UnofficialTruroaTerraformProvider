package truora

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	truora "terraform-provider-truora/truora/client"
)

func resourceFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlowCreate,
		ReadContext:   resourceFlowRead,
		UpdateContext: resourceFlowUpdate,
		DeleteContext: resourceFlowDelete,
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Type:     schema.TypeString,
				Computed: true,
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
				Optional: true,
				Default:  "permanent",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:       schema.TypeList,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeBlock,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lang": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "es",
						},
						"enable_desktop_flow": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"verification": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"verification_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"logic": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
										Required: true,
									},
									"expected_inputs": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"response_options": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
															"alias": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"title": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flowFromResourceData(d *schema.ResourceData) *truora.IdentityProcessFlow {
	flow := truora.IdentityProcessFlow{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	if v, ok := d.GetOk("config"); ok {

		configMapList := v.([]interface{})

		mapConfig := configMapList[0].(map[string]interface{})

		flow.Config = &truora.IdentityFlowConfig{
			Lang:              mapConfig["lang"].(string),
			EnableDesktopFlow: mapConfig["enable_desktop_flow"].(bool),
		}
	}

	if v, ok := d.GetOk("verification"); ok {
		verificationsList := v.([]interface{})

		verifications := make([]*truora.IdentityVerification, len(verificationsList))

		for i, verificationMap := range verificationsList {
			verificationMap := verificationMap.(map[string]interface{})

			verification := &truora.IdentityVerification{
				Name: verificationMap["name"].(string),
			}

			if v, ok := verificationMap["logic"]; ok {
				logicList := v.([]interface{})

				verification.Logic = make([]string, len(logicList))

				for i, logic := range logicList {
					verification.Logic[i] = logic.(string)
				}
			}

			if v, ok := verificationMap["steps"]; ok {
				stepsMapList := v.([]interface{})
				steps := make([]*truora.Step, len(stepsMapList))

				for j, stepMap := range stepsMapList {
					stepMap := stepMap.(map[string]interface{})

					step := truora.Step{
						Type:        stepMap["type"].(string),
						Title:       stepMap["title"].(string),
						Description: stepMap["description"].(string),
					}

					if v, ok := stepMap["expected_inputs"]; ok {
						expectedInputs := v.([]interface{})
						step.ExpectedInputs = make([]*truora.Input, len(expectedInputs))

						for j, expectedInputMap := range expectedInputs {
							expectedInputMap := expectedInputMap.(map[string]interface{})

							expectedInput := truora.Input{
								Type: expectedInputMap["type"].(string),
								Name: expectedInputMap["name"].(string),
							}

							if v, ok := expectedInputMap["response_options"]; ok {
								responseOptions := v.([]interface{})
								expectedInput.ResponseOptions = make([]*truora.ResponseOption, len(responseOptions))

								for k, responseOptionMap := range responseOptions {
									responseOptionMap := responseOptionMap.(map[string]interface{})

									responseOption := truora.ResponseOption{
										Value: responseOptionMap["value"].(string),
										Alias: responseOptionMap["alias"].(string),
									}

									expectedInput.ResponseOptions[k] = &responseOption
								}
							}

							step.ExpectedInputs[j] = &expectedInput
						}
					}

					steps[j] = &step
				}

				verification.Steps = steps
			}

			verifications[i] = verification
		}

		flow.IdentityVerifications = verifications
	}

	return &flow
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*truora.TruoraClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	flow := flowFromResourceData(d)

	flowResponse, err := c.CreateFlow(ctx, flow)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(flowResponse.FlowID)

	resourceFlowRead(ctx, d, m)

	return diags
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*truora.TruoraClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	flowID := d.Id()

	flow, err := c.GetFlow(flowID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("flow_id", flow.FlowID); err != nil {
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

		if err = d.Set("verification", mapVerification); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*truora.TruoraClient)

	flowID := d.Id()

	flow := flowFromResourceData(d)

	_, err := client.UpdateFlow(flowID, flow)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*truora.TruoraClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	flowID := d.Id()

	err := client.DeleteFlow(flowID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
