import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { MainLayout } from './components/Layout/MainLayout';
import { DashboardPage } from './pages/DashboardPage';
import { CharacterLibraryPage } from './pages/CharacterLibraryPage';
import { CharacterFormPage } from './pages/CharacterFormPage';

function App() {
  return (
    <BrowserRouter>
      <MainLayout>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/characters" element={<CharacterLibraryPage />} />
          <Route path="/characters/new" element={<CharacterFormPage />} />
          <Route path="/characters/:id" element={<CharacterFormPage />} />
        </Routes>
      </MainLayout>
    </BrowserRouter>
  );
}

export default App;