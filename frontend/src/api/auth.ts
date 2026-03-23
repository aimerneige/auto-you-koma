import { apiClient } from './client';

export const authApi = {
  register: (data: any) => apiClient.post('/auth/register', data),
  login: (data: any) => apiClient.post('/auth/login', data),
  setup2FA: () => apiClient.post('/auth/2fa/setup'),
  verify2FA: (data: { passcode: string }) => apiClient.post('/auth/2fa/verify', data),
  getMe: () => apiClient.get('/auth/me'),
};
