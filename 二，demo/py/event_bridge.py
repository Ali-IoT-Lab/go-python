import requests
import logging
import json
import voluptuous as vol
from datetime import timedelta
from homeassistant.components.sensor import PLATFORM_SCHEMA
from homeassistant.helper.entity import Entity
import homeassistant.helpers.config_validation as cv
from homeassistant.util import Throttle
from requests.exceptions import (
        ConnectionError as ConnectError, HTTPError, Timeout)
from sseclient import SSEClient

_LOGGER = logging.getLogger(__name__)

TIME_BETWEEN_UPDATES = timedelta(minutes=30)


CONF_RESOURCE = 'resource'

PLATFORM_SCHEMA = PLATFORM_SCHEMA.extend(
    {
        vol.Required(CONF_RESOURCE): cv.string,
    }        
)
def setup_platform(hass, config, add_entities, discovery_info=None):
    #_LOGGER.info(config.get(CONF_RESOURCE))
    auth = {'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NTE0MjA0NzIsImV4cCI6MTg2Njc4MDQ3MiwiaXNzIjoiMDhiYzQ0YzYzMDQ4NGMyOTg4MzIzMWFkZjRkMmY2ZTMifQ.A6cbNYwsqqXFFrG83jebJ1LzQ8VZBs8JiytoLolZb70'}
    messages = SSEClient('http://localhost:8123/api/stream', headers=auth)
    for msg in messages:
        outputMsg = msg.data
        if outputMsg != 'ping':
            outputJS = json.loads(outputMsg)
            requests.post(
                config.get(CONF_RESOURCE),
                data=json.dumps(outputJS['data']),
                headers={'content-type': 'application/json'}
            )
    dev = []
    add_entities(dev, True)
