data "archive_file" "this" {
  type        = "zip"
  source_file = "${path.module}/files/bin/${var.name}"
  output_path = "${path.module}/files/bin/${var.name}.zip"
}

resource "aws_lambda_function" "this" {
  filename         = data.archive_file.this.output_path
  function_name    = var.name
  description      = var.description
  role             = aws_iam_role.this.arn
  handler          = var.name
  source_code_hash = data.archive_file.this.output_base64sha256
  runtime          = "go1.x"

  memory_size = 128
  timeout     = 3

  environment {
    variables = {
      SlackChannelName = var.slack_channel_name
    }
  }
}

/* lambdaのiam作成 */
resource "aws_iam_role" "this" {
  name               = format("lambda-%s", var.name)
  assume_role_policy = data.aws_iam_policy_document.this_assume.json
}

data "aws_iam_policy_document" "this_assume" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "this" {
  name = format("lambda-%s", var.name)
  role = aws_iam_role.this.id

  policy = data.aws_iam_policy_document.this.json
}

data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    effect = "Allow"

    resources = [
      "arn:aws:logs:ap-northeast-1:${var.aws_account_id}:log-group:/aws/lambda/${var.name}",
      "arn:aws:logs:ap-northeast-1:${var.aws_account_id}:log-group:/aws/lambda/${var.name}:*",
    ]
  }

  statement {
    actions = [
      "cloudwatch:Describe*",
      "cloudwatch:Get*",
      "cloudwatch:List*",
      "logs:Get*",
      "logs:List*",
      "logs:Describe*",
      "logs:TestMetricFilter",
    ]

    effect = "Allow"

    resources = [
      "*",
    ]
  }

  // slackのtokenをgetするための許可
  statement {
    actions = [
      "ssm:DescribeParameters",
      "ssm:GetParameters",
      "ssm:GetParameter",
      "ssm:GetParametersByPath",
    ]

    effect = "Allow"

    resources = [
      "*",
    ]
  }

  # kmsの権限を付与
  statement {
    actions = [
      "kms:Decrypt",
      "kms:DescribeKey",
    ]

    resources = var.kms_arns
  }
}
