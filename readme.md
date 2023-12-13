# Chat App w/ Webhook and Redis
### ~ A Small Tutorial Project ~

## Dependency
Docker redis:7.2.3-alpine

## Summary
Redis is not necessarily needed to implement Websocket. But with Redis, we could overcome the Websocket's bottleneck of limited simultaneous connection in a single server. By using Redis, communication between users across server could be done.
