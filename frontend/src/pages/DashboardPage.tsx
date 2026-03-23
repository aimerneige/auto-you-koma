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
      </section>
    </div>
  );
};
