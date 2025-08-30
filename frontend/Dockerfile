FROM node:24.7.0-alpine3.21 AS deps
WORKDIR /app
COPY package.json package-lock.json .
RUN npm i

FROM node:24.7.0-alpine3.21 AS builder
WORKDIR /app
COPY --from=deps /app/node_modules node_modules
COPY vite.config.ts tsconfig.json index.html .
COPY src src
RUN npx vite build

FROM nginx:1.29.1-alpine3.22-slim
COPY default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]
