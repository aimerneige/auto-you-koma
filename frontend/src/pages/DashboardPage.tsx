import { useNavigate } from 'react-router-dom';

export function DashboardPage() {
  const navigate = useNavigate();

  return (
    <div className="page-container">
      <h1 className="title-comic">Welcome to Auto Yon Koma!</h1>

      <div className="dashboard-grid">
        <div className="dashboard-card" onClick={() => navigate('/characters')}>
          <div className="card-icon">👤</div>
          <h3>Characters</h3>
          <p>Manage your character library</p>
        </div>

        <div className="dashboard-card" onClick={() => navigate('/projects/new')}>
          <div className="card-icon">📝</div>
          <h3>New Project</h3>
          <p>Create a new 4-panel comic</p>
        </div>

        <div className="dashboard-card" onClick={() => navigate('/projects')}>
          <div className="card-icon">📚</div>
          <h3>Projects</h3>
          <p>View your ongoing projects</p>
        </div>
      </div>

      <div className="quick-start">
        <h2>Quick Start</h2>
        <p>Ready to create your first comic? Follow these steps:</p>
        <ol className="steps-list">
          <li>Create or select characters in the Character Library</li>
          <li>Generate reference sheets for your characters (HITL Node 1)</li>
          <li>Create a new project and write your story idea</li>
          <li>Review and edit the generated script (HITL Node 2)</li>
          <li>Preview and adjust the storyboard (HITL Node 3)</li>
          <li>Generate and approve images (HITL Node 4)</li>
          <li>Export your finished comic!</li>
        </ol>
        <button className="btn" onClick={() => navigate('/projects/new')}>
          Start Creating!
        </button>
      </div>
    </div>
  );
}