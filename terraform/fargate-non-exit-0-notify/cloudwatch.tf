resource "aws_cloudwatch_event_rule" "this" {
  name        = "lambda-${var.name}"
  description = var.description

  # task のイベントのみ見るようにする
  event_pattern = <<PATTERN
{
  "source": [
    "aws.ecs"
  ],
  "detail-type": [
    "ECS Task State Change"
  ],
  "detail": {
    "clusterArn": [
      "${var.ecs_cluster_arn}"
    ]
  }
}
PATTERN

}

resource "aws_cloudwatch_event_target" "this" {
  target_id = aws_cloudwatch_event_rule.this.name
  rule      = aws_cloudwatch_event_rule.this.name
  arn       = aws_lambda_function.this.arn
}

resource "aws_lambda_permission" "this" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.this.arn
}

resource "aws_cloudwatch_log_group" "lambda_log" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = 3
}