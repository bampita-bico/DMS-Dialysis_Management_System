import api from './api';

export const adminService = {
  async listHospitals() {
    const response = await api.get('/hospitals');
    return response.data;
  },

  async createHospital(payload) {
    const response = await api.post('/hospitals', payload);
    return response.data;
  },

  async updateHospital(hospitalId, payload) {
    const response = await api.patch(`/hospitals/${hospitalId}`, payload);
    return response.data;
  },

  async listUsers(hospitalId) {
    const response = await api.get('/users', {
      params: hospitalId ? { hospital_id: hospitalId } : {},
    });
    return response.data;
  },

  async createUser(payload) {
    const response = await api.post('/users', payload);
    return response.data;
  },

  async listRoles(hospitalId) {
    const response = await api.get('/roles', {
      params: hospitalId ? { hospital_id: hospitalId } : {},
    });
    return response.data;
  },

  async listUserRoles(userId) {
    const response = await api.get(`/users/${userId}/roles`);
    return response.data;
  },

  async assignUserRole(userId, payload) {
    const response = await api.post(`/users/${userId}/roles`, payload);
    return response.data;
  },

  async revokeUserRole(userId, roleId) {
    await api.delete(`/users/${userId}/roles/${roleId}`);
  },

  async resetUserPassword(userId, password) {
    await api.patch(`/users/${userId}/password`, { password });
  },
};
