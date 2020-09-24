import json
import base64
import logging
import requests
import sys
from faktory import Worker

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

stream = logging.StreamHandler(sys.stdout)
stream.setLevel(logging.INFO)
logger.addHandler(stream)


def subtitute_path_variables_in_url(url, path_params):
    """
    """
    for path_p in path_params:
        url = url.replace(f':{path_p}', path_params[path_p])

    return url


def send_request(url, filter_exp, payload):
    """
    """
    logger.info(f"[Switch_Consumer] recieved new job for url {url} with filter {filter_exp}")
    try:
        if not filter_exp or not len(filter_exp):
            filter_exp = "True"

        is_satisfied = eval(filter_exp)
        if not is_satisfied:
            return
    except Exception as e:
        logger.error(
            f'[Switch_Consumer] Failed to evaluate filter for sub (filter: {filter_exp})  | error {e}')
        return

    try:
        data = json.loads(base64.b64decode(payload).decode('utf-8'))
    except Exception as e:
        logger.error(
            f'[Switch_Consumer] Failed to parse payload into json| error {e}')
        return

    logger.info("[Switch_Consumer] Sucessfully parsed job for repective subscriber")
    body = data.get('body') or {}
    headers = data.get('headers') or {}
    query_params = data.get('query_params') or {}
    path_params = data.get('path_params') or {}
    http_method = data.get('http_method') or 'POST'

    url = subtitute_path_variables_in_url(url, path_params)

    resp = requests.request(method=http_method,
                            url=url,
                            params=query_params,
                            headers=headers,
                            data=body)
    resp.raise_for_status()
    logger.info("[Switch_Consumer] Job Executed Sucessfully")

w = Worker(queues=['default'], concurrency=1)
w.register('my_queue', send_request)
logger.info('[Switch_Consumer] Sucessfully registered consumer on task queue, waiting for new jobs ....')
w.run()