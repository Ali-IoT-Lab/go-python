DOMAIN = 'event_bridge'

import requests
import json

#server_url = 'http://47.97.210.118/push_event'
headers = {'content-type': 'application/json'}

def setup(hass, config):
    """Set up is called when Home Assistant is loading our component."""
    count = 0
    a=config
    # Listener to handle fired events
    def handle_event(event):
        print(event)
        new_state = event.data.get('new_state')
        requests.post(config.resource, data=json.dumps({'some': 'data'}), headers=headers)

    # Listen for when state_changed is fired
    hass.bus.listen('s11tate_changed', handle_event)
          return True
