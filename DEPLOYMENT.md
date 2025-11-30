# 🚀 Panduan Deployment ke VPS

Panduan lengkap untuk deploy RukunOS ke VPS server.

## 📋 Prerequisites

Sebelum memulai, pastikan VPS Anda sudah memiliki:
- Ubuntu 20.04+ atau Debian 11+
- Docker & Docker Compose terinstall
- Git terinstall
- Akses root atau user dengan sudo privileges

## 🔧 Setup Awal VPS

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

# Add user to docker group (ganti 'your-user' dengan username Anda)
sudo usermod -aG docker your-user

# Verifikasi instalasi
docker --version
docker-compose --version
```

### 2. Clone Repository

```bash
# Buat direktori untuk aplikasi
sudo mkdir -p /var/www
cd /var/www

# Clone repository
sudo git clone https://github.com/Join-Project/rukunOS.git rukunos-app
cd rukunos-app

# Set ownership (ganti 'your-user' dengan username Anda)
sudo chown -R your-user:your-user /var/www/rukunos-app
```

## ⚙️ Konfigurasi Environment

### 1. Buat File .env

```bash
cd /var/www/rukunos-app
cp .env.example .env
nano .env
```

### 2. Isi Environment Variables

```env
# Database Configuration
DB_HOST=db
DB_USER=rukunos_user
DB_PASSWORD=your_secure_password_here
DB_NAME=rukunos_db
DB_PORT=5432

# Application
PORT=8080
ENV=production
JWT_SECRET=your_jwt_secret_here_use_openssl_rand_base64_32

# Google OAuth (Optional)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback

# Frontend URL
FRONTEND_URL=https://yourdomain.com
NUXT_PUBLIC_API_BASE=https://yourdomain.com/api
```

**Generate JWT Secret:**
```bash
openssl rand -base64 32
```

## 🗄️ Setup Database & Migrations

### 1. Start Database Container

```bash
cd /var/www/rukunos-app
docker-compose -f docker-compose.prod.yml up -d db
```

Tunggu beberapa detik hingga database siap.

### 2. Jalankan Migrations

```bash
# Gunakan script migrasi otomatis
chmod +x scripts/deploy.sh
./scripts/deploy.sh migrate

# Atau manual:
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/001_create_core_tables.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/001_init_schema.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/002_create_rbac_tables.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/003_create_units_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/004_seed_default_roles.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/005_seed_demo_tenant.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/006_add_super_admin_and_fix_roles.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/007_create_bills_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/008_alter_bills_due_date_nullable.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/008_create_billing_templates_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/009_create_announcements_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/010_create_visitor_logs_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/011_create_panic_alerts_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/012_create_complaints_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/013_create_document_requests_table.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/014_update_billing_templates_add_fields.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/015_create_billing_template_amount_rules.sql
docker-compose -f docker-compose.prod.yml exec -T db psql -U rukunos_user -d rukunos_db < backend/migrations/016_add_bill_number_to_bills.sql
```

## 🚀 Start Aplikasi

### 1. Build dan Start Semua Services

```bash
cd /var/www/rukunos-app
docker-compose -f docker-compose.prod.yml up -d --build
```

### 2. Cek Status Services

```bash
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f
```

### 3. Verifikasi Aplikasi

```bash
# Cek API health
curl http://localhost:8086/health

# Cek frontend
curl http://localhost:3000
```

## 🔄 Update Aplikasi

### 1. Pull Latest Code

```bash
cd /var/www/rukunos-app
git pull origin main
```

### 2. Rebuild dan Restart

```bash
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build
```

### 3. Jalankan Migrations Baru (jika ada)

```bash
./scripts/deploy.sh migrate
```

## 🌐 Setup Nginx Reverse Proxy (Optional)

### 1. Install Nginx

```bash
sudo apt install nginx -y
```

### 2. Buat Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/rukunos
```

Isi dengan:

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Backend API
    location /api {
        proxy_pass http://localhost:8086;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. Enable Site

```bash
sudo ln -s /etc/nginx/sites-available/rukunos /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 4. Setup SSL dengan Let's Encrypt

```bash
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d yourdomain.com
```

## 📊 Monitoring & Maintenance

### Cek Logs

```bash
# Semua services
docker-compose -f docker-compose.prod.yml logs -f

# Specific service
docker-compose -f docker-compose.prod.yml logs -f api
docker-compose -f docker-compose.prod.yml logs -f app
docker-compose -f docker-compose.prod.yml logs -f db
```

### Backup Database

```bash
# Manual backup
docker-compose -f docker-compose.prod.yml exec db pg_dump -U rukunos_user rukunos_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Atau gunakan script
./scripts/backup-db.sh
```

### Restart Services

```bash
docker-compose -f docker-compose.prod.yml restart
```

### Stop Services

```bash
docker-compose -f docker-compose.prod.yml down
```

## 🔐 Security Checklist

- [ ] Ganti semua password default
- [ ] Generate JWT_SECRET yang kuat
- [ ] Setup firewall (UFW)
- [ ] Enable SSL/HTTPS
- [ ] Restrict database port (hanya localhost)
- [ ] Setup regular backups
- [ ] Monitor logs untuk suspicious activity
- [ ] Update sistem secara berkala

## 🆘 Troubleshooting

### Database connection error
```bash
# Cek apakah database container running
docker-compose -f docker-compose.prod.yml ps db

# Cek database logs
docker-compose -f docker-compose.prod.yml logs db

# Test connection
docker-compose -f docker-compose.prod.yml exec db psql -U rukunos_user -d rukunos_db -c "SELECT 1;"
```

### Port already in use
```bash
# Cek port yang digunakan
sudo netstat -tulpn | grep :8086
sudo netstat -tulpn | grep :3000
sudo netstat -tulpn | grep :5433

# Stop service yang menggunakan port
sudo systemctl stop <service-name>
```

### Container tidak start
```bash
# Cek logs
docker-compose -f docker-compose.prod.yml logs

# Rebuild containers
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml build --no-cache
docker-compose -f docker-compose.prod.yml up -d
```

## 📞 Support

Jika mengalami masalah, cek:
1. Logs aplikasi: `docker-compose -f docker-compose.prod.yml logs`
2. Status containers: `docker-compose -f docker-compose.prod.yml ps`
3. Resource usage: `docker stats`



