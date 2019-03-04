DOMAIN = 'event_bridge'

import requests
import json
import traceback

#server_url = 'http://47.97.210.118/push_event'
headers = {'content-type': 'application/json'}
url = 'http://yanglao.moja-lab.com:3000/push_event'

def setup(hass, config):
    """Set up is called when Home Assistant is loading our component."""
    count = 0
    a=config
    # Listener to handle fired events
    def handle_event(event):
        print('=======================')
        #print(event)
        new_state = event.data.get('new_state')
        print(event.data.get('entity_id'))

        print('-----------------------')
        
        print(type(new_state))

        print('!!!!!!!!!!!!!!!!!!!!!!!')
        print(new_state.list_all_member())
        #try:
        #    b = json.dumps(new_state)
        #except Exception:
        #    print(Exception)
        #    print(str(Exception))
        #else:
        #    print ('success')
        #for a in new_state:
        #    print(a)
        #print(type(b))
       # requests.post(url, data=new_state, headers=headers)

    # Listen for when state_changed is fired
    hass.bus.listen('state_changed', handle_event)
    return True
