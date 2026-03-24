import { useState } from 'react';
import { useCharacterStore } from '../store/characterStore';
import { useNavigate } from 'react-router-dom';

export function CharacterLibraryPage() {
  const { characters, loading, error, fetchCharacters, deleteCharacter } = useCharacterStore();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');

  // Fetch characters on mount
  useState(() => {
    fetchCharacters();
  });

  const filteredCharacters = characters.filter((char) =>
    char.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    char.name_jp?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="page-container">
      <div className="page-header">
        <h1 className="title-comic">Character Library</h1>
        <button className="btn" onClick={() => navigate('/characters/new')}>
          + New Character
        </button>
      </div>

      <div className="search-bar">
        <input
          type="text"
          placeholder="Search characters..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="search-input"
        />
      </div>

      {loading && <div className="loading">Loading...</div>}
      {error && <div className="error-message">{error}</div>}

      <div className="character-grid">
        {filteredCharacters.map((character) => (
          <div
            key={character.id}
            className="character-card"
            onClick={() => navigate(`/characters/${character.id}`)}
          >
            <div className="character-avatar">
              {character.reference_sheet_url ? (
                <img src={character.reference_sheet_url} alt={character.name} />
              ) : (
                <span className="avatar-placeholder">{character.name[0]}</span>
              )}
            </div>
            <div className="character-info">
              <h3>{character.name}</h3>
              {character.name_jp && <p className="jp-name">{character.name_jp}</p>}
              {character.category && <span className="category-tag">{character.category}</span>}
            </div>
            <button
              className="delete-btn"
              onClick={(e) => {
                e.stopPropagation();
                if (confirm('Delete this character?')) {
                  deleteCharacter(character.id);
                }
              }}
            >
              ×
            </button>
          </div>
        ))}
      </div>

      {filteredCharacters.length === 0 && !loading && (
        <div className="empty-state">
          <p>No characters yet. Create your first character!</p>
        </div>
      )}
    </div>
  );
}