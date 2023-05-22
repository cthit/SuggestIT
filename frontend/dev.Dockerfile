FROM node:16-alpine

WORKDIR /usr

COPY package.json .

RUN yarn install
RUN yarn global add react-scripts

COPY . .

CMD yarn start
