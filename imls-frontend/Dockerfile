FROM node:18-alpine as base

COPY package*.json /
EXPOSE 3000

FROM base as dev
WORKDIR /
ENV NODE_ENV=development
RUN npm install -g nodemon && npm install
COPY . /
RUN npm run build