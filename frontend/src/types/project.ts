export interface Project {
  id: string;
  user_id: string;
  title: string;
  mode: string;
  status: string;
  synopsis: string;
  created_at: string;
  updated_at: string;
}

export interface CreateProjectRequest {
  title: string;
  mode?: string;
  synopsis: string;
}

export interface ScriptPanel {
  panel_number: number;
  structure: string;
  scene: string;
  characters: string;
  dialogue: string;
  narration?: string;
}

export interface ScriptContent {
  project_id: string;
  title: string;
  synopsis: string;
  mode: string;
  panels: ScriptPanel[];
}

export interface Script {
  id: string;
  project_id: string;
  episode_number: number;
  content: string;
  version: number;
  created_at: string;
  updated_at: string;
}

export interface StoryboardPanel {
  panel_number: number;
  shot_type: string;
  angle: string;
  atmosphere: string;
  characters: string;
  action: string;
  expression: string;
  positive_prompt: string;
  negative_prompt: string;
  dialogue: string;
  seed: number;
}

export interface StoryboardContent {
  storyboard_id: string;
  script_id: string;
  panels: StoryboardPanel[];
}

export interface Storyboard {
  id: string;
  script_id: string;
  content: string;
  version: number;
  created_at: string;
  updated_at: string;
}