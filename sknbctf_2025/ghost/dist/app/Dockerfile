FROM node:24-slim

WORKDIR /app
COPY server.js .

EXPOSE 8080
CMD ["node", "--insecure-http-parser", "server.js"]
