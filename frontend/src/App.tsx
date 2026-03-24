import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { MainLayout } from './components/Layout/MainLayout';
import { DashboardPage } from './pages/DashboardPage';
import { CharacterLibraryPage } from './pages/CharacterLibraryPage';
import { CharacterFormPage } from './pages/CharacterFormPage';
import { ProjectListPage } from './pages/ProjectListPage';
import { ProjectCreatePage } from './pages/ProjectCreatePage';
import { ScriptEditorPage } from './pages/ScriptEditorPage';

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
          <Route path="/projects" element={<ProjectListPage />} />
          <Route path="/projects/new" element={<ProjectCreatePage />} />
          <Route path="/projects/:id" element={<ProjectCreatePage />} />
          <Route path="/projects/:id/script" element={<ScriptEditorPage />} />
        </Routes>
      </MainLayout>
    </BrowserRouter>
  );
}

export default App;