import { useAuthStore } from '../store/authStore';
import { useNavigate, Link } from 'react-router-dom';
import { LogOut, Shield } from 'lucide-react';

export const DashboardPage = () => {
  const { email, clearAuth } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    clearAuth();
    navigate('/login');
  };

  return (
    <div style={{ padding: 40, maxWidth: 800, margin: '0 auto' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #eee', paddingBottom: 20, marginBottom: 40 }}>
        <h2>Auto Yon Koma Dashboard</h2>
        <div style={{ display: 'flex', gap: 15, alignItems: 'center' }}>
          <span>Welcome, {email}</span>
          <button onClick={handleLogout} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '5px 10px' }}><LogOut size={16} /> Logout</button>
        </div>
      </header>

      <section style={{ display: 'flex', gap: 20 }}>
        <div style={{ border: '1px solid #ddd', padding: 20, borderRadius: 8, flex: 1 }}>
          <h3>Security Settings</h3>
          <p>Protect your account with Two-Factor Authentication.</p>
          <Link to="/2fa-setup" style={{ display: 'inline-flex', alignItems: 'center', gap: 5, marginTop: 10, padding: '8px 12px', background: '#0070f3', color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            <Shield size={16} /> Setup / Manage 2FA
          </Link>
        </div>

        <div style={{ border: '1px solid #ddd', padding: 20, borderRadius: 8, flex: 1 }}>
          <h3>Character Library</h3>
          <p>Manage your AI generated character assets and lore.</p>
          <Link to="/characters" style={{ display: 'inline-flex', alignItems: 'center', gap: 5, marginTop: 10, padding: '8px 12px', background: '#222', color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            Go to Library
          </Link>
        </div>

        <div style={{ border: '1px solid #ddd', padding: 20, borderRadius: 8, flex: 1 }}>
          <h3>Script Studio</h3>
          <p>Write and parse stories into storyboard panels.</p>
          <Link to="/scripts/new" style={{ display: 'inline-flex', alignItems: 'center', gap: 5, marginTop: 10, padding: '8px 12px', background: '#10b981', color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            New Script Hub
          </Link>
        </div>

        <div style={{ border: '1px solid #ddd', padding: 20, borderRadius: 8, flex: 1 }}>
          <h3>Series Universe</h3>
          <p>Maintain character states and continuity lore over time.</p>
          <Link to="/series" style={{ display: 'inline-flex', alignItems: 'center', gap: 5, marginTop: 10, padding: '8px 12px', background: '#8b5cf6', color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            Series Dashboard
          </Link>
        </div>
      </section>
    </div>
  );
};
