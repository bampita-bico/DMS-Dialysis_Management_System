# 👋 Welcome to DMS - Dialysis Management System

**Uganda-Localized • Offline-First • Production-Ready**

---

## 🎯 What You Need to Do RIGHT NOW

### 1️⃣ Configure Supabase (5 minutes)
📖 **Read:** `SUPABASE_SETUP.md`

**Quick version:**
1. Get Supabase connection string from dashboard
2. Create `backend/.env` with credentials
3. Seed database with Uganda data

### 2️⃣ Start the System
```bash
./START.sh
```

### 3️⃣ Open Browser
👉 **http://localhost:5173**

---

## 📚 Documentation Files

| File | Purpose | When to Read |
|------|---------|--------------|
| **SUPABASE_SETUP.md** | Connect to Supabase, seed database | **READ FIRST** |
| **QUICK_START.md** | 5-minute start guide | After Supabase setup |
| **UGANDA_COMPLETION_SUMMARY.md** | What's included, features | Reference |
| **START.sh** | Automatic startup script | Just run it! |

---

## ✅ System Features

### 🇺🇬 Uganda-Specific
- ✅ 75 locally-available medications (Heparin, Eprex, Venofer)
- ✅ Lung ultrasound for pulmonary edema detection
- ✅ Ugandan insurance schemes (NHIS, AAR, Britam)
- ✅ Local pricing in UGX/KES

### 📡 Offline-First
- ✅ Works without internet
- ✅ Auto-sync when online
- ✅ Conflict detection
- ✅ IndexedDB storage

### 🏥 Clinical Features
- ✅ Patient management
- ✅ Session scheduling
- ✅ Lab orders & results (54 tests)
- ✅ Billing & invoices
- ✅ Staff & shift management
- ✅ Module toggles (enable/disable features)

### 📊 Real Data Ready
- ✅ Import script for MS Access data
- ✅ 33 real patients ready to import
- ✅ 330 real sessions with vitals

---

## 🚀 Quick Start (One Command)

```bash
cd /home/bampita/Projects/My-apps/DMS-Dialysis_Management_System
./START.sh
```

**What it does:**
1. Checks dependencies
2. Installs packages
3. Starts backend (port 8080)
4. Starts frontend (port 5173)
5. Opens browser automatically

---

## 🔗 Access URLs (After Starting)

| Service | URL | Description |
|---------|-----|-------------|
| **Frontend** | http://localhost:5173 | Main application |
| **Backend API** | http://localhost:8080/api/v1 | REST API |
| **Health Check** | http://localhost:8080/health | Server status |
| **Supabase** | https://supabase.com/dashboard | Database dashboard |

---

## 🗄️ Update Supabase Database

### Option 1: Via Supabase Dashboard (Easiest)
1. Go to Supabase → SQL Editor
2. Get hospital UUID: `SELECT id FROM hospitals LIMIT 1;`
3. Open each seed file in `backend/seeds/`
4. Replace `:hospital_id` with your UUID
5. Copy SQL → Paste in SQL Editor → Run

**Seed files (in order):**
```
001_lab_tests.sql              → 54 dialysis lab tests
002_lab_panels.sql             → 12 pre-built panels
003_lab_reference_ranges.sql   → Normal ranges
004_medications_uganda.sql     → 75 Uganda medications 🇺🇬
005_drug_interactions.sql      → Critical interactions
006_consumables.sql            → 35 dialysis consumables
007_insurance_schemes.sql      → 16 insurance schemes
008_price_lists.sql            → 33 service prices
```

### Option 2: Via Command Line
```bash
cd backend/scripts

# Set credentials
export DB_HOST=aws-0-[region].pooler.supabase.com
export DB_PORT=6543
export DB_NAME=postgres
export DB_USER=postgres.[project-ref]
export DB_PASSWORD=your_password

# Run seeding
./seed_all.sh <hospital-uuid>
```

### Import Real Patient Data
```bash
cd backend/scripts
export DB_PASSWORD=your_password
python3 import_access_data.py <hospital-uuid>
```

---

## 🎓 First-Time User Guide

### 1. Create Admin User
Go to Supabase → SQL Editor:
```sql
INSERT INTO users (id, hospital_id, email, password_hash, full_name, role, is_active)
VALUES (
  gen_random_uuid(),
  (SELECT id FROM hospitals LIMIT 1),
  'admin@hospital.com',
  crypt('admin123', gen_salt('bf')),
  'Admin User',
  'admin',
  true
);
```

### 2. Login
- Email: `admin@hospital.com`
- Password: `admin123`

### 3. Configure System
- Go to **Settings** → Toggle modules you need
- Module options:
  - Lab Management
  - Full Pharmacy
  - HR Management
  - Inventory Tracking
  - Advanced Billing
  - Imaging Integration
  - CHW Program
  - Outcomes Reporting
  - Offline Sync (keep ON)

### 4. Add Data
- **Patients:** Add first patient or import MS Access data
- **Sessions:** Schedule dialysis session
- **Lab:** Create lab order (lung ultrasound available!)
- **Billing:** Create invoice
- **Staff:** Add staff and assign shifts

### 5. Test Offline Mode
1. Open DevTools (F12)
2. Network tab → Select "Offline"
3. Navigate pages → Still works!
4. Create patient offline → Saves locally
5. Go online → Watch sync indicator (bottom-right)

---

## 🛠️ System Requirements

### Backend
- Go 1.21+ (check: `go version`)
- PostgreSQL 15+ (Supabase)
- Port 8080 available

### Frontend
- Node.js 18+ (check: `node --version`)
- npm 9+ (check: `npm --version`)
- Port 5173 available
- Modern browser (Chrome, Firefox, Edge)

---

## 🐛 Troubleshooting

### Backend won't start
```bash
# Check .env exists
ls backend/.env

# Check Go installed
go version

# Check port available
lsof -i :8080
```

### Frontend won't start
```bash
# Check Node installed
node --version

# Install dependencies
cd frontend && npm install

# Check port available
lsof -i :5173
```

### Can't connect to Supabase
```bash
# Test connection
psql "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require" -c "SELECT 1;"
```

### No data showing
- Check if database seeded (see seed counts in Supabase)
- Check browser console (F12) for errors
- Check backend logs for database errors

---

## 📞 Support & Documentation

| Question | Documentation |
|----------|---------------|
| How to connect Supabase? | `SUPABASE_SETUP.md` |
| How to start system? | `QUICK_START.md` or `./START.sh` |
| What's included? | `UGANDA_COMPLETION_SUMMARY.md` |
| What's new? | `PHASE_5_STATUS_AND_NEXT_STEPS.md` |
| Technical details? | `backend/README.md`, `frontend/README.md` |

---

## 📊 What's Included

### Database
- ✅ 68 tables (fully normalized)
- ✅ Row-level security (RLS)
- ✅ Multi-tenancy support
- ✅ Sync tracking tables
- ✅ Audit logging

### Backend (Go)
- ✅ REST API (Gin framework)
- ✅ JWT authentication
- ✅ Background sync worker
- ✅ Conflict detection
- ✅ 50+ API endpoints

### Frontend (React)
- ✅ React 19 + Vite
- ✅ Tailwind CSS
- ✅ Offline-first (IndexedDB)
- ✅ Module toggles
- ✅ Sync indicator
- ✅ 7 pages fully integrated

### Seeds & Data
- ✅ 8 SQL seed files
- ✅ 225+ reference items
- ✅ Import script for MS Access
- ✅ 33 real patients ready

---

## 🎉 You're Ready!

### Next Steps:
1. ✅ Read `SUPABASE_SETUP.md`
2. ✅ Run `./START.sh`
3. ✅ Open http://localhost:5173
4. ✅ Login and explore

**Welcome to your Uganda-localized Dialysis Management System! 🇺🇬**

---

### 🔥 Pro Tips

💡 **Backup your .env file** - Contains important credentials  
💡 **Test offline mode first** - Main feature of system  
💡 **Import real data early** - Better for testing  
💡 **Enable only modules you need** - Simpler UI  
💡 **Check sync indicator** - Bottom-right corner  

---

**Built with ❤️ for Ugandan dialysis centers**
