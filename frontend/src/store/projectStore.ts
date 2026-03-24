import { create } from 'zustand';
import { Project, CreateProjectRequest, Script, ScriptContent, Storyboard, StoryboardContent, PanelRenderResult } from '../types/project';
import { RenderTask } from '../types/render';
import { projectApi } from '../api/projects';

interface ProjectState {
  projects: Project[];
  selectedProject: Project | null;
  currentScript: Script | null;
  scriptContent: ScriptContent | null;
  currentStoryboard: Storyboard | null;
  storyboardContent: StoryboardContent | null;
  renderTask: RenderTask | null;
  renderResults: PanelRenderResult[];
  loading: boolean;
  generating: boolean;
  rendering: boolean;
  error: string | null;
  fetchProjects: () => Promise<void>;
  fetchProject: (id: string) => Promise<void>;
  createProject: (data: CreateProjectRequest) => Promise<Project>;
  updateProject: (id: string, data: CreateProjectRequest) => Promise<Project>;
  deleteProject: (id: string) => Promise<void>;
  generateScript: (id: string) => Promise<void>;
  fetchScript: (id: string) => Promise<void>;
  updateScript: (id: string, content: string) => Promise<void>;
  generateStoryboard: (id: string) => Promise<void>;
  fetchStoryboard: (id: string) => Promise<void>;
  updateStoryboard: (id: string, content: string) => Promise<void>;
  startRender: (id: string, options?: any) => Promise<void>;
  fetchRenderStatus: (id: string) => Promise<void>;
  regeneratePanel: (id: string, panelNumber: number) => Promise<void>;
  confirmRender: (id: string) => Promise<void>;
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  projects: [],
  selectedProject: null,
  currentScript: null,
  scriptContent: null,
  currentStoryboard: null,
  storyboardContent: null,
  renderTask: null,
  renderResults: [],
  loading: false,
  generating: false,
  rendering: false,
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

  generateStoryboard: async (id: string) => {
    set({ generating: true, error: null });
    try {
      await projectApi.generateStoryboard(id);
      const response = await projectApi.getStoryboard(id);
      set({
        currentStoryboard: response.data,
        storyboardContent: JSON.parse(response.data.content),
        generating: false,
      });
    } catch (error: any) {
      set({ generating: false, error: error.message });
      throw error;
    }
  },

  fetchStoryboard: async (id: string) => {
    set({ loading: true, error: null });
    try {
      const response = await projectApi.getStoryboard(id);
      set({
        currentStoryboard: response.data,
        storyboardContent: JSON.parse(response.data.content),
        loading: false,
      });
    } catch (error: any) {
      set({ error: error.message, loading: false });
    }
  },

  updateStoryboard: async (id: string, content: string) => {
    set({ loading: true, error: null });
    try {
      await projectApi.updateStoryboard(id, content);
      set({
        storyboardContent: JSON.parse(content),
        loading: false,
      });
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  startRender: async (id: string, options?: any) => {
    set({ rendering: true, error: null });
    try {
      await projectApi.startRender(id, options);
      // Poll for status
      const response = await projectApi.getRenderStatus(id);
      set({
        renderTask: response.data,
        renderResults: JSON.parse(response.data.output_paths || '[]'),
        rendering: false,
      });
    } catch (error: any) {
      set({ rendering: false, error: error.message });
      throw error;
    }
  },

  fetchRenderStatus: async (id: string) => {
    set({ loading: true, error: null });
    try {
      const response = await projectApi.getRenderStatus(id);
      set({
        renderTask: response.data,
        renderResults: JSON.parse(response.data.output_paths || '[]'),
        loading: false,
      });
    } catch (error: any) {
      set({ error: error.message, loading: false });
    }
  },

  regeneratePanel: async (id: string, panelNumber: number) => {
    set({ rendering: true, error: null });
    try {
      const result = await projectApi.regeneratePanel(id, panelNumber);
      // Update the specific panel in renderResults
      set((state) => ({
        renderResults: state.renderResults.map((r) =>
          r.panel_number === panelNumber ? result.data : r
        ),
        rendering: false,
      }));
    } catch (error: any) {
      set({ rendering: false, error: error.message });
      throw error;
    }
  },

  confirmRender: async (id: string) => {
    set({ loading: true, error: null });
    try {
      await projectApi.confirmRender(id);
      set({ loading: false });
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },
}));