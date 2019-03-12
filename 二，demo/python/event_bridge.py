DOMAIN = 'event_bridge'

import requests
import json
import traceback
import homeassistant.helpers.config_validation as cv

server_url = 'http://47.97.210.118/push_event'
headers = {'content-type': 'application/json'}

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
        device_list = ['opening', 'humidity', 'temperature', 'illuminance', 'motion']
        for item in device_list:
            if (item == device_class):
                flag = True
                break
        if (flag == False):
            checkKey = new_state.object_id.split('_')
            if (checkKey[0] == 'vibration' or checkKey[0] == 'plug'):
                flag = True
        if (flag == True):
            requests.post(server_url, data=json.dumps(msg), headers=headers)
    #a Listen for when state_changed is fired
    hass.bus.listen('state_changed', handle_event)
    return True
