terraform {
  required_providers {
    npm = {
      version = "0.1.0"
      source  = "hashicorp.com/mohamed-gara/npm"
    }
  }
}

provider "npm" {
  url = "http://localhost:4873"
}
