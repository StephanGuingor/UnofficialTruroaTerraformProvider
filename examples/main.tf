terraform {
  required_providers {
    truora = {
      version = "0.2"
      source  = "truora.com/local/truora"
    }
  }
}

provider "truora" {
  api_server = "https://api.identity.truora.com"
}


module "flows" {
  source = "./flows"
}

output "root" {
  value = module.flows 
}
