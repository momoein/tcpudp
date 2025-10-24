## 🚀 Usage

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/momoein/tcpudp.git
cd tcpudp
```

### 2️⃣ Start the TCP and UDP Servers

Use Docker Compose to build and start both servers in detached mode:

```bash
docker compose up -d
```

This will spin up:
• TCP server on port 1241 → 11000
• UDP server on port 1242 → 11000 (UDP)

Check running services:

```bash
docker compose ps
```

### 3️⃣ Run the Client (from your host)

🔹 TCP Client:

```bash
go run main.go --mode=client --proto=tcp --addr=127.0.0.1:1241
```

🔹 UDP Client:

```bash
go run main.go --mode=client --proto=udp --addr=127.0.0.1:1242
```

### 4️⃣ Stop the Servers

```bash
docker compose down
```

### 💡 Notes

- --mode controls whether the program runs as "server" or "client".
- --proto selects the protocol ("tcp" or "udp").
- --addr sets the IP address and port to listen/connect.
- Optionally use --file to specify a file to send or save:
  --file image.jpg
