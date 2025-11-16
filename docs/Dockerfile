FROM node:24.7.0-alpine3.21 AS deps
WORKDIR /app
COPY package.json package-lock.json .
RUN npm i

FROM node:24.7.0-alpine3.21 AS builder
WORKDIR /app
COPY --from=deps /app/node_modules node_modules
COPY .vitepress .vitepress
COPY src src
COPY .env* .
RUN npx vitepress build

FROM nginx:1.29.1-alpine3.22-slim
WORKDIR app
COPY nginx.conf.template .
COPY --from=builder /app/.vitepress/dist /usr/share/nginx/html
EXPOSE 32782
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
ENTRYPOINT /app/entrypoint.sh
