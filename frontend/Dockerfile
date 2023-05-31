FROM node:16-alpine as build

COPY package.json yarn.lock /

RUN yarn install --production --network-timeout 1000000000
RUN yarn global add react-scripts 

COPY . .
RUN yarn build

FROM nginx:alpine
COPY default.conf /etc/nginx/conf.d/default.conf
COPY --from=build ./build /usr/share/nginx/html