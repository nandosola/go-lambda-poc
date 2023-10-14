module "api_gateway" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "~> 2.0"

  name          = "${random_pet.this.id}-http"
  description   = "My awesome HTTP API Gateway"
  protocol_type = "HTTP"

  domain_name                 = "${local.subdomain}.${local.domain_name}"
  domain_name_certificate_arn = module.acm.acm_certificate_arn

  default_stage_access_log_destination_arn = aws_cloudwatch_log_group.api_gw.arn
  default_stage_access_log_format          = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.integrationErrorMessage"


  integrations = {
    "GET /hello/{username}" = {
      integration_type       = "AWS_PROXY"
      lambda_arn             = module.lambda_get.lambda_function_arn
      payload_format_version = "2.0"
    }

    "$default" = {
      lambda_arn = module.lambda_get.lambda_function_arn
    }
  }
}

resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/apigw/${random_pet.this.id}-http"

  retention_in_days = 30
}

