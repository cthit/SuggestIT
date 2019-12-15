FROM node:latest

WORKDIR /usr

COPY package.json .

RUN yarn install
RUN yarn global add react-scripts

COPY . .

EXPOSE 3000


CMD yarn start
