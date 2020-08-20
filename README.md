# Terraform Provider Travis (Unofficial)

A Terraform provider to interact with [Travis CI](https://travis-ci.com/) resources.

https://registry.terraform.io/providers/kamatama41/utravis/latest

## Support
The Latest version supports Terraform >= 0.13.0. It might work on Terraform 0.12 but not be guaranteed.

## Installation
### For Terraform 0.13 users
You can install the provider via `terraform init` with the following configuration

```
terraform {
  required_providers {
    utravis = {
      source = "kamatama41/utravis"
      version = "~> 0.0"
    }
  }
}
```

### For Terraform 0.12 users
As this is a Third party plugin, so you have to install it manually.

Download plugin from GitHub releases and unarchive it. That's all!
```sh
$ latest=$(curl -s https://api.github.com/repos/kamatama41/terraform-provider-utravis/releases/latest | jq -r ".name") 
$ os=$(uname | tr '[:upper:]' '[:lower:]')
$ archive_name=terraform-provider-utravis_${latest//v}_${os}_amd64.zip
$ curl -LO https://github.com/kamatama41/terraform-provider-utravis/releases/download/${latest}/${archive_name}
$ unzip ${archive_name} && rm ${archive_name}
```

(Optional) If you want to use the plugin for other Terraform projects, place the binary into `~/.terraform.d/plugins` (`%APPDATA%\terraform.d\plugins` for Windows users)

## Configuration
`base_url` and `token` are required. You can use the environment variable `TRAVIS_BASE_URL` and `TRAVIS_API_TOKEN` instead of them.

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

# Add a private key to the repository
resource "utravis_key_pair" "my-repo" {
  slug = "myuser/my-repository"
  value = "${file("~/.ssh/id_travis_rsa")}"
}
```

## Supported resources
- `utravis_env_var`: https://developer.travis-ci.com/resource/env_var
- `utravis_key_pair`: https://developer.travis-ci.com/resource/key_pair
