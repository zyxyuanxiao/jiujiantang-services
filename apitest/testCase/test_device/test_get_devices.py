import json

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class GetDevices(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def get_devices(self):
        url = self.host + "owner/organizations/" + Util().get_organization_id() + "/devices"
        # url = urljoin(self.host, 'owner/organizations/', Util().get_organization_id(), '/devices')
        LOG.info("请求url:%s" % url)
        data = {}
        res = requests.get(url=url, json=json.dumps(data), headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_get_devices_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_devices()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
