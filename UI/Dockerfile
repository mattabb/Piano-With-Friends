FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY ./ .

RUN npm run build

FROM nginx as production-stage
RUN mkdir /app
COPY --from=0 /app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf