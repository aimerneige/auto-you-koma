import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { scriptApi } from '../api/scripts';
import { Sparkles, Save, ListTree, ArrowLeft, Lock, Unlock, Edit3, X } from 'lucide-react';

export const ScriptEditorPage = () => {
  const { id } = useParams<{ id: string }>(); // script id, or "new"
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [prompt, setPrompt] = useState('');
  const [parsedData, setParsedData] = useState<any[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [focusedPanelIndex, setFocusedPanelIndex] = useState<number | null>(null);
  const [panelInstructions, setPanelInstructions] = useState('');

  const fetchScript = async () => {
    if (id && id !== 'new') {
      const res = await scriptApi.get(id);
      setTitle(res.data.title);
      setContent(res.data.content);
      if (res.data.parsed_data) {
        try {
          setParsedData(JSON.parse(res.data.parsed_data));
        } catch(e) {}
      }
    }
  };

  useEffect(() => {
    fetchScript();
  }, [id]);

  const handleSave = async () => {
    setLoading(true);
    try {
      if (id === 'new') {
        const res = await scriptApi.create({ title, content, project_id: 'default' });
        navigate(`/scripts/${res.data.id}`);
      } else {
        await scriptApi.update(id as string, { title, content });
      }
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async () => {
    if (!prompt) return;
    setLoading(true);
    try {
      const res = await scriptApi.generate(prompt);
      setContent(content + (content ? '\n\n' : '') + res.data.content);
    } finally {
      setLoading(false);
    }
  };

  const handleParse = async () => {
    if (!id || id === 'new') return alert('Please save the script first below');
    setLoading(true);
    try {
      await scriptApi.update(id as string, { title, content });
      const res = await scriptApi.parse(id);
      try {
        setParsedData(JSON.parse(res.data.parsed_data));
      } catch(e) {
        alert("Failed to parse the LLM's JSON output. Please tweak prompt.");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleToggleLock = async (index: number) => {
    if (!parsedData || !id) return;
    const newPanels = [...parsedData];
    newPanels[index].locked = !newPanels[index].locked;
    setParsedData(newPanels);
    await scriptApi.updatePanel(id, index, newPanels[index]);
  };

  const handlePanelContentChange = async (index: number, field: 'visual_desc' | 'dialog', value: string) => {
    if (!parsedData || !id) return;
    const newPanels = [...parsedData];
    newPanels[index][field] = value;
    setParsedData(newPanels);
    // Debounced save would be better, but saving directly for prototype
    await scriptApi.updatePanel(id, index, newPanels[index]);
  };

  const handleRegeneratePanel = async () => {
    if (focusedPanelIndex === null || !id || !panelInstructions) return;
    setLoading(true);
    try {
      await scriptApi.regeneratePanel(id, focusedPanelIndex, panelInstructions);
      await fetchScript(); // refresh data
      setPanelInstructions('');
      setFocusedPanelIndex(null);
    } catch(e: any) {
      alert("Error regenerating panel: " + e.response?.data?.error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', height: '100vh', background: '#f5f5f5' }}>
      {/* Sidebar: AI Tools */}
      <div style={{ width: 320, background: '#fff', borderRight: '1px solid #ddd', padding: 20, display: 'flex', flexDirection: 'column', gap: 15 }}>
        <button onClick={() => navigate('/dashboard')} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: 'transparent', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', marginBottom: 20 }}>
          <ArrowLeft size={16} /> Dashboard
        </button>

        {focusedPanelIndex === null ? (
          <>
            <h3>AI Story Assistant</h3>
            <textarea 
              placeholder="例如：'两个小学生发现学校有外星人基地的搞笑短文'" 
              value={prompt} 
              onChange={e => setPrompt(e.target.value)} 
              style={{ height: 100, padding: 8, fontSize: '0.9em' }} 
            />
            <button onClick={handleGenerate} disabled={loading} style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 5, padding: 10, background: '#0070f3', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', opacity: loading ? 0.7 : 1 }}>
              <Sparkles size={16} /> {loading ? 'Thinking...' : 'Generate Story'}
            </button>

            <hr style={{ margin: '20px 0', border: 'none', borderTop: '1px solid #eee' }} />
            
            <p style={{ fontSize: '0.9em', color: '#666' }}>一键利用 AI 理解上述内容，拆分出明确的画面场景描写和气泡对话。</p>
            <button onClick={handleParse} disabled={loading || id === 'new'} style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 5, padding: 10, background: id === 'new' ? '#ccc' : '#10b981', color: '#fff', border: 'none', borderRadius: 4, cursor: id === 'new' ? 'not-allowed' : 'pointer' }}>
              <ListTree size={16} /> Parse into Panels
            </button>
          </>
        ) : (
          <>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <h3 style={{ margin: 0, color: '#10b981' }}>Director Chat (P.{focusedPanelIndex + 1})</h3>
              <button onClick={() => setFocusedPanelIndex(null)} style={{ background: 'transparent', border: 'none', cursor: 'pointer' }}><X size={20} /></button>
            </div>
            <p style={{ fontSize: '0.9em', color: '#666' }}>对当前选择的格子发送修改意见，AI Agent 将替你重写。</p>
            <textarea 
              placeholder="例如：把这段对话改得更暴躁一点；或是背景改成夜晚。" 
              value={panelInstructions} 
              onChange={e => setPanelInstructions(e.target.value)} 
              style={{ height: 150, padding: 8, fontSize: '0.9em' }} 
            />
            <button onClick={handleRegeneratePanel} disabled={loading || !panelInstructions} style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 5, padding: 10, background: '#10b981', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', opacity: loading ? 0.7 : 1 }}>
              <Edit3 size={16} /> {loading ? 'Regenerating...' : 'Regenerate Panel'}
            </button>
          </>
        )}
      </div>

      {/* Main content: Editor & Panels */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        <header style={{ padding: '15px 20px', background: '#fff', borderBottom: '1px solid #ddd', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <input 
            type="text" 
            value={title} 
            onChange={e => setTitle(e.target.value)} 
            placeholder="Untitled Script" 
            style={{ fontSize: '1.2em', fontWeight: 'bold', border: 'none', outline: 'none', background: 'transparent' }} 
          />
          <div style={{ display: 'flex', gap: 10 }}>
            <button onClick={handleSave} disabled={loading} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', borderRadius: 4, border: '1px solid #ddd', cursor: 'pointer', background: '#fff' }}>
              <Save size={16} /> Save
            </button>
            <button onClick={() => {
              if (id === 'new') alert('Please save first');
              else navigate(`/generate/${id}`);
            }} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', borderRadius: 4, border: 'none', cursor: 'pointer', background: '#000', color: '#fff', fontWeight: 'bold' }}>
              Proceed to Render &rarr;
            </button>
          </div>
        </header>

        <div style={{ flex: 1, overflowY: 'auto', padding: 20, display: 'flex', gap: 20 }}>
          {/* Text Editor */}
          <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
            <h4 style={{ margin: '0 0 10px 0' }}>Manuscript (Story Document)</h4>
            <textarea 
              value={content} 
              onChange={e => setContent(e.target.value)} 
              placeholder="在这里编写原始长生剧本，或使用左侧 AI 扩写。"
              style={{ flex: 1, padding: 15, border: '1px solid #ddd', borderRadius: 8, fontSize: '1.1em', resize: 'none', lineHeight: 1.6, outline: 'none' }}
            />
          </div>

          {/* Parsed Panels View - Previz */}
          {parsedData && (
            <div style={{ width: 500, display: 'flex', flexDirection: 'column' }}>
              <h4 style={{ margin: '0 0 10px 0' }}>Interactive Previz Panels</h4>
              <div style={{ flex: 1, overflowY: 'auto', paddingRight: 5 }}>
                {parsedData.map((panel, idx) => (
                  <div 
                    key={idx} 
                    style={{ 
                      background: '#fff', 
                      border: focusedPanelIndex === idx ? '2px solid #10b981' : '1px solid #ccc', 
                      padding: 15, 
                      borderRadius: 8, 
                      marginBottom: 15, 
                      boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
                      transition: 'border 0.2s',
                      position: 'relative'
                    }}
                  >
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #eee', paddingBottom: 5, marginBottom: 10 }}>
                      <h5 style={{ margin: 0, color: focusedPanelIndex === idx ? '#10b981' : '#0070f3', cursor: 'pointer' }} onClick={() => setFocusedPanelIndex(idx)}>
                        P.{panel.panel || (idx + 1)} {panel.locked && '(Locked)'}
                      </h5>
                      <div style={{ display: 'flex', gap: 10 }}>
                        <button onClick={() => setFocusedPanelIndex(idx)} style={{ background: 'transparent', border: '1px solid #10b981', color: '#10b981', borderRadius: 4, padding: '4px 8px', cursor: 'pointer', fontSize: '0.8em' }}>Chat AI</button>
                        <button onClick={() => handleToggleLock(idx)} style={{ background: 'transparent', border: 'none', cursor: 'pointer', color: panel.locked ? '#f43f5e' : '#aaa' }}>
                          {panel.locked ? <Lock size={18} /> : <Unlock size={18} />}
                        </button>
                      </div>
                    </div>
                    
                    <div style={{ marginBottom: 10 }}>
                      <span style={{ fontSize: '0.8em', textTransform: 'uppercase', color: '#888', fontWeight: 'bold' }}>Visual Description</span>
                      <textarea 
                        value={panel.visual_desc || ''} 
                        onChange={e => handlePanelContentChange(idx, 'visual_desc', e.target.value)}
                        disabled={panel.locked}
                        style={{ width: '100%', minHeight: 60, marginTop: 5, padding: 8, boxSizing: 'border-box', border: '1px solid #eee', borderRadius: 4, fontSize: '0.95em', resize: 'vertical' }}
                      />
                    </div>
                    <div>
                      <span style={{ fontSize: '0.8em', textTransform: 'uppercase', color: '#888', fontWeight: 'bold' }}>Dialog</span>
                      <textarea 
                        value={panel.dialog || ''} 
                        onChange={e => handlePanelContentChange(idx, 'dialog', e.target.value)}
                        disabled={panel.locked}
                        style={{ width: '100%', minHeight: 60, marginTop: 5, padding: 8, boxSizing: 'border-box', border: '1px solid #eee', borderRadius: 4, fontSize: '0.95em', fontStyle: 'italic', resize: 'vertical' }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
