import unittest

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class Devices(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def device(self, mac):
        data = {'mac': mac}
        url = self.host + 'owner/organizations/' + Util().get_organization_id() + '/devices/is_bindable'
        # url = urljoin(self.host, 'owner/organizations/', Util().get_organization_id(), 'devices/is_bindable')
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_device_correct_parameters(self):
        u"""正确参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.device(data.mac)
        LOG.info('获取测试结果：%s' % result)
        self.assertOkResult(result)
        LOG.info('------pass!---------')

    def test_device_error_parameters(self):
        u"""错误参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.device('')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrInvalidDevice)
        LOG.info('------pass!---------')


if __name__ == '__main__':
    unittest.TestCase()
