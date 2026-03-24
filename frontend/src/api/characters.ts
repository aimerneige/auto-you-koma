import client from './client';
import { Character, CreateCharacterRequest } from '../types/character';

export const characterApi = {
  list: () => client.get<Character[]>('/characters'),

  get: (id: string) => client.get<Character>(`/characters/${id}`),

  create: (data: CreateCharacterRequest) => client.post<Character>('/characters', data),

  update: (id: string, data: CreateCharacterRequest) => client.put<Character>(`/characters/${id}`, data),

  delete: (id: string) => client.delete(`/characters/${id}`),

  generateReference: (id: string) => client.post<{ image_url: string; message: string }>(`/characters/${id}/generate-reference`, {}),

  confirmReference: (id: string, imageUrl: string) => client.post(`/characters/${id}/confirm-reference`, { image_url: imageUrl }),
};