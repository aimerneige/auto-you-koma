import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { scriptApi } from '../api/scripts';
import { Sparkles, Save, ListTree, ArrowLeft } from 'lucide-react';

export const ScriptEditorPage = () => {
  const { id } = useParams<{ id: string }>(); // script id, or "new"
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [prompt, setPrompt] = useState('');
  const [parsedData, setParsedData] = useState<any[] | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (id && id !== 'new') {
      scriptApi.get(id).then(res => {
        setTitle(res.data.title);
        setContent(res.data.content);
        if (res.data.parsed_data) {
          try {
            setParsedData(JSON.parse(res.data.parsed_data));
          } catch(e) {}
        }
      });
    }
  }, [id]);

  const handleSave = async () => {
    setLoading(true);
    try {
      if (id === 'new') {
        const res = await scriptApi.create({ title, content, project_id: 'default' }); // temporary unified project
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
      await scriptApi.update(id as string, { title, content }); // save automatically
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

  return (
    <div style={{ display: 'flex', height: '100vh', background: '#f5f5f5' }}>
      {/* Sidebar: AI Tools */}
      <div style={{ width: 300, background: '#fff', borderRight: '1px solid #ddd', padding: 20, display: 'flex', flexDirection: 'column', gap: 15 }}>
        <button onClick={() => navigate('/dashboard')} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: 'transparent', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', marginBottom: 20 }}>
          <ArrowLeft size={16} /> Dashboard
        </button>

        <h3>AI Assistant</h3>
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
        
        <p style={{ fontSize: '0.9em', color: '#666' }}>一键利用 AI 理解上述内容，拆分出 4 个画面的场景描写和气泡对话。</p>
        <button onClick={handleParse} disabled={loading || id === 'new'} style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 5, padding: 10, background: id === 'new' ? '#ccc' : '#10b981', color: '#fff', border: 'none', borderRadius: 4, cursor: id === 'new' ? 'not-allowed' : 'pointer' }}>
          <ListTree size={16} /> Parse into Panels
        </button>
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
          <button onClick={handleSave} disabled={loading} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', borderRadius: 4, border: '1px solid #ddd', cursor: 'pointer', background: '#fff' }}>
            <Save size={16} /> Save
          </button>
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

          {/* Parsed Panels View */}
          {parsedData && (
            <div style={{ width: 450, display: 'flex', flexDirection: 'column' }}>
              <h4 style={{ margin: '0 0 10px 0' }}>Grid Storyboard (Pans)</h4>
              <div style={{ flex: 1, overflowY: 'auto', paddingRight: 5 }}>
                {parsedData.map((panel, idx) => (
                  <div key={idx} style={{ background: '#fff', border: '1px solid #ccc', padding: 15, borderRadius: 8, marginBottom: 15, boxShadow: '0 2px 4px rgba(0,0,0,0.05)' }}>
                    <h5 style={{ margin: '0 0 10px 0', borderBottom: '1px solid #eee', paddingBottom: 5, color: '#0070f3' }}>
                      P.{panel.panel || (idx + 1)}
                    </h5>
                    <div style={{ marginBottom: 10 }}>
                      <span style={{ fontSize: '0.8em', textTransform: 'uppercase', color: '#888', fontWeight: 'bold' }}>Visual Description</span>
                      <p style={{ margin: '5px 0 0 0', fontSize: '0.95em' }}>{panel.visual_desc}</p>
                    </div>
                    <div>
                      <span style={{ fontSize: '0.8em', textTransform: 'uppercase', color: '#888', fontWeight: 'bold' }}>Dialog</span>
                      <p style={{ margin: '5px 0 0 0', fontSize: '0.95em', fontStyle: 'italic' }}>"{panel.dialog}"</p>
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
