import client from './client';
import { Project, CreateProjectRequest, Script, Storyboard, RenderTask, PanelRenderResult } from '../types/project';
import { RenderTask as RenderTaskType, PanelRenderResult as PanelRenderResultType } from '../types/render';

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

  startRender: (id: string, data: { export_type?: string; layout?: string; image_width?: number; image_height?: number }) =>
    client.post<{ message: string; task: RenderTaskType }>(`/projects/${id}/render`, data),

  getRenderStatus: (id: string) => client.get<RenderTaskType>(`/projects/${id}/render-status`),

  regeneratePanel: (id: string, panelNumber: number) =>
    client.post<PanelRenderResultType>(`/projects/${id}/render/regenerate`, { panel_number: panelNumber }),

  confirmRender: (id: string) => client.post(`/projects/${id}/render/confirm`, {}),
};