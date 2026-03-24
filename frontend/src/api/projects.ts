import client from './client';
import { Project, CreateProjectRequest, Script, Storyboard } from '../types/project';

export const projectApi = {
  list: () => client.get<Project[]>('/projects'),

  get: (id: string) => client.get<Project>(`/projects/${id}`),

  create: (data: CreateProjectRequest) => client.post<Project>('/projects', data),

  update: (id: string, data: CreateProjectRequest) => client.put<Project>(`/projects/${id}`, data),

  delete: (id: string) => client.delete(`/projects/${id}`),

  generateScript: (id: string) => client.post<{ message: string; script: Script }>(`/projects/${id}/generate-script`, {}),

  getScript: (id: string) => client.get<Script>(`/projects/${id}/script`),

  updateScript: (id: string, content: string) => client.put<Script>(`/projects/${id}/script`, { content }),

  generateStoryboard: (id: string) => client.post<{ message: string; storyboard: Storyboard }>(`/projects/${id}/generate-storyboard`, {}),

  getStoryboard: (id: string) => client.get<Storyboard>(`/projects/${id}/storyboard`),

  updateStoryboard: (id: string, content: string) => client.put<Storyboard>(`/projects/${id}/storyboard`, { content }),
};