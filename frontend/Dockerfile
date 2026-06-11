# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM node:22-alpine AS build
WORKDIR /app

# Enable pnpm via corepack. pnpm resolves Nuxt 4's large dependency tree
# (incl. the many oxc native bindings) far more reliably than npm.
RUN corepack enable && corepack prepare pnpm@9 --activate

ENV NODE_OPTIONS="--max-old-space-size=2048"

# 1. Copy manifest and install deps (skip scripts so nuxt prepare runs later).
COPY package.json ./
RUN pnpm install --no-frozen-lockfile --ignore-scripts

# 2. Copy the rest of the project source.
COPY . .

# 3. Generate Nuxt types, then build.
RUN pnpm exec nuxt prepare
RUN pnpm run build

# ---- Runtime stage ----
FROM node:22-alpine
WORKDIR /app

ENV NODE_ENV=production \
    HOST=0.0.0.0 \
    PORT=3000 \
    NUXT_PUBLIC_API_BASE=http://localhost:5000/api

COPY --from=build /app/.output ./.output

EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
