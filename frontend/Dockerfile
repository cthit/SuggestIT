FROM node:latest as build

COPY /package.json /

RUN yarn install --production --network-timeout 1000000000
RUN yarn global add react-scripts 

COPY . .
RUN yarn build

FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build ./build /usr/share/nginx/html