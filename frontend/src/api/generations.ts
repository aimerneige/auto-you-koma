import { apiClient } from './client';

export const generationApi = {
  start: (projectId: string, scriptId: string, layout: string) => 
    apiClient.post('/generations', { project_id: projectId, script_id: scriptId, layout }),
  get: (id: string) => apiClient.get(`/generations/${id}`)
};
