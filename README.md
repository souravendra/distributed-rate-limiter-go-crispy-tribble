# Distributed Rate Limiter in Go (Token Bucket + Redis)

This is an implementation of a **distributed rate limiter** using the **Token Bucket algorithm**, backed by **Redis**, and built in **Go**. It's designed to be efficient, scalable, and production-ready.

---

## Highlights

- Token Bucket algorithm using Redis atomic operations
- Clean architecture using Go interfaces and patterns
- Middleware-friendly HTTP integration
- Flexible configuration with Functional Options pattern
- Singleton limiter instance for safe shared use
- Works out-of-the-box with any Redis instance

---

## 📦 Design Patterns Used

| Pattern                   | Purpose                                                                 |
|--------------------------|-------------------------------------------------------------------------|
| Strategy Pattern         | Abstracts different rate limiting algorithms                            |
| Adapter Pattern          | Encapsulates Redis operations through a `Store` interface               |
| Functional Options       | Clean and extensible configuration of the limiter                       |
| Singleton                | Ensures only one limiter instance exists across the app                 |
| Middleware               | Integrates rate limiting logic into HTTP stack cleanly                  |

---

## 🛠️ Technologies

- **Go** 1.20+
- **Redis** 6+
- **go-redis** client

---

## 📁 Project Structure

```
.
├── main.go               # HTTP server + limiter wiring
├── limiter.go            # RateLimiter logic
├── README.md             # You're here
```

---

## 🧪 Rate Limiting Logic

- **Rate:** 2 requests/sec
- **Burst:** 2 requests
- Keyed by a hardcoded client ID (can be extended to IP/user ID)

Redis stores count keys like `rate:limiter:test-client` and uses `INCR` and `EXPIRE` to count and auto-reset.

---

## ▶️ Running the Project

### 1. **Start Redis Locally** (no Docker)

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

Now open your browser or run:
```bash
curl localhost:8080
```
Make more than 2 requests per second to get `429 Too Many Requests`.

---

## 📈 Example Output

```
Request allowed: 2025-04-01T18:26:15+05:30
Request allowed: 2025-04-01T18:26:16+05:30
Rate limit exceeded
```

---

## 💡 Possible Improvements / Extensions

- 🔁 Sliding window or Leaky bucket algorithms
- 🔑 Per-IP or API-key based throttling
- 🔧 Admin dashboard for live metrics
- 🌐 Redis Cluster or Sentinel support
- 📊 Prometheus metrics + Grafana dashboard


---

## 👨‍💻 Author


Crafted with care by me – designed to be readable and extensible. Will probably build on top of it later

---

