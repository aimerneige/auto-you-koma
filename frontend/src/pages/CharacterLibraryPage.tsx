import { useEffect, useState } from 'react';
import { characterApi } from '../api/characters';
import { Link } from 'react-router-dom';
import { Users, Plus } from 'lucide-react';

export const CharacterLibraryPage = () => {
  const [characters, setCharacters] = useState<any[]>([]);
  const [q, setQ] = useState('');

  const loadDocs = async (query = '') => {
    try {
      const res = await characterApi.list(query ? { q: query } : {});
      setCharacters(res.data || []);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    loadDocs();
  }, []);

  return (
    <div style={{ padding: 40, maxWidth: 1200, margin: '0 auto' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 20 }}>
        <h2 style={{ display: 'flex', alignItems: 'center', gap: 10 }}><Users /> Character Library</h2>
        <Link to="/characters/new" style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', background: '#0070f3', color: '#fff', textDecoration: 'none', borderRadius: 4 }}>
          <Plus size={16} /> New Character
        </Link>
      </header>
      
      <div style={{ marginBottom: 20 }}>
        <input 
          type="text" 
          placeholder="Search characters by name or tags..." 
          value={q} 
          onChange={e => setQ(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && loadDocs(q)}
          style={{ width: '100%', maxWidth: 400, padding: 10, borderRadius: 4, border: '1px solid #ccc' }}
        />
        <button onClick={() => loadDocs(q)} style={{ marginLeft: 10, padding: 10, cursor: 'pointer' }}>Search</button>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: 20 }}>
        {characters.map(char => {
          const primaryImage = char.images?.find((i: any) => i.is_primary) || char.images?.[0];
          return (
            <Link key={char.id} to={`/characters/${char.id}`} style={{ border: '1px solid #eaeaea', padding: 15, borderRadius: 8, textDecoration: 'none', color: 'inherit', display: 'block', transition: 'box-shadow 0.2s', boxShadow: '0 2px 4px rgba(0,0,0,0.05)' }}>
              {primaryImage ? (
                <img src={`/api/v1${primaryImage.file_path}`} alt={char.name} style={{ width: '100%', height: 200, objectFit: 'cover', borderRadius: 4, background: '#f5f5f5' }} />
              ) : (
                <div style={{ width: '100%', height: 200, background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center', borderRadius: 4, color: '#aaa' }}>No Image</div>
              )}
              <h3 style={{ margin: '15px 0 5px' }}>{char.name}</h3>
              <p style={{ margin: 0, color: '#666', fontSize: '0.9em' }}>{char.category || 'Uncategorized'}</p>
            </Link>
          );
        })}
        {characters.length === 0 && (
          <div style={{ gridColumn: '1 / -1', padding: 40, textAlign: 'center', color: '#888' }}>
            No characters found. Create one!
          </div>
        )}
      </div>
    </div>
  );
};
