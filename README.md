# Distributed Rate Limiter in Go (Token Bucket + Redis)

![Go CI](https://github.com/souravendra/distributed-rate-limiter-go-crispy-tribble/actions/workflows/ci.yml/badge.svg)
<!-- [![codecov](https://codecov.io/gh/souravendra/distributed-rate-limiter-go-crispy-tribble/branch/main/graph/badge.svg)](https://codecov.io/gh/souravendra/distributed-rate-limiter-go-crispy-tribble) -->


This is an implementation of a **distributed rate limiter** using the **Token Bucket algorithm**, backed by **Redis**, and built in **Go**. It's designed to be efficient, scalable, and production-ready.

---

## Highlights

- Token Bucket algorithm using Redis atomic operations
- Clean architecture using Go interfaces and patterns
- Middleware-friendly HTTP integration
- Flexible configuration with Functional Options pattern
- Singleton limiter instance for safe shared use
- Works out-of-the-box with any Redis instance
- Unit tests using `stretchr/testify`
- Code coverage tracking via Codecov

---

## Design Patterns Used

| Pattern                   | Purpose                                                                 |
|--------------------------|-------------------------------------------------------------------------|
| Strategy Pattern         | Abstracts different rate limiting algorithms                            |
| Adapter Pattern          | Encapsulates Redis operations through a `Store` interface               |
| Functional Options       | Clean and extensible configuration of the limiter                       |
| Singleton                | Ensures only one limiter instance exists across the app                 |
| Middleware               | Integrates rate limiting logic into HTTP stack cleanly                  |

---

## Technologies

- **Go** 1.20+
- **Redis** 6+
- **go-redis** client
- **Testify** for unit testing
- **Codecov** for real-time test coverage

---

## Project Structure

```
.
├── main.go                   → Entrypoint for HTTP server
├── limiter/
│   ├── limiter.go           → Limiter struct & functional options
│   ├── strategy.go          → TokenBucket + RateLimiter interface
│   ├── store.go             → RedisStore + Store interface
│   └── limiter_test.go      → Unit tests for limiter
├── middleware/
│   └── ratelimit.go         → HTTP middleware
└── .github/workflows/ci.yml → GitHub Actions CI pipeline
└── Taskfile.yml             → Task runner for linting and testing
```
---

## Rate Limiting Logic

- **Rate:** 2 requests/sec
- **Burst:** 2 requests
- Keyed by a hardcoded client ID (can be extended to IP/user ID)

Redis stores count keys like `rate:limiter:test-client` and uses `INCR` and `EXPIRE` to count and auto-reset.

---

## Running Locally

### 1. **Start Redis**

(I wasnt using Docker for Redis)
#### macOS:
```bash
brew install redis
brew services start redis
```

---

### 2. **Run the Server**

```bash
task run # using Taskfile.yaml
```
### 3. **Spamming Requests & Other Commands**

 (can use a script if you want, I kept the window short to be able to trigger it manually):
```bash
curl localhost:8080 # spam
redis-cli FLUSHALL # clearing Redis
curl localhost:6379 # checking 
```
Make more than 2 requests per second to get `429 Too Many Requests`.

### 4. **Run Lint + Tests + Coverage**

```bash
task lint       # static analysis
```
```bash
task test       # unit tests with coverage
```
```bash
task coverage   # HTML coverage report
```

---

## Example Output

```
Request allowed: 2025-04-01T18:26:15+05:30
Request allowed: 2025-04-01T18:26:16+05:30
Rate limit exceeded
```

## Example Debug Logs

```
task: [run] go run main.go
Server running on :8080
New TTL set: 1s
Key: rate:limiter:test-client, Count: 1
New TTL set: 1s
Key: rate:limiter:test-client, Count: 1
Key: rate:limiter:test-client, Count: 2
Key: rate:limiter:test-client, Count: 3
Key: rate:limiter:test-client, Count: 4
Key: rate:limiter:test-client, Count: 5
New TTL set: 1s
Key: rate:limiter:test-client, Count: 1
Key: rate:limiter:test-client, Count: 2
```

---

## Possible Improvements / Extensions

- Sliding window or Leaky bucket algorithms
- Per-IP or API-key based throttling
- Admin dashboard for live metrics
- Redis Cluster or Sentinel support
- Prometheus metrics + Grafana dashboard


_Crafted with care by me – designed to be readable and extensible. Will probably build on top of it later_

---

