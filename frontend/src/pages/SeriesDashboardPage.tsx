import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { seriesApi } from '../api/series';
import { Book, Plus, ArrowLeft } from 'lucide-react';

export const SeriesDashboardPage = () => {
  const navigate = useNavigate();
  const [seriesList, setSeriesList] = useState<any[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');

  const fetchSeries = async () => {
    try {
      const res = await seriesApi.list();
      setSeriesList(res.data);
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    fetchSeries();
  }, []);

  const handleCreate = async () => {
    if (!title) return;
    try {
      await seriesApi.create({ title, description: desc });
      setShowModal(false);
      setTitle('');
      setDesc('');
      fetchSeries();
    } catch(e) {
      alert("Error creating series");
    }
  };

  return (
    <div style={{ padding: 40, background: '#f5f5f5', minHeight: '100vh' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 40 }}>
        <div style={{ display: 'flex', gap: 15, alignItems: 'center' }}>
          <button onClick={() => navigate('/dashboard')} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: '#fff', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer' }}>
            <ArrowLeft size={16} /> Back
          </button>
          <h1 style={{ margin: 0, display: 'flex', alignItems: 'center', gap: 10 }}>
            <Book size={28} /> Series Universe
          </h1>
        </div>
        <button onClick={() => setShowModal(true)} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '10px 20px', background: '#0070f3', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontWeight: 'bold' }}>
          <Plus size={18} /> New Series
        </button>
      </header>

      {seriesList.length === 0 ? (
        <div style={{ textAlign: 'center', color: '#888', marginTop: 100 }}>
          <Book size={64} style={{ opacity: 0.2, marginBottom: 20 }} />
          <h2>No Series created yet</h2>
          <p>Group your comic episodes logically to share memory and context.</p>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: 20 }}>
          {seriesList.map(s => (
            <div key={s.id} onClick={() => alert("Continuity Tracker UI coming later Phase 2")} style={{ background: '#fff', padding: 25, borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.05)', cursor: 'pointer', border: '1px solid transparent', transition: 'border 0.2s' }}>
              <h3 style={{ margin: '0 0 10px 0', color: '#0070f3' }}>{s.title}</h3>
              <p style={{ color: '#666', fontSize: '0.9em', margin: 0 }}>{s.description || 'No description'}</p>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ background: '#fff', padding: 30, borderRadius: 12, width: 400 }}>
            <h3>Create a new Series Universe</h3>
            <input 
              placeholder="Series Title" 
              value={title} 
              onChange={e => setTitle(e.target.value)} 
              style={{ width: '100%', padding: 10, margin: '15px 0', boxSizing: 'border-box', border: '1px solid #ddd', borderRadius: 4 }}
            />
            <textarea 
              placeholder="Synopsis / Canon guidelines" 
              value={desc} 
              onChange={e => setDesc(e.target.value)} 
              style={{ width: '100%', padding: 10, height: 100, boxSizing: 'border-box', border: '1px solid #ddd', borderRadius: 4, resize: 'none' }}
            />
            <div style={{ display: 'flex', gap: 10, marginTop: 20, justifyContent: 'flex-end' }}>
              <button onClick={() => setShowModal(false)} style={{ padding: '8px 16px', border: '1px solid #ccc', background: '#fff', borderRadius: 4, cursor: 'pointer' }}>Cancel</button>
              <button onClick={handleCreate} style={{ padding: '8px 16px', background: '#0070f3', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}>Create</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
