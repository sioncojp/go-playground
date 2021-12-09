import logging
import django_datadog_logger.formatters.datadog

"""ログ出力しないuri"""
silent_logging_uris = [
    "/healthz",
    "/healthcheck",
]


class Main(logging.Filter):
    def filter(self, record: logging.LogRecord) -> bool:
        """
        出力されたlogにフィルターをかけ、出力するorしないのジャッジをする

        :param record:  djangoが受け取るlog
        :return: bool。Falseなら出力させない
        """
        # config/base.pyにセットされた、django_datadog_logger.middleware
        # からgetしてきたhttp周りのログをセットする
        ddl = django_datadog_logger.formatters.datadog.DataDogJSONFormatter()
        d = ddl.get_wsgi_request()

        # ログ出力させないuri
        for v in silent_logging_uris:
            if d.path_info == v:
                return False

        return True
