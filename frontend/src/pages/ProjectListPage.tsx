import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProjectStore } from '../store/projectStore';

export function ProjectListPage() {
  const navigate = useNavigate();
  const { projects, loading, error, fetchProjects, deleteProject } = useProjectStore();

  useEffect(() => {
    fetchProjects();
  }, []);

  const handleDelete = async (id: string) => {
    if (confirm('Delete this project?')) {
      await deleteProject(id);
    }
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h1 className="title-comic">Projects</h1>
        <button className="btn" onClick={() => navigate('/projects/new')}>
          + New Project
        </button>
      </div>

      {loading && <div className="loading">Loading...</div>}
      {error && <div className="error-message">{error}</div>}

      <div className="project-grid">
        {projects.map((project) => (
          <div
            key={project.id}
            className="project-card"
            onClick={() => navigate(`/projects/${project.id}`)}
          >
            <div className="project-info">
              <h3>{project.title}</h3>
              <p className="project-synopsis">{project.synopsis}</p>
              <div className="project-meta">
                <span className={`status-badge status-${project.status}`}>{project.status}</span>
                <span className="mode-tag">{project.mode || 'standalone'}</span>
              </div>
            </div>
            <button
              className="delete-btn"
              onClick={(e) => {
                e.stopPropagation();
                handleDelete(project.id);
              }}
            >
              ×
            </button>
          </div>
        ))}
      </div>

      {projects.length === 0 && !loading && (
        <div className="empty-state">
          <p>No projects yet. Create your first comic!</p>
        </div>
      )}
    </div>
  );
}