.PHONY: start fix-connection-refused fix-proxy check-ports

# Default command to start the frontend
start:
	@echo "🚀 Starting React frontend..."
	npm start

endf:
# Command to stop the frontend
endf:
	@echo "🛑 Stopping React frontend...
	@lsof -ti:3000 | xargs kill -9 2>/dev/null || echo "No React process found"
	@echo "🚀 React frontend has stopped."

# Fix connection refused errors
fix-connection-refused:
	@echo "🔧 Fixing Connection Refused Errors (Frontend)"
	@echo "============================================="
	@echo ""
	@echo "Checking if servers are running:"
	@lsof -ti:8080 >/dev/null && echo "✅ API server is running on port 8080" || echo "❌ API server is NOT running"
	@lsof -ti:3000 >/dev/null && echo "✅ Frontend server is running on port 3000" || echo "❌ Frontend server is NOT running"
	@echo ""
	@echo "🔍 Checking proxy configuration in package.json:"
	@if grep -q "\"proxy\":" package.json; then \
		echo "✅ Proxy configuration found:"; \
		grep -A 1 "\"proxy\":" package.json; \
	else \
		echo "❌ No proxy configuration found in package.json"; \
		echo ""; \
		echo "🔧 SOLUTION: Add this line to package.json (at the top level):"; \
		echo "   \"proxy\": \"http://localhost:8080\","; \
		echo ""; \
		echo "   This will automatically forward API requests to the backend"; \
	fi
	@echo ""
	@echo "🔍 Common issues:"
	@echo "1. API server not running - Start it with: cd .. && make server"
	@echo "2. Missing proxy config - Add \"proxy\": \"http://localhost:8080\" to package.json"
	@echo "3. React using wrong URL - Make sure API calls use relative paths like /api/..."
	@echo "4. Browser navigating to API - Always use http://localhost:3000 in browser"
	@echo ""
	@echo "🔧 Quick fix steps:"
	@echo "1. Add proxy configuration to package.json:"
	@echo "   \"proxy\": \"http://localhost:8080\","
	@echo "2. Kill and restart the frontend server:"
	@echo "   make fix-proxy"
	@echo ""
	@echo "💡 TIP: Always navigate using React Router links, not by typing URLs"

# Fix proxy configuration and restart
fix-proxy:
	@echo "🔧 Fixing proxy configuration..."
	@if ! grep -q "\"proxy\":" package.json; then \
		echo "Adding proxy to package.json..."; \
		sed -i '' 's/\("name": "frontend"\)/\1,\n  "proxy": "http:\/\/localhost:8080"/' package.json; \
		echo "✅ Added proxy configuration"; \
	else \
		echo "✅ Proxy already configured"; \
	fi
	@echo "🔄 Stopping any existing React processes..."
	@lsof -ti:3000 | xargs kill -9 2>/dev/null || echo "No React process found"
	@echo "🔄 Clearing React cache..."
	@rm -rf node_modules/.cache 2>/dev/null || echo "No cache to clear"
	@echo "🚀 Restarting React frontend..."
	@npm start

# Check ports
check-ports:
	@echo "🔍 Checking port status:"
	@echo "Frontend (3000): $(shell lsof -ti:3000 >/dev/null && echo '✅ Running' || echo '❌ Not running')"
	@echo "API (8080): $(shell lsof -ti:8080 >/dev/null && echo '✅ Running' || echo '❌ Not running')"
	@echo "MongoDB (27017): $(shell lsof -ti:27017 >/dev/null && echo '✅ Running' || echo '❌ Not running')"
	@echo "Weaviate (8081): $(shell lsof -ti:8081 >/dev/null && echo '✅ Running' || echo '❌ Not running')"
