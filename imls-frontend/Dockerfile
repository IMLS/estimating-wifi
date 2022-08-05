FROM node:18-alpine
RUN mkdir -p /app
COPY package*.json /app/

WORKDIR /app
ENV NODE_ENV=development
COPY . /app/
RUN npm install
CMD ["/bin/sh -c 'npm install && npm run dev'"]
