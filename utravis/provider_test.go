package utravis

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testSlug = os.Getenv("TRAVIS_SLUG")

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"utravis": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TRAVIS_API_TOKEN"); v == "" {
		t.Fatal("TRAVIS_API_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("TRAVIS_BASE_URL"); v == "" {
		t.Fatal("TRAVIS_BASE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("TRAVIS_SLUG"); v == "" {
		t.Fatal("TRAVIS_SLUG must be set for acceptance tests")
	}
}
