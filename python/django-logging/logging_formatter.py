from django.utils import timezone
import logging
import json
import django_datadog_logger.formatters.datadog
from django.http.request import split_domain_port

"""出力しないLOG FIELD"""
IGNORE_LOG_FIELD = {
    'args',
    'created',
    'exc_info',
    'exc_text',
    'funcName',
    'module',
    'msg',
    'levelname'
    'filename',
    'pathname',
    'process',
    'processName',
    'relativeCreated',
    'stack_info',
    'thread',
    'threadName',
    'server_time',
    'msecs',
    'lineno',
    'levelno',
    'error.kind',
}

"""LOG FIELD名を変更する。datadogのlogsに寄せる。before -> afterの書き方"""
CUSTOMIZE_LOG_FIELD = {
    'name': 'middleware',
    'http.status_code': 'status',
}


class Main(logging.Formatter):
    def format(self, record: logging.LogRecord) -> json:
        """
        受け取ったログに独自のログフィールドを追加し、jsonにフォーマットする

        :param record: djangoが受け取るlog
        :return: json型
        """
        #  "GET /v1/chat/health HTTP/1.1" 200 0
        json_record = self.new_record(record)
        result = self.recursive_to_json_format(json_record)

        return self.to_json(result)

    @staticmethod
    def to_json(record: dict) -> json:
        """
        jsonに変換
        :params record: jsonにしたいdict
        :return: json型
        """
        return json.dumps(record)

    def recursive_to_json_format(self, extra: dict) -> dict:
        """
        dictやlistやその他で書かれてる値を再帰的にjson形式に変換できる形にする

        :param extra: ユーザ定義したレコード e.g. `logger.info('Sign up', extra={'referral_code': '52d6ce'})`.
        :return: json libraryに渡すdict
        """
        result = {}
        if type(extra) is dict:
            result = {}
            for key, value in extra.items():
                if type(value) in [dict, list, tuple]:
                    res = self.recursive_to_json_format(value)
                    if res is not None:
                        result[key] = res
                else:
                    try:
                        json.dumps({"a": value})
                    except TypeError:
                        pass
                    else:
                        result[key] = value
        elif type(extra) in [list, tuple]:
            result = []
            for value in extra:
                if type(value) in [dict, list, tuple]:
                    res = self.recursive_to_json_format(value)
                    if res is not None:
                        result.append(res)
                else:
                    try:
                        json.dumps({"a": value})
                    except TypeError:
                        pass
                    else:
                        result.append(value)
        else:
            try:
                json.dumps({"a": extra})
            except TypeError:
                pass
            else:
                return extra

        return result

    def new_record(self, record: logging.LogRecord) -> dict:
        """
        出力するログを生成
        recordの中でIGNORE_LOG_FIELDに含まれないフィールドを返す
        log fieldのkeyをCUSTOMIZE_LOG_FIELDに書かれてるものに変更する
        また最低限必要なログ項目を拡張として設定する

        :param record: recordは下記の値
        {
            'name': 'django.server', 'msg': '"%s" %s %s', 'args': ('GET /v1/chat/health HTTP/1.1', '200', '0'),
            'levelname': 'INFO', 'levelno': 20,
            'pathname': '/home/python/.local/lib/python3.10/site-packages/django/core/servers/basehttp.py',
            'filename': 'basehttp.py', 'module': 'basehttp', 'exc_info': None, 'exc_text': None,
            'stack_info': None, 'lineno': 161, 'funcName': 'log_message', 'created': 1638518824.289155,
            'msecs': 289.1550064086914, 'relativeCreated': 2702.704668045044, 'thread': 281473646694880,
            'threadName': 'Thread-1', 'processName': 'MainProcess', 'process': 11,
            'request':
                <socket.socket fd=5, family=AddressFamily.AF_INET, type=SocketKind.SOCK_STREAM, proto=0,
                laddr=('172.23.0.3', 8000), raddr=('172.23.0.1', 62918)>,
            'server_time': '03/Dec/2021 08:07:04', 'status_code': 200
        }
        :return: json libraryに渡すdict
        """
        result = {}

        # recordから必要なものをセット
        for log_field, log_value in record.__dict__.items():
            if log_field not in IGNORE_LOG_FIELD:
                if log_field in CUSTOMIZE_LOG_FIELD.keys():
                    result[CUSTOMIZE_LOG_FIELD[log_field]] = log_value
                else:
                    result[log_field] = log_value

        # ここからdefault拡張ログ設定
        # config/base.pyにセットされた、django_datadog_logger.middleware
        # からgetしてきたhttp周りのログをセットする
        ddl = django_datadog_logger.formatters.datadog.DataDogJSONFormatter()
        d = ddl.get_wsgi_request()
        if d is not None:
            domain, port = split_domain_port(d.get_host())
            result["host"] = domain
            result["port"] = int(port) if port else None
            result["protocol"] = d.scheme
            result["method"] = d.method
            result["path"] = d.path_info
            result["query_string"] = d.GET.dict()
            result["request_body"] = d.body.decode('utf-8')
            result["user-agent"] = d.META.get("HTTP_USER_AGENT")
            result["remote_ip"] = get_client_ip(d)
            result["referer"] = d.META.get("HTTP_REFERER")

        # YYYY-mm-dd HH:MM:SS.zzzzzz+HH::MM / timezone(UTC/JST)を出力
        time = timezone.now()
        result['time'] = time.isoformat()
        result['timezone'] = time.tzinfo

        # datadog logsに合わせて、exceptionとcriticalはemergency扱いにする
        if str(record.levelname).lower() == 'exception' or str(record.levelname).lower() == 'critical':
            result['level'] = "EMERGENCY"
        return result


def get_client_ip(request: any) -> str:
    """
    x-forwarded-for or remote ipを返す

    :param request:
    :return: IPアドレスの文字列
    """
    x_forwarded_for = request.META.get("HTTP_X_FORWARDED_FOR")
    if x_forwarded_for:
        return x_forwarded_for.split(",")[0] or None
    else:
        return request.META.get("REMOTE_ADDR") or None
