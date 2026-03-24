import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { generationApi } from '../api/generations';
import { ArrowLeft, Download, Layers } from 'lucide-react';

export const ComicViewerPage = () => {
  const { genId } = useParams<{ genId: string }>();
  const navigate = useNavigate();
  const [gen, setGen] = useState<any>(null);

  useEffect(() => {
    if (genId) {
      generationApi.get(genId).then(res => setGen(res.data));
    }
  }, [genId]);

  if (!gen) return <div style={{ padding: 40, textAlign: 'center' }}>Loading output...</div>;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100vh', background: '#ececec' }}>
      <header style={{ padding: '15px 20px', background: '#fff', borderBottom: '1px solid #ddd', display: 'flex', justifyContent: 'space-between', alignItems: 'center', boxShadow: '0 1px 3px rgba(0,0,0,0.05)' }}>
        <button onClick={() => navigate('/dashboard')} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: 'transparent', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer' }}>
          <ArrowLeft size={16} /> Dashboard
        </button>
        <h3 style={{ margin: 0, display: 'flex', alignItems: 'center', gap: 8 }}><Layers size={20} color="#0070f3" /> Final Artboard</h3>
        <div style={{ display: 'flex', gap: 10 }}>
          <button onClick={() => navigate(`/viewer/advanced/${genId}`)} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', background: '#8b5cf6', color: '#fff', textDecoration: 'none', border: 'none', borderRadius: 4, fontWeight: 'bold', cursor: 'pointer' }}>
            <Layers size={16} /> Freeform Layout Sandbox
          </button>
          <a href={gen.result_image_url} download={`comic_${genId}.jpg`} target="_blank" rel="noreferrer" style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', background: '#10b981', color: '#fff', textDecoration: 'none', borderRadius: 4, fontWeight: 'bold' }}>
            <Download size={16} /> Save Final
          </a>
        </div>
      </header>

      <div style={{ flex: 1, overflowY: 'auto', padding: 40, display: 'flex', justifyContent: 'center', alignItems: 'flex-start' }}>
        <div style={{ background: '#fff', padding: 20, borderRadius: 12, boxShadow: '0 8px 30px rgba(0,0,0,0.1)', display: 'inline-block', maxWidth: '100%' }}>
           <img src={gen.result_image_url} alt="Generated Comic" style={{ maxWidth: '100%', maxHeight: '80vh', objectFit: 'contain', display: 'block', borderRadius: 4 }} />
        </div>
      </div>
    </div>
  );
};
