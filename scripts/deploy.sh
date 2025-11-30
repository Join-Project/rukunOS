#!/bin/bash

# Script Deployment untuk VPS
# Usage: ./deploy.sh [command]
#   setup    - Setup awal (clone, install dependencies)
#   migrate  - Jalankan database migrations
#   start    - Start semua services
#   stop     - Stop semua services
#   restart  - Restart semua services
#   update   - Update code dan restart
#   logs     - Tampilkan logs
#   backup   - Backup database

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
cd "$PROJECT_DIR"

COMPOSE_FILE="docker-compose.prod.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    log_error "docker-compose tidak ditemukan. Silakan install Docker Compose terlebih dahulu."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    log_warn "File .env tidak ditemukan. Membuat dari .env.example..."
    if [ -f .env.example ]; then
        cp .env.example .env
        log_warn "Silakan edit file .env dengan konfigurasi yang sesuai sebelum melanjutkan."
        exit 1
    else
        log_error "File .env.example tidak ditemukan."
        exit 1
    fi
fi

# Setup function
setup() {
    log_info "Setup awal aplikasi..."
    
    # Check if Docker is running
    if ! docker info &> /dev/null; then
        log_error "Docker tidak berjalan. Silakan start Docker terlebih dahulu."
        exit 1
    fi
    
    log_info "Docker sudah berjalan ✓"
    
    # Start database first
    log_info "Starting database container..."
    docker-compose -f $COMPOSE_FILE up -d db
    
    log_info "Menunggu database siap..."
    sleep 10
    
    log_info "Setup selesai! Selanjutnya jalankan: ./deploy.sh migrate"
}

# Migrate function
migrate() {
    log_info "Menjalankan database migrations..."
    
    # Check if database is running
    if ! docker-compose -f $COMPOSE_FILE ps db | grep -q "Up"; then
        log_error "Database container tidak berjalan. Jalankan: ./deploy.sh setup"
        exit 1
    fi
    
    MIGRATIONS_DIR="$PROJECT_DIR/backend/migrations"
    
    # List of migrations in order
    MIGRATIONS=(
        "001_create_core_tables.sql"
        "001_init_schema.sql"
        "002_create_rbac_tables.sql"
        "003_create_units_table.sql"
        "004_seed_default_roles.sql"
        "005_seed_demo_tenant.sql"
        "006_add_super_admin_and_fix_roles.sql"
        "007_create_bills_table.sql"
        "008_alter_bills_due_date_nullable.sql"
        "008_create_billing_templates_table.sql"
        "009_create_announcements_table.sql"
        "010_create_visitor_logs_table.sql"
        "011_create_panic_alerts_table.sql"
        "012_create_complaints_table.sql"
        "013_create_document_requests_table.sql"
        "014_update_billing_templates_add_fields.sql"
        "015_create_billing_template_amount_rules.sql"
        "016_add_bill_number_to_bills.sql"
    )
    
    # Load environment variables
    set -a
    source .env
    set +a
    
    DB_USER=${DB_USER:-rukunos_user}
    DB_NAME=${DB_NAME:-rukunos_db}
    
    for migration in "${MIGRATIONS[@]}"; do
        MIGRATION_FILE="$MIGRATIONS_DIR/$migration"
        
        if [ ! -f "$MIGRATION_FILE" ]; then
            log_warn "Migration file tidak ditemukan: $migration (skipping...)"
            continue
        fi
        
        log_info "Running migration: $migration"
        
        if docker-compose -f $COMPOSE_FILE exec -T db psql -U "$DB_USER" -d "$DB_NAME" < "$MIGRATION_FILE" 2>&1 | grep -q "ERROR"; then
            # Check if it's a "already exists" error (which is OK)
            if docker-compose -f $COMPOSE_FILE exec -T db psql -U "$DB_USER" -d "$DB_NAME" < "$MIGRATION_FILE" 2>&1 | grep -q "already exists\|duplicate\|already exists"; then
                log_warn "  Migration sudah diterapkan, skipping..."
            else
                log_error "  Migration $migration gagal!"
                exit 1
            fi
        else
            log_info "  ✓ Migration $migration selesai"
        fi
    done
    
    log_info "Semua migrations selesai!"
}

# Start function
start() {
    log_info "Starting semua services..."
    docker-compose -f $COMPOSE_FILE up -d --build
    log_info "Services sudah di-start!"
    log_info "Cek status: docker-compose -f $COMPOSE_FILE ps"
    log_info "Cek logs: docker-compose -f $COMPOSE_FILE logs -f"
}

# Stop function
stop() {
    log_info "Stopping semua services..."
    docker-compose -f $COMPOSE_FILE down
    log_info "Services sudah di-stop!"
}

# Restart function
restart() {
    log_info "Restarting semua services..."
    docker-compose -f $COMPOSE_FILE restart
    log_info "Services sudah di-restart!"
}

# Update function
update() {
    log_info "Updating aplikasi..."
    
    # Pull latest code
    if [ -d .git ]; then
        log_info "Pulling latest code dari Git..."
        git pull origin main || log_warn "Git pull gagal atau tidak ada perubahan"
    else
        log_warn "Direktori ini bukan Git repository"
    fi
    
    # Rebuild and restart
    log_info "Rebuilding containers..."
    docker-compose -f $COMPOSE_FILE down
    docker-compose -f $COMPOSE_FILE up -d --build
    
    log_info "Update selesai!"
}

# Logs function
logs() {
    docker-compose -f $COMPOSE_FILE logs -f
}

# Backup function
backup() {
    log_info "Backing up database..."
    
    set -a
    source .env
    set +a
    
    DB_USER=${DB_USER:-rukunos_user}
    DB_NAME=${DB_NAME:-rukunos_db}
    BACKUP_DIR="$PROJECT_DIR/backups"
    
    mkdir -p "$BACKUP_DIR"
    
    BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"
    
    docker-compose -f $COMPOSE_FILE exec -T db pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
        log_info "Backup berhasil: $BACKUP_FILE"
        
        # Compress backup
        gzip "$BACKUP_FILE"
        log_info "Backup dikompres: ${BACKUP_FILE}.gz"
    else
        log_error "Backup gagal!"
        exit 1
    fi
}

# Main command handler
case "${1:-}" in
    setup)
        setup
        ;;
    migrate)
        migrate
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    update)
        update
        ;;
    logs)
        logs
        ;;
    backup)
        backup
        ;;
    *)
        echo "Usage: $0 {setup|migrate|start|stop|restart|update|logs|backup}"
        echo ""
        echo "Commands:"
        echo "  setup    - Setup awal (start database)"
        echo "  migrate  - Jalankan database migrations"
        echo "  start    - Start semua services"
        echo "  stop     - Stop semua services"
        echo "  restart  - Restart semua services"
        echo "  update   - Update code dan restart"
        echo "  logs     - Tampilkan logs"
        echo "  backup   - Backup database"
        exit 1
        ;;
esac

