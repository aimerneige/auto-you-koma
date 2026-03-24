export interface RenderTask {
  id: string;
  project_id: string;
  storyboard_id: string;
  export_type: string;
  layout: string;
  image_width: number;
  image_height: number;
  status: string;
  output_paths: string;
  error_message?: string;
  created_at: string;
  updated_at: string;
}

export interface PanelRenderResult {
  panel_number: number;
  image_url: string;
  seed: number;
  success: boolean;
  error?: string;
}