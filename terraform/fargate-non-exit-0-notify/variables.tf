variable "name" {
  default = "fargate-non-exit-0-notify"
}

variable "description" {
  default = ""
}

variable "ecs_cluster_arn" {
  type = string
}

variable "slack_channel_name" {
  type = string
}

variable "aws_account_id" {
  type = string
}

variable "kms_arns" {
  type = list(string)
}

