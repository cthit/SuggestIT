import requests
import config
import json

def sendMail(_to, _from, _subject, _body):
    header = {'Authorization': 'pre-shared: %s' % config.GOTIFY_PRE_SHARED_KEY}
    data = json.dumps({
        "to": _to,
        "from": _from,
        "subject": _subject,
        "body": _body
    })
    return requests.post('%s/mail' % config.GOTIFY_URL,headers=header ,data=data).ok
