import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProjectStore } from '../store/projectStore';

export function ProjectCreatePage() {
  const navigate = useNavigate();
  const { createProject, loading } = useProjectStore();

  const [formData, setFormData] = useState({
    title: '',
    synopsis: '',
    mode: 'standalone',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const project = await createProject(formData);
      navigate(`/projects/${project.id}`);
    } catch (error) {
      console.error('Failed to create project:', error);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <div className="page-container">
      <h1 className="title-comic">Create New Project</h1>

      <form onSubmit={handleSubmit} className="project-form">
        <div className="form-section">
          <h3>Project Info</h3>

          <div className="input-box">
            <label>Title *</label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
              required
              placeholder="Your comic's title"
            />
          </div>

          <div className="input-box">
            <label>Synopsis *</label>
            <textarea
              name="synopsis"
              value={formData.synopsis}
              onChange={handleChange}
              required
              rows={4}
              placeholder="What happens in your comic? (e.g., A magical girl forgets her transformation words...)"
            />
          </div>

          <div className="input-box">
            <label>Mode</label>
            <select name="mode" value={formData.mode} onChange={handleChange}>
              <option value="standalone">Standalone (Single 4-panel comic)</option>
              <option value="serialized">Serialized (Multiple episodes)</option>
            </select>
          </div>
        </div>

        <div className="form-actions">
          <button type="button" className="btn btn-alt" onClick={() => navigate('/projects')}>
            Cancel
          </button>
          <button type="submit" className="btn" disabled={loading}>
            {loading ? 'Creating...' : 'Create Project'}
          </button>
        </div>
      </form>
    </div>
  );
}