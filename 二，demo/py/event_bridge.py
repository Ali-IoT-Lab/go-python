from sseclient import SSEClient
import requests
import json
url = 'http://192.168.1.5:8123/api/stream'
server_url = 'http://47.97.210.118/push_event'
payload = {'some': 'data'}
headers = {'content-type': 'application/json'}

auth = {'Authorization': 'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiI5YjRjNDI5Yzg2MmY0OGFjOWJmODk0NjdjYjRhNjI0OSIsImV4cCI6MTg2Njk3NzkyNywiaWF0IjoxNTUxNjE3OTI3fQ.lWNLKS18WCvsP_SocEA7lf1-ZyjVxEbU3slvZtZhNtI'}

def setup(hass, config=None):
    messages = SSEClient(url, headers=auth)
    for msg in messages:
      outputMsg = msg.data
      if outputMsg != 'ping':
        ret_server = requests.post(url, data=json.dumps(payload), headers=headers)
    return True 