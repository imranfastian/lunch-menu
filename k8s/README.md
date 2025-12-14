# Kubernetes Deployment Guide for Lunch Menu API

This guide explains how to deploy the Lunch Menu API and PostgreSQL database on Kubernetes, using best practices for configuration, security, and local development.  
It covers the purpose and usage of each file in the `k8s/` directory, environment setup, secrets/configmaps, and testing instructions.

---

## Directory Structure

```
k8s/
├── api-deployment.yaml        # Deploys your Go API Pod(s)
├── api-service.yaml           # Internal service to expose your API inside cluster
├── configmap.yaml             # Non-sensitive configuration for the API
├── ingress.yaml               # (Optional) Routes external traffic into your cluster
├── postgres-deployment.yaml   # Deploys the Postgres Pod (container)
├── postgres-pvc.yaml          # Persistent volume claim for Postgres data
├── postgres-service.yaml      # Exposes Postgres internally to other Pods
├── secret.yaml                # Sensitive credentials for the API (DB creds, JWT secret)
```

---

## File Responsibilities

| File                         | Purpose                                                                                       |
| ---------------------------- | --------------------------------------------------------------------------------------------- |
| **postgres-deployment.yaml** | Runs PostgreSQL container with attached PVC and configures readiness/liveness probes.         |
| **postgres-service.yaml**    | Exposes PostgreSQL as a `ClusterIP` service (internal access only).                           |
| **postgres-pvc.yaml**        | Ensures persistent storage for the Postgres database data.                                    |
| **api-deployment.yaml**      | Deploys your Go API (replicas, container image, env vars, etc.).                              |
| **api-service.yaml**         | Exposes your API to the cluster (often `ClusterIP` for Ingress or `LoadBalancer` for public). |
| **configmap.yaml**           | Defines environment variables and configuration used by API deployment.                       |
| **secret.yaml**              | Contains sensitive info (passwords, API keys, tokens).                                        |
| **ingress.yaml**             | Defines routing rules (e.g. `/api/` → `api-service`) and optionally TLS certificates.         |

---

## ConfigMap vs Secret

- **ConfigMap**: Stores non-sensitive configuration (e.g., DB name, host, cookie settings).  
  Used for environment variables that are safe to expose.
- **Secret**: Stores sensitive data (e.g., DB password, JWT secret).  
  Values are base64-encoded and mounted as environment variables.

**How to create base64 values (Windows PowerShell):**

```powershell
[Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("your-value"))
```

---

## Request Lifecycle Diagram

```
[ Browser / Postman ]
        │
        ▼
[ Ingress Controller (NGINX) ]
        │
        ▼
[ api-service (ClusterIP) ]
        │
        ▼
[ api-deployment Pod(s) ]
        │
        ▼
[ postgres-service ] ─► [ postgres-deployment + PVC ]
```

---

## Kubernetes Commands

```sh
# Build and push Docker image
docker build --no-cache -t imranfastian1982/lunch-menu-api:latest .
docker push imranfastian1982/lunch-menu-api:latest

# Apply Kubernetes manifests
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml
kubectl apply -f k8s/api-deployment.yaml
kubectl apply -f k8s/api-service.yaml
kubectl apply -f k8s/ingress.yaml

# Delete deployments
kubectl delete deployment lunch-menu-api
kubectl delete deployment lunch-menu-postgres

# Check status
kubectl get pods
kubectl get svc
kubectl get ingress

# Port-forward for local testing
kubectl port-forward svc/lunch-menu-api 8000:8000

# Restart deployment after image update
kubectl rollout restart deployment/lunch-menu-api
```

---

## Testing the API

- Access Swagger UI at [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html) (adjust host/port as needed).
- Use tools like Postman or curl to test endpoints.
- `kubectl port-forward` exposes your service on localhost for quick testing, but is not production-ready.
- Ingress provides routing, TLS, and load balancing for external access.
- For Kubernetes, use `configmap.yaml` and `secret.yaml` instead of `.env`.
- For local Docker Compose, `.env` is still used.

---

## Additional Recommendations

- **Security:**

  - Use Secrets for all sensitive data.
  - Set `COOKIE_SECURE=true` and `COOKIE_HTTPONLY=true` for production.
  - Restrict CORS to trusted domains.
  - Use resource limits in deployments.
  - Add liveness/readiness probes.
  - Use NetworkPolicies for pod communication.
  - Keep images updated and scan for vulnerabilities.

- **Scaling:**

  - Increase `replicas` in `api-deployment.yaml` for high availability.
  - Use a LoadBalancer service or Ingress for public access.

- **Troubleshooting:**
  - Check pod logs with `kubectl logs <pod-name>`.
  - Use `kubectl describe` for detailed resource info.
  - Ensure your ingress controller is installed and running.

---

## Summary Table

| File                         | Purpose                            |
| ---------------------------- | ---------------------------------- |
| k8s/configmap.yaml           | Non-sensitive config for API/DB    |
| k8s/secret.yaml              | Sensitive secrets (passwords, JWT) |
| k8s/postgres-deployment.yaml | PostgreSQL deployment              |
| k8s/postgres-service.yaml    | PostgreSQL service                 |
| k8s/postgres-pvc.yaml        | Persistent volume for Postgres     |
| k8s/api-deployment.yaml      | API deployment                     |
| k8s/api-service.yaml         | API service                        |
| k8s/ingress.yaml             | Ingress for external access        |

---

## Self-Signed Certificate for HTTPS Ingress

To enable HTTPS for your API in Kubernetes, you can use a self-signed TLS certificate for local development. This allows secure access to your endpoints via `https://localhost`, and is required for testing secure cookies and production-like setups.

---

### How to Create a Self-Signed Certificate

You can customize the certificate's validity period (e.g., 2 years) and subject details.  
**Example command for a 2-year (730 days) certificate:**

```sh
openssl req -x509 -nodes -days 730 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=localhost/O=localdev"
```

- **Common Name (CN):** `localhost`
- **Organization (O):** `localdev`
- **Expires On:** 2 years from creation (customizable with `-days`)

---

### Create Kubernetes TLS Secret

```sh
kubectl create secret tls lunch-menu-tls --cert=tls.crt --key=tls.key
```

---

### Update Ingress for TLS

Add the following to your `k8s/ingress.yaml`:

```yaml
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - localhost
      secretName: lunch-menu-tls
  rules:
    - host: localhost
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: lunch-menu-api
                port:
                  number: 8000
```

Apply the changes:

```sh
kubectl apply -f k8s/ingress.yaml
```

---

### Test in Browser

- Go to [https://localhost/swagger/index.html](https://localhost/swagger/index.html)
- You will see a browser warning about the self-signed certificate. Accept the warning to proceed.
- Certificate details (as shown in browser):
  - **Common Name (CN):** `localhost`
  - **Organization (O):** `localdev`
  - **Issued On:** (e.g.) Thursday, October 16, 2025
  - **Expires On:** (e.g.) Friday, October 16, 2026 (or as set by `-days`)

---

### Notes

- For longer validity, increase the `-days` value in the OpenSSL command.
- For production, use a certificate from a trusted CA (e.g., Let's Encrypt) and set your real domain name.
- Self-signed certificates are suitable for local development and testing only.

---

**Tip:**  
You can view certificate details in your browser by clicking the padlock icon next to the URL and selecting "Certificate" or "View Certificate".

---
