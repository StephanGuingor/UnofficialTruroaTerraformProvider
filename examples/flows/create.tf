resource "truora_flow" "new_automated_flow" {
    name = "welpsjj :D"
    type = "permanent"
    
    config {
        lang = "en"
        enable_desktop_flow = true
    }

    verification {
        name = "email_verification"
    }

    verification {
        name = "phone_verification"
    }
}

output "new_automated_flow" {
    value = {
        flow_id = truora_flow.new_automated_flow.flow_id
        name = truora_flow.new_automated_flow.name
        type = truora_flow.new_automated_flow.type
        lang = truora_flow.new_automated_flow.config[0].lang
    }
  
}