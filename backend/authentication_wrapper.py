from functools import wraps
from flask import request
import config

def login_required(f):
  @wraps(f)
  def decorated_function(*args, **kwargs):
      if request.headers.get('Authorization') != config.PRIT_AUTH_KEY:
          return 'You are not P.R.I.T.', 401
      return f(*args, **kwargs)
  return decorated_function
