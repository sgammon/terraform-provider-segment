package segment

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fenderdigital/segment-apis-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSegmentTrackingPlan() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"rules": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				StateFunc: func(val interface{}) string {
					s := segment.Rules{}
					json.Unmarshal([]byte(val.(string)), &s)
					result, _ := json.MarshalIndent(s, "", "  ")
					return string(result)
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
		Create: resourceSegmentTrackingPlanCreate,
		Read:   resourceSegmentTrackingPlanRead,
		Delete: resourceSegmentTrackingPlanDelete,
		Update: resourceSegmentTrackingPlanUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceSegmentTrackingPlanImport,
		},
	}
}

func resourceSegmentTrackingPlanCreate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	displayName := r.Get("display_name").(string)
	rules := r.Get("rules").(string)
	s := segment.Rules{}
	json.Unmarshal([]byte(rules), &s)
	fmt.Printf("%+v\n", s)
	trackingPlan, err := client.CreateTrackingPlan(displayName, s)
	if err != nil {
		return fmt.Errorf("ERROR Creating Tracking Plan!! DisplayName: %q; err: %v", displayName, err)
	}

	planName := parseNameID(trackingPlan.Name)
	r.SetId(planName)
	return resourceSegmentTrackingPlanRead(r, meta)
}

func resourceSegmentTrackingPlanRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Id()
	names, err0 := getNameIDs(client)
	if err0 != nil {
		return err0
	}
	if _, ok := names[planName]; !ok {
		r.SetId("")
		return nil
	}
	trackingPlan, err := client.GetTrackingPlan(planName)
	if err != nil {
		return fmt.Errorf("ERROR Reading Tracking Plan!! PlanName: %q; err: %v", planName, err)
	}
	stringRules, err := json.MarshalIndent(trackingPlan.Rules, "", "  ")
	if err != nil {
		return err
	}
	r.Set("display_name", trackingPlan.DisplayName)
	r.Set("rules", string(stringRules))
	r.Set("name", planName)
	return nil
}

func resourceSegmentTrackingPlanDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Id()
	err := client.DeleteTrackingPlan(planName)
	if err != nil {
		return fmt.Errorf("ERROR Deleting Tracking Plan!! PlanName: %q; err: %v", planName, err)
	}

	return nil
}

func resourceSegmentTrackingPlanUpdate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Id()
	rules := r.Get("rules").(string)
	displayName := r.Get("display_name").(string)

	paths := []string{"tracking_plan.display_name", "tracking_plan.rules"}

	s := segment.Rules{}
	json.Unmarshal([]byte(rules), &s)
	updatedPlan := segment.TrackingPlan{
		DisplayName: displayName,
		Rules:       s,
	}
	_, err := client.UpdateTrackingPlan(planName, paths, updatedPlan)
	if err != nil {
		return fmt.Errorf("ERROR Updating Tracking Plan!! PlanName: %q; err: %v", planName, err)
	}
	return resourceSegmentTrackingPlanRead(r, meta)
}

func resourceSegmentTrackingPlanImport(r *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*segment.Client)
	s, err := client.GetTrackingPlan(r.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid tracking plan: %q; err: %v", r.Id(), err)
	}
	stringRules, err := json.Marshal(s.Rules)
	if err != nil {
		return nil, err
	}
	planName := parseNameID(s.Name)
	r.SetId(planName)
	r.Set("name", planName)
	r.Set("display_name", s.DisplayName)
	r.Set("rules", stringRules)

	results := make([]*schema.ResourceData, 1)
	results[0] = r

	return results, nil
}

func parseNameID(name string) string {
	nameSplit := strings.Split(name, "/")
	return nameSplit[len(nameSplit)-1]
}

func getNameIDs(client *segment.Client) (map[string]string, error) {
	plans, err := client.ListTrackingPlans()
	if err != nil {
		return nil, err
	}
	names := make(map[string]string)
	for _, element := range plans.TrackingPlans {
		id := parseNameID(element.Name)
		names[id] = element.Name
	}
	return names, nil
}
