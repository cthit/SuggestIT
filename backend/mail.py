import requests
import config
import json

def sendMail(_to, _from, _titel, _body):
    header = {'Authorization': 'pre-shared: ' + config.GOTIFY_PRE_SHARED_KEY}
    data = json.dumps({
        "to": _to,
        "from": _from,
        "titel": _titel,
        "body": _body
    })
    return requests.post(config.GOTIFY_URL + '/mail',headers=header ,data=data).ok
