module "lambda_get" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 6.0"

  function_name = "${random_pet.this.id}-lambda-get"
  description   = "Greet user"
  handler       = "main"
  runtime       = "provided.al2"
  publish       = true

  create_package         = false
  local_existing_package = "${local.lambda_zip_path}/lambda-get.zip"

  environment_variables = {
    DYNAMODB_TABLE = "${module.dynamodb_table.dynamodb_table_id}"
  }

  attach_tracing_policy    = true
  attach_policy_statements = true

  policy_statements = {
    dynamodb_read = {
      effect    = "Allow",
      actions   = ["dynamodb:GetItem"],
      resources = [module.dynamodb_table.dynamodb_table_arn]
    }
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*/*/*"
    }
  }
}

module "lambda_put" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 6.0"

  function_name = "${random_pet.this.id}-lambda-put"
  description   = "Add user birthday"
  handler       = "main"
  runtime       = "provided.al2"
  publish       = true

  create_package         = false
  local_existing_package = "${local.lambda_zip_path}/lambda-put.zip"

  environment_variables = {
    DYNAMODB_TABLE = "${module.dynamodb_table.dynamodb_table_id}"
  }

  attach_tracing_policy    = true
  attach_policy_statements = true

  policy_statements = {
    dynamodb_read = {
      effect    = "Allow",
      actions   = ["dynamodb:PutItem"],
      resources = [module.dynamodb_table.dynamodb_table_arn]
    }
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*/*/*"
    }
  }
}

