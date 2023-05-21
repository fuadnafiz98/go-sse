# Golang SSE (Server Send Events)

## Testing

- On Server:
  - `curl http://localhost:8888/?id=1`
    - creates new client
    - sends messages to frontend.
- On Client:
  - `const client = new EventSource("http://localhost:8888/events?id=1")`
  - `client.onmessage = function (msg) { console.log(msg) }`
  - `client.onerror = function (msg) { console.log(msg) }`
  - this will start listening for messages from backend.
