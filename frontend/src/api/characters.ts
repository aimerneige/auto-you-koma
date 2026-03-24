import { apiClient } from './client';

export interface CharacterFilter {
  category?: string;
  tags?: string[];
  q?: string;
}

export const characterApi = {
  list: (params?: CharacterFilter) => apiClient.get('/characters', { params }),
  get: (id: string) => apiClient.get(`/characters/${id}`),
  create: (data: any) => apiClient.post('/characters', data),
  addVariant: (id: string, data: any) => apiClient.post(`/characters/${id}/variants`, data),
  uploadImage: (id: string, formData: FormData) => apiClient.post(`/characters/${id}/images`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
};
