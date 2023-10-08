data "truora_flow_document" "new_automated_flow_document" {
    name = "Epic new flow :D"
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

resource "truora_flow" "new_automated_flow_inline" {
    document = data.truora_flow_document.new_automated_flow_document.json
}

output "new_automated_flow" {
    value = {
        flow_id = truora_flow.new_automated_flow_inline.flow_id
        name = truora_flow.new_automated_flow_inline.name
    }
}