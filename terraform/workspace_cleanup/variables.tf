variable "enabled" {
  default     = false
  description = "Enable/Disable the module"
}

locals {
  enabled = var.enabled ? 1 : 0
}