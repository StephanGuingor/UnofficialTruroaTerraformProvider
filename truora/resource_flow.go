package truora

import (
	"context"
	"encoding/json"

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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"document": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*truora.TruoraClient)

	jsonFlow := d.Get("document").(string)

	var flow truora.IdentityProcessFlow
	err := json.Unmarshal([]byte(jsonFlow), &flow)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateFlow(ctx, &flow)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.FlowID)

	return resourceFlowRead(ctx, d, m)
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*truora.TruoraClient)

	flowID := d.Id()

	flow, err := client.GetFlow(flowID)
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

	if err = d.Set("name", flow.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*truora.TruoraClient)

	jsonFlow := d.Get("document").(string)

	var flow truora.IdentityProcessFlow
	err := json.Unmarshal([]byte(jsonFlow), &flow)
	if err != nil {
		return diag.FromErr(err)
	}

	flowID := d.Id()

	_, err = client.UpdateFlow(ctx, flowID, &flow)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*truora.TruoraClient)

	flowID := d.Id()

	err := client.DeleteFlow(flowID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
