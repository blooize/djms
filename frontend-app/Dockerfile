FROM node:lts-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json ./

# Install dependencies
RUN npm i

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Expose port 3000
EXPOSE 3000

# Start the development server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]
