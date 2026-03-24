import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useProjectStore } from '../store/projectStore';
import { StoryboardPanel } from '../types/project';

export function StoryboardPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const {
    selectedProject,
    storyboardContent,
    generating,
    error,
    fetchProject,
    generateStoryboard,
    fetchStoryboard,
    updateStoryboard,
  } = useProjectStore();

  const [localPanels, setLocalPanels] = useState<StoryboardPanel[]>([]);
  const [selectedPanel, setSelectedPanel] = useState<number | null>(null);

  useEffect(() => {
    if (id) {
      fetchProject(id);
      fetchStoryboard(id);
    }
  }, [id]);

  useEffect(() => {
    if (storyboardContent) {
      setLocalPanels(storyboardContent.panels);
    }
  }, [storyboardContent]);

  const handleGenerate = async () => {
    if (!id) return;
    try {
      await generateStoryboard(id);
    } catch (error) {
      console.error('Failed to generate storyboard:', error);
    }
  };

  const handlePanelChange = (index: number, field: keyof StoryboardPanel, value: string) => {
    const newPanels = [...localPanels];
    newPanels[index] = { ...newPanels[index], [field]: value };
    setLocalPanels(newPanels);
  };

  const handleSave = async () => {
    if (!id || !storyboardContent) return;
    const updatedContent = { ...storyboardContent, panels: localPanels };
    try {
      await updateStoryboard(id, JSON.stringify(updatedContent));
      alert('Storyboard saved!');
    } catch (error) {
      console.error('Failed to save storyboard:', error);
    }
  };

  const handleConfirm = () => {
    handleSave();
    navigate(`/projects/${id}/render`);
  };

  if (!id) {
    return <div>Project not found</div>;
  }

  const shotTypes = ['wide_shot', 'medium_shot', 'close_up', 'extreme_close_up', 'over_shoulder', 'pov'];
  const angles = ['eye_level', 'high_angle', 'low_angle', 'dutch_angle', 'birds_eye', 'worms_eye'];

  return (
    <div className="page-container">
      <div className="page-header">
        <h1 className="title-comic">Storyboard Preview (HITL Node 3)</h1>
        <button className="btn btn-alt" onClick={() => navigate(`/projects/${id}/script`)}>
          ← Back to Script
        </button>
      </div>

      <div className="project-info-bar">
        <h2>{selectedProject?.title || 'Loading...'}</h2>
      </div>

      {error && <div className="error-message">{error}</div>}

      {!storyboardContent ? (
        <div className="generate-section">
          <p>Generate a detailed storyboard with visual parameters from your script.</p>
          <button className="btn" onClick={handleGenerate} disabled={generating}>
            {generating ? 'Generating Storyboard...' : 'Generate Storyboard'}
          </button>
        </div>
      ) : (
        <div className="previz-container">
          {/* Chat-style preview */}
          <div className="previz-chat">
            <h3>Preview</h3>
            {localPanels.map((panel, index) => (
              <div key={index} className="chat-panel-group">
                <div className="panel-divider">
                  <span>Panel {panel.panel_number}</span>
                </div>

                <div className="chat-bubble user">
                  <div className="bubble-header">
                    <span className="bubble-label">Scene:</span>
                  </div>
                  <p className="bubble-text">{panel.atmosphere}</p>
                </div>

                <div className="chat-bubble system">
                  <div className="bubble-header">
                    <span className="bubble-label">Shot:</span>
                    <span className="bubble-value">{panel.shot_type}</span>
                  </div>
                  <div className="bubble-header">
                    <span className="bubble-label">Angle:</span>
                    <span className="bubble-value">{panel.angle}</span>
                  </div>
                </div>

                {panel.dialogue && (
                  <div className="chat-bubble character">
                    <p className="bubble-text">{panel.dialogue}</p>
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Editor panel */}
          <div className="previz-editor">
            <h3>Edit Panel Details</h3>
            {selectedPanel !== null ? (
              <div className="panel-edit-form">
                <div className="form-section">
                  <h4>Panel {localPanels[selectedPanel].panel_number}</h4>

                  <div className="input-box">
                    <label>Shot Type</label>
                    <select
                      value={localPanels[selectedPanel].shot_type}
                      onChange={(e) => handlePanelChange(selectedPanel, 'shot_type', e.target.value)}
                    >
                      {shotTypes.map((type) => (
                        <option key={type} value={type}>{type.replace('_', ' ')}</option>
                      ))}
                    </select>
                  </div>

                  <div className="input-box">
                    <label>Angle</label>
                    <select
                      value={localPanels[selectedPanel].angle}
                      onChange={(e) => handlePanelChange(selectedPanel, 'angle', e.target.value)}
                    >
                      {angles.map((angle) => (
                        <option key={angle} value={angle}>{angle.replace('_', ' ')}</option>
                      ))}
                    </select>
                  </div>

                  <div className="input-box">
                    <label>Atmosphere</label>
                    <input
                      type="text"
                      value={localPanels[selectedPanel].atmosphere}
                      onChange={(e) => handlePanelChange(selectedPanel, 'atmosphere', e.target.value)}
                    />
                  </div>

                  <div className="input-box">
                    <label>Positive Prompt (for AI)</label>
                    <textarea
                      value={localPanels[selectedPanel].positive_prompt}
                      onChange={(e) => handlePanelChange(selectedPanel, 'positive_prompt', e.target.value)}
                      rows={3}
                    />
                  </div>

                  <div className="input-box">
                    <label>Dialogue</label>
                    <textarea
                      value={localPanels[selectedPanel].dialogue}
                      onChange={(e) => handlePanelChange(selectedPanel, 'dialogue', e.target.value)}
                      rows={2}
                    />
                  </div>

                  <button className="btn btn-alt" onClick={() => setSelectedPanel(null)}>
                    Done Editing
                  </button>
                </div>
              </div>
            ) : (
              <div className="panel-select">
                <p>Select a panel to edit:</p>
                <div className="panel-buttons">
                  {localPanels.map((panel, index) => (
                    <button
                      key={index}
                      className="btn"
                      onClick={() => setSelectedPanel(index)}
                    >
                      Panel {panel.panel_number}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {storyboardContent && (
        <div className="form-actions">
          <button className="btn" onClick={handleGenerate} disabled={generating}>
            {generating ? 'Regenerating...' : 'Regenerate Storyboard'}
          </button>
          <button className="btn btn-alt" onClick={handleSave}>
            Save
          </button>
          <button className="btn" onClick={handleConfirm}>
            Confirm & Continue to Render →
          </button>
        </div>
      )}
    </div>
  );
}