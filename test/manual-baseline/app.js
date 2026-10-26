const http = require('http');

const PORT = process.env.PORT || 8080;

const server = http.createServer((req, res) => {
  if (req.url === '/healthz' || req.url === '/ready') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ status: 'ok', timestamp: new Date().toISOString() }));
  } else if (req.url === '/metrics') {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('# HELP http_requests_total Total number of HTTP requests\n# TYPE http_requests_total counter\nhttp_requests_total{status="200"} 1\n');
  } else {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ message: 'Hello from manual baseline orders API', version: '1.0.0' }));
  }
});

server.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});

