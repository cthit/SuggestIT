FROM node:latest

RUN mkdir -p /usr/src/app/
RUN chown -R node /usr/src/app

USER node


WORKDIR /usr/src/app

COPY ./src/ .
COPY package.json .
COPY ./public/ .

RUN yarn install
RUN yarn global add react-scripts

EXPOSE 3000


CMD yarn start
