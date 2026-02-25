#!/bin/bash
# Kompakter Real-Time Log Generator für rsyslog

# Konfiguration
DB_CONF="-u ${DB_USER:-rsyslog} -p${DB_PASS:-password} ${DB_NAME:-Syslog}"
INTERVAL=10
LOGS_PER_BURST=3

# Daten-Pools für dynamische Sätze
ACT=("Access" "Connection" "Update" "Backup" "Request" "Login" "Sync")
OBJ=("database" "user session" "upstream API" "config file" "cache cluster" "firewall rule")
REASON=("timeout" "permission denied" "checksum mismatch" "invalid token" "success" "resource exhaustion")
HOSTS=("web-server" "db-server" "app" "mail-server" "fw" "lb")
TAGS=("sshd" "nginx" "mysqld" "postfix" "docker")

# Datenbank-Check (Einzeiler)
until mysqladmin ping -h ${DB_HOST:-localhost} --silent; do sleep 1; done

generate_log() {
    local sev=$(( (RANDOM % 6) ))
    local fac=$(( (RANDOM % 24) ))
    local msg="$(slice "${ACT[@]}") $(slice "${OBJ[@]}"): $(slice "${REASON[@]}") ($((1 + RANDOM % 999))ms)"
    local host="$(slice "${HOSTS[@]}")-0$((1 + RANDOM % 5))"
    local tag="$(slice "${TAGS[@]}")"
    local ra="$(date '+%Y-%m-%d %H:%M:%S')"
    local drt="$(date --date "$((RANDOM % 9)) seconds ago" '+%Y-%m-%d %H:%M:%S')"

    # Kompakter Aufruf via -e
    mysql $DB_CONF -e "INSERT INTO SystemEvents (ReceivedAt, DeviceReportedTime, Facility, Priority, FromHost, Message, EventSource, EventUser, SysLogTag, EventID, CustomerID, Importance) VALUES ('$ra', '$drt', $fac, $sev, '$host', '$msg', 'service-$((RANDOM%5))', 'user-$((RANDOM%10))', '$tag', $((RANDOM%5000)), 1, $((6-prio)));"

    echo "[$ra] $host $tag: [$sev|$fac] - $msg"
}

# Hilfsfunktion für Zufallswahl
slice() { local a=("$@"); echo "${a[$RANDOM % ${#a[@]}]}"; }

echo "Generator läuft... (Intervall: ${INTERVAL}s, Burst: $LOGS_PER_BURST)"

while true; do
    for i in $(seq 1 $LOGS_PER_BURST); do
        generate_log
        sleep 0.2 # Minimale Verzögerung für Realismus
    done
    sleep $INTERVAL
done
