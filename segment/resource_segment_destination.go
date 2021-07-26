package segment

import (
	"fmt"
	"log"
	"strings"

	"github.com/fenderdigital/segment-apis-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSegmentDestination() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"source_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"configs": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
		},
		Create: resourceSegmentDestinationCreate,
		Read:   resourceSegmentDestinationRead,
		Update: resourceSegmentDestinationUpdate,
		Delete: resourceSegmentDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSegmentDestinationImport,
		},
	}
}

func resourceSegmentDestinationCreate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	destName := r.Get("destination_name").(string)
	connMode := r.Get("connection_mode").(string)
	enabled := r.Get("enabled").(bool)
	configs := r.Get("configs").(*schema.Set)

	dest, err := client.CreateDestination(srcName, destName, connMode, enabled, extractConfigs(configs))
	if err != nil {
		return fmt.Errorf("ERROR Creating Destination!! Source: %q; Destination: %q; err: %v", srcName, destName, err)
	}

	r.SetId(dest.Name)

	return resourceSegmentDestinationRead(r, meta)
}

func resourceSegmentDestinationRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	id := r.Id()
	destName := idToName(id)

	d, err := client.GetDestination(srcName, destName)
	if err != nil {
		return fmt.Errorf("ERROR Reading Destination!! Source: %q; Destination: %q; err: %v", srcName, destName, err)
	}

	r.Set("enabled", d.Enabled)
	r.Set("configs", d.Configs)
	r.Set("connection_mode", d.ConnectionMode)

	return nil
}

func resourceSegmentDestinationUpdate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	configs := r.Get("configs").(*schema.Set)
	enabled := r.Get("enabled").(bool)
	id := r.Id()
	destName := idToName(id)

	_, err := client.UpdateDestination(srcName, destName, enabled, extractConfigs(configs))
	if err != nil {
		return fmt.Errorf("ERROR Updating Destination!! Source: %q; Destination: %q; err: %v", srcName, destName, err)
	}

	return resourceSegmentDestinationRead(r, meta)
}

func resourceSegmentDestinationDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	id := r.Id()
	destName := idToName(id)

	err := client.DeleteDestination(srcName, destName)
	if err != nil {
		return fmt.Errorf("ERROR Deleting Destination!! Source: %q; Destination: %q; err: %v", srcName, destName, err)
	}

	return nil
}

func resourceSegmentDestinationImport(r *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*segment.Client)
	s := strings.SplitN(r.Id(), "/", 2)
	if len(s) != 2 {
		return nil, fmt.Errorf(
			"invalid destination import format: %s (expected <SOURCE-NAME>/<DESTINATION-NAME>)",
			r.Id(),
		)
	}

	srcName := s[0]
	destName := s[1]

	d, err := client.GetDestination(srcName, destName)
	if err != nil {
		return nil, fmt.Errorf("invalid destination: %q; err: %v", r.Id(), err)
	}

	r.SetId(d.Name)
	r.Set("source_name", srcName)
	r.Set("destination_name", destName)
	r.Set("enabled", d.Enabled)
	r.Set("connection_mode", d.ConnectionMode)

	y := make([]interface{}, len(d.Configs))
	for i, v := range d.Configs {
		y[i] = map[string]interface{}{
			"id":    v.Name,
			"name":  v.DisplayName,
			"type":  v.Type,
			"value": fmt.Sprintf("%v", v.Value),
		}
	}

	confs := schema.NewSet(func(val interface{}) int {
		config := val.(map[string]interface{})
		log.Printf("[DEBUG] Found config value %v", config)
		return schema.HashString(config["id"])
	}, y)

	r.Set("configs", confs)

	results := make([]*schema.ResourceData, 1)
	results[0] = r

	return results, nil
}

func extractConfigs(s *schema.Set) []segment.DestinationConfig {
	configs := []segment.DestinationConfig{}

	if s != nil {
		for _, config := range s.List() {
			c := segment.DestinationConfig{
				Name:        config.(map[string]interface{})["id"].(string),
				Type:        config.(map[string]interface{})["type"].(string),
				Value:       config.(map[string]interface{})["value"],
				DisplayName: config.(map[string]interface{})["name"].(string),
			}
			configs = append(configs, c)
		}
	}

	return configs
}
