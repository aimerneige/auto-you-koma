import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { characterApi } from '../api/characters';

export const CharacterDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const [char, setChar] = useState<any>(null);
  const [uploading, setUploading] = useState(false);

  const loadChar = async () => {
    try {
      if (!id) return;
      const res = await characterApi.get(id);
      setChar(res.data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    loadChar();
  }, [id]);

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || !e.target.files[0] || !id) return;
    
    const file = e.target.files[0];
    const formData = new FormData();
    formData.append('image', file);
    formData.append('image_type', 'avatar');
    // Set first image as primary automatically for simplicity
    formData.append('is_primary', char?.images?.length === 0 ? 'true' : 'false');

    setUploading(true);
    try {
      await characterApi.uploadImage(id, formData);
      await loadChar(); // reload
    } catch (err) {
      alert("Failed to upload image");
    } finally {
      setUploading(false);
      e.target.value = ''; // reset input
    }
  };

  if (!char) return <div style={{ padding: 40 }}>Loading...</div>;

  return (
    <div style={{ padding: 40, maxWidth: 1000, margin: '0 auto', display: 'flex', gap: 40 }}>
      {/* Left Sidebar: Images */}
      <div style={{ width: 300 }}>
        <div style={{ background: '#f5f5f5', height: 400, borderRadius: 8, overflow: 'hidden', marginBottom: 20 }}>
           {char.images?.find((i:any) => i.is_primary) ? (
             <img src={`/api/v1${char.images.find((i:any) => i.is_primary).file_path}`} style={{ width:'100%', height:'100%', objectFit: 'cover' }} />
           ) : (
             <div style={{ display:'flex', alignItems:'center', justifyContent:'center', height:'100%', color:'#aaa' }}>No Image</div>
           )}
        </div>
        
        <div>
          <h4>Upload Image</h4>
          <input type="file" accept="image/*" onChange={handleImageUpload} disabled={uploading} />
          {uploading && <span style={{ fontSize: '0.8em', color: '#666' }}>Uploading...</span>}
        </div>

        <div style={{ marginTop: 20, display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10 }}>
          {char.images?.filter((i:any) => !i.is_primary).map((img: any) => (
            <img key={img.id} src={`/api/v1${img.file_path}`} style={{ width: '100%', aspectRatio: '1/1', objectFit: 'cover', borderRadius: 4 }} />
          ))}
        </div>
      </div>

      {/* Right Content: Details */}
      <div style={{ flex: 1 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline' }}>
          <h1 style={{ margin: 0 }}>{char.name} <span style={{ fontSize: '0.6em', color: '#888' }}>{char.name_jp}</span></h1>
          <Link to="/characters" style={{ textDecoration: 'none', color: '#0070f3' }}>&larr; Back to Library</Link>
        </div>
        <div style={{ display: 'flex', gap: 10, margin: '15px 0' }}>
          <span style={{ padding: '4px 8px', background: '#eee', borderRadius: 4, fontSize: '0.8em' }}>{char.category || 'Uncategorized'}</span>
          {char.gender && <span style={{ padding: '4px 8px', background: '#eee', borderRadius: 4, fontSize: '0.8em' }}>{char.gender}</span>}
          {char.age && <span style={{ padding: '4px 8px', background: '#eee', borderRadius: 4, fontSize: '0.8em' }}>{char.age}</span>}
        </div>
        
        <div style={{ marginTop: 30 }}>
          <h3>Visual Prompt</h3>
          <p style={{ background: '#f9f9f9', padding: 15, borderRadius: 8, fontFamily: 'monospace' }}>{char.visual_prompt || 'N/A'}</p>
        </div>

        <div style={{ marginTop: 30 }}>
          <h3>Backstory</h3>
          <p style={{ lineHeight: 1.6 }}>{char.backstory || 'N/A'}</p>
        </div>
      </div>
    </div>
  );
};
