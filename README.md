# 🪙 Kasi - The Privacy-First Budget Manager

![Go Version](https://img.shields.io/badge/Go-1.21-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)
![Size](https://img.shields.io/badge/Image_Size-~15MB-orange)

> **Simple. Self-Hosted. Secure.**

**Kasi** is a lightweight, open-source expense tracking solution designed for individuals, freelancers, and families who value data privacy. Unlike traditional finance apps, Kasi is designed to be **Self-Hosted**, meaning your financial data stays on your server, 100% under your control.

---

## 📸 Screenshots

| Dashboard (Mobile) | Reports (Desktop) |
|:---:|:---:|
| ![Dashboard](https://via.placeholder.com/300x600?text=Upload+Mobile+Screenshot) | ![Report](https://via.placeholder.com/600x400?text=Upload+Desktop+Screenshot) |
*(Upload screenshots to your repo and replace links above)*

---

## 🚀 Key Features

- **🔒 Privacy-First:** No tracking, no ads. Your data belongs to you.
- **📱 PWA Ready:** Install on **Android & iOS** as a native app (No App Store needed).
- **🚀 Ultra-Lightweight:** Built with **Go & SQLite**. Runs on <20MB RAM.
- **🌍 Global:** Multi-currency support (USD, EUR, LKR, GBP, etc.) & Multi-language.
- **👥 Multi-User:** Create accounts for family members or project partners.
- **📊 Reporting:** Generate clean, printable PDF reports.
- **💾 Backup & Restore:** Built-in tools to safeguard your data.

---

## 🐳 Quick Start (Docker)

Get Kasi running in seconds using Docker.

### 1. Run the Container
```bash
docker run -d \
  -p 8080:8080 \
  -v ./kasi-data:/app/data \
  -e SESSION_SECRET="ChangeThisToSomethingSecret" \
  --name kasi \
  --restart unless-stopped \
  amilakothalawalasolo/kasi:latest
