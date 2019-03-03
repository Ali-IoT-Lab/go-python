DOMAIN = 'event_bridge'
import requests
import json

server_url = 'http://47.97.210.118/push_event'
headers = {'content-type': 'application/json'}

def setup(hass, config):
    """Set up is called when Home Assistant is loading our component."""
    count = 0

    # Listener to handle fired events
    def handle_event(event):
        nonlocal count
        count += 1
        print('Total events received:', count)
        new_state = event.data.get('new_state')
        requests.post(server_url, data=json.dumps(new_state), headers=headers)

    # Listen for when my_cool_event is fired
    hass.bus.listen('state_changed', handle_event)