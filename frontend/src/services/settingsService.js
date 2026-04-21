import api from './api';

export const settingsService = {
  // Update enabled modules (admin only)
  async updateModules(modules) {
    const response = await api.put('/subscription/modules', modules);
    return response.data;
  },

  // Get hospital settings
  async getHospitalSettings() {
    const response = await api.get('/hospitals/settings');
    return response.data;
  },

  // Update hospital setting
  async updateSetting(key, value) {
    const response = await api.put('/hospitals/settings', { key, value });
    return response.data;
  },
};
