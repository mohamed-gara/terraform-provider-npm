terraform {
  required_providers {
    npm = {
      version = "0.2.0"
      source  = "mohamed-gara/npm"
    }
  }
}

provider "npm" {
  url = "https://registry.npmjs.org"
}
