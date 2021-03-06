---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "npm Provider"
subcategory: ""
description: |-
  
---

# npm Provider

A Terraform plugin for downloading and accessing npm package files.

## Example Usage

```terraform
terraform {
  required_providers {
    npm = {
      version = "0.3.0"
      source  = "mohamed-gara/npm"
    }
  }
}

provider "npm" {
  url = "https://registry.npmjs.org"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **url** (String) The npm registry URL. The value can be defined using the TERRAFORM_PROVIDER_NPM_URL environment variable.

### Optional

- **password** (String, Sensitive) The password to use for authentication. The password should be defined with the TERRAFORM_PROVIDER_NPM_PASSWORD environment variable.
- **username** (String) The username to use for authentication. The username should be defined with the TERRAFORM_PROVIDER_NPM_USERNAME environment variable.
