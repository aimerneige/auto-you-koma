import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useProjectStore } from '../store/projectStore';
import { projectApi } from '../api/projects';

export function RenderPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const {
    selectedProject,
    renderTask,
    renderResults,
    rendering,
    error,
    fetchProject,
    startRender,
    fetchRenderStatus,
    regeneratePanel,
    confirmRender,
  } = useProjectStore();

  const [renderOptions, setRenderOptions] = useState({
    exportType: 'clean_plate',
    layout: '2x2',
    imageWidth: 1024,
    imageHeight: 1024,
  });

  useEffect(() => {
    if (id) {
      fetchProject(id);
      fetchRenderStatus(id).catch(() => {
        // No render task yet, that's ok
      });
    }
  }, [id]);

  const handleStartRender = async () => {
    if (!id) return;
    try {
      await startRender(id, renderOptions);
    } catch (error) {
      console.error('Failed to start render:', error);
    }
  };

  const handleRegenerate = async (panelNumber: number) => {
    if (!id) return;
    try {
      await regeneratePanel(id, panelNumber);
    } catch (error) {
      console.error('Failed to regenerate panel:', error);
    }
  };

  const handleConfirm = async () => {
    if (!id) return;
    try {
      // Call composite to create final image
      await projectApi.composite(id, {
        export_type: renderOptions.exportType,
        layout: renderOptions.layout,
      });
      await confirmRender(id);
      alert('Comic is done! You can view and download it.');
    } catch (error) {
      console.error('Failed to confirm render:', error);
    }
  };

  const handleExport = async () => {
    if (!id) return;
    try {
      const response = await projectApi.exportProject(id);
      // In a real app, this would trigger a file download
      console.log('Export data:', response.data);
      alert('Export initiated! Check console for data.');
    } catch (error) {
      console.error('Failed to export:', error);
    }
  };

  if (!id) {
    return <div>Project not found</div>;
  }

  return (
    <div className="page-container">
      <div className="page-header">
        <h1 className="title-comic">Render & Export (HITL Node 4)</h1>
        <button className="btn btn-alt" onClick={() => navigate(`/projects/${id}/storyboard`)}>
          ← Back to Storyboard
        </button>
      </div>

      <div className="project-info-bar">
        <h2>{selectedProject?.title || 'Loading...'}</h2>
        <p>Review and approve generated panels</p>
      </div>

      {error && <div className="error-message">{error}</div>}

      {/* Render Settings */}
      {!renderTask && (
        <div className="render-settings">
          <h3>Render Settings</h3>
          <div className="settings-grid">
            <div className="input-box">
              <label>Export Type</label>
              <select
                value={renderOptions.exportType}
                onChange={(e) => setRenderOptions({ ...renderOptions, exportType: e.target.value })}
              >
                <option value="clean_plate">Clean Plate (without text)</option>
                <option value="native_text">With Text</option>
              </select>
            </div>

            <div className="input-box">
              <label>Layout</label>
              <select
                value={renderOptions.layout}
                onChange={(e) => setRenderOptions({ ...renderOptions, layout: e.target.value })}
              >
                <option value="2x2">2x2 Grid</option>
                <option value="1x4">1x4 Vertical Strip</option>
              </select>
            </div>

            <div className="input-box">
              <label>Image Width</label>
              <select
                value={renderOptions.imageWidth}
                onChange={(e) => setRenderOptions({ ...renderOptions, imageWidth: parseInt(e.target.value) })}
              >
                <option value="512">512px</option>
                <option value="768">768px</option>
                <option value="1024">1024px</option>
                <option value="2048">2048px</option>
              </select>
            </div>

            <div className="input-box">
              <label>Image Height</label>
              <select
                value={renderOptions.imageHeight}
                onChange={(e) => setRenderOptions({ ...renderOptions, imageHeight: parseInt(e.target.value) })}
              >
                <option value="512">512px</option>
                <option value="768">768px</option>
                <option value="1024">1024px</option>
                <option value="2048">2048px</option>
              </select>
            </div>
          </div>

          <button className="btn" onClick={handleStartRender} disabled={rendering}>
            {rendering ? 'Rendering...' : 'Start Rendering'}
          </button>
        </div>
      )}

      {/* Render Results / Quality Check */}
      {renderTask && (
        <div className="render-results">
          <h3>Generated Panels - Quality Check</h3>
          <p className="hint">Click "Regenerate" on any panel that needs improvement</p>

          <div className="panel-grid">
            {renderResults.map((result, index) => (
              <div key={index} className="render-panel">
                <div className="panel-label">Panel {result.panel_number}</div>
                <div className="panel-image">
                  {result.success ? (
                    <img src={result.image_url} alt={`Panel ${result.panel_number}`} />
                  ) : (
                    <div className="error-placeholder">Failed to generate</div>
                  )}
                </div>
                <div className="panel-actions">
                  <button
                    className="btn btn-alt"
                    onClick={() => handleRegenerate(result.panel_number)}
                    disabled={rendering}
                  >
                    {rendering ? 'Regenerating...' : 'Regenerate'}
                  </button>
                </div>
                {result.error && <p className="panel-error">{result.error}</p>}
              </div>
            ))}
          </div>

          <div className="form-actions">
            <button className="btn btn-alt" onClick={() => setRenderOptions({ ...renderOptions })}>
              Change Settings
            </button>
            <button className="btn btn-alt" onClick={handleExport}>
              Export
            </button>
            <button className="btn" onClick={handleConfirm}>
              Confirm & Export →
            </button>
          </div>
        </div>
      )}
    </div>
  );
}