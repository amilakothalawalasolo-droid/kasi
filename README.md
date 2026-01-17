# 🪙 Kasi - The Privacy-First Budget Manager

![Go Version](https://img.shields.io/badge/Go-1.21-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Stable-success)

> **"Your Money. Your Server. Your Rules."**

**Kasi** is a lightweight, self-hosted expense tracking solution built for individuals, freelancers, and families who prioritize data privacy. Unlike traditional finance apps that store your data on third-party clouds, Kasi allows you to host your own financial data securely.

It is built with **Go (Golang)** and **SQLite**, ensuring it runs efficiently on low-resource hardware like a 512MB VPS, Raspberry Pi, or your local laptop.

🌐 **Official Website & Managed Services:** [www.amilakothalawala.work](https://www.amilakothalawala.work)

---

---

## ✨ Why Choose Kasi?

Most finance apps sell your data or require monthly subscriptions. **Kasi** is different:

- **🔒 100% Privacy:** Your data never leaves your server.
- **📱 PWA Native Experience:** Install on **iOS & Android** directly from the browser (No App Store required).
- **⚡ Ultra-Lightweight:** Docker image is small (~15MB) and uses minimal RAM (<20MB).
- **🌍 Global Support:** Supports **Multi-Currency** (USD, EUR, LKR, GBP, etc.) & Multi-Language.
- **📄 Smart Reports:** Generate printable PDF-ready reports for Projects, Weddings, or Monthly Expenses.
- **👥 Multi-User Support:** Perfect for couples, families, or small project teams.
- **💾 Auto Backups:** Built-in tools to download your database instantly.

---

## 🛠️ Tech Stack

Built for performance and simplicity:
- **Backend:** Go (Golang) 1.21
- **Database:** SQLite (Embedded, zero-config)
- **Frontend:** HTML/CSS (Server Side Rendered) + Vanilla JS
- **Deployment:** Docker & Docker Compose

---

## 🛡️ Security & Transparency

Security is a top priority for Kasi. The Docker image is scanned regularly using **Trivy** to ensure safety.

**Latest Scan Results:**

| Component | Status | Vulnerabilities |
| :--- | :--- | :--- |
| **Base OS (Alpine Linux)** | ✅ **Clean** | 0 Critical, 0 High, 0 Medium |
| **Application Logic** | ✅ **Safe** | 0 Critical, 0 High, 1 Medium* |

> *Note: The single "Medium" vulnerability (`CVE-2025-47909`) is an upstream issue in the `gorilla/csrf` library. There is currently no fix availab>

---


## 🚀 Getting Started (Self-Hosted)

You can run Kasi in seconds using Docker.

### Option A: Quick Run (Docker CLI)

```bash
docker run -d \
  -p 8080:8080 \
  -v ./kasi-data:/app/data \
  -e SESSION_SECRET="ReplaceWithStrongPassword" \
  --name kasi \
  --restart unless-stopped \
  amilakothalawalasolo/kasi:latest
