import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class WeekMeasurements(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def weekmeasurements(self, user_id):
        url = self.host + 'owner/week_measurements?'+str(user_id)
        data = {}
        res = requests.get(url=url, params=data, headers=Util().get_token())
        LOG.info("url是%s" % url)
        LOG.info("data是%s" % data)
        LOG.info("结果是%s" % res.json())
        return res.json()

    def test_get_week_measurements_data(self):
        result = self.weekmeasurements(Util().get_user_id())
        self.assertOkResult(result)

    def test_get_week_measurements_error_user_id_type(self):
        result = self.weekmeasurements("44981148494984944984894849848498484984984")
        self.assertErrorResult(result, const.ErrInvalidValue)

    def test_get_week_measurements_error_user_id_value(self):
        result = self.weekmeasurements(44981148494984944984894849848498484984984)
        self.assertErrorResult(result, const.ErrInvalidValue)
