module "dynamodb_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "~> 3.3"

  name      = random_pet.this.id
  hash_key  = "Id"

  attributes = [
    {
      name = "Id"
      type = "S"
    },
  ]
}
