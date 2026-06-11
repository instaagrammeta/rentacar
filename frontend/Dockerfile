# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM node:22-alpine AS build
WORKDIR /app

# Increase Node memory for large dependency trees
ENV NODE_OPTIONS="--max-old-space-size=4096"

# 1. Copy package files and .npmrc.
COPY package.json .npmrc ./

# 2. Install deps. Use --prefer-online to avoid stale caches and
#    --maxsockets 5 to prevent connection pool exhaustion on Alpine.
RUN npm install --prefer-online --maxsockets 5 --ignore-scripts 2>&1; \
    EXIT_CODE=$?; \
    if [ $EXIT_CODE -ne 0 ]; then \
      echo "=== NPM DEBUG LOG (last 200 lines) ==="; \
      tail -200 /root/.npm/_logs/*-debug-0.log 2>/dev/null; \
      exit $EXIT_CODE; \
    fi

# 3. Copy the rest of the project source.
COPY . .

# 4. Run nuxt prepare (generates .nuxt types) then build.
RUN npx nuxt prepare
RUN npm run build

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
