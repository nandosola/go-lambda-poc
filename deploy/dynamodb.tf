module "dynamodb_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "~> 3.3"

  name      = random_pet.this.id
  hash_key  = "id"

  attributes = [
    {
      name = "id"  # SHA256sum, just a uid, max 12 chars
      type = "S"
    },
  ]
}
