import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { generationApi } from '../api/generations';
import { Loader2, Clapperboard, ArrowLeft } from 'lucide-react';

export const ComicGenerationPage = () => {
  const { scriptId } = useParams<{ scriptId: string }>();
  const navigate = useNavigate();
  const [layout, setLayout] = useState('2x2');
  const [status, setStatus] = useState<string>('idle');
  const [jobId, setJobId] = useState<string | null>(null);

  const startGeneration = async () => {
    try {
      setStatus('launching');
      const res = await generationApi.start('default', scriptId as string, layout);
      setJobId(res.data.id);
      setStatus('processing');
    } catch (e) {
      alert("Failed to start formatting: " + e);
      setStatus('idle');
    }
  };

  useEffect(() => {
    if (!jobId || status !== 'processing') return;
    
    const interval = setInterval(async () => {
      try {
        const res = await generationApi.get(jobId);
        if (res.data.status === 'done') {
          setStatus('done');
          clearInterval(interval);
          navigate(`/viewer/${res.data.id}`);
        } else if (res.data.status === 'failed') {
          setStatus('failed');
          alert("Generation failed: " + res.data.error);
          clearInterval(interval);
        }
      } catch (e) {
        console.error("Polling error", e);
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [jobId, status, navigate]);

  return (
    <div style={{ padding: 40, maxWidth: 600, margin: '0 auto', textAlign: 'center' }}>
      <button onClick={() => navigate(-1)} style={{ position: 'absolute', top: 20, left: 20, display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: 'transparent', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer' }}>
        <ArrowLeft size={16} /> Back
      </button>

      <h1>Comic Generation Setup</h1>
      <p style={{ color: '#666' }}>Your storyboard is ready. Choose a layout for the final render.</p>

      {status === 'idle' && (
        <div style={{ marginTop: 40, textAlign: 'left', background: '#fff', border: '1px solid #eee', padding: 30, borderRadius: 12, boxShadow: '0 4px 12px rgba(0,0,0,0.05)' }}>
          <label style={{ display: 'block', marginBottom: 10, fontWeight: 'bold' }}>Select Layout</label>
          <select value={layout} onChange={e => setLayout(e.target.value)} style={{ padding: 12, width: '100%', marginBottom: 30, borderRadius: 6, border: '1px solid #ccc', fontSize: '1em' }}>
            <option value="2x2">2x2 Grid (Instagram Square Format)</option>
            <option value="1x4">1x4 Vertical (Webtoon Strip Format)</option>
          </select>
          <button onClick={startGeneration} style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 10, padding: '14px 20px', width: '100%', background: '#000', color: '#fff', border: 'none', borderRadius: 8, cursor: 'pointer', fontSize: '1.2em' }}>
            <Clapperboard size={20} /> Render Comic
          </button>
        </div>
      )}

      {(status === 'launching' || status === 'processing') && (
        <div style={{ marginTop: 80 }}>
          <Loader2 size={54} className="spin" style={{ color: '#0070f3', marginBottom: 30, animation: 'spin 2s linear infinite' }} />
          <h3>AI is painting your story...</h3>
          <p style={{ color: '#666', lineHeight: 1.6 }}>This usually takes 15 to 30 seconds as we orchestrate the models and synthesize the final image.<br/>Please do not close this page.</p>
        </div>
      )}

      {status === 'failed' && (
         <div style={{ marginTop: 80, color: 'red' }}>
            <h3>Render Error</h3>
            <button onClick={() => setStatus('idle')} style={{ padding: '10px 20px', background: '#000', color: 'white', borderRadius: 6, border: 'none', cursor: 'pointer' }}>Try Again</button>
         </div>
      )}

      <style>{`
        @keyframes spin { 100% { transform: rotate(360deg); } }
      `}</style>
    </div>
  );
};
