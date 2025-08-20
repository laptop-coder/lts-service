FROM node:24.6.0-bookworm-slim AS build
WORKDIR /frontend
COPY package.json package-lock.json .
RUN npm i
COPY vite.config.ts tsconfig.json index.html .
COPY src src
RUN npm run build

FROM nginx:1.29.1-alpine-slim AS production
COPY default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /frontend/dist /usr/share/nginx/html
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]
