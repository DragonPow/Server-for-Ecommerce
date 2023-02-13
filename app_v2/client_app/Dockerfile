FROM node:16-alpine as builder
RUN node -v
ENV NODE_ENV=production
WORKDIR /app
COPY ["package.json", "package-lock.json", "./"]
RUN npm i -g npm@9.4.2
RUN npm install --legacy-peer-deps --production --silent && mv node_modules ../
COPY . .
RUN npm run build

FROM nginx:stable-alpine as publish
ENV NODE_ENV production
WORKDIR /usr/share/nginx/html
COPY ./conf/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/build/ .
EXPOSE 3000
ENTRYPOINT [ "nginx", "-g", "daemon off;" ]
