import { useEffect, useState } from 'react';
import { adminService } from '../../services/adminService';

const hospitalTiers = [
  { value: 'national', label: 'National' },
  { value: 'regional', label: 'Regional' },
  { value: 'district', label: 'District' },
  { value: 'private', label: 'Private' },
];

const initialHospitalForm = {
  name: '',
  short_code: '',
  tier: 'private',
  region: 'Central',
  country: 'Uganda',
  address: '',
  phone: '',
  email: '',
  license_no: '',
};

const initialUserForm = {
  full_name: '',
  email: '',
  phone: '',
  password: '',
  role_name: 'doctor',
  is_active: true,
};

const hospitalToEditForm = (hospital) => ({
  name: hospital?.name || '',
  tier: hospital?.tier || 'private',
  region: hospital?.region || '',
  address: hospital?.address || '',
  phone: hospital?.phone || '',
  email: hospital?.email || '',
});

export default function AdminPage() {
  const [hospitals, setHospitals] = useState([]);
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [userRoles, setUserRoles] = useState({});
  const [selectedHospitalId, setSelectedHospitalId] = useState('');
  const [hospitalForm, setHospitalForm] = useState(initialHospitalForm);
  const [hospitalEditForm, setHospitalEditForm] = useState(hospitalToEditForm());
  const [userForm, setUserForm] = useState(initialUserForm);
  const [assignRole, setAssignRole] = useState({});
  const [passwordInputs, setPasswordInputs] = useState({});
  const [userSearch, setUserSearch] = useState('');
  const [loading, setLoading] = useState(true);
  const [savingHospital, setSavingHospital] = useState(false);
  const [savingHospitalEdit, setSavingHospitalEdit] = useState(false);
  const [savingUser, setSavingUser] = useState(false);
  const [assigningRoleFor, setAssigningRoleFor] = useState('');
  const [revokingRoleFor, setRevokingRoleFor] = useState('');
  const [resettingPasswordFor, setResettingPasswordFor] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    loadHospitals();
  }, []);

  useEffect(() => {
    if (!selectedHospitalId) return;
    loadHospitalAccessData(selectedHospitalId);
  }, [selectedHospitalId]);

  useEffect(() => {
    const selected = hospitals.find((hospital) => hospital.id === selectedHospitalId);
    setHospitalEditForm(hospitalToEditForm(selected));
  }, [hospitals, selectedHospitalId]);

  const loadHospitals = async () => {
    setLoading(true);
    setError('');
    try {
      const hospitalData = await adminService.listHospitals();
      setHospitals(hospitalData);
      if (hospitalData.length > 0) {
        setSelectedHospitalId((current) => current || hospitalData[0].id);
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load hospitals');
    } finally {
      setLoading(false);
    }
  };

  const loadHospitalAccessData = async (hospitalId) => {
    setError('');
    try {
      if (hospitalId === 'all') {
        const userGroups = await Promise.all(
          hospitals.map(async (hospital) => {
            const hospitalUsers = await adminService.listUsers(hospital.id);
            return hospitalUsers.map((user) => ({
              ...user,
              hospital_name: hospital.name,
              hospital_short_code: hospital.short_code,
            }));
          })
        );
        const usersData = userGroups.flat();
        setUsers(usersData);
        setRoles([]);

        const roleEntries = await Promise.all(
          usersData.map(async (user) => {
            try {
              const assigned = await adminService.listUserRoles(user.id);
              return [user.id, assigned];
            } catch {
              return [user.id, []];
            }
          })
        );

        setUserRoles(Object.fromEntries(roleEntries));
        setAssignRole({});
        return;
      }

      const [usersData, rolesData] = await Promise.all([
        adminService.listUsers(hospitalId),
        adminService.listRoles(hospitalId),
      ]);

      const selected = hospitals.find((hospital) => hospital.id === hospitalId);
      setUsers(usersData.map((user) => ({
        ...user,
        hospital_name: selected?.name || '',
        hospital_short_code: selected?.short_code || '',
      })));
      setRoles(rolesData);

      const roleEntries = await Promise.all(
        usersData.map(async (user) => {
          try {
            const assigned = await adminService.listUserRoles(user.id);
            return [user.id, assigned];
          } catch {
            return [user.id, []];
          }
        })
      );

      const roleMap = Object.fromEntries(roleEntries);
      setUserRoles(roleMap);

      const firstRole = rolesData[0]?.name || '';
      const nextAssign = {};
      usersData.forEach((user) => {
        nextAssign[user.id] = firstRole;
      });
      setAssignRole(nextAssign);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load hospital access data');
    }
  };

  const handleHospitalFormChange = (e) => {
    const { name, value } = e.target;
    setHospitalForm((current) => ({ ...current, [name]: value }));
  };

  const handleHospitalEditFormChange = (e) => {
    const { name, value } = e.target;
    setHospitalEditForm((current) => ({ ...current, [name]: value }));
  };

  const handleUserFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setUserForm((current) => ({
      ...current,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleCreateHospital = async (e) => {
    e.preventDefault();
    setSavingHospital(true);
    setError('');
    setSuccess('');
    try {
      const created = await adminService.createHospital(hospitalForm);
      setSuccess(`Hospital created: ${created.name}`);
      setHospitalForm(initialHospitalForm);
      await loadHospitals();
      setSelectedHospitalId(created.id);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create hospital');
    } finally {
      setSavingHospital(false);
    }
  };

  const handleUpdateHospital = async (e) => {
    e.preventDefault();
    if (!selectedHospitalId || selectedHospitalId === 'all') return;

    setSavingHospitalEdit(true);
    setError('');
    setSuccess('');
    try {
      const updated = await adminService.updateHospital(selectedHospitalId, hospitalEditForm);
      setSuccess(`Hospital updated: ${updated.name}`);
      setHospitals((current) => current.map((hospital) => (
        hospital.id === updated.id ? updated : hospital
      )));
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update hospital');
    } finally {
      setSavingHospitalEdit(false);
    }
  };

  const handleCreateUser = async (e) => {
    e.preventDefault();
    if (!selectedHospitalId) return;

    setSavingUser(true);
    setError('');
    setSuccess('');
    try {
      await adminService.createUser({
        ...userForm,
        hospital_id: selectedHospitalId,
      });
      setSuccess(`User created: ${userForm.email}`);
      setUserForm(initialUserForm);
      await loadHospitalAccessData(selectedHospitalId);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create user');
    } finally {
      setSavingUser(false);
    }
  };

  const handleAssignRole = async (userId) => {
    const roleName = assignRole[userId];
    if (!roleName || !selectedHospitalId) return;

    setAssigningRoleFor(userId);
    setError('');
    setSuccess('');
    try {
      await adminService.assignUserRole(userId, {
        hospital_id: selectedHospitalId,
        role_name: roleName,
      });
      setSuccess(`Granted ${roleName} access`);
      await loadHospitalAccessData(selectedHospitalId);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to assign role');
    } finally {
      setAssigningRoleFor('');
    }
  };

  const handleRevokeRole = async (userId, role) => {
    if (!role?.role_id) return;

    const revokeKey = `${userId}:${role.role_id}`;
    setRevokingRoleFor(revokeKey);
    setError('');
    setSuccess('');
    try {
      await adminService.revokeUserRole(userId, role.role_id);
      setSuccess(`Revoked ${role.role_name} access`);
      await loadHospitalAccessData(selectedHospitalId);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to revoke role');
    } finally {
      setRevokingRoleFor('');
    }
  };

  const handleResetPassword = async (userId, email) => {
    const password = passwordInputs[userId] || '';
    if (password.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }

    setResettingPasswordFor(userId);
    setError('');
    setSuccess('');
    try {
      await adminService.resetUserPassword(userId, password);
      setSuccess(`Password reset for ${email}`);
      setPasswordInputs((current) => ({ ...current, [userId]: '' }));
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to reset password');
    } finally {
      setResettingPasswordFor('');
    }
  };

  const selectedHospital = hospitals.find((hospital) => hospital.id === selectedHospitalId);
  const isAllHospitals = selectedHospitalId === 'all';
  const normalizedSearch = userSearch.trim().toLowerCase();
  const filteredUsers = normalizedSearch
    ? users.filter((user) => [
      user.full_name,
      user.email,
      user.hospital_name,
      user.hospital_short_code,
      (userRoles[user.id] || []).map((role) => role.role_name).join(' '),
    ].some((value) => String(value || '').toLowerCase().includes(normalizedSearch)))
    : users;

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <h1 className="text-3xl font-serif font-bold text-gray-900 tracking-tight">Platform Control</h1>
          <p className="mt-2 text-sm text-gray-600">
            Internal DMS console for hospital onboarding, subscriptions, and cross-hospital access control.
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-6 py-8 space-y-6">
        {error && (
          <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        {success && (
          <div className="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
            {success}
          </div>
        )}

        <div className="grid grid-cols-1 gap-6 xl:grid-cols-2">
          <section className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="text-lg font-semibold text-gray-900">Create Hospital</h2>
            <p className="mt-1 text-sm text-gray-500">New hospitals get the default admin and clinical roles automatically.</p>

            <form className="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2" onSubmit={handleCreateHospital}>
              <input name="name" value={hospitalForm.name} onChange={handleHospitalFormChange} placeholder="Hospital name" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
              <input name="short_code" value={hospitalForm.short_code} onChange={handleHospitalFormChange} placeholder="Short code" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
              <select name="tier" value={hospitalForm.tier} onChange={handleHospitalFormChange} className="rounded-lg border border-gray-300 px-4 py-2.5">
                {hospitalTiers.map((tier) => (
                  <option key={tier.value} value={tier.value}>{tier.label}</option>
                ))}
              </select>
              <input name="region" value={hospitalForm.region} onChange={handleHospitalFormChange} placeholder="Region" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
              <input name="country" value={hospitalForm.country} onChange={handleHospitalFormChange} placeholder="Country" className="rounded-lg border border-gray-300 px-4 py-2.5" />
              <input name="license_no" value={hospitalForm.license_no} onChange={handleHospitalFormChange} placeholder="License number" className="rounded-lg border border-gray-300 px-4 py-2.5" />
              <input name="phone" value={hospitalForm.phone} onChange={handleHospitalFormChange} placeholder="Phone" className="rounded-lg border border-gray-300 px-4 py-2.5" />
              <input name="email" value={hospitalForm.email} onChange={handleHospitalFormChange} placeholder="Email" className="rounded-lg border border-gray-300 px-4 py-2.5" />
              <textarea name="address" value={hospitalForm.address} onChange={handleHospitalFormChange} placeholder="Address" className="rounded-lg border border-gray-300 px-4 py-2.5 md:col-span-2" rows="3" />
              <div className="md:col-span-2">
                <button type="submit" disabled={savingHospital} className="rounded-lg bg-sky-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-sky-700 disabled:opacity-50">
                  {savingHospital ? 'Creating...' : 'Create Hospital'}
                </button>
              </div>
            </form>
          </section>

          <section className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="text-lg font-semibold text-gray-900">Hospital Access</h2>
            <p className="mt-1 text-sm text-gray-500">Choose a hospital, create a user, then assign or elevate roles.</p>

            <div className="mt-6">
              <label className="mb-2 block text-sm font-medium text-gray-700">Selected Hospital</label>
              <select
                value={selectedHospitalId}
                onChange={(e) => setSelectedHospitalId(e.target.value)}
                className="w-full rounded-lg border border-gray-300 px-4 py-2.5"
                disabled={loading || hospitals.length === 0}
              >
                <option value="all">All hospitals</option>
                {hospitals.map((hospital) => (
                  <option key={hospital.id} value={hospital.id}>
                    {hospital.name} ({hospital.short_code})
                  </option>
                ))}
              </select>
            </div>

            {selectedHospital && (
              <div className="mt-4 rounded-lg bg-gray-50 px-4 py-3 text-sm text-gray-600">
                <span className="font-medium text-gray-900">{selectedHospital.name}</span>
                {' '}• {selectedHospital.region} • {selectedHospital.subscription_plan}
              </div>
            )}

            {isAllHospitals && (
              <div className="mt-4 rounded-lg bg-amber-50 px-4 py-3 text-sm text-amber-800">
                All-hospital mode is for searching and reviewing users. Select one hospital to create users or grant roles.
              </div>
            )}

            {!isAllHospitals && selectedHospital && (
              <form className="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2" onSubmit={handleUpdateHospital}>
                <input name="name" value={hospitalEditForm.name} onChange={handleHospitalEditFormChange} placeholder="Hospital name" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
                <select name="tier" value={hospitalEditForm.tier} onChange={handleHospitalEditFormChange} className="rounded-lg border border-gray-300 px-4 py-2.5">
                  {hospitalTiers.map((tier) => (
                    <option key={tier.value} value={tier.value}>{tier.label}</option>
                  ))}
                </select>
                <input name="region" value={hospitalEditForm.region} onChange={handleHospitalEditFormChange} placeholder="Region" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
                <input name="phone" value={hospitalEditForm.phone} onChange={handleHospitalEditFormChange} placeholder="Phone" className="rounded-lg border border-gray-300 px-4 py-2.5" />
                <input name="email" value={hospitalEditForm.email} onChange={handleHospitalEditFormChange} placeholder="Email" className="rounded-lg border border-gray-300 px-4 py-2.5" />
                <textarea name="address" value={hospitalEditForm.address} onChange={handleHospitalEditFormChange} placeholder="Address" className="rounded-lg border border-gray-300 px-4 py-2.5 md:col-span-2" rows="2" />
                <div className="md:col-span-2">
                  <button type="submit" disabled={savingHospitalEdit} className="rounded-lg border border-gray-300 bg-white px-5 py-2.5 text-sm font-medium text-gray-800 hover:bg-gray-50 disabled:opacity-50">
                    {savingHospitalEdit ? 'Saving...' : 'Save Hospital Details'}
                  </button>
                </div>
              </form>
            )}

            <form className="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2" onSubmit={handleCreateUser}>
              <input name="full_name" value={userForm.full_name} onChange={handleUserFormChange} placeholder="Full name" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
              <input name="email" type="email" value={userForm.email} onChange={handleUserFormChange} placeholder="Email" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
              <input name="phone" value={userForm.phone} onChange={handleUserFormChange} placeholder="Phone" className="rounded-lg border border-gray-300 px-4 py-2.5" />
              <input name="password" type="password" value={userForm.password} onChange={handleUserFormChange} placeholder="Temporary password" className="rounded-lg border border-gray-300 px-4 py-2.5" required />
              <select name="role_name" value={userForm.role_name} onChange={handleUserFormChange} className="rounded-lg border border-gray-300 px-4 py-2.5">
                {roles.map((role) => (
                  <option key={role.id} value={role.name}>{role.name}</option>
                ))}
              </select>
              <label className="flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2.5 text-sm text-gray-700">
                <input type="checkbox" name="is_active" checked={userForm.is_active} onChange={handleUserFormChange} />
                Active account
              </label>
              <div className="md:col-span-2">
                <button type="submit" disabled={savingUser || !selectedHospitalId || isAllHospitals} className="rounded-lg bg-gray-900 px-5 py-2.5 text-sm font-medium text-white hover:bg-gray-800 disabled:opacity-50">
                  {savingUser ? 'Creating...' : 'Create User'}
                </button>
              </div>
            </form>
          </section>
        </div>

        <section className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-lg font-semibold text-gray-900">Users and Roles</h2>
              <p className="mt-1 text-sm text-gray-500">Grant, revoke, reset passwords, or filter users across hospitals.</p>
            </div>
            <div className="w-72">
              <input
                value={userSearch}
                onChange={(e) => setUserSearch(e.target.value)}
                placeholder="Search users, roles, hospitals"
                className="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm"
              />
            </div>
          </div>

          <div className="mt-6 overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm">
              <thead>
                <tr className="text-left text-gray-500">
                  <th className="pb-3 pr-4 font-medium">User</th>
                  <th className="pb-3 pr-4 font-medium">Email</th>
                  <th className="pb-3 pr-4 font-medium">Hospital</th>
                  <th className="pb-3 pr-4 font-medium">Current Roles</th>
                  <th className="pb-3 pr-4 font-medium">Grant Role</th>
                  <th className="pb-3 pr-4 font-medium">Reset Password</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {filteredUsers.map((user) => (
                  <tr key={user.id}>
                    <td className="py-3 pr-4 font-medium text-gray-900">{user.full_name}</td>
                    <td className="py-3 pr-4 text-gray-600">{user.email}</td>
                    <td className="py-3 pr-4 text-gray-600">{user.hospital_short_code || user.hospital_name || '-'}</td>
                    <td className="py-3 pr-4 text-gray-600">
                      {(userRoles[user.id] || []).length > 0 ? (
                        <div className="flex flex-wrap gap-2">
                          {userRoles[user.id].map((role) => {
                            const revokeKey = `${user.id}:${role.role_id}`;
                            return (
                              <span key={role.id || role.role_id} className="inline-flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1 text-xs text-gray-700">
                                {role.role_name}
                                <button
                                  type="button"
                                  onClick={() => handleRevokeRole(user.id, role)}
                                  disabled={revokingRoleFor === revokeKey}
                                  className="font-semibold text-red-600 hover:text-red-700 disabled:opacity-50"
                                >
                                  {revokingRoleFor === revokeKey ? '...' : 'Revoke'}
                                </button>
                              </span>
                            );
                          })}
                        </div>
                      ) : 'No roles'}
                    </td>
                    <td className="py-3 pr-4">
                      <div className="flex gap-2">
                        <select
                          value={assignRole[user.id] || ''}
                          onChange={(e) => setAssignRole((current) => ({ ...current, [user.id]: e.target.value }))}
                          className="rounded-lg border border-gray-300 px-3 py-2"
                        >
                          <option value="">Select role</option>
                          {roles.map((role) => (
                            <option key={role.id} value={role.name}>{role.name}</option>
                          ))}
                        </select>
                        <button
                          type="button"
                          onClick={() => handleAssignRole(user.id)}
                          disabled={assigningRoleFor === user.id || !assignRole[user.id] || isAllHospitals}
                          className="rounded-lg border border-sky-200 bg-sky-50 px-3 py-2 font-medium text-sky-700 hover:bg-sky-100 disabled:opacity-50"
                        >
                          {assigningRoleFor === user.id ? 'Granting...' : 'Grant'}
                        </button>
                      </div>
                    </td>
                    <td className="py-3 pr-4">
                      <div className="flex gap-2">
                        <input
                          type="password"
                          value={passwordInputs[user.id] || ''}
                          onChange={(e) => setPasswordInputs((current) => ({ ...current, [user.id]: e.target.value }))}
                          placeholder="New password"
                          className="w-36 rounded-lg border border-gray-300 px-3 py-2"
                        />
                        <button
                          type="button"
                          onClick={() => handleResetPassword(user.id, user.email)}
                          disabled={resettingPasswordFor === user.id || !(passwordInputs[user.id] || '')}
                          className="rounded-lg border border-gray-300 bg-white px-3 py-2 font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                        >
                          {resettingPasswordFor === user.id ? 'Resetting...' : 'Reset'}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {!loading && filteredUsers.length === 0 && (
              <div className="py-8 text-center text-sm text-gray-500">
                No users found for this hospital yet.
              </div>
            )}
          </div>
        </section>
      </div>
    </div>
  );
}
