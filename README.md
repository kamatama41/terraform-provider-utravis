# Terraform Provider Unofficial Travis

A Terraform provider to interact with [Travis CI](https://travis-ci.com/) resources.

## Prerequisites
- Terraform (tested on 0.11.3)

## Installation
This is a [Third party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins), so you have to install it manually.

Download plugin from GitHub releases and unarchive it. That's it!
```sh
$ latest=$(curl -s https://api.github.com/repos/kamatama41/terraform-provider-unofficial-travis/releases/latest | jq -r ".name")
$ os=$(uname | tr '[:upper:]' '[:lower:]')
$ curl -LO https://github.com/kamatama41/terraform-provider-unofficial-travis/releases/download/${latest}/terraform-provider-utravis_${latest}_${os}_amd64.zip
$ unzip terraform-provider-utravis_${latest}_${os}_amd64.zip && rm terraform-provider-utravis_${latest}_${os}_amd64.zip
```

(Optional) If you want to use the plugin for other Terraform projects, place it to `~/.terraform.d/plugins` (`%APPDATA%\terraform.d\plugins` for Windows users)

## Configuration
```hcl
# Configure the unofficial Travis Provider (utravis)
provider "utravis" {
  base_url = "https://api.travis-ci.com/"
  token = "${var.travis_api_token}"
}

# Add an environment variable to the repository
resource "utravis_env_var" "my-repo" {
  slug = "myuser/my-repository"
  name = "FOO"
  value = "bar"
  public = true
}
```
