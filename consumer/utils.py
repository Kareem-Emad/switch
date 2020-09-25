import json
import base64
import logging
import requests

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

tag = '[Switch_Consumer]'

def subtitute_path_variables_in_url(url, path_params):
    """
    Args:
        url: (string)url of the subscriber containing placeholders
        path_params: (dict) a map of each placeholder and its corresponding value

    Returns:
        (string) the url after subtituting all placeholders in url.
    """
    for path_p in path_params:
        url = url.replace(f':{path_p}', path_params[path_p])

    return url

def parse_message(payload):
    """parses a base64 string into its original json format
    Args:
        payload: (string) the data contianing all info to be passed in request(headers,body,..)
                          encoded as base64 string
    Returns:
        (dict) map containing the data decoded from the base64 string or false if failed to parse
    """
    try:
        return json.loads(base64.b64decode(payload).decode('utf-8'))
    except Exception as e:
        logger.error(
            f'{tag} Failed to parse payload into json| error {e}')
        return False

def check_filter_expression_satisfied(filter_exp, data):
    """check if  boolean expression for the subscriber is satified for the current request params or not 
    Args:
        filter_exp: (string) a string containig the condition on which the job is accepted for this subscriber
                    example data['body']['blabla'] == 'anything' or data['headers']['kaza'] == 'kaz_tani'
        data: (dict) a dict contaning all request params (body, headers, ...)
    Returns:
        (bool) whether the expression final evalution is true or false
    """
    try:
        if not filter_exp or not len(filter_exp):
            filter_exp = "True"
        return eval(filter_exp)

    except Exception as e:
        logger.error(
            f'{tag} Failed to evaluate filter for sub (filter: {filter_exp})  | error {e}')
        return False

def process_job(url, filter_exp, payload):
    """processes the job data and sends the request to the respective url with all params as specified
    Args:
        url: (string)url of the subscriber containing placeholders
        filter_exp: (string) a string containig the condition on which the job is accepted for this subscriber
                    example data['body']['blabla'] == 'anything' or data['headers']['kaza'] == 'kaz_tani'
        payload: (string) encoded request data including (headers, body, ...)
    Returns:
        (bool) whether the expression final evalution is true or false
    Raises:
        HTTP_ERROR: if the request failed or got any respone code from 4XX, 5XX,..
    """

    logger.info(f"{tag} recieved new job for url {url} with filter {filter_exp}")

    data = parse_message(payload)

    if data and check_filter_expression_satisfied(filter_exp=filter_exp, data=data):

        logger.info(f"{tag} Sucessfully parsed job data | filter expression for subscriber is satsified")
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
        logger.info(f"{tag} Job Executed Sucessfully")
    else:
        logger.info(f"{tag} Job Skipped due to irrecoverable error in request data || filter_expression")
