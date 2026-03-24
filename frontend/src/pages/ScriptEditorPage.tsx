import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useProjectStore } from '../store/projectStore';
import { ScriptPanel } from '../types/project';

export function ScriptEditorPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { selectedProject, scriptContent, generating, error, fetchProject, generateScript, fetchScript, updateScript } = useProjectStore();
  const [localPanels, setLocalPanels] = useState<ScriptPanel[]>([]);

  useEffect(() => {
    if (id) {
      fetchProject(id);
      fetchScript(id);
    }
  }, [id]);

  useEffect(() => {
    if (scriptContent) {
      setLocalPanels(scriptContent.panels);
    }
  }, [scriptContent]);

  const handleGenerate = async () => {
    if (!id) return;
    try {
      await generateScript(id);
    } catch (error) {
      console.error('Failed to generate script:', error);
    }
  };

  const handlePanelChange = (index: number, field: keyof ScriptPanel, value: string) => {
    const newPanels = [...localPanels];
    newPanels[index] = { ...newPanels[index], [field]: value };
    setLocalPanels(newPanels);
  };

  const handleSave = async () => {
    if (!id || !scriptContent) return;
    const updatedContent = { ...scriptContent, panels: localPanels };
    try {
      await updateScript(id, JSON.stringify(updatedContent));
      alert('Script saved successfully!');
    } catch (error) {
      console.error('Failed to save script:', error);
    }
  };

  const handleConfirm = () => {
    handleSave();
    navigate(`/projects/${id}/storyboard`);
  };

  if (!id) {
    return <div>Project not found</div>;
  }

  const structureLabels: Record<string, string> = {
    '起': 'Opening',
    '承': 'Development',
    '转': 'Twist',
    '合': 'Conclusion',
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h1 className="title-comic">Script Editor (HITL Node 2)</h1>
        <button className="btn btn-alt" onClick={() => navigate(`/projects/${id}`)}>
          ← Back
        </button>
      </div>

      <div className="project-info-bar">
        <h2>{selectedProject?.title || 'Loading...'}</h2>
        <p>{selectedProject?.synopsis}</p>
      </div>

      {error && <div className="error-message">{error}</div>}

      {!scriptContent ? (
        <div className="generate-section">
          <p>Generate a 4-panel script based on your project synopsis.</p>
          <button className="btn" onClick={handleGenerate} disabled={generating}>
            {generating ? 'Generating Script...' : 'Generate Script'}
          </button>
        </div>
      ) : (
        <>
          <div className="script-editor">
            {localPanels.map((panel, index) => (
              <div key={index} className="panel-editor">
                <div className="panel-header">
                  <span className="panel-number">Panel {panel.panel_number}</span>
                  <span className="panel-structure">{structureLabels[panel.structure] || panel.structure}</span>
                </div>

                <div className="panel-fields">
                  <div className="input-box">
                    <label>Scene</label>
                    <input
                      type="text"
                      value={panel.scene}
                      onChange={(e) => handlePanelChange(index, 'scene', e.target.value)}
                    />
                  </div>

                  <div className="input-box">
                    <label>Characters</label>
                    <input
                      type="text"
                      value={panel.characters}
                      onChange={(e) => handlePanelChange(index, 'characters', e.target.value)}
                    />
                  </div>

                  <div className="input-box">
                    <label>Dialogue</label>
                    <textarea
                      value={panel.dialogue}
                      onChange={(e) => handlePanelChange(index, 'dialogue', e.target.value)}
                      rows={3}
                    />
                  </div>

                  <div className="input-box">
                    <label>Narration (optional)</label>
                    <textarea
                      value={panel.narration || ''}
                      onChange={(e) => handlePanelChange(index, 'narration', e.target.value)}
                      rows={2}
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>

          <div className="form-actions">
            <button className="btn" onClick={handleGenerate} disabled={generating}>
              {generating ? 'Regenerating...' : 'Regenerate Script'}
            </button>
            <button className="btn btn-alt" onClick={handleSave}>
              Save
            </button>
            <button className="btn" onClick={handleConfirm}>
              Confirm & Continue to Storyboard →
            </button>
          </div>
        </>
      )}
    </div>
  );
}