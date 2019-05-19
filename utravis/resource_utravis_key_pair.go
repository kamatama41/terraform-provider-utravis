package utravis

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/shuheiktgw/go-travis"
)

func resourceTravisKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyPairCreate,
		Read:   resourceKeyPairRead,
		Update: resourceKeyPairUpdate,
		Delete: resourceKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				StateFunc: func(value interface{}) string {
					return hashString(value.(string))
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Created by Terraform",
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createKeyPairBody(d *schema.ResourceData) *travis.KeyPairBody {
	return &travis.KeyPairBody{
		Value:       d.Get("value").(string),
		Description: d.Get("description").(string),
	}
}

func resourceKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	slug := d.Get("slug").(string)
	body := createKeyPairBody(d)
	config := meta.(*config)

	_, _, err := config.client.KeyPair.CreateByRepoSlug(context.Background(), slug, body)
	if err != nil {
		return err
	}

	d.SetId(slug)
	return resourceKeyPairRead(d, meta)
}

func resourceKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	slug := d.Id()
	config := meta.(*config)

	keyPair, _, err := config.client.KeyPair.FindByRepoSlug(context.Background(), slug)
	if err != nil {
		e, ok := err.(*travis.ErrorResponse)
		if ok && e.ErrorType == "not_found" {
			d.SetId("")
			return nil
		}
		return err
	}

	// value will be saved only when creating
	d.Set("description", *keyPair.Description)
	d.Set("public_key", *keyPair.PublicKey)
	d.Set("fingerprint", *keyPair.Fingerprint)
	return nil
}

func resourceKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	slug := d.Id()
	body := createKeyPairBody(d)
	config := meta.(*config)

	_, _, err := config.client.KeyPair.UpdateByRepoSlug(context.Background(), slug, body)
	if err != nil {
		return err
	}
	return resourceKeyPairRead(d, meta)
}

func resourceKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	slug := d.Id()
	config := meta.(*config)

	_, err := config.client.KeyPair.DeleteByRepoSlug(context.Background(), slug)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
