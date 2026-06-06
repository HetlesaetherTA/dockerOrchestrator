# dockerOrchestrator

### Start:

#### Warning!:
- Dev runs in containerized environment, it can change docker containers state, but not alter host fs
- Prod runs directly on host, use with caution.

```bash
# Dev environment (logs: visit http://localhost:42069)
docker compose up -d

# Prod environment (logs: journalctl -u dockerOrchestrator -f) 
sudo chmod +x ./deploy.sh
./deploy.sh
```
