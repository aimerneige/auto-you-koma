import { apiClient } from './client';

export const seriesApi = {
  create: (data: { title: string; description: string }) => apiClient.post('/series', data),
  list: () => apiClient.get('/series'),
  synthesizeMemory: (seriesId: string, scriptId: string, charId: string) =>
    apiClient.post(`/series/${seriesId}/continuity/${scriptId}/${charId}`)
};
