package utravis

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/shuheiktgw/go-travis"
)

func resourceTravisEnvVar() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvVarCreate,
		Read:   resourceEnvVarRead,
		Update: resourceEnvVarUpdate,
		Delete: resourceEnvVarDelete,

		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
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
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
		},
	}
}

func createEnvVarBody(d *schema.ResourceData) *travis.EnvVarBody {
	return &travis.EnvVarBody{
		Name:   d.Get("name").(string),
		Value:  d.Get("value").(string),
		Public: d.Get("public").(bool),
	}
}

func resourceEnvVarCreate(d *schema.ResourceData, meta interface{}) error {
	slug := d.Get("slug").(string)
	body := createEnvVarBody(d)
	config := meta.(*config)

	config.lock()
	envVar, _, err := config.client.EnvVars.CreateByRepoSlug(context.Background(), slug, body)
	config.unlock()
	if err != nil {
		return err
	}

	d.SetId(*envVar.Id)
	return resourceEnvVarRead(d, meta)
}

func resourceEnvVarRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	slug := d.Get("slug").(string)
	config := meta.(*config)

	envVar, _, err := config.client.EnvVars.FindByRepoSlug(context.Background(), slug, id)
	if err != nil {
		e, ok := err.(*travis.ErrorResponse)
		if ok && e.ErrorType == "not_found" {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(*envVar.Id)
	d.Set("name", *envVar.Name)
	if envVar.Value != nil {
		// Value returns nil if the env var is private
		d.Set("value", hashString(*envVar.Value))
	}
	d.Set("public", *envVar.Public)

	return nil
}

func resourceEnvVarUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	slug := d.Get("slug").(string)
	body := createEnvVarBody(d)
	config := meta.(*config)

	config.lock()
	_, _, err := config.client.EnvVars.UpdateByRepoSlug(context.Background(), slug, id, body)
	config.unlock()
	if err != nil {
		return err
	}
	return resourceEnvVarRead(d, meta)
}

func resourceEnvVarDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	slug := d.Get("slug").(string)
	config := meta.(*config)

	config.lock()
	_, err := config.client.EnvVars.DeleteByRepoSlug(context.Background(), slug, id)
	config.unlock()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func hashString(str string) string {
	hash := sha256.Sum256([]byte(str))
	return base64.StdEncoding.EncodeToString(hash[:])
}
