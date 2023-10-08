terraform {
  required_providers {
    truora = {
      version = "0.2"
      source  = "truora.com/local/truora"
    }
  }
}

data "truora_flow" "onboarding_flow" {
  flow_id = "IPFf33fc222f4eb13e55422f07158edb7b3"
}

output "onboarding_flow" {
  value = {
    flow_id = data.truora_flow.onboarding_flow.flow_id
    name = data.truora_flow.onboarding_flow.name
    type = data.truora_flow.onboarding_flow.type
    lang = data.truora_flow.onboarding_flow.config.lang
  }
}