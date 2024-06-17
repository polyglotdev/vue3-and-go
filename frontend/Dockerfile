# Use the specified Node.js version
FROM node:20.14.0

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and yarn.lock files
COPY package.json yarn.lock ./

# Install dependencies using yarn
RUN yarn install

# Copy the rest of the application code
COPY . .

# Build the Vue.js application
RUN yarn build

# Expose the port the application will run on
EXPOSE 4173

# Start the application using the 'preview' script
CMD ["yarn", "preview"]
