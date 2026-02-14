#!/bin/bash
# Advanced Live Log Generator for rsyslog REST API
# Generates realistic syslog entries with ALL fields filled

set -e

# Configuration
INTERVAL=10
LOGS_PER_BURST=3
DB_USER="rsyslog"
DB_PASS="password"
DB_NAME="Syslog"

# Arrays for generation
HOSTS=("webserver01" "webserver02" "dbserver01" "appserver01" "mailserver01" "firewall01")
TAGS=("sshd" "nginx" "mysqld" "node" "postfix" "iptables" "systemd" "docker")
USERS=("admin" "deploy" "www-data" "root" "postgres" "nginx" "app-user")
EVENT_SOURCES=("auth-service" "web-service" "db-service" "app-service" "mail-service" "firewall")

# Helper functions
random_element() {
    local arr=("$@")
    echo "${arr[$RANDOM % ${#arr[@]}]}"
}

random_range() {
    echo $((RANDOM % ($2 - $1 + 1) + $1))
}

# Generate message based on priority
generate_message() {
    local priority=$1
    local msg=""

    case $priority in
        6) # Info
            local variants=(
                "User login successful"
                "HTTP request processed: 200 OK"
                "Database query completed in $(random_range 10 500)ms"
                "Service started on port $(random_range 3000 9000)"
                "Email sent successfully"
                "Backup completed"
                "Configuration reloaded"
                "Cache cleared"
            )
            msg=$(random_element "${variants[@]}")
            ;;
        5) # Notice
            local variants=(
                "Service restarted"
                "Certificate expires in $(random_range 7 30) days"
                "Disk usage: $(random_range 60 85)%"
                "Memory usage: $(random_range 60 90)%"
                "Connection pool size increased"
                "Rate limit reached"
            )
            msg=$(random_element "${variants[@]}")
            ;;
        4) # Warning
            local variants=(
                "Slow query detected: $(random_range 1000 5000)ms"
                "Failed login attempt"
                "Connection timeout"
                "Retry attempt failed"
                "Queue size growing"
                "Response time degraded"
            )
            msg=$(random_element "${variants[@]}")
            ;;
        3) # Error
            local variants=(
                "Connection refused"
                "Database error"
                "Failed to write to disk"
                "Service crashed"
                "Authentication failed"
                "Cannot allocate memory"
            )
            msg=$(random_element "${variants[@]}")
            ;;
        2) # Critical
            local variants=(
                "CRITICAL: Service down"
                "CRITICAL: Disk full"
                "CRITICAL: Database connection lost"
                "CRITICAL: Security breach detected"
            )
            msg=$(random_element "${variants[@]}")
            ;;
    esac

    echo "$msg"
}

# Generate realistic log entry
generate_log() {
    local priority=$1
    local host=$(random_element "${HOSTS[@]}")
    local tag=$(random_element "${TAGS[@]}")
    local user=$(random_element "${USERS[@]}")
    local event_source=$(random_element "${EVENT_SOURCES[@]}")
    local facility=1  # user facility

    # Generate timestamps with slight offset
    local received_at=$(date '+%Y-%m-%d %H:%M:%S')
    local device_time=$(date -d "2 seconds ago" '+%Y-%m-%d %H:%M:%S')

    # Generate message
    local message=$(generate_message $priority)

    # Generate realistic extended fields
    local customer_id=$((RANDOM % 100 + 1))
    local nt_severity=$((priority * 1000))
    local importance=$((6 - priority + 1))
    local event_category=$((RANDOM % 10 + 1))
    local event_id=0

    # Realistic Event IDs based on tag
    case $tag in
        sshd) event_id=$((4624 + RANDOM % 10)) ;;
        nginx) event_id=$((RANDOM % 300 + 200)) ;;
        mysqld) event_id=$((1000 + RANDOM % 100)) ;;
        *) event_id=$((RANDOM % 1000 + 1000)) ;;
    esac

    local system_id=$((RANDOM % 10 + 1))
    local info_unit_id=$((RANDOM % 5 + 1))

    # Escape message for SQL
    local escaped_message="${message//\'/\'\'}"

    # Insert into database with ALL fields
    mysql -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" 2>/dev/null <<SQLEOF
INSERT INTO SystemEvents (
    CustomerID, ReceivedAt, DeviceReportedTime, Facility, Priority,
    FromHost, Message, NTSeverity, Importance, EventSource, EventUser,
    EventCategory, EventID, SysLogTag, InfoUnitID, SystemID
) VALUES (
    $customer_id,
    '$received_at',
    '$device_time',
    $facility,
    $priority,
    '$host',
    '$escaped_message',
    $nt_severity,
    $importance,
    '$event_source',
    '$user',
    $event_category,
    $event_id,
    '$tag',
    $info_unit_id,
    $system_id
);
SQLEOF

    # Log to console
    local priority_label=""
    case $priority in
        6) priority_label="INFO" ;;
        5) priority_label="NOTICE" ;;
        4) priority_label="WARNING" ;;
        3) priority_label="ERROR" ;;
        2) priority_label="CRITICAL" ;;
    esac

    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$host] [$priority_label] [$tag] $message"
}

echo "=========================================="
echo "Advanced Live Log Generator Started"
echo "=========================================="
echo "Database: $DB_NAME"
echo "Interval: ${INTERVAL}s"
echo "Logs per burst: $LOGS_PER_BURST"
echo "Hosts: ${HOSTS[*]}"
echo ""
echo "Features:"
echo "  ✓ All 25+ fields filled"
echo "  ✓ Realistic Event IDs"
echo "  ✓ Extended columns populated"
echo "  ✓ Varied priorities (weighted)"
echo ""

# Wait for database
for i in {1..30}; do
    if mysql -u "$DB_USER" -p"$DB_PASS" -e "USE $DB_NAME" 2>/dev/null; then
        echo "✓ Database connection established"
        break
    fi
    [ $i -eq 30 ] && echo "✗ Database connection failed!" && exit 1
    sleep 1
done

echo ""
echo "Generating realistic logs..."
echo ""

# Weighted priorities (more info, less critical)
PRIORITIES=(6 6 6 6 5 5 4 3)

# Main loop
while true; do
    for i in $(seq 1 $LOGS_PER_BURST); do
        # Pick random weighted priority
        priority=${PRIORITIES[$RANDOM % ${#PRIORITIES[@]}]}
        generate_log $priority
        sleep 0.5
    done

    sleep $INTERVAL
done
