# RukunOS - Platform Manajemen Komunitas Dinamis

Platform manajemen komunitas dengan Multi-Tenant Architecture dan Dynamic RBAC.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.20+ (untuk development lokal)
- Node.js 22+ (untuk development lokal)

### Development Setup

1. **Clone dan setup environment**
```bash
cd rukunos-app
cp .env.example .env
# Edit .env dengan konfigurasi yang sesuai
```

2. **Start services dengan Docker**
```bash
docker compose up -d
```

3. **Run migrations**
```bash
docker compose exec api sh -c "cd /app && psql -h db -U rukunos_user -d rukunos_db -f migrations/001_create_core_tables.sql"
docker compose exec api sh -c "cd /app && psql -h db -U rukunos_user -d rukunos_db -f migrations/002_create_rbac_tables.sql"
docker compose exec api sh -c "cd /app && psql -h db -U rukunos_user -d rukunos_db -f migrations/003_create_units_table.sql"
docker compose exec api sh -c "cd /app && psql -h db -U rukunos_user -d rukunos_db -f migrations/004_seed_default_roles.sql"
```

4. **Access aplikasi**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Health: http://localhost:8080/health

## 📁 Project Structure

```
rukunos-app/
├── backend/           # Go backend (Echo framework)
│   ├── db/           # Database connection
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # Middleware (JWT, Tenant, Permission)
│   ├── models/       # Data models
│   ├── migrations/   # Database migrations
│   └── main.go       # Entry point
├── frontend/         # Nuxt.js 3 frontend
│   ├── pages/        # Pages (auto-routing)
│   ├── components/   # Vue components
│   ├── stores/       # Pinia stores
│   └── composables/  # Composables
└── docker-compose.yml
```

## 🔧 Development

### Backend
```bash
cd backend
go mod download
go run main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## 📚 Documentation

Lihat folder `docs/` untuk dokumentasi lengkap:
- URS.md - User Requirements Specification
- API_CONTRACT.md - API Documentation
- DATABASE_SCHEMA.md - Database Schema
- TECHNICAL_ARCHITECTURE.md - Technical Architecture
- DEVELOPMENT_ROADMAP.md - Development Roadmap
- UI_UX_DESIGN.md - UI/UX Design & Wireframes

## 🗄️ Database Migrations

Migrations ada di `backend/migrations/`. Untuk run migrations:

```bash
# Via Docker
docker compose exec api sh -c "cd /app && psql -h db -U rukunos_user -d rukunos_db -f migrations/XXX_migration_name.sql"
```

## 🔐 Environment Variables

Lihat `.env.example` untuk daftar environment variables yang diperlukan.

## 📝 License

MIT License











