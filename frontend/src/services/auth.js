import api from './api';

export const authService = {
  async login(email, password) {
    const response = await api.post('/auth/login', { email, password });
    const { token, user } = response.data;
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    return { token, user };
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
  },

  getCurrentUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },

  hasRole(roleName) {
    const user = this.getCurrentUser();
    const roles = user?.role_names || [];
    return roles.some(role => String(role).toLowerCase() === String(roleName).toLowerCase());
  },

  isHospitalAdmin() {
    const user = this.getCurrentUser();
    return Boolean(user?.is_hospital_admin || user?.is_admin || this.hasRole('admin') || this.hasRole('super_admin'));
  },

  isPlatformAdmin() {
    const user = this.getCurrentUser();
    return Boolean(user?.is_platform_admin || this.hasRole('super_admin'));
  },

  isAuthenticated() {
    return !!localStorage.getItem('token');
  },
};
