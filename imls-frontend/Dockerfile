FROM node:18-buster-slim
RUN mkdir -p /app && chown -R 1000 /app

COPY package*.json /app/
WORKDIR /app

RUN npm install

COPY --chown=1000 . /app/
ENV NODE_ENV=development

USER 1000
CMD ["/bin/sh -c 'npm run dev'"]
