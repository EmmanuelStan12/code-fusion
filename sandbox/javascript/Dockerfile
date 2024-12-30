FROM node:22.12.0-alpine3.20

RUN apk add --no-cache \
    python3 \
    make \
    g++

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY ./src/ ./src/

COPY index.js .

CMD ["node", "index.js"]