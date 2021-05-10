FROM node:latest

WORKDIR /usr

COPY package.json yarn.lock ./

RUN yarn install
RUN yarn global add react-scripts

COPY . .

CMD yarn start
