import unittest
import responses

from requests import ConnectionError
from utils import check_filter_expression_satisfied, parse_message, process_job


class TestUtils(unittest.TestCase):
    def test_filters_should_include_body_data_in_expression_happy(self):
        """
        filters should allow you to use the data from the body/headrs/query_params
        as varaibles in the boolean expression
        happy scenario should equate to true
        """
        data = {
            'body': {
                'test_true': True,
                'test_false': False
            },
            'headers': {
                'test_number': 1,
                'test_string': 'yy'
            }
        }

        test_truth = "data['body']['test_true'] and data['headers']['test_number'] == 1"
        self.assertEqual(check_filter_expression_satisfied(test_truth, data),
                         True)

    def test_filters_should_include_body_data_in_expression_fail(self):
        """
        filters should allow you to use the data from the body/headrs/query_params
        as varaibles in the boolean expression
        fail scenario should equate to false
        """
        data = {
            'body': {
                'test_true': True,
                'test_false': False
            },
            'headers': {
                'test_number': 1,
                'test_string': 'yy'
            }
        }

        test_truth = "data['body']['test_false'] and data['headers']['test_string'] == 'yyyy'"
        self.assertEqual(check_filter_expression_satisfied(test_truth, data),
                         False)

    def test_filters_should_include_body_data_in_expression_raise_error(self):
        """
        passing an invalid filter string should not break the function
        only return false if string is not to be processed or evaluated to false
        """
        data = {}

        test_truth = "data['body']['test_false'] and data['headers']['test_string'] == 'yyyy'"
        self.assertEqual(check_filter_expression_satisfied(test_truth, data),
                         False)

    def test_parsing_should_return_original_message_as_json_happy(self):
        """
        a base64 string passed to this function should be transformed back to its
        original json format
        """
        encoded_data = "eyJib2R5IjogeyJtZXNzYWdlIjogIlR\oZSBmb3JjZSBpcyBzdHJvbmcgd2l0aCB0aGlzIG9uZS4uLiJ9LCAiaGVhZGVycyI6IHsiYWxwaGEiOiAiYmV0YSIsICJDb250ZW50LVR5cGUiOiAiYXBwbGljYXRpb24vanNvbiJ9LCAicXVlcnlfcGFyYW1zIjogeyJtaSI6ICJwaWFjaXRvIn0sICJwYXRoX3BhcmFtcyI6IHsiZG9tYWluIjogInBpcGVkcmVhbSJ9fQ=="
        decoded_data = parse_message(encoded_data)

        expected_data = {
            'body': {
                'message': 'The force is strong with this one...'
            },
            'headers': {
                'alpha': 'beta',
                'Content-Type': 'application/json'
            },
            'query_params': {
                'mi': 'piacito'
            },
            'path_params': {
                'domain': 'pipedream'
            }
        }

        self.assertEqual(decoded_data, expected_data)

    def test_parsing_should_return_original_message_as_json_fail(self):
        """
        passing an invalid base64 string should not break the function
        should only return false
        """
        encoded_data = "fake_base64_string_will_not_parse_to_json_isa"
        decoded_data = parse_message(encoded_data)

        self.assertEqual(decoded_data, False)

    @responses.activate
    def test_process_job_happy(self):
        """
        process job should pass filtering and data decoding phases
        and sucessfully send request to mocked url
        """
        encoded_data = "eyJib2R5IjogeyJtZXNzYWdlIjogIlR\oZSBmb3JjZSBpcyBzdHJvbmcgd2l0aCB0aGlzIG9uZS4uLiJ9LCAiaGVhZGVycyI6IHsiYWxwaGEiOiAiYmV0YSIsICJDb250ZW50LVR5cGUiOiAiYXBwbGljYXRpb24vanNvbiJ9LCAicXVlcnlfcGFyYW1zIjogeyJtaSI6ICJwaWFjaXRvIn0sICJwYXRoX3BhcmFtcyI6IHsiZG9tYWluIjogInBpcGVkcmVhbSJ9fQ=="
        filter_exp = "data['headers']['alpha'] == 'beta' and data['query_params']['mi'] == 'piacito'"
        url = 'https://leadrole.cage'

        responses.add(responses.POST,
                      f'{url}/?mi=piacito',
                      json={'sucess': 'thank you'},
                      status=200)

        process_job(url=url, filter_exp=filter_exp, payload=encoded_data)

    def test_process_job_fail(self):
        """
        process job should pass filtering and data decoding phases
        but fail to send request as url is invalid
        """
        encoded_data = "eyJib2R5IjogeyJtZXNzYWdlIjogIlR\oZSBmb3JjZSBpcyBzdHJvbmcgd2l0aCB0aGlzIG9uZS4uLiJ9LCAiaGVhZGVycyI6IHsiYWxwaGEiOiAiYmV0YSIsICJDb250ZW50LVR5cGUiOiAiYXBwbGljYXRpb24vanNvbiJ9LCAicXVlcnlfcGFyYW1zIjogeyJtaSI6ICJwaWFjaXRvIn0sICJwYXRoX3BhcmFtcyI6IHsiZG9tYWluIjogInBpcGVkcmVhbSJ9fQ=="
        filter_exp = "data['headers']['alpha'] == 'beta' and data['query_params']['mi'] == 'piacito'"
        url = 'https://leadrole.cage'

        with self.assertRaises(ConnectionError):
            process_job(url=url, filter_exp=filter_exp, payload=encoded_data)

    def test_process_job_fail_decode(self):
        """
        failure in decoding should still pass sucessfully(skip job as no meaning in repeating it)
        """
        encoded_data = "malformed_base64"
        filter_exp = "data['headers']['alpha'] == 'beta' and data['query_params']['mi'] == 'piacito'"
        url = 'https://leadrole.cage'

        process_job(url=url, filter_exp=filter_exp, payload=encoded_data)

    def test_process_job_fail_filter(self):
        """
        failure in filtering should still pass sucessfully(skip job as no meaning in repeating it)
        """
        encoded_data = "eyJib2R5IjogeyJtZXNzYWdlIjogIlR\oZSBmb3JjZSBpcyBzdHJvbmcgd2l0aCB0aGlzIG9uZS4uLiJ9LCAiaGVhZGVycyI6IHsiYWxwaGEiOiAiYmV0YSIsICJDb250ZW50LVR5cGUiOiAiYXBwbGljYXRpb24vanNvbiJ9LCAicXVlcnlfcGFyYW1zIjogeyJtaSI6ICJwaWFjaXRvIn0sICJwYXRoX3BhcmFtcyI6IHsiZG9tYWluIjogInBpcGVkcmVhbSJ9fQ=="
        filter_exp = "data['headers']['alpha'] == 'beta' and data['query_params']['mi'] == 'piacito' and unkown"
        url = 'https://leadrole.cage'

        process_job(url=url, filter_exp=filter_exp, payload=encoded_data)

if __name__ == '__main__':
    unittest.main()