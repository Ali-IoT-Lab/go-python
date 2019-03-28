DOMAIN = 'event_bridge'

import requests
import json
import traceback
import homeassistant.helpers.config_validation as cv
from datetime import datetime
import hmac
import base64
import hashlib
import platform

accountid = "1385629681127691"
accessKeyID = "LTAIANf6x7njUgCS"
accessKeySecret = "32YOvOX9fT9nZVqn6kVNFN94Tz6dIM"
method = "POST"
version = "2016-08-15"
path = "/2016-08-15/proxy/gaofeng-test/httptest/"
region = "cn-hangzhou"
host = "1385629681127691.cn-hangzhou.fc.aliyuncs.com"
server_url = 'https://1385629681127691.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/gaofeng-test/httptest/'
GMT_FORMAT = "%a, %d %b %Y %H:%M:%S GMT"
__version__ = "2.0.11"
user_agent = \
            'aliyun-fc-sdk-v{0}.python-{1}.{2}-{3}-{4}'.\
            format(__version__, platform.python_version(),
                    platform.system(), platform.release(), platform.machine())
headers = {
    "content-type": "application/json",
    "accept": "application/json",
    "host": host,
    "x-fc-account-id": accountid,
    "user-agent": user_agent
}

def build_canonical_headers(headers):
    canonical_headers = []
    for k, v in headers.items():
        lower_key = k.lower()
        if lower_key.startswith('x-fc-'):
            canonical_headers.append((lower_key, v))
    canonical_headers.sort(key=lambda x: x[0])
    if canonical_headers:
        return '\n'.join(k + ':' + v for k, v in canonical_headers) + '\n'
    else:
        return ''

def get_sign_resource(unescaped_path, unescaped_queries):
    if not isinstance(unescaped_queries, dict):
        raise TypeError("`dict` type required for queries")
    params = []
    for key, values in unescaped_queries.items():
        if isinstance(values, str):
            params.append('%s=%s' % (key, values))
            continue
        if len(values) > 0:
            for value in values:
                params.append('%s=%s' % (key, value))
        else:
            params.append('%s' % key)
    params.sort()
    resource = unescaped_path + '\n' + '\n'.join(params)
    return resource

def sign_request(accessKeyID, accessKeySecret, method, unescaped_path, headers, unescaped_queries=None):
    content_md5 = headers.get('content-md5', '')
    content_type = headers.get('content-type', '')
    date = headers.get('date', '')
    canonical_headers = build_canonical_headers(headers)
    canonical_resource = unescaped_path
    if isinstance(unescaped_queries, dict):
        canonical_resource = get_sign_resource(unescaped_path, unescaped_queries)
    string_to_sign = '\n'.join(
            [method.upper(), content_md5, content_type, date, canonical_headers + canonical_resource])
    h = hmac.new(accessKeySecret.encode('utf-8'), string_to_sign.encode('utf-8'), hashlib.sha256)
    signature = 'FC ' + accessKeyID + ':' + base64.b64encode(h.digest()).decode('utf-8')
    return signature

def setup(hass, config):
    """Set up is called when Home Assistant is loading our component."""
    # Listener to handle fired events
    def handle_event(event):
        new_state = event.data.get('new_state')
        msg = {}
        msg.update(state = new_state.state)
        msg.update(object_id = new_state.object_id)
        msg.update(name = new_state.name)
        msg.update(time = new_state.last_updated.timestamp())
        msg.update(battery_level = new_state.attributes.get('battery_level'))
        device_class = new_state.attributes.get('device_class')
        flag = False
        device_list = ['opening', 'humidity', 'temperature', 'motion']
        for item in device_list:
            if (item == device_class):
                flag = True
                break
        if (flag == False):
            checkKey = new_state.object_id.split('_')
            if (checkKey[0] == 'vibration'):
                flag = True
            if (checkKey[0] == 'plug'):
                flag = True
                msg.update(load_power = new_state.attributes.get('load_power'))
        if (flag == True):
            send_data = json.dumps(msg).encode('utf-8')
            headers["content-length"] = str(len(send_data))
            headers["date"] = datetime.utcnow().strftime(GMT_FORMAT)
            headers["authorization"] = sign_request(accessKeyID, accessKeySecret, method, path, headers, {})
            requests.post(server_url, data=send_data, headers=headers)
    #a Listen for when state_changed is fired
    hass.bus.listen('state_changed', handle_event)
    return True
