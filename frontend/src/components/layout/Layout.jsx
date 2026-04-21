import { Link, useLocation, useNavigate } from 'react-router-dom';
import { authService } from '../../services/auth';
import EmojiIcon from '../ui/EmojiIcon';

export default function Layout({ children }) {
  const location = useLocation();
  const navigate = useNavigate();
  const user = authService.getCurrentUser();

  const navigation = [
    { name: 'Dashboard', path: '/dashboard', icon: '📊', dark: false },
    { name: 'Patients', path: '/patients', icon: '👥', dark: false },
    { name: 'Sessions', path: '/sessions', icon: '🩺', dark: false },
    { name: 'Laboratory', path: '/lab', icon: '🧪', dark: false },
    { name: 'Billing', path: '/billing', icon: '💰', dark: true },
    { name: 'Staff', path: '/staff', icon: '👨‍⚕️', dark: true },
  ];

  const handleLogout = () => {
    authService.logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Top Navigation Bar */}
      <nav className="bg-gray-50 shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            {/* Logo & Nav Links */}
            <div className="flex">
              <div className="flex-shrink-0 flex items-center">
                <span className="text-2xl font-bold text-sky-600">DMS</span>
                <span className="ml-2 text-sm text-gray-600">Kiruddu Hospital</span>
              </div>
              <div className="hidden sm:ml-6 sm:flex sm:space-x-1">
                {navigation.map((item) => {
                  const isActive = location.pathname.startsWith(item.path);
                  return (
                    <Link
                      key={item.path}
                      to={item.path}
                      className={`inline-flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                        isActive
                          ? 'bg-sky-50 text-sky-700 border-b-2 border-sky-600'
                          : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                      }`}
                    >
                      <span className="mr-2">
                        <EmojiIcon dark={item.dark} size="sm">{item.icon}</EmojiIcon>
                      </span>
                      {item.name}
                    </Link>
                  );
                })}
              </div>
            </div>

            {/* Right side - User menu */}
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-700">
                {user?.full_name || 'User'}
              </span>
              <Link
                to="/settings"
                className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md transition"
                title="Settings"
              >
                ⚙️
              </Link>
              <button
                onClick={handleLogout}
                className="px-3 py-1.5 text-sm font-medium text-red-600 hover:text-red-700 hover:bg-red-50 rounded-md transition"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main>
        {children}
      </main>
    </div>
  );
}
