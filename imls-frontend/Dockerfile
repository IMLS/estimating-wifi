FROM node:18-alpine
RUN mkdir -p /app
COPY package*.json /app/

WORKDIR /app
ENV NODE_ENV=development
RUN npm -g install
#COPY . /app/
#RUN npm run build
