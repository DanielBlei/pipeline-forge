# -----------------------------------------------
# Variables
# -----------------------------------------------

variable "postgres_password" {
  description = "Docker Composer PostgreSQL password"
  type        = string
  default     = "pipeline_forge_sandbox"
  sensitive   = true
}

variable "mysql_password" {
  description = "Docker Composer MySQL password"
  type        = string
  default     = "pipeline_forge_sandbox"
  sensitive   = true
}
