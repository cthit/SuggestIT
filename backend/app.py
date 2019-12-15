from datetime import datetime
from swagger_blueprint import SWAGGERUI_BLUEPRINT, SWAGGER_URL

from db import Suggestion, suggestion_to_json
import config
from mail import sendMail
from authentication_wrapper import login_required

from requests.exceptions import RequestException
from flask import Flask, request, jsonify
from flask_restful import Resource, Api
from pony.orm import db_session, ObjectNotFound, commit

from flask_cors import CORS

app = Flask(__name__)
api = Api(app)

PRIT_MAIL = "prit@chalmers.it"
SENDER_MAIL = "no-reply@chalmers.it"

if config.IS_PROD == 'true':
  app.config['DEBUG'] = False
  app.config['TESTING'] = False

cors = CORS(app, resources={r"/*": {"origins":"*"}})
app.register_blueprint(SWAGGERUI_BLUEPRINT, url_prefix=SWAGGER_URL)

@db_session
def add_new_suggestion(title_inp, text_inp, author_inp):
    return Suggestion(timestamp=datetime.utcnow(), title=title_inp, text=text_inp, author=author_inp), 201


class SuggestionRes(Resource):
    @db_session
    @login_required
    def get(self, id):
        return suggestion_to_json(Suggestion[id])


    @db_session
    @login_required
    def put(self, id):
        try:
            suggestion = Suggestion[id]
            suggestion.title = request.json["title"]
            suggestion.text = request.json["text"]
            suggestion.author = request.json["author"]
            commit()
        except ObjectNotFound:
            add_new_suggestion(request.json["title"], request.json["text"], request.json["author"])

class RemoveSuggestion(Resource):
  @db_session
  @login_required
  def delete(self, id):
    Suggestion[id].delete()

class RemoveSuggestions(Resource):
  @db_session
  @login_required
  def put(self):
    for id in request.json['ids']:
      try:
        Suggestion[id].delete()
      except ObjectNotFound:
        return 'Cound not find the suggestion with id %s' % (id), 400

class SuggestionResList(Resource):
    @db_session
    def post(self):
        try:
          add_new_suggestion(request.json["title"], request.json["text"], request.json["author"])

          if config.ENABLE_EMAIL == "True":
              respons = sendMail( \
                  PRIT_MAIL, \
                  SENDER_MAIL, \
                  "SuggestIT: %s" % (request.json["title"]), \
                  "%s\n\nPosted by %s" % (request.json["text"], request.json["author"]))          
              
              respons.raise_for_status()
              print("A mail has been sent to P.R.I.T.")
        except RequestException as err:
          print("Failed to send mail: \n %s" % err)
        except Exception as err:
          print("Failed to create suggestion: \n %s" % err)
          return "Could not create suggestion", 400
        
        return "Suggestion created", 201
    @db_session
    @login_required
    def get(self):
        return jsonify([suggestion_to_json(s) for s in Suggestion.select(lambda t : True)])

class Authentication(Resource):
    def put(self):
        if request.json["password"] == config.PRIT_PASSWORD:
            return jsonify(key=config.PRIT_AUTH_KEY)

        return "You are not P.R.I.T.", 401
    
    @login_required
    def get(self):
        return 'You are P.R.I.T.', 200

api.add_resource(SuggestionRes, '/api/<string:id>')
api.add_resource(Authentication, '/api/authenticate')
api.add_resource(SuggestionResList, '/api/')
api.add_resource(RemoveSuggestion, '/api/delete/<string:id>')
api.add_resource(RemoveSuggestions, '/api/delete')

if __name__ == '__main__':
    app.run(host='0.0.0.0')
