package utravis

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/shuheiktgw/go-travis"
)

func TestAccTravisEnvVar_basic(t *testing.T) {
	var envVar travis.EnvVar
	randString := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("KEY_%s", randString)
	updatedName := fmt.Sprintf("KEY_updated_%s", randString)
	value := "val"
	updatedValue := "val - updated"
	public := false
	updatedPublic := true

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTravisEnvVarDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTravisEnvVarConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTravisEnvVarExists("utravis_env_var.foo", &envVar),
					testAccCheckTravisEnvVarAttributes(&envVar, name, value, public),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "slug", testSlug),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "name", name),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "value", hashString(value)),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "public", strconv.FormatBool(public)),
				),
			},
			{
				Config: testAccTravisEnvVarUpdateConfig(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTravisEnvVarExists("utravis_env_var.foo", &envVar),
					testAccCheckTravisEnvVarAttributes(&envVar, updatedName, updatedValue, updatedPublic),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "slug", testSlug),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "name", updatedName),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "value", hashString(updatedValue)),
					resource.TestCheckResourceAttr("utravis_env_var.foo", "public", strconv.FormatBool(updatedPublic)),
				),
			},
		},
	})
}

func testAccCheckTravisEnvVarExists(n string, envVar *travis.EnvVar) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Env Var ID is set")
		}

		conn := testAccProvider.Meta().(*config).client
		id := rs.Primary.ID

		e, _, err := conn.EnvVars.FindByRepoSlug(context.TODO(), testSlug, id)
		if err != nil {
			return err
		}
		*envVar = *e
		return nil
	}
}

func testAccCheckTravisEnvVarAttributes(envVar *travis.EnvVar, name, value string, public bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *envVar.Name != name {
			return fmt.Errorf("Env name does not match: %s, %s", *envVar.Name, name)
		}

		if public {
			if *envVar.Value != value {
				return fmt.Errorf("Env value does not match: %s, %s", *envVar.Value, value)
			}
		} else {
			// Private env var doesn't return value, so should be nil
			if envVar.Value != nil {
				return fmt.Errorf("Env value is not nil: %s", *envVar.Value)
			}
		}

		if *envVar.Public != public {
			return fmt.Errorf("Env visiblity does not match: %t, %t", *envVar.Public, public)
		}

		return nil
	}
}

func testAccCheckTravisEnvVarDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*config).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "utravis_env_var" {
			continue
		}

		id := rs.Primary.ID
		slug := rs.Primary.Attributes["slug"]
		envVar, resp, err := conn.EnvVars.FindByRepoSlug(context.TODO(), slug, id)
		if err == nil && envVar != nil {
			return fmt.Errorf("env var still exists")
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccTravisEnvVarConfig(name string) string {
	return fmt.Sprintf(`
resource "utravis_env_var" "foo" {
    slug = "%s"
	name = "%s"
	value = "val"
	public = false
}
`, testSlug, name)
}

func testAccTravisEnvVarUpdateConfig(updatedName string) string {
	return fmt.Sprintf(`
resource "utravis_env_var" "foo" {
    slug = "%s"
	name = "%s"
	value = "val - updated"
	public = true
}
`, testSlug, updatedName)
}
