const http = require('http');
const fs = require('fs');
const path = require('path');

const FLAG = process.env.FLAG || 'sknb{dummy}';

http.createServer((req, res) => {
  console.log(`${req.method} ${req.url}`);

  switch (req.url) {
    case '/flag':
      if (!req.headers['x-from-proxy']) {
        res.end(FLAG);
      } else {
        res.writeHead(403, { 'Content-Type': 'text/plain' });
        res.end('Forbidden');
      }
      return;
    
    case '/':
      res.end('Hello');
      return;

    default:
      res.writeHead(404, { 'Content-Type': 'text/plain' });
      res.end('Not Found');
      return;
  }

}).listen(8080, () => console.log('Backend on :8080'));
