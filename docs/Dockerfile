FROM node:24.4.1-bookworm-slim AS build
WORKDIR /docs
COPY package.json package-lock.json .
RUN npm i
COPY .vitepress .vitepress
COPY src src
RUN npm run docs:build

FROM nginx:1.29.0-alpine-slim AS production
COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /docs/.vitepress/dist /usr/share/nginx/html
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]
