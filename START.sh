#!/bin/bash
# DMS Quick Start Script
# Run this after configuring backend/.env with Supabase credentials

set -e

echo "========================================="
echo "🚀 DMS - Dialysis Management System"
echo "========================================="
echo ""

# Check if backend .env exists
if [ ! -f "backend/.env" ]; then
  echo "❌ Error: backend/.env not found!"
  echo ""
  echo "📝 Please create backend/.env with your Supabase credentials:"
  echo "   cp backend/.env.example backend/.env"
  echo "   nano backend/.env"
  echo ""
  echo "See SUPABASE_SETUP.md for detailed instructions"
  exit 1
fi

echo "✅ Found backend/.env"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
  echo "❌ Error: Go not installed"
  echo "Install: https://go.dev/doc/install"
  exit 1
fi

# Check if Node is installed
if ! command -v node &> /dev/null; then
  echo "❌ Error: Node.js not installed"
  echo "Install: https://nodejs.org/"
  exit 1
fi

echo "🔧 Installing dependencies..."
echo ""

# Install backend dependencies
cd backend
echo "→ Installing Go modules..."
go mod download
cd ..

# Install frontend dependencies
cd frontend
echo "→ Installing npm packages..."
npm install
cd ..

echo ""
echo "========================================="
echo "✅ Ready to start!"
echo "========================================="
echo ""
echo "Opening 2 terminal windows:"
echo "  Terminal 1: Backend (Go server on :8080)"
echo "  Terminal 2: Frontend (Vite dev server on :5173)"
echo ""
echo "Press Ctrl+C in each terminal to stop"
echo ""
echo "🌐 Access URLs:"
echo "   Frontend: http://localhost:5173"
echo "   Backend:  http://localhost:8080"
echo "   Health:   http://localhost:8080/health"
echo ""
echo "========================================="
echo ""

# Start backend in new terminal
gnome-terminal --title="DMS Backend" -- bash -c "cd backend && echo '🟢 Starting Backend Server...' && go run cmd/api/main.go; exec bash" 2>/dev/null || \
xterm -title "DMS Backend" -e "cd backend && echo '🟢 Starting Backend Server...' && go run cmd/api/main.go; bash" 2>/dev/null || \
(cd backend && go run cmd/api/main.go &)

sleep 2

# Start frontend in new terminal
gnome-terminal --title="DMS Frontend" -- bash -c "cd frontend && echo '🟢 Starting Frontend Dev Server...' && npm run dev; exec bash" 2>/dev/null || \
xterm -title "DMS Frontend" -e "cd frontend && echo '🟢 Starting Frontend Dev Server...' && npm run dev; bash" 2>/dev/null || \
(cd frontend && npm run dev &)

sleep 3

echo "🎉 DMS is starting up..."
echo ""
echo "Waiting for servers to initialize..."
sleep 5

# Try to open browser
if command -v xdg-open &> /dev/null; then
  xdg-open http://localhost:5173 2>/dev/null &
elif command -v open &> /dev/null; then
  open http://localhost:5173 2>/dev/null &
else
  echo "👉 Open your browser manually: http://localhost:5173"
fi

echo ""
echo "✅ DMS is running!"
echo ""
echo "📚 Quick Start Guide:"
echo "   1. Login with your credentials"
echo "   2. Go to Settings → Toggle modules"
echo "   3. Import patient data (see SUPABASE_SETUP.md)"
echo "   4. Create test patient"
echo "   5. Test offline mode (DevTools → Network → Offline)"
echo ""
echo "🛑 To stop: Press Ctrl+C in backend and frontend terminals"
echo ""
