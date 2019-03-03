DOMAIN = 'event_bridge'
from datetime import timedelta
import functools as ft

from homeassistant.loader import bind_hass
from homeassistant.helpers.sun import get_astral_event_next
from ..core import HomeAssistant, callback
from ..util.async_ import run_callback_threadsafe
import requests
import json

server_url = 'http://47.97.210.118/push_event'
headers = {'content-type': 'application/json'}

def threaded_listener_factory(async_factory):
    """Convert an async event helper to a threaded one."""
    @ft.wraps(async_factory)
    def factory(*args, **kwargs):
        """Call async event helper safely."""
        hass = args[0]

        if not isinstance(hass, HomeAssistant):
            raise TypeError('First parameter needs to be a hass instance')

        async_remove = run_callback_threadsafe(
            hass.loop, ft.partial(async_factory, *args, **kwargs)).result()

        def remove():
            """Threadsafe removal."""
            run_callback_threadsafe(hass.loop, async_remove).result()

        return remove

    return factory

@callback
@bind_hass
def async_track_state_change(hass, config):
    """Set up is called when Home Assistant is loading our component."""
    count = 0

    # Listener to handle fired events
    def handle_event(event):
        nonlocal count
        count += 1
        print('Total events received:', count)
        new_state = event.data.get('new_state')
        requests.post(server_url, data=json.dumps({'some': 'data'}), headers=headers)

    # Listen for when my_cool_event is fired
    return hass.bus.async_listen_once('state_changed', handle_event)
track_state_change = threaded_listener_factory(async_track_state_change)