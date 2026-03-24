import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useCharacterStore } from '../store/characterStore';

export function CharacterFormPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { selectedCharacter, fetchCharacter, createCharacter, updateCharacter, loading } = useCharacterStore();

  const [formData, setFormData] = useState({
    name: '',
    name_jp: '',
    gender: '',
    age: '',
    personality: '',
    backstory: '',
    visual_prompt: '',
    tags: '',
    category: '',
  });

  useEffect(() => {
    if (id) {
      fetchCharacter(id);
    }
  }, [id]);

  useEffect(() => {
    if (selectedCharacter && id) {
      setFormData({
        name: selectedCharacter.name,
        name_jp: selectedCharacter.name_jp || '',
        gender: selectedCharacter.gender || '',
        age: selectedCharacter.age || '',
        personality: selectedCharacter.personality || '',
        backstory: selectedCharacter.backstory || '',
        visual_prompt: selectedCharacter.visual_prompt || '',
        tags: selectedCharacter.tags || '',
        category: selectedCharacter.category || '',
      });
    }
  }, [selectedCharacter, id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (id) {
        await updateCharacter(id, formData as any);
      } else {
        await createCharacter(formData as any);
      }
      navigate('/characters');
    } catch (error) {
      console.error('Failed to save character:', error);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <div className="page-container">
      <h1 className="title-comic">{id ? 'Edit Character' : 'Create Character'}</h1>

      <form onSubmit={handleSubmit} className="character-form">
        <div className="form-section">
          <h3>Basic Info</h3>
          <div className="form-row">
            <div className="input-box">
              <label>Name *</label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
              />
            </div>
            <div className="input-box">
              <label>Name (Japanese)</label>
              <input
                type="text"
                name="name_jp"
                value={formData.name_jp}
                onChange={handleChange}
              />
            </div>
          </div>

          <div className="form-row">
            <div className="input-box">
              <label>Gender</label>
              <select name="gender" value={formData.gender} onChange={handleChange}>
                <option value="">Select...</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
              </select>
            </div>
            <div className="input-box">
              <label>Age</label>
              <input
                type="text"
                name="age"
                value={formData.age}
                onChange={handleChange}
                placeholder="e.g., 16 years old"
              />
            </div>
          </div>

          <div className="input-box">
            <label>Category</label>
            <input
              type="text"
              name="category"
              value={formData.category}
              onChange={handleChange}
              placeholder="e.g., Original, LoveLive, Fanmade"
            />
          </div>

          <div className="input-box">
            <label>Tags (comma separated)</label>
            <input
              type="text"
              name="tags"
              value={formData.tags}
              onChange={handleChange}
              placeholder="e.g., student, magical girl, leader"
            />
          </div>
        </div>

        <div className="form-section">
          <h3>Personality & Background</h3>
          <div className="input-box">
            <label>Personality Keywords</label>
            <textarea
              name="personality"
              value={formData.personality}
              onChange={handleChange}
              rows={3}
              placeholder="Describe the character's personality..."
            />
          </div>

          <div className="input-box">
            <label>Backstory</label>
            <textarea
              name="backstory"
              value={formData.backstory}
              onChange={handleChange}
              rows={5}
              placeholder="The character's background story..."
            />
          </div>
        </div>

        <div className="form-section">
          <h3>Visual Description</h3>
          <div className="input-box">
            <label>Visual Prompt (for AI generation)</label>
            <textarea
              name="visual_prompt"
              value={formData.visual_prompt}
              onChange={handleChange}
              rows={4}
              placeholder="Describe the character's appearance for AI image generation..."
            />
          </div>
        </div>

        <div className="form-actions">
          <button type="button" className="btn btn-alt" onClick={() => navigate('/characters')}>
            Cancel
          </button>
          <button type="submit" className="btn" disabled={loading}>
            {loading ? 'Saving...' : 'Save Character'}
          </button>
        </div>
      </form>
    </div>
  );
}