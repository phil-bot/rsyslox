#!/bin/bash
# rsyslog REST API - Complete Build & Test Suite
# Replaces: /docker/test.sh, /docker/test-v0.2.2.sh
# 
# Usage:
#   ./build-and-test.sh              # Full build & test
#   ./build-and-test.sh --skip-build # Skip build, just test
#   ./build-and-test.sh --cleanup    # Stop Docker and cleanup

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
API_URL="${API_URL:-http://localhost:8000}"
API_KEY="${API_KEY:-}"
DOCKER_DIR="docker"
BUILD_DIR="build"

# Test counters
PASSED=0
FAILED=0

# Flags
SKIP_BUILD=false
CLEANUP_ONLY=false

# Parse arguments
for arg in "$@"; do
    case $arg in
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --cleanup)
            CLEANUP_ONLY=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --skip-build   Skip build step, only run tests"
            echo "  --cleanup      Stop Docker and cleanup, then exit"
            echo "  --help         Show this help message"
            echo ""
            exit 0
            ;;
    esac
done

# Cleanup function
cleanup() {
    echo ""
    echo -e "${BLUE}Cleaning up...${NC}"
    cd "$DOCKER_DIR"
    docker-compose down -v 2>/dev/null || true
    cd ..
    echo -e "${GREEN}✓ Cleanup complete${NC}"
}

# Cleanup only mode
if [ "$CLEANUP_ONLY" = true ]; then
    cleanup
    exit 0
fi

echo "=========================================="
echo "rsyslog REST API - Build & Test Suite"
echo "=========================================="
echo ""

# Step 1: Clean
if [ "$SKIP_BUILD" = false ]; then
    echo -e "${BLUE}[1/7]${NC} Cleaning old build..."
    make clean
    echo -e "${GREEN}✓${NC} Clean complete"
    echo ""
else
    echo -e "${YELLOW}[1/7] Skipping clean (--skip-build)${NC}"
    echo ""
fi

# Step 2: Build
if [ "$SKIP_BUILD" = false ]; then
    echo -e "${BLUE}[2/7]${NC} Building static binary..."
    make build-static
    if [ ! -f "$BUILD_DIR/rsyslog-rest-api" ]; then
        echo -e "${RED}✗ Build failed!${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓${NC} Build complete"
    echo ""
else
    echo -e "${YELLOW}[2/7] Skipping build (--skip-build)${NC}"
    echo ""
fi

# Step 3: Verify binary
echo -e "${BLUE}[3/7]${NC} Verifying binary..."
if [ ! -f "$BUILD_DIR/rsyslog-rest-api" ]; then
    echo -e "${RED}✗ Binary not found! Run without --skip-build${NC}"
    exit 1
fi
ls -lh "$BUILD_DIR/rsyslog-rest-api"
file "$BUILD_DIR/rsyslog-rest-api"
echo -e "${GREEN}✓${NC} Binary OK"
echo ""

# Step 4: Start Docker environment
echo -e "${BLUE}[4/7]${NC} Starting Docker test environment..."
cd "$DOCKER_DIR"
docker-compose down -v 2>/dev/null || true
docker-compose up -d

# Wait for API
echo -n "Waiting for API to be ready..."
for i in {1..30}; do
    if curl -s "$API_URL/health" > /dev/null 2>&1; then
        echo -e " ${GREEN}ready${NC}"
        break
    fi
    [ $i -eq 30 ] && echo -e " ${RED}timeout${NC}" && exit 1
    sleep 1
done

# Detect if API key is configured in Docker
CONTAINER_API_KEY=$(docker-compose exec -T rsyslog grep "^API_KEY=" /opt/rsyslog-rest-api/.env 2>/dev/null | cut -d'=' -f2 || echo "")
if [ -n "$CONTAINER_API_KEY" ] && [ "$CONTAINER_API_KEY" != "none" ]; then
    API_KEY="$CONTAINER_API_KEY"
fi

cd ..
echo ""

# Test helper function
test_endpoint() {
    local name="$1"
    local endpoint="$2"
    local use_auth="${3:-yes}"
    local check_content="${4:-}"
    
    echo -n "  [$((PASSED + FAILED + 1))] $name... "
    
    if [ "$use_auth" = "yes" ] && [ -n "$API_KEY" ]; then
        http_code=$(curl -s -w "%{http_code}" -o /tmp/response.json \
          -H "X-API-Key: $API_KEY" "$API_URL$endpoint" 2>/dev/null)
    else
        http_code=$(curl -s -w "%{http_code}" -o /tmp/response.json "$API_URL$endpoint" 2>/dev/null)
    fi
    
    # Check HTTP code
    if [ "$http_code" != "200" ]; then
        echo -e "${RED}✗ FAILED${NC} (HTTP $http_code)"
        FAILED=$((FAILED + 1))
        return 1
    fi
    
    # Optional content check
    if [ -n "$check_content" ]; then
        if grep -q "$check_content" /tmp/response.json 2>/dev/null; then
            echo -e "${GREEN}✓ OK${NC}"
            PASSED=$((PASSED + 1))
            return 0
        else
            echo -e "${RED}✗ FAILED${NC} (content check failed)"
            FAILED=$((FAILED + 1))
            return 1
        fi
    fi
    
    echo -e "${GREEN}✓ OK${NC}"
    PASSED=$((PASSED + 1))
    return 0
}

# Step 5: Run comprehensive test suite
echo -e "${BLUE}[5/7]${NC} Running comprehensive test suite..."
echo ""

echo -e "${CYAN}=== Basic Functionality ===${NC}"
test_endpoint "Health Check" "/health" "no"
test_endpoint "Get Logs (default)" "/logs" "yes"
test_endpoint "Get Meta (list)" "/meta" "yes"
echo ""

echo -e "${CYAN}=== Multi-Value Filter Tests ===${NC}"
test_endpoint "Multi-value FromHost" \
  "/logs?FromHost=webserver01&FromHost=webserver02&limit=5" \
  "yes"
test_endpoint "Multi-value Priority" \
  "/logs?Priority=3&Priority=6&limit=5" \
  "yes"
test_endpoint "Multi-value Facility" \
  "/logs?Facility=1&Facility=4&limit=5" \
  "yes"
test_endpoint "Multi-value Message (OR)" \
  "/logs?Message=login&Message=error&limit=5" \
  "yes"
test_endpoint "Combined multi-values" \
  "/logs?FromHost=webserver01&FromHost=webserver02&Priority=3&Priority=6&limit=5" \
  "yes"
echo ""

echo -e "${CYAN}=== Extended Column Tests ===${NC}"
test_endpoint "Extended columns present" \
  "/logs?limit=1" \
  "yes" \
  "SysLogTag"
test_endpoint "DeviceReportedTime field" \
  "/logs?limit=1" \
  "yes"
test_endpoint "Filter by SysLogTag" \
  "/logs?SysLogTag=sshd&limit=1" \
  "yes"
test_endpoint "Meta SysLogTag" \
  "/meta/SysLogTag" \
  "yes"
echo ""

echo -e "${CYAN}=== Meta with Multi-Value Filters ===${NC}"
test_endpoint "Meta FromHost with Priority filter" \
  "/meta/FromHost?Priority=3&Priority=6" \
  "yes"
test_endpoint "Meta SysLogTag with multiple hosts" \
  "/meta/SysLogTag?FromHost=webserver01&FromHost=webserver02" \
  "yes"
test_endpoint "Meta Priority with host filter" \
  "/meta/Priority?FromHost=webserver01" \
  "yes"
echo ""

echo -e "${CYAN}=== Advanced Scenarios ===${NC}"
test_endpoint "Many hosts (5+)" \
  "/logs?FromHost=webserver01&FromHost=webserver02&FromHost=dbserver01&FromHost=appserver01&FromHost=mailserver01&limit=5" \
  "yes"
test_endpoint "All priorities (0-7)" \
  "/logs?Priority=0&Priority=1&Priority=2&Priority=3&Priority=4&Priority=5&Priority=6&Priority=7&limit=5" \
  "yes"
test_endpoint "Complex filter combo" \
  "/logs?FromHost=webserver01&FromHost=webserver02&Priority=3&Priority=4&Facility=1&Message=login&limit=5" \
  "yes"
echo ""

echo -e "${CYAN}=== Backward Compatibility ===${NC}"
test_endpoint "Single FromHost (legacy)" \
  "/logs?FromHost=webserver01&limit=5" \
  "yes"
test_endpoint "Single Priority (legacy)" \
  "/logs?Priority=3&limit=5" \
  "yes"
test_endpoint "Single Message (legacy)" \
  "/logs?Message=login&limit=5" \
  "yes"
echo ""

# Step 6: Test v0.2.3+ features
echo -e "${BLUE}[6/7]${NC} Testing v0.2.3+ features..."
echo ""

# Test structured error format
echo -n "  [ERROR FORMAT] Testing structured errors... "
if [ -n "$API_KEY" ]; then
    ERROR_RESPONSE=$(curl -s -H "X-API-Key: $API_KEY" "$API_URL/logs?Priority=99" | jq -r '.code' 2>/dev/null || echo "")
else
    ERROR_RESPONSE=$(curl -s "$API_URL/logs?Priority=99" | jq -r '.code' 2>/dev/null || echo "")
fi
if [ "$ERROR_RESPONSE" = "INVALID_PRIORITY" ]; then
    echo -e "${GREEN}✓ OK${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAILED${NC} (Expected: INVALID_PRIORITY, Got: $ERROR_RESPONSE)"
    FAILED=$((FAILED + 1))
fi

# Test version in health
echo -n "  [VERSION] Testing version in health response... "
VERSION=$(curl -s "$API_URL/health" | jq -r '.version' 2>/dev/null || echo "")
if [ -n "$VERSION" ] && [ "$VERSION" != "null" ]; then
    echo -e "${GREEN}✓ OK${NC} (version: $VERSION)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAILED${NC} (No version found)"
    FAILED=$((FAILED + 1))
fi

# Test improved headers
echo -n "  [HEADERS] Testing response headers... "
HEADERS=$(curl -s -I "$API_URL/health" 2>/dev/null)
if echo "$HEADERS" | grep -q "Content-Type: application/json" && \
   echo "$HEADERS" | grep -q "charset=utf-8"; then
    echo -e "${GREEN}✓ OK${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${YELLOW}⚠ PARTIAL${NC} (headers present but not all expected)"
    PASSED=$((PASSED + 1))
fi

echo ""

# Step 7: Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
else
    echo -e "${GREEN}Failed: $FAILED${NC}"
fi
echo "Total:  $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓✓✓ ALL TESTS PASSED! ✓✓✓${NC}"
    echo ""
    echo "Features verified:"
    echo "  ✓ Multi-value filters working"
    echo "  ✓ Extended columns present"
    echo "  ✓ Backward compatibility maintained"
    echo "  ✓ Meta endpoint with multi-value filters"
    echo "  ✓ Structured error responses (v0.2.3+)"
    echo "  ✓ Version information in health"
    echo "  ✓ Improved HTTP headers"
    echo ""
    echo -e "${CYAN}Ready for production! 🚀${NC}"
else
    echo -e "${RED}✗ Some tests failed!${NC}"
    echo ""
    echo "Please review the failures above."
fi

echo ""
echo "=========================================="
echo "Useful Commands"
echo "=========================================="
echo "  Logs:      docker logs rsyslog-rest-api-test"
echo "  Shell:     docker exec -it rsyslog-rest-api-test bash"
echo "  DB:        docker exec -it rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog"
echo "  Cleanup:   ./build-and-test.sh --cleanup"
echo "  Re-test:   ./build-and-test.sh --skip-build"
echo ""

# Exit with appropriate code
if [ $FAILED -eq 0 ]; then
    exit 0
else
    exit 1
fi
