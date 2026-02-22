# ==============================================
# Stage 1: Install dependencies
# ==============================================
FROM node:20-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci && npm cache clean --force

# ==============================================
# Stage 2: Build the Next.js app
# ==============================================
FROM node:20-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Build arguments for environment variables (baked at build time)
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_WS_URL
ARG NEXT_PUBLIC_META_APP_ID
ARG NEXT_PUBLIC_TIKTOK_CLIENT_KEY

ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_WS_URL=$NEXT_PUBLIC_WS_URL
ENV NEXT_PUBLIC_META_APP_ID=$NEXT_PUBLIC_META_APP_ID
ENV NEXT_PUBLIC_TIKTOK_CLIENT_KEY=$NEXT_PUBLIC_TIKTOK_CLIENT_KEY
ENV NEXT_TELEMETRY_DISABLED=1

# Ensure public directory exists
RUN mkdir -p public

RUN npm run build

# ==============================================
# Stage 3: Production image (standalone)
# ==============================================
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create non-root user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy standalone build output
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Standalone mode uses server.js
CMD ["node", "server.js"]
