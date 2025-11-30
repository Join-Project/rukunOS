# 🚀 Quick Deploy Guide - VPS

Panduan cepat untuk deploy RukunOS ke VPS.

## Langkah-langkah di VPS

### 1. Install Docker & Docker Compose

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verifikasi
docker --version
docker-compose --version
```

### 2. Clone Repository

```bash
cd /var/www
sudo git clone https://github.com/Join-Project/rukunOS.git rukunos-app
cd rukunos-app
sudo chown -R $USER:$USER /var/www/rukunos-app
```

### 3. Setup Environment

```bash
cp .env.example .env
nano .env
```

**Isi .env dengan:**
- `DB_PASSWORD`: Password database yang kuat
- `JWT_SECRET`: Generate dengan `openssl rand -base64 32`
- `FRONTEND_URL`: URL domain Anda (misal: https://yourdomain.com)
- `NUXT_PUBLIC_API_BASE`: URL API (misal: https://yourdomain.com/api)

### 4. Deploy Aplikasi

```bash
# Berikan permission execute
chmod +x scripts/deploy.sh

# Setup database
./scripts/deploy.sh setup

# Jalankan migrations
./scripts/deploy.sh migrate

# Start aplikasi
./scripts/deploy.sh start
```

### 5. Verifikasi

```bash
# Cek status
docker-compose -f docker-compose.prod.yml ps

# Cek logs
docker-compose -f docker-compose.prod.yml logs -f

# Test API
curl http://localhost:8086/health
```

## Update Aplikasi

```bash
cd /var/www/rukunos-app
git pull origin main
./scripts/deploy.sh update
```

## Backup Database

```bash
./scripts/deploy.sh backup
```

## Troubleshooting

**Database tidak connect:**
```bash
docker-compose -f docker-compose.prod.yml logs db
```

**Port sudah digunakan:**
```bash
sudo netstat -tulpn | grep :8086
```

**Lihat semua logs:**
```bash
./scripts/deploy.sh logs
```

## Dokumentasi Lengkap

Lihat [DEPLOYMENT.md](./DEPLOYMENT.md) untuk dokumentasi lengkap termasuk setup Nginx dan SSL.
