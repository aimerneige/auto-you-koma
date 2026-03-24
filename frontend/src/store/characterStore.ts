import { create } from 'zustand';
import { Character } from '../types/character';
import { characterApi } from '../api/characters';

interface CharacterState {
  characters: Character[];
  selectedCharacter: Character | null;
  loading: boolean;
  error: string | null;
  fetchCharacters: () => Promise<void>;
  fetchCharacter: (id: string) => Promise<void>;
  createCharacter: (data: Partial<Character>) => Promise<Character>;
  updateCharacter: (id: string, data: Partial<Character>) => Promise<Character>;
  deleteCharacter: (id: string) => Promise<void>;
}

export const useCharacterStore = create<CharacterState>((set, get) => ({
  characters: [],
  selectedCharacter: null,
  loading: false,
  error: null,

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
}));