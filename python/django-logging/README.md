# django-logging
## Overview
- filter: ログを出力したくないuriを設定
- formatter: 下記のようなログを出力する
```shell
$ curl -X POST -H "Content-Type: application/json" -d '{"Name":"hogehoge", "fugarfuga": "fugaaaaaa"}' localhost:8000/hoge
{
	"middleware": "django_datadog_logger.middleware.request_log",
	"levelname": "WARNING",
	"filename": "request_log.py",
	"status": 404,
	"duration": 28986692.428588867,
	"error.message": "Not Found",
	"host": "localhost",
	"port": 8000,
	"protocol": "http",
	"method": "POST",
	"path": "/hoge",
	"query_string": {},
	"request_body": "{\"Name\":\"hogehoge\", \"fugarfuga\": \"fugaaaaaa\"}",
	"user-agent": "curl/7.77.0",
	"remote_ip": "172.23.0.1",
	"referer": null,
	"time": "2021-12-08T12:44:50.447804+00:00"
}
```

## Usage
```shell
# config/logging_filter.py
# config/logging_formatter.py
# を配置

# requirements.txt
django-datadog-logger==0.5.0


# config/base.py
MIDDLEWARE = [
    "django_datadog_logger.middleware.request_id.RequestIdMiddleware",
    .
    .
    .
    "django_datadog_logger.middleware.error_log.ErrorLoggingMiddleware",
    "django_datadog_logger.middleware.request_log.RequestLoggingMiddleware",
]

# config/environemnt名.py
LOGGING = {
    "version": 1,
    "disable_existing_loggers": False,
    'filters': {
        'custom': {
            '()': 'config.logging_filter.Main',
        },
    },
    "formatters": {
        "json": {"()": "config.logging_formatter.Main"},
    },
    "handlers": {
        "console": {"level": "INFO", "class": "logging.StreamHandler", "formatter": "json", 'filters': ['custom']},
    },
    "loggers": {
        'django.server': {'level': 'ERROR'},
        "django_datadog_logger.middleware.request_log": {"handlers": ["console"], "level": "INFO", "propagate": False},
    },
}
DEBUG: False
```

## Motivation
- healthcheck用uriなど、ログを出力させたくない要件があった
- djangoの標準ログは質素すぎて足りない部分が多かった
- datadog logsでjsonフィールドを認識させたかった。標準フォーマットだと、json型じゃないので認識しなかった
- datadog logsで検索しやすいように、意図的にフィールド名を変更したかった

## 実装周り解説
### overview
- django標準ログだと質素で必要なデータがとれない + datadogでjsonとしてジャッジされないのでカスタマイズしてます
- required: https://github.com/namespace-ee/django-datadog-logger
- 参考実装
  - https://github.com/namespace-ee/django-datadog-logger
  - https://github.com/eieste/django-datadog-logger
  
### `config/base.py`
- django-datadog-loggerをmiddlewareに登録
- logging_formatterでラッパーし、httpの詳細ログを取得するため

### `config/logging_formatter.py`
- format用に作成 
- django-datadog-loggerでも足りない必要なログを追加
- datadog logsに合わせてフィールド名を変更
- 最後にjson型に変更して出力

### `config/logging_filter.py`
- filter用。ログの出力をするorしないをジャッジ
- 主にヘルスチェックなど、有効なURIだがロギングしたくないものに対して行ってます

### `config/environemnt名.py`
- formatに上記のlogging_formatter.pyを設定
- filterは上記のlogging_filter.pyを設定
- handlerはconsole（コンソールにstdoutとして出力用）の1つだけ作成し、filter/formatterは上記を選択
- loggersは2つ
  - `django.server`...djangoが生成する標準ログ。質素で使えない（`"GET / HTTP/1.1" 404 179` のようなログ）ため出力させないようにしてます
    - level: ERRORにすることで、 上記のログを出さないようにしています
    - handlerも設定不要
  - `django_datadog_logger.middleware.request_log`...httpリクエストを解析したログも出力されるやつ。これをベースにログを生成します
- `DEBUG: True` にすることで、ルーティングしてないURIアクセス時の無駄なログ（ `Not Found: /` ）を消しています。