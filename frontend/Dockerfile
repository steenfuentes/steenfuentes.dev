# Build stage
FROM --platform=linux/amd64 node:23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies (including dev dependencies for build)
RUN npm ci

# Copy all frontend files
COPY . .



# Build the application
RUN npm run build

# Production stage
FROM node:23-alpine

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Set working directory
WORKDIR /app

# Copy package files and install production dependencies
COPY package*.json ./
RUN npm ci --only=production && \
    npm cache clean --force

# Copy built application from builder stage
COPY --from=builder --chown=nodejs:nodejs /app/build ./build
COPY --from=builder --chown=nodejs:nodejs /app/package.json ./package.json

# Switch to non-root user
USER nodejs

# Expose port (Cloud Run will set PORT env var)
EXPOSE 8080

# Set environment variables
ENV HOST=0.0.0.0

# Start the application
CMD ["node", "build/index.js"]