FROM node:14-alpine AS build

WORKDIR /app

ARG REACT_APP_API_HOST
ENV REACT_APP_API_HOST=$REACT_APP_API_HOST

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# Block npm ads
ENV OPEN_SOURCE_CONTRIBUTOR=true

COPY ./package*.json /app/

# Better npm install https://docs.npmjs.com/cli/ci.html
RUN npm ci

COPY ./ /app/

RUN npm run build

FROM nginx:1.21.0-alpine
COPY --from=build /app/build/ /var/www/
COPY ./nginx.conf /etc/nginx/nginx.conf
EXPOSE 80