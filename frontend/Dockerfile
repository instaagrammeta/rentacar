# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM node:20-alpine AS build
WORKDIR /app

# Ensure npm can reach the registry (helps in restrictive Docker networks).
RUN npm config set registry https://registry.npmjs.org/ && \
    npm config set fetch-retries 5 && \
    npm config set fetch-retry-mintimeout 20000 && \
    npm config set fetch-retry-maxtimeout 120000

# 1. Copy package files and .npmrc, then install deps.
COPY package.json .npmrc ./
RUN npm install --ignore-scripts || \
    (echo "=== NPM DEBUG LOG ===" && cat /root/.npm/_logs/*.log 2>/dev/null && exit 1)

# 2. Copy the rest of the project source.
COPY . .

# 3. Run nuxt prepare (generates .nuxt types) then build.
RUN npx nuxt prepare
RUN npm run build

# ---- Runtime stage ----
FROM node:20-alpine
WORKDIR /app

ENV NODE_ENV=production \
    HOST=0.0.0.0 \
    PORT=3000 \
    NUXT_PUBLIC_API_BASE=http://localhost:5000/api

COPY --from=build /app/.output ./.output

EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
