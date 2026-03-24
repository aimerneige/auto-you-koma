import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { characterApi } from '../api/characters';

export const CharacterForm = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    name: '',
    name_jp: '',
    gender: '',
    age: '',
    category: '',
    visual_prompt: '',
    backstory: ''
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await characterApi.create(form);
      navigate(`/characters/${res.data.id}`);
    } catch (err) {
      alert('Failed to create character');
    }
  };

  return (
    <div style={{ padding: 40, maxWidth: 600, margin: '0 auto' }}>
      <h2>Create New Character</h2>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 15, marginTop: 20 }}>
        <div><label>Name</label><input type="text" value={form.name} onChange={e => setForm({...form, name: e.target.value})} required style={{ width: '100%', padding: 8 }} /></div>
        <div><label>Japanese Name</label><input type="text" value={form.name_jp} onChange={e => setForm({...form, name_jp: e.target.value})} style={{ width: '100%', padding: 8 }} /></div>
        <div style={{ display: 'flex', gap: 10 }}>
          <div style={{ flex: 1 }}><label>Gender</label><input type="text" value={form.gender} onChange={e => setForm({...form, gender: e.target.value})} style={{ width: '100%', padding: 8 }} /></div>
          <div style={{ flex: 1 }}><label>Age</label><input type="text" value={form.age} onChange={e => setForm({...form, age: e.target.value})} style={{ width: '100%', padding: 8 }} /></div>
        </div>
        <div><label>Category (e.g. Original, LoveLive)</label><input type="text" value={form.category} onChange={e => setForm({...form, category: e.target.value})} style={{ width: '100%', padding: 8 }} /></div>
        <div><label>Visual Prompt (SD / AI Generation tags)</label><textarea value={form.visual_prompt} onChange={e => setForm({...form, visual_prompt: e.target.value})} style={{ width: '100%', padding: 8, height: 80 }} /></div>
        <div><label>Backstory</label><textarea value={form.backstory} onChange={e => setForm({...form, backstory: e.target.value})} style={{ width: '100%', padding: 8, height: 100 }} /></div>
        <button type="submit" style={{ padding: 12, background: '#000', color: '#fff', border: 'none', cursor: 'pointer', borderRadius: 4 }}>Save Character</button>
      </form>
    </div>
  );
};
