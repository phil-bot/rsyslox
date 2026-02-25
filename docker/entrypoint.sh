#!/bin/bash
# Docker test environment entrypoint.
# Sets up MariaDB, seeds test data, then starts rsyslox WITHOUT a config file
# so the setup wizard runs and can be tested end-to-end.
#
# RSYSLOX_ALLOW_REMOTE_SETUP=true (set in docker-compose.yml) allows the
# wizard to be reached from the Docker host, not just from localhost.

set -e

echo "================================================"
echo "rsyslox - Test Environment"
echo "================================================"
echo ""

# ── Check binary ─────────────────────────────────────────────────────────────
if [ ! -f /host-build/rsyslox ]; then
    echo "✗ ERROR: Binary not found at /host-build/rsyslox"
    echo ""
    echo "Build first:  make all"
    exit 1
fi

echo "[1/5] Installing API binary..."
cp /host-build/rsyslox /opt/rsyslox/rsyslox
chmod +x /opt/rsyslox/rsyslox
echo "✓ Binary installed ($(ls -lh /opt/rsyslox/rsyslox | awk '{print $5}'))"

# ── MariaDB ───────────────────────────────────────────────────────────────────
echo "[2/5] Starting MariaDB..."
mysqld_safe --datadir=/var/lib/mysql --user=mysql &
for i in {1..30}; do
    mysqladmin ping --silent 2>/dev/null && { echo "✓ MariaDB ready"; break; }
    [ $i -eq 30 ] && echo "✗ MariaDB timeout!" && exit 1
    sleep 1
done

# ── Database + table ─────────────────────────────────────────────────────────
DB_NAME="${DB_NAME:-Syslog}"
DB_USER="${DB_USER:-rsyslog}"
DB_PASS="${DB_PASS:-password}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
SERVER_PORT="${SERVER_PORT:-8000}"

echo "[3/5] Creating database, user and table..."
mysql <<SQL
CREATE DATABASE IF NOT EXISTS ${DB_NAME};
CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASS}';
GRANT ALL ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;
SQL

mysql "${DB_NAME}" <<'SQL'
CREATE TABLE IF NOT EXISTS SystemEvents (
    ID            int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CustomerID    bigint,
    ReceivedAt    datetime NULL,
    DeviceReportedTime datetime NULL,
    Facility      smallint NULL,
    Priority      smallint NULL,
    FromHost      varchar(60) NULL,
    Message       text,
    NTSeverity    int NULL,
    Importance    int NULL,
    EventSource   varchar(60),
    EventUser     varchar(60) NULL,
    EventCategory int NULL,
    EventID       int NULL,
    EventBinaryData text NULL,
    MaxAvailable  int NULL,
    CurrUsage     int NULL,
    MinUsage      int NULL,
    MaxUsage      int NULL,
    InfoUnitID    int NULL,
    SysLogTag     varchar(60),
    EventLogType  varchar(60),
    GenericFileName varchar(60),
    SystemID      int NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SQL
echo "✓ Database ready"

# ── Start rsyslox without config → setup wizard ───────────────────────────────
echo "[4/5] Starting rsyslox in setup wizard mode..."

/opt/rsyslox/rsyslox > /var/log/rsyslox.log 2>&1 &
API_PID=$!
sleep 2

if ! kill -0 $API_PID 2>/dev/null; then
    echo "✗ rsyslox failed to start!"
    cat /var/log/rsyslox.log
    exit 1
fi
echo "✓ rsyslox started (PID: $API_PID)"

echo ""
echo "================================================"
echo "✓ Environment Ready — Setup Wizard"
echo "================================================"
echo ""
echo "  Open in your browser:"
echo "  → http://localhost:${SERVER_PORT}"
echo ""
echo "  Use these database credentials in the wizard:"
echo "    DB Host:     ${DB_HOST}"
echo "    DB Port:     ${DB_PORT}"
echo "    DB Name:     ${DB_NAME}"
echo "    DB User:     ${DB_USER}"
echo "    DB Password: ${DB_PASS}"
echo ""
echo "[5/5] Starting log-generator.sh"
echo ""
/host-scripts/log-generator.sh
