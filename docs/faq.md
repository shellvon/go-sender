# FAQ

**Q: How do I add a new provider?**  
A: Implement the Provider interface and register it with the sender. See [advanced.md](./advanced.md).

**Q: How do I debug failed sends?**  
A: Check error messages, enable debug logging, and consult provider API docs.

**Q: Can I use go-sender in a microservice?**  
A: Yes, go-sender is designed for both monoliths and microservices.

**Q: How do I send messages asynchronously?**  
A: Use the built-in queue middleware. See [middleware.md](./middleware.md).

**Q: How do I set a custom HTTP client?**  
A: Use `core.WithSendHTTPClient(myClient)` when sending.

**Q: How do I add rate limiting or retries?**  
A: Use the built-in middleware. See [middleware.md](./middleware.md).

**Q: Where can I find more examples?**  
A: See [examples.md](./examples.md).
