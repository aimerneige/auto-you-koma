export interface Character {
  id: string;
  user_id: string;
  name: string;
  name_jp: string;
  gender: string;
  age: string;
  personality: string;
  backstory: string;
  visual_prompt: string;
  tags: string;
  category: string;
  reference_sheet_url: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCharacterRequest {
  name: string;
  name_jp?: string;
  gender?: string;
  age?: string;
  personality?: string;
  backstory?: string;
  visual_prompt?: string;
  tags?: string;
  category?: string;
}