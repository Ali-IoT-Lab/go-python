import requests
import json
from sseclient import SSEClient

url = 'http://47.97.210.118/push_event'
payload = {'some': 'data'}
headers = {'content-type': 'application/json'}

auth = {'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NTE0MjA0NzIsImV4cCI6MTg2Njc4MDQ3MiwiaXNzIjoiMDhiYzQ0YzYzMDQ4NGMyOTg4MzIzMWFkZjRkMmY2ZTMifQ.A6cbNYwsqqXFFrG83jebJ1LzQ8VZBs8JiytoLolZb70'}

def setup(hass, config):
    messages = SSEClient('http://192.168.31.94:8123/api/stream', headers=auth)

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
      






