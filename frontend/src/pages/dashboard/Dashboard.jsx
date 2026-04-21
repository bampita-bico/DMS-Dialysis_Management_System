import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { authService } from '../../services/auth';
import { useModules } from '../../contexts/ModuleContext';
import useOfflineData from '../../hooks/useOfflineData';
import api from '../../services/api';
import EmojiIcon from '../../components/ui/EmojiIcon';

export default function Dashboard() {
  const navigate = useNavigate();
  const { isModuleEnabled } = useModules();
  const [stats, setStats] = useState({
    todaySessions: 0,
    activeSessions: 0,
    pendingAlerts: 0,
    overdueInvoices: 0,
  });
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState(null);

  // Load real data using offline-first approach
  const today = new Date().toISOString().split('T')[0];
  const { data: todaySessions } = useOfflineData('dialysis_sessions', { date: today });
  const { data: activeSessions } = useOfflineData('dialysis_sessions', { status: 'in_progress' });
  const { data: alerts } = useOfflineData('lab_critical_alerts', { acknowledged: false });
  const { data: invoices } = useOfflineData('invoices', { invoice_status: 'overdue' });

  useEffect(() => {
    setUser(authService.getCurrentUser());
    loadDashboardStats();
  }, [todaySessions, activeSessions, alerts, invoices]);

  const loadDashboardStats = () => {
    setStats({
      todaySessions: todaySessions?.length || 0,
      activeSessions: activeSessions?.length || 0,
      pendingAlerts: alerts?.length || 0,
      overdueInvoices: invoices?.length || 0,
    });
    setLoading(false);
  };

  const handleLogout = () => {
    authService.logout();
    navigate('/login');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
              <p className="text-sm text-gray-500 mt-1">Welcome back, {user?.full_name || 'Doctor'}</p>
            </div>
            <div className="flex items-center space-x-4">
              <Link
                to="/settings"
                className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-md transition"
              >
                ⚙️ Settings
              </Link>
              <button
                onClick={handleLogout}
                className="px-4 py-2 text-sm font-medium text-red-600 hover:text-red-700 hover:bg-red-50 rounded-md transition"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Navigation */}
      <nav className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex space-x-8 overflow-x-auto py-4">
            <Link to="/patients" className="text-gray-700 hover:text-blue-600 whitespace-nowrap">
              👤 Patients
            </Link>
            <Link to="/sessions" className="text-gray-700 hover:text-blue-600 whitespace-nowrap">
              💉 Sessions
            </Link>
            {isModuleEnabled('lab_management') && (
              <Link to="/lab" className="text-gray-700 hover:text-blue-600 whitespace-nowrap">
                🧪 Lab
              </Link>
            )}
            {isModuleEnabled('advanced_billing') && (
              <Link to="/billing" className="text-gray-700 hover:text-blue-600 whitespace-nowrap">
                <EmojiIcon dark size="sm">💰</EmojiIcon> Billing
              </Link>
            )}
            {isModuleEnabled('hr_management') && (
              <Link to="/staff" className="text-gray-700 hover:text-blue-600 whitespace-nowrap">
                <EmojiIcon dark size="sm">👨‍⚕️</EmojiIcon> Staff
              </Link>
            )}
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatCard
            title="Today's Sessions"
            value={stats.todaySessions}
            icon="📅"
            color="blue"
            description="Scheduled for today"
            onClick={() => navigate('/sessions')}
          />
          <StatCard
            title="Active Now"
            value={stats.activeSessions}
            icon="💉"
            color="green"
            description="Currently in progress"
            onClick={() => navigate('/sessions')}
          />
          {isModuleEnabled('lab_management') && (
            <StatCard
              title="Critical Alerts"
              value={stats.pendingAlerts}
              icon="⚠️"
              color="red"
              description="Requiring attention"
              onClick={() => navigate('/lab')}
            />
          )}
          {isModuleEnabled('advanced_billing') && (
            <StatCard
              title="Overdue Invoices"
              value={stats.overdueInvoices}
              icon="💰"
              color="yellow"
              description="Payment pending"
              onClick={() => navigate('/billing')}
              darkEmoji={true}
            />
          )}
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <QuickAction
              icon="👤"
              label="Add Patient"
              onClick={() => navigate('/patients')}
            />
            <QuickAction
              icon="📋"
              label="Schedule Session"
              onClick={() => navigate('/sessions')}
            />
            {isModuleEnabled('lab_management') && (
              <QuickAction
                icon="🧪"
                label="Lab Results"
                onClick={() => navigate('/lab')}
              />
            )}
            {isModuleEnabled('advanced_billing') && (
              <QuickAction
                icon="📄"
                label="Create Invoice"
                onClick={() => navigate('/billing')}
              />
            )}
          </div>
        </div>

        {/* Recent Activity */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Recent Activity</h2>
          <div className="text-center text-gray-500 py-8">
            <p>Recent dialysis sessions and activities</p>
            <Link to="/sessions" className="text-blue-600 hover:underline mt-2 inline-block">
              View all sessions →
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
}

function StatCard({ title, value, icon, color, description, onClick, darkEmoji = false }) {
  const colors = {
    blue: 'from-blue-500 to-blue-600',
    green: 'from-green-500 to-green-600',
    red: 'from-red-500 to-red-600',
    yellow: 'from-yellow-500 to-yellow-600',
  };

  return (
    <button
      onClick={onClick}
      className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition text-left w-full"
    >
      <div className={`h-2 bg-gradient-to-r ${colors[color]}`}></div>
      <div className="p-6">
        <div className="flex items-center justify-between mb-2">
          <p className="text-gray-500 text-sm font-medium">{title}</p>
          <EmojiIcon dark={darkEmoji} size="lg">{icon}</EmojiIcon>
        </div>
        <p className="text-3xl font-bold text-gray-900 mb-1">{value}</p>
        <p className="text-xs text-gray-500">{description}</p>
      </div>
    </button>
  );
}

function QuickAction({ icon, label, onClick, darkEmoji = false }) {
  return (
    <button
      onClick={onClick}
      className="flex flex-col items-center justify-center p-4 bg-gray-100 hover:bg-blue-50 rounded-lg transition group"
    >
      <span className="text-3xl mb-2 group-hover:scale-110 transition">
        <EmojiIcon dark={darkEmoji} size="xl">{icon}</EmojiIcon>
      </span>
      <span className="text-sm font-medium text-gray-700 group-hover:text-blue-600">{label}</span>
    </button>
  );
}
