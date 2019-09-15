FROM node:latest as build

COPY /package.json /

RUN yarn install --production --network-timeout 1000000000
RUN yarn global add react-scripts 

COPY . .
RUN yarn build

RUN yarn global add serve
CMD ["sh","build-startup.sh"]                                                                                                                                            130 ↵