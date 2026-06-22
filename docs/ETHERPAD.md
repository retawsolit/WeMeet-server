# Etherpad integration – troubleshooting

## "no active etherpad host found" / API code 4 ("no or wrong API Key")

### 1. Config: use `api_key` (recommended)

WeMeet supports two ways to authenticate with Etherpad:

- **API key** (Etherpad standard): set `api_key` in `config.yaml` under `shared_notepad.etherpad_hosts[]`. The server will send it as the `apikey` query parameter. The value must match the content of `APIKEY.txt` inside the Etherpad container.
- **OIDC**: set `client_id` and `client_secret` to match your Etherpad SSO/oidc config (e.g. in `settings.json`). The server will obtain a Bearer token from `/oidc/token` and use it for API calls.

If you only have an API key (e.g. from `APIKEY.txt`), configure it like this and you can leave `client_id` / `client_secret` commented out:

```yaml
shared_notepad:
  enabled: true
  etherpad_hosts:
    - id: "node_01"
      host: "http://192.168.1.119:9001"
      api_key: "YOUR_KEY_FROM_APIKEY_TXT"
      # client_id: "WeMeet"
      # client_secret: "..."
```

### 2. File permissions for `APIKEY.txt` in the container

Etherpad runs as user `etherpad`. If `APIKEY.txt` is owned by `root`, the process may not be able to read it and the HTTP API will reject requests (e.g. code 4).

**Fix:** ensure the file is readable by the etherpad user. For example, after starting the container:

```bash
docker exec -u root <etherpad_container_name> chown etherpad:etherpad /opt/etherpad-lite/APIKEY.txt
docker restart <etherpad_container_name>
```

Or in your image/build: create or copy `APIKEY.txt` before switching to `USER etherpad`, and ensure it is owned by `etherpad:etherpad`.

### 3. Network / host

- `config.yaml` uses `host: "http://192.168.1.119:9001"`. The WeMeet server (e.g. `wemeet-api` container) must be able to reach this address (same Docker network or reachable IP).
- If you run everything on one host, you can use `http://etherpad:9001` when the Etherpad service name in Docker Compose is `etherpad` and they share a network. Using the host IP (e.g. `192.168.1.119`) is fine if the server is not in the same Compose network.

### 4. Build / image

- Changes to `APIKEY.txt` on the host only take effect if the container uses the updated file (e.g. bind mount) or if you rebuild the image and recreate the container.
- After fixing config or permissions, restart the WeMeet server and the Etherpad container so both use the new settings.
