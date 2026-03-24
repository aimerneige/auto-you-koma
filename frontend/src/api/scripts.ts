import { apiClient } from './client';

export const scriptApi = {
  list: (projectId: string) => apiClient.get('/scripts', { params: { project_id: projectId } }),
  get: (id: string) => apiClient.get(`/scripts/${id}`),
  create: (data: any) => apiClient.post('/scripts', data),
  update: (id: string, data: any) => apiClient.put(`/scripts/${id}`, data),
  generate: (prompt: string) => apiClient.post('/scripts/generate', { prompt }),
  parse: (id: string) => apiClient.post(`/scripts/${id}/parse`)
};
