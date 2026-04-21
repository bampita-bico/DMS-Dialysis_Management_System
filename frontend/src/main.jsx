import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'
import db, { initializeMetadata } from './db/schema'
import { seedDemoData } from './utils/demoDataSeeder'

// Initialize IndexedDB and seed demo data
(async () => {
  try {
    await db.open();
    await initializeMetadata();
    console.log('💾 IndexedDB initialized successfully');

    // Seed demo data for investor presentations
    // (Only runs if DB is empty)
    await seedDemoData();
  } catch (error) {
    console.error('Failed to initialize IndexedDB:', error);
  }
})();

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
