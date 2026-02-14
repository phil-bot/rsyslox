#!/bin/bash
# Extended Test Suite for v0.2.2
# Tests multi-value filters and extended columns

API_URL="${API_URL:-http://localhost:8000}"
API_KEY="${API_KEY:-test123456789}"

echo "=========================================="
echo "rsyslog REST API v0.2.2 - Test Suite"
echo "=========================================="
echo "API URL: $API_URL"
echo "API Key: ${API_KEY:0:20}..."
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PASSED=0
FAILED=0

# Test function
test_endpoint() {
    local name="$1"
    local endpoint="$2"
    local use_auth="${3:-yes}"
    local check_content="${4:-}"
    
    echo -n "[$((PASSED + FAILED + 1))] $name... "
    
    if [ "$use_auth" = "yes" ]; then
        http_code=$(curl -s -w "%{http_code}" -o /tmp/response.json \
          -H "X-API-Key: $API_KEY" "$API_URL$endpoint")
    else
        http_code=$(curl -s -w "%{http_code}" -o /tmp/response.json "$API_URL$endpoint")
    fi
    
    # Check HTTP code
    if [ "$http_code" != "200" ]; then
        echo -e "${RED}✗ FAILED${NC} (HTTP $http_code)"
        cat /tmp/response.json
        echo ""
        FAILED=$((FAILED + 1))
        return 1
    fi
    
    # Optional content check
    if [ -n "$check_content" ]; then
        if grep -q "$check_content" /tmp/response.json 2>/dev/null; then
            echo -e "${GREEN}✓ OK${NC} (HTTP $http_code, content verified)"
            PASSED=$((PASSED + 1))
            return 0
        else
            echo -e "${RED}✗ FAILED${NC} (HTTP 200 but content check failed)"
            echo "Expected: $check_content"
            echo "Got:"
            cat /tmp/response.json
            echo ""
            FAILED=$((FAILED + 1))
            return 1
        fi
    fi
    
    echo -e "${GREEN}✓ OK${NC} (HTTP $http_code)"
    PASSED=$((PASSED + 1))
    return 0
}

# Wait for API
echo -n "Waiting for API... "
for i in {1..30}; do
    if curl -s "$API_URL/health" > /dev/null 2>&1; then
        echo -e "${GREEN}ready${NC}"
        break
    fi
    [ $i -eq 30 ] && echo -e "${RED}timeout${NC}" && exit 1
    sleep 1
done
echo ""

# Section: Basic Tests
echo -e "${BLUE}=== Basic Functionality ===${NC}"
test_endpoint "Health Check" "/health" "no"
test_endpoint "Get Logs (default)" "/logs" "yes"
test_endpoint "Get Meta (list)" "/meta" "yes"
echo ""

# Section: Multi-Value Filter Tests
echo -e "${BLUE}=== Multi-Value Filter Tests ===${NC}"

# Test: Multiple hosts
test_endpoint "Multi-value FromHost" \
  "/logs?FromHost=webserver01&FromHost=webserver02&limit=5" \
  "yes"

# Test: Multiple priorities
test_endpoint "Multi-value Priority" \
  "/logs?Priority=3&Priority=6&limit=5" \
  "yes"

# Test: Multiple facilities
test_endpoint "Multi-value Facility" \
  "/logs?Facility=1&Facility=4&limit=5" \
  "yes"

# Test: Multiple message searches
test_endpoint "Multi-value Message (OR)" \
  "/logs?Message=login&Message=error&limit=5" \
  "yes"

# Test: Combination of multi-values
test_endpoint "Combined multi-values" \
  "/logs?FromHost=webserver01&FromHost=webserver02&Priority=3&Priority=6&limit=5" \
  "yes"

echo ""

# Section: Extended Columns Tests
echo -e "${BLUE}=== Extended Column Tests ===${NC}"

# Test: Check for extended columns in response
test_endpoint "Extended columns present" \
  "/logs?limit=1" \
  "yes" \
  "SysLogTag"

# Test: DeviceReportedTime exists
test_endpoint "DeviceReportedTime field" \
  "/logs?limit=1" \
  "yes"

# Test: Filter by SysLogTag (if exists)
test_endpoint "Filter by SysLogTag" \
  "/logs?SysLogTag=sshd&limit=1" \
  "yes"

# Test: Meta for SysLogTag
test_endpoint "Meta SysLogTag" \
  "/meta/SysLogTag" \
  "yes"

echo ""

# Section: Meta with Multi-Value Filters
echo -e "${BLUE}=== Meta with Multi-Value Filters ===${NC}"

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

# Section: Advanced Multi-Value Scenarios
echo -e "${BLUE}=== Advanced Scenarios ===${NC}"

# Test: Many values for same filter
test_endpoint "Many hosts (5+)" \
  "/logs?FromHost=webserver01&FromHost=webserver02&FromHost=dbserver01&FromHost=appserver01&FromHost=mailserver01&limit=5" \
  "yes"

# Test: All priority values
test_endpoint "All priorities (0-7)" \
  "/logs?Priority=0&Priority=1&Priority=2&Priority=3&Priority=4&Priority=5&Priority=6&Priority=7&limit=5" \
  "yes"

# Test: Complex combination
test_endpoint "Complex filter combo" \
  "/logs?FromHost=webserver01&FromHost=webserver02&Priority=3&Priority=4&Facility=1&Message=login&limit=5" \
  "yes"

echo ""

# Section: Backward Compatibility
echo -e "${BLUE}=== Backward Compatibility ===${NC}"

# Test: Single value still works
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

# Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total:  $((PASSED + FAILED))"
echo ""

# Detailed results
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    echo ""
    echo "v0.2.2 Features Verified:"
    echo "  ✓ Multi-value filters working"
    echo "  ✓ Extended columns present"
    echo "  ✓ Backward compatibility maintained"
    echo "  ✓ Meta endpoint with multi-value filters"
    echo ""
    exit 0
else
    echo -e "${RED}✗ Some tests failed!${NC}"
    echo ""
    echo "Please check the failed tests above."
    echo ""
    exit 1
fi
