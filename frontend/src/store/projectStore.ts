import { create } from 'zustand';
import { Project, CreateProjectRequest, Script, ScriptContent } from '../types/project';
import { projectApi } from '../api/projects';

interface ProjectState {
  projects: Project[];
  selectedProject: Project | null;
  currentScript: Script | null;
  scriptContent: ScriptContent | null;
  loading: boolean;
  generating: boolean;
  error: string | null;
  fetchProjects: () => Promise<void>;
  fetchProject: (id: string) => Promise<void>;
  createProject: (data: CreateProjectRequest) => Promise<Project>;
  updateProject: (id: string, data: CreateProjectRequest) => Promise<Project>;
  deleteProject: (id: string) => Promise<void>;
  generateScript: (id: string) => Promise<void>;
  fetchScript: (id: string) => Promise<void>;
  updateScript: (id: string, content: string) => Promise<void>;
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  projects: [],
  selectedProject: null,
  currentScript: null,
  scriptContent: null,
  loading: false,
  generating: false,
  error: null,

  fetchProjects: async () => {
    set({ loading: true, error: null });
    try {
      const response = await projectApi.list();
      set({ projects: response.data, loading: false });
    } catch (error: any) {
      set({ error: error.message, loading: false });
    }
  },

  fetchProject: async (id: string) => {
    set({ loading: true, error: null });
    try {
      const response = await projectApi.get(id);
      set({ selectedProject: response.data, loading: false });
    } catch (error: any) {
      set({ error: error.message, loading: false });
    }
  },

  createProject: async (data: CreateProjectRequest) => {
    set({ loading: true, error: null });
    try {
      const response = await projectApi.create(data);
      set((state) => ({
        projects: [...state.projects, response.data],
        loading: false,
      }));
      return response.data;
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  updateProject: async (id: string, data: CreateProjectRequest) => {
    set({ loading: true, error: null });
    try {
      const response = await projectApi.update(id, data);
      set((state) => ({
        projects: state.projects.map((p) => (p.id === id ? response.data : p)),
        selectedProject: response.data,
        loading: false,
      }));
      return response.data;
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  deleteProject: async (id: string) => {
    set({ loading: true, error: null });
    try {
      await projectApi.delete(id);
      set((state) => ({
        projects: state.projects.filter((p) => p.id !== id),
        loading: false,
      }));
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  generateScript: async (id: string) => {
    set({ generating: true, error: null });
    try {
      await projectApi.generateScript(id);
      // Fetch the generated script
      const response = await projectApi.getScript(id);
      set({
        currentScript: response.data,
        scriptContent: JSON.parse(response.data.content),
        generating: false,
      });
    } catch (error: any) {
      set({ generating: false, error: error.message });
      throw error;
    }
  },

  fetchScript: async (id: string) => {
    set({ loading: true, error: null });
    try {
      const response = await projectApi.getScript(id);
      set({
        currentScript: response.data,
        scriptContent: JSON.parse(response.data.content),
        loading: false,
      });
    } catch (error: any) {
      set({ error: error.message, loading: false });
    }
  },

  updateScript: async (id: string, content: string) => {
    set({ loading: true, error: null });
    try {
      await projectApi.updateScript(id, content);
      set({
        scriptContent: JSON.parse(content),
        loading: false,
      });
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },
}));