package truora

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func requestVerificationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func requestFlowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
					"continue_flow_in_new_device": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"enable_follow_up": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"follow_up_delay": {
						Type:     schema.TypeInt,
						Optional: true,
						RequiredWith: []string{
							"config.enable_follow_up",
						},
					},
					"follow_up_message": {
						Type:     schema.TypeString,
						Optional: true,
						RequiredWith: []string{
							"config.enable_follow_up",
						},
					},
					"start_business_hours": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"end_business_hours": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"messages": {
						Type:       schema.TypeList,
						Optional:   true,
						MaxItems:   1,
						ConfigMode: schema.SchemaConfigModeBlock,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"failure_message": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"success_message": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"pending_message": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"exit_message": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"waiting_for_results_message": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"custom_messages": {
									Type:     schema.TypeList,
									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"message": {
												Type:     schema.TypeString,
												Required: true,
											},
											"status": {
												Type:     schema.TypeString,
												Required: true,
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
		"verification": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: requestVerificationSchema(),
			},
		},
	}
}

func dataSourceFlowDocument() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFlowDocumentRead,
		Schema: mergeSchemaMaps(
			map[string]*schema.Schema{
				"json": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
			requestFlowSchema(),
		),
	}
}

func dataFlowDocumentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	flow := flowFromResourceData(d)

	flowMarshal, err := json.Marshal(flow)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("json", string(flowMarshal))

	d.SetId(flow.Name)

	return diags
}
