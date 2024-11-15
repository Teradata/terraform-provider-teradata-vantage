terraform {
  required_providers {
    teradata-vantage = {
      source = "hashicorp.com/edu/teradata-vantage"
    }
  }
}

provider "teradata-vantage" {
  db_host = ""
  db_user = ""
  db_password = ""
}

resource "teradata-vantage_computecluster" "edu1" {
  compute_profile_name = "test_profile"
}

