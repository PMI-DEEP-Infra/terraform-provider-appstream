package appstream

import (
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/service/appstream"
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
    	"log"
	"strings"
	"time"
)


func resourceAppstreamFleet() *schema.Resource {
	return &schema.Resource {
        Create: resourceAppstreamFleetCreate,
        Read:   resourceAppstreamFleetRead,
        Update: resourceAppstreamFleetUpdate,
        Delete: resourceAppstreamFleetDelete,
        Importer: &schema.ResourceImporter {
            State: schema.ImportStatePassthrough,
        },

        Schema: map[string]*schema.Schema{
            "compute_capacity": {
                Type:         schema.TypeList,
                Required:     true,
                Elem:         &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "desired_instances": {
                            Type:       schema.TypeInt,
                            Required:   true,
                        },
                    },
                },
            },

            "description": {
                Type:         schema.TypeString,
                Optional:     true,
            },

            "disconnect_timeout": {
                Type:         schema.TypeInt,
                Optional:     true,
            },

            "display_name": {
                Type:         schema.TypeString,
                Optional:     true,
            },

            "domain_info": {
                Type:         schema.TypeList,
                Optional:     true,
                Elem:         &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "directory_name": {
                            Type:       schema.TypeString,
                            Optional:   true,
                        },
                        "organizational_unit_distinguished_name": {
                            Type:       schema.TypeString,
                            Optional:   true,
                        },
                    },
                },
            },

            "enable_default_internet_access": {
                Type:         schema.TypeBool,
                Optional:     true,
            },

            "fleet_type": {
                Type:         schema.TypeString,
                Optional:     true,
            },
			
			"image_arn": {
                Type:         schema.TypeString,
				Required:     true,
				ForceNew:	  true,
            },

            "instance_type": {
                Type:         schema.TypeString,
                Required:     true,
            },

            "max_user_duration": {
                Type:         schema.TypeInt,
                Optional:     true,
            },

            "name": {
                Type:       schema.TypeString,
                Required:   true,
            },

	    "stack_name": {
                Type:       schema.TypeString,
                Optional:   true,
            },

	    "state": {
		    Type:	schema.TypeString,
		    Optional:	true,
	    },

            "vpc_config": {
                Type:         schema.TypeList,
                Optional:     true,
                Elem:         &schema.Resource {
                    Schema: map[string]*schema.Schema {
                        "security_group_ids": {
                            Type:       schema.TypeString,
                            Optional:   true,
                        },
                        "subnet_ids":   {
                            Type:       schema.TypeString,
                            Optional:   true,
                        },
                    },
                },
            },
	    "tags": {
		    Type:	schema.TypeMap,
		    Optional:	true,
	    },
        },
    }
}

func resourceAppstreamFleetCreate(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn
	CreateFleetInputOpts := &appstream.CreateFleetInput{}

	if v, ok := d.GetOk("name"); ok {
		CreateFleetInputOpts.Name = aws.String(v.(string))
	}

	ComputeConfig := &appstream.ComputeCapacity{}

	if a, ok := d.GetOk("compute_capacity"); ok {
		ComputeAttributes := a.([]interface{})
		attr := ComputeAttributes[0].(map[string]interface{})
		if v, ok := attr["desired_instances"]; ok {
			ComputeConfig.DesiredInstances = aws.Int64(int64(v.(int)))
		}
		CreateFleetInputOpts.ComputeCapacity = ComputeConfig
	}


	if v, ok := d.GetOk("description"); ok {
		CreateFleetInputOpts.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("disconnect_timeout"); ok {
		CreateFleetInputOpts.DisconnectTimeoutInSeconds = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("display_name"); ok {
		CreateFleetInputOpts.DisplayName = aws.String(v.(string))
	}

	DomainJoinInfoConfig := &appstream.DomainJoinInfo{}

	if dom, ok := d.GetOk("domain_info"); ok {
		DomainAttributes := dom.([]interface{})
		attr := DomainAttributes[0].(map[string]interface{})
		if v, ok := attr["directory_name"]; ok {
			DomainJoinInfoConfig.DirectoryName = aws.String(v.(string))
		}
		if v, ok := attr["organizational_unit_distinguished_name"]; ok {
			DomainJoinInfoConfig.OrganizationalUnitDistinguishedName = aws.String(v.(string))
		}
		CreateFleetInputOpts.DomainJoinInfo = DomainJoinInfoConfig
	}

	if v, ok := d.GetOk("enable_default_internet_access"); ok {
		CreateFleetInputOpts.EnableDefaultInternetAccess = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("fleet_type"); ok {
		CreateFleetInputOpts.FleetType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("image_arn"); ok {
		CreateFleetInputOpts.ImageArn = aws.String(v.(string))
	}


	if v, ok := d.GetOk("instance_type"); ok {
		CreateFleetInputOpts.InstanceType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("max_user_duration"); ok {
		CreateFleetInputOpts.MaxUserDurationInSeconds = aws.Int64(int64(v.(int)))
	}

	VpcConfigConfig := & appstream.VpcConfig{}

	if vpc, ok := d.GetOk("vpc_config"); ok {
        	VpcAttributes := vpc.([]interface{})
        	attr := VpcAttributes[0].(map[string]interface{})
        	if v, ok := attr["security_group_ids"]; ok {
            		strSlice := strings.Split(v.(string), ",")
                	for i, s := range strSlice {
		    		strSlice[i] = strings.TrimSpace(s)
			}
            		VpcConfigConfig.SecurityGroupIds = aws.StringSlice(strSlice)
        	}
        	if v, ok := attr["subnet_ids"]; ok {
            		strSlice := strings.Split(v.(string), ",")
                	for i, s := range strSlice {
		    		strSlice[i] = strings.TrimSpace(s)
			}
            		VpcConfigConfig.SubnetIds = aws.StringSlice(strSlice)
    		}
        	CreateFleetInputOpts.VpcConfig = VpcConfigConfig
    	}

	log.Printf("[DEBUG] Run configuration: %s", CreateFleetInputOpts)
	resp, err := svc.CreateFleet(CreateFleetInputOpts)

	if err != nil {
		log.Printf("[ERROR] Error creating Appstream Fleet: %s", err)
		return err
	}

	log.Printf("[DEBUG] %s", resp)
	time.Sleep(2 * time.Second)
	if v, ok := d.GetOk("tags"); ok {

		data_tags := v.(map[string]interface{})

		attr := make(map[string]string)

		for k, v := range data_tags {
		    attr[k] = v.(string)
		}

		tags := aws.StringMap(attr)

		fleet_name := aws.StringValue(CreateFleetInputOpts.Name)
		get, err := svc.DescribeFleets(&appstream.DescribeFleetsInput {
		    Names:   aws.StringSlice([]string{fleet_name}),
		})
		if err != nil {
		    log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
		    return err
		}
		if get.Fleets == nil {
		    log.Printf("[DEBUG] Apsstream Fleet (%s) not found", d.Id())
		}

		fleetArn := get.Fleets[0].Arn

		tag, err := svc.TagResource(&appstream.TagResourceInput{
		    ResourceArn:    fleetArn,
		    Tags:           tags,
		})
		if err != nil {
		    log.Printf("[ERROR] Error tagging Appstream Stack: %s", err)
		    return err
		}
		log.Printf("[DEBUG] %s", tag)
	}


	if v, ok := d.GetOk("stack_name"); ok {
		AssociateFleetInputOpts := &appstream.AssociateFleetInput{}
		AssociateFleetInputOpts.FleetName = CreateFleetInputOpts.Name
		AssociateFleetInputOpts.StackName = aws.String(v.(string))
		resp, err := svc.AssociateFleet(AssociateFleetInputOpts)
		if err != nil {
			log.Printf("[ERROR] Error associating Appstream Fleet: %s", err)
			return err
		}

		log.Printf("[DEBUG] %s", resp)
	}

	if v, ok := d.GetOk("state"); ok {
		if v == "RUNNING" {
			desired_state := v
			resp, err := svc.StartFleet(&appstream.StartFleetInput{
				Name: CreateFleetInputOpts.Name,
			})

			if err != nil {
				log.Printf("[ERROR] Error satrting Appstream Fleet: %s", err)
				return err
			}
			log.Printf("[DEBUG] %s", resp)

			for {

			    resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
				Names:  aws.StringSlice([]string{*CreateFleetInputOpts.Name}),
			    })

			    if err != nil {
				log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
				return err
			    }

			    curr_state := resp.Fleets[0].State
			    if aws.StringValue(curr_state) == desired_state{
				break
			    }
			    if aws.StringValue(curr_state) != desired_state {
				time.Sleep(20 * time.Second)
				continue
			    }

			}
		}
	}

	d.SetId(*CreateFleetInputOpts.Name)

	return resourceAppstreamFleetRead(d, meta)
}

func resourceAppstreamFleetRead(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn

	resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{})

	if err != nil {
		log.Printf("[ERROR] Error reading Appstream Fleet: %s", err)
		return err
	}
	for _, v := range resp.Fleets {
		if aws.StringValue(v.Name) == d.Get("name") {

			d.Set("name", v.Name)

			if v.ComputeCapacityStatus != nil {
				comp_attr := map[string]interface{}{}
				comp_attr["desired_instances"] = aws.Int64Value(v.ComputeCapacityStatus.Desired)
				d.Set("compute_capacity", comp_attr)
			}

			d.Set("description", v.Description)
			d.Set("display_name", v.DisplayName)
			d.Set("disconnect_timeout", v.DisconnectTimeoutInSeconds)
			d.Set("enable_default_internet_access", v.EnableDefaultInternetAccess)
			d.Set("fleet_type", v.FleetType)
			d.Set("image_arn", v.ImageArn)
			d.Set("instance_type", v.InstanceType)
			d.Set("max_user_duration", v.MaxUserDurationInSeconds)

			if v.VpcConfig != nil {
				vpc_attr := map[string]interface{}{}
				vpc_config_sg := aws.StringValueSlice(v.VpcConfig.SecurityGroupIds)
				vpc_config_sub := aws.StringValueSlice(v.VpcConfig.SubnetIds)
				vpc_attr["security_group_ids"] = aws.String(strings.Join(vpc_config_sg, ","))
				vpc_attr["subnet_ids"] = aws.String(strings.Join(vpc_config_sub, ","))
				d.Set("vpc_config", vpc_attr)
			}
			tg, err := svc.ListTagsForResource(&appstream.ListTagsForResourceInput{
        			ResourceArn: v.Arn,
    			})
    			if err != nil {
				log.Printf("[ERROR] Error listing stack tags: %s", err)
				return err
    			}
			if tg.Tags == nil {
				log.Printf("[DEBUG] Apsstream Stack tags (%s) not found", d.Id())
				return nil
			}

			if len(tg.Tags) > 0 {
				tags_attr := make(map[string]string)
				tags := tg.Tags
				for k, v := range tags {
					tags_attr[k] = aws.StringValue(v)
				}
				d.Set("tags", tags_attr)
			}

			d.Set("state", v.State)

			return nil

		}
	}
	d.SetId("")
	return nil
}

func resourceAppstreamFleetUpdate(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn
    UpdateFleetInputOpts := &appstream.UpdateFleetInput{}

    d.Partial(true)

    if v, ok := d.GetOk("name"); ok {
        UpdateFleetInputOpts.Name = aws.String(v.(string))
    }

    if d.HasChange("description") {
    	d.SetPartial("description")
        log.Printf("[DEBUG] Modify Fleet")
        description :=d.Get("description").(string)
        UpdateFleetInputOpts.Description = aws.String(description)
    }

    if d.HasChange("disconnect_timeout") {
        d.SetPartial("disconnect_timeout")
        log.Printf("[DEBUG] Modify Fleet")
        disconnect_timeout := d.Get("disconnect_timeout").(int)
        UpdateFleetInputOpts.DisconnectTimeoutInSeconds = aws.Int64(int64(disconnect_timeout))
    }

    if d.HasChange("display_name") {
        d.SetPartial("display_name")
        log.Printf("[DEBUG] Modify Fleet")
        display_name :=d.Get("display_name").(string)
        UpdateFleetInputOpts.DisplayName = aws.String(display_name)
    }

	if d.HasChange("image_arn") {
        d.SetPartial("image_arn")
        log.Printf("[DEBUG] Modify Fleet")
        image_arn :=d.Get("image_arn").(string)
        UpdateFleetInputOpts.ImageArn = aws.String(image_arn)
    }

    if d.HasChange("instance_type") {
        d.SetPartial("instance_type")
        log.Printf("[DEBUG] Modify Fleet")
        instance_type := d.Get("instance_type").(string)
        UpdateFleetInputOpts.InstanceType = aws.String(instance_type)
    }

    if d.HasChange("max_user_duration") {
        d.SetPartial("max_user_duration")
        log.Printf("[DEBUG] Modify Fleet")
        max_user_duration :=d.Get("max_user_duration").(int)
        UpdateFleetInputOpts.MaxUserDurationInSeconds = aws.Int64(int64(max_user_duration))
    }

    resp, err := svc.UpdateFleet(UpdateFleetInputOpts)

    if err != nil {
        log.Printf("[ERROR] Error updating Appstream Fleet: %s", err)
	return err
    }
    log.Printf("[DEBUG] %s", resp)
    desired_state := d.Get("state")
    if d.HasChange("state") {
        d.SetPartial("state")
	if desired_state == "STOPPED" {
            svc.StopFleet(&appstream.StopFleetInput{
		    Name: aws.String(d.Id()),
	    })
	    for {

		    resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
			Names:  aws.StringSlice([]string{*UpdateFleetInputOpts.Name}),
		    })
		    if err != nil {
			log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
			return err
		    }

		    curr_state := resp.Fleets[0].State
		    if aws.StringValue(curr_state) == desired_state{
			break
		    }
		    if aws.StringValue(curr_state) != desired_state {
			time.Sleep(20 * time.Second)
			continue
		    }

	    }
        } else if desired_state == "RUNNING" {
            svc.StartFleet(&appstream.StartFleetInput{
		    Name: aws.String(d.Id()),
	    })
	    for {

		    resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
			Names:  aws.StringSlice([]string{*UpdateFleetInputOpts.Name}),
		    })
		    if err != nil {
			log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
			return err
		    }

		    curr_state := resp.Fleets[0].State
		    if aws.StringValue(curr_state) == desired_state{
			break
		    }
		    if aws.StringValue(curr_state) != desired_state {
			time.Sleep(20 * time.Second)
			continue
		    }

	    }
        }
    }
    d.Partial(false)
    return resourceAppstreamFleetRead(d, meta)

}

func resourceAppstreamFleetDelete(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

    resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
	    Names: aws.StringSlice([]string{*aws.String(d.Id())}),
    })

    if err != nil {
	log.Printf("[ERROR] Error reading Appstream Fleet: %s", err)
	return err
    }

    curr_state := aws.StringValue(resp.Fleets[0].State)

    if  curr_state == "RUNNING" {
	    desired_state := "STOPPED"
	    svc.StopFleet(&appstream.StopFleetInput{
		    Name: aws.String(d.Id()),
	    })
	    for {

		    resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
			Names:  aws.StringSlice([]string{*aws.String(d.Id())}),
		    })
		    if err != nil {
			log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
			return err
		    }

		    curr_state := resp.Fleets[0].State
		    if aws.StringValue(curr_state) == desired_state{
			break
		    }
		    if aws.StringValue(curr_state) != desired_state {
			time.Sleep(20 * time.Second)
			continue
		    }

	    }

    }



    dis, err := svc.DisassociateFleet(&appstream.DisassociateFleetInput{
	    FleetName: aws.String(d.Id()),
	    StackName: aws.String(d.Get("stack_name").(string)),
    })
    if err != nil {
        log.Printf("[ERROR] Error deleting Appstream Fleet: %s", err)
	return err
    }
    log.Printf("[DEBUG] %s", dis)

    del, err := svc.DeleteFleet(&appstream.DeleteFleetInput{
        Name:   aws.String(d.Id()),
    })
    if err != nil {
        log.Printf("[ERROR] Error deleting Appstream Fleet: %s", err)
	return err
    }
    log.Printf("[DEBUG] %s", del)
    return nil

}
