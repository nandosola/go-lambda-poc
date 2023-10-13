terraform {
  backend "s3" {
    encrypt = true
    bucket = "d29a2323db9d4"
    key = "terraform.tfstate"
    dynamodb_table = "tf-statelock"
    region = "eu-west-1"
    profile = "playground_iac"
  }
}

