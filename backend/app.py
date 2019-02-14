from datetime import datetime

from db import Suggestion, suggestion_to_json

from flask import Flask, request, jsonify
from flask_restful import Resource, Api
from pony.orm import db_session, ObjectNotFound, commit

app = Flask(__name__)
api = Api(app)

app.config.setdefault('WTF_CSRF_METHODS', [])

@db_session
def add_new_suggestion(title_inp, text_inp, author_inp):
    return Suggestion(timestamp=datetime.utcnow(), title=title_inp, text=text_inp, author=author_inp)


class SuggestionRes(Resource):
    @db_session
    def get(self, id):
        return suggestion_to_json(Suggestion[id])

    @db_session
    def delete(self, id):
        Suggestion[id].delete()


    @db_session
    def put(self, id):
        try:
            suggestion = Suggestion[id]
            suggestion.title = request.json["title"]
            suggestion.text = request.json["text"]
            suggestion.author = request.json["author"]
            commit()
        except ObjectNotFound:
            add_new_suggestion(request.json["title"], request.json["text"], request.json["author"])


class SuggestionResList(Resource):
    @db_session
    def post(self):
        return suggestion_to_json(add_new_suggestion(request.json["title"], request.json["text"], request.json["author"]))


    @db_session
    def get(self):
        return jsonify([suggestion_to_json(s) for s in Suggestion.select(lambda t : True)])



api.add_resource(SuggestionRes, '/<string:id>')
api.add_resource(SuggestionResList, '/')


if __name__ == '__main__':
    app.run(host='0.0.0.0')
