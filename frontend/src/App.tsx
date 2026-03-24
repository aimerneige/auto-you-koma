import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './store/authStore';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { DashboardPage } from './pages/DashboardPage';
import { TwoFactorSetup } from './pages/TwoFactorSetup';
import { CharacterLibraryPage } from './pages/CharacterLibraryPage';
import { CharacterForm } from './pages/CharacterForm';
import { CharacterDetailPage } from './pages/CharacterDetailPage';
import { ScriptEditorPage } from './pages/ScriptEditorPage';
import { ComicGenerationPage } from './pages/ComicGenerationPage';
import { ComicViewerPage } from './pages/ComicViewerPage';
import { ComicAdvancedViewer } from './pages/ComicAdvancedViewer';
import { SeriesDashboardPage } from './pages/SeriesDashboardPage';

const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const token = useAuthStore((state) => state.token);
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return children;
};

export const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <DashboardPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/2fa-setup"
          element={
            <ProtectedRoute>
              <TwoFactorSetup />
            </ProtectedRoute>
          }
        />
        <Route
          path="/characters"
          element={
            <ProtectedRoute>
              <CharacterLibraryPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/characters/new"
          element={
            <ProtectedRoute>
              <CharacterForm />
            </ProtectedRoute>
          }
        />
        <Route
          path="/characters/:id"
          element={
            <ProtectedRoute>
              <CharacterDetailPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/scripts/:id"
          element={
            <ProtectedRoute>
              <ScriptEditorPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/generate/:scriptId"
          element={
            <ProtectedRoute>
              <ComicGenerationPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/viewer/:genId"
          element={
            <ProtectedRoute>
              <ComicViewerPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/viewer/advanced/:genId"
          element={
            <ProtectedRoute>
              <ComicAdvancedViewer />
            </ProtectedRoute>
          }
        />
        <Route
          path="/series"
          element={
            <ProtectedRoute>
              <SeriesDashboardPage />
            </ProtectedRoute>
          }
        />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
