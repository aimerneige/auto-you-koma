import { create } from 'zustand';
import { Character } from '../types/character';
import { characterApi } from '../api/characters';

interface CharacterState {
  characters: Character[];
  selectedCharacter: Character | null;
  loading: boolean;
  error: string | null;
  generatingReference: boolean;
  generatedReferenceUrl: string | null;
  fetchCharacters: () => Promise<void>;
  fetchCharacter: (id: string) => Promise<void>;
  createCharacter: (data: Partial<Character>) => Promise<Character>;
  updateCharacter: (id: string, data: Partial<Character>) => Promise<Character>;
  deleteCharacter: (id: string) => Promise<void>;
  generateReferenceSheet: (id: string) => Promise<string>;
  confirmReferenceSheet: (id: string, url: string) => Promise<void>;
}

export const useCharacterStore = create<CharacterState>((set, get) => ({
  characters: [],
  selectedCharacter: null,
  loading: false,
  error: null,
  generatingReference: false,
  generatedReferenceUrl: null,

  fetchCharacters: async () => {
    set({ loading: true, error: null });
    try {
      const response = await characterApi.list();
      set({ characters: response.data, loading: false });
    } catch (error: any) {
      set({ error: error.message, loading: false });
    }
  },

  fetchCharacter: async (id: string) => {
    set({ loading: true, error: null });
    try {
      const response = await characterApi.get(id);
      set({ selectedCharacter: response.data, loading: false });
    } catch (error: any) {
      set({ error: error.message, loading: false });
    }
  },

  createCharacter: async (data: Partial<Character>) => {
    set({ loading: true, error: null });
    try {
      const response = await characterApi.create(data as any);
      set((state) => ({
        characters: [...state.characters, response.data],
        loading: false,
      }));
      return response.data;
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  updateCharacter: async (id: string, data: Partial<Character>) => {
    set({ loading: true, error: null });
    try {
      const response = await characterApi.update(id, data as any);
      set((state) => ({
        characters: state.characters.map((c) => (c.id === id ? response.data : c)),
        selectedCharacter: response.data,
        loading: false,
      }));
      return response.data;
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  deleteCharacter: async (id: string) => {
    set({ loading: true, error: null });
    try {
      await characterApi.delete(id);
      set((state) => ({
        characters: state.characters.filter((c) => c.id !== id),
        loading: false,
      }));
    } catch (error: any) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  generateReferenceSheet: async (id: string) => {
    set({ generatingReference: true, error: null, generatedReferenceUrl: null });
    try {
      const response = await characterApi.generateReference(id);
      set({ generatingReference: false, generatedReferenceUrl: response.data.image_url });
      return response.data.image_url;
    } catch (error: any) {
      set({ generatingReference: false, error: error.message });
      throw error;
    }
  },

  confirmReferenceSheet: async (id: string, url: string) => {
    set({ loading: true, error: null });
    try {
      await characterApi.confirmReference(id, url);
      // Refresh character data
      const response = await characterApi.get(id);
      set((state) => ({
        characters: state.characters.map((c) => (c.id === id ? response.data : c)),
        selectedCharacter: response.data,
        generatedReferenceUrl: null,
        loading: false,
      }));
    } catch (error: any) {
      set({ loading: false, error: error.message });
      throw error;
    }
  },
}));