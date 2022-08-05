FROM node:18-alpine
RUN mkdir -p /app
COPY package*.json /app/

WORKDIR /app
ENV NODE_ENV=development
COPY . /app/
RUN npm install

# This CMD will attempt to make sure the node_modules are
# present which may be problematic if there is no connection
# to the Internet (npm install will hang); therefore, in
# such a use-case, specify a CMD with 'npm run dev' instead.
CMD ["/bin/sh -c 'npm install && npm run dev'"]
