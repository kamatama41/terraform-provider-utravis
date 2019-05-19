package utravis

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/shuheiktgw/go-travis"
)

func TestAccTravisKeyPair_basic(t *testing.T) {
	t.Skip("The resource is available on Travis CI Enterprise. I don't have the Enterprise account for CI testing so far.")

	var keyPair travis.KeyPair
	description := "Test"
	updatedDescription := "Test - updated"
	publicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzL18oVN9cRKbqUKNCR7o
3cGeCafRlvpAOs2ApjU6gkhMwXfnReykVe3NhYN9v29NWO/28mubn7/sUsUtg89t
IKPqe6G0bkBbHwJMaesHN+A4XXylOR/xcP3JnRizAaGWfe5B0QE9IoggCu7EX0Jv
h602z9P+C2NlHIbHVC/TwW/rBDHGKLovtl7LHtm0At99010WUNnjZ4kfUu9Iyq4h
exfLqJFwy3tVMXRN2Z3nvG5MRR0uQu//oCXynTvOZgSDsW0tqBUFpEvK0h9ekS9a
CAz9SIB7UQ2J6Bq2mWsYy9nhU7FUWaPUfGyR7rBhvC6x4bww4uIsbMaXPEWDrEky
AwIDAQAB
-----END PUBLIC KEY-----
`
	updatedPublicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAojZ9LuMENdlS0mAfdkOe
+Lo8G9702KTMXAJG1mZ+/eQhxNd+6SCjF23EBBW8Ru/+IDBFJepM++Wt2yg/Uk+v
ZfUuB2962ES08Z5ceLp9GF8t3YqMEq7hDwS7hsjiGRFrw0813ZonSbuuPnq9B8Vg
o9Wq/E0I7VtYK4OA28+kiatXJxesSyHuZeEZlq13tlzM97M8l+hFtk55kbgABFrj
UL/qy6F4m4udnt7n0MJ4OaWGwhYLpNUHDrrst9U48Ta5riR+Auj8qYDJvNZxD6/a
uvTUgJdax+I7C4VK8ywzHY15Wf2EGN/DI8BdEAFmBF+P7ziFVVb/RwiV+ZZQi7XJ
mwIDAQAB
-----END PUBLIC KEY-----
`
	fingerprint := "f4:3d:1b:23:ca:91:df:35:36:d1:40:04:8a:b6:41:ee"
	updatedFingerprint := "cc:72:b8:f9:05:d0:c4:f4:f7:39:ff:ed:c1:02:b7:b7"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTravisKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testAccDeleteTravisKeyPair()
				},
				Config: testAccTravisKeyPairConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTravisKeyPairExists("utravis_key_pair.foo", &keyPair),
					testAccCheckTravisKeyPairAttributes(&keyPair, description, publicKey, fingerprint),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "slug", testSlug),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "description", description),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "public_key", publicKey),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "fingerprint", fingerprint),
				),
			},
			{
				Config: testAccTravisKeyPairUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTravisKeyPairExists("utravis_key_pair.foo", &keyPair),
					testAccCheckTravisKeyPairAttributes(&keyPair, updatedDescription, updatedPublicKey, updatedFingerprint),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "slug", testSlug),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "description", updatedDescription),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "public_key", updatedPublicKey),
					resource.TestCheckResourceAttr("utravis_key_pair.foo", "fingerprint", updatedFingerprint),
				),
			},
		},
	})
}

func testAccDeleteTravisKeyPair() error {
	client := travis.NewClient(os.Getenv("TRAVIS_BASE_URL"), os.Getenv("TRAVIS_API_TOKEN"))
	res, err := client.KeyPair.DeleteByRepoSlug(context.TODO(), testSlug)
	if err != nil {
		return err
	}
	if res.StatusCode != 204 && res.StatusCode != 404 {
		return fmt.Errorf("unexpected response: %d", res.StatusCode)
	}
	return nil
}

func testAccCheckTravisKeyPairExists(n string, keyPair *travis.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Key Pair ID is set")
		}

		conn := testAccProvider.Meta().(*config).client
		slug := rs.Primary.ID

		k, _, err := conn.KeyPair.FindByRepoSlug(context.TODO(), slug)
		if err != nil {
			return err
		}
		*keyPair = *k
		return nil
	}
}

func testAccCheckTravisKeyPairAttributes(keyPair *travis.KeyPair, description, publicKey, fingerprint string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *keyPair.Description != description {
			return fmt.Errorf("Key description does not match: %s, %s", *keyPair.Description, description)
		}

		if *keyPair.PublicKey != publicKey {
			return fmt.Errorf("Public key does not match: %s, %s", *keyPair.PublicKey, publicKey)
		}

		if *keyPair.Fingerprint != fingerprint {
			return fmt.Errorf("Fingerprint does not match: %s, %s", *keyPair.Fingerprint, fingerprint)
		}
		return nil
	}
}

func testAccCheckTravisKeyPairDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*config).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "utravis_key_pair" {
			continue
		}

		slug := rs.Primary.ID
		keyPair, resp, err := conn.KeyPair.FindByRepoSlug(context.TODO(), slug)
		if err == nil && keyPair != nil {
			return fmt.Errorf("key pair still exists")
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccTravisKeyPairConfig() string {
	return fmt.Sprintf(`
resource "utravis_key_pair" "foo" {
    slug = "%s"
	description = "Test"
	value = "${file("test-fixtures/id_rsa")}"
}
`, testSlug)
}

func testAccTravisKeyPairUpdateConfig() string {
	return fmt.Sprintf(`
resource "utravis_key_pair" "foo" {
    slug = "%s"
	description = "Test - updated"
	value = "${file("test-fixtures/id_rsa_updated")}"
}
`, testSlug)
}
