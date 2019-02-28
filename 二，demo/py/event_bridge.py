import requests
import json
from sseclient import SSEClient

url = 'http://47.97.210.118/push_event'
payload = {'some': 'data'}
headers = {'content-type': 'application/json'}

auth = {'Authorization': 'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE1NTEzMjMxNTUsImV4cCI6MTg2NjY4MzE1NSwiaXNzIjoiN2EyNjI2YzE0ZWYwNDgzZjk4ZmM4MDkzNDExMjFkODIifQ.zJ8a9W9rkYoIjQPM822tVB8aeGZ39FviBL8Iz7Mxys4'}

def setup(hass, config):
    messages = SSEClient('http://192.168.31.95:8123/api/stream', headers=auth)

    for msg in messages:
       outputMsg = msg.data
       #print(outputMsg)
       if outputMsg != 'ping':
          print('---------------------')
          outputJS = json.loads(outputMsg)
          #print( FilterName, outputJS[FilterName] )
          print(outputJS['data'])
          ret = requests.post(url, data=json.dumps(outputJS['data']), headers=headers)
          print(ret)
      






