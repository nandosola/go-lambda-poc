locals {
  domain_name         = "playground-4fd1.net"
  subdomain           = "api-684b"
  lambda_zip_path     = "../dist"
  lambda_insights_arn = "arn:aws:lambda:eu-west-1:580247275435:layer:LambdaInsightsExtension:38"
}

provider "aws" {
  region  = "eu-west-1"
  profile = "playground_iac"

  # Make it faster
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_credentials_validation = true

  # skip_requesting_account_id should be disabled to generate valid ARN in apigatewayv2_api_execution_arn
  skip_requesting_account_id = false
}

resource "random_pet" "this" {
  length = 2
}

