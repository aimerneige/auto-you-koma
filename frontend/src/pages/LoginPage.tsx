import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../api/auth';
import { useAuthStore } from '../store/authStore';
import { LogIn } from 'lucide-react';

export const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [passcode, setPasscode] = useState('');
  const [error, setError] = useState('');
  const [needs2FA, setNeeds2FA] = useState(false);
  const setAuth = useAuthStore((state) => state.setAuth);
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setError('');
      const res = await authApi.login({ email, password, passcode: passcode || undefined });
      const { token } = res.data;
      // Fetch user info with token
      // Temp hack: interceptor uses store state but here it's not set yet. 
      // Better to manually pass token for this one call, or set store first.
      useAuthStore.getState().setAuth(token, '', ''); // Set temporarily
      const meRes = await authApi.getMe();
      setAuth(token, meRes.data.user_id, meRes.data.email);
      navigate('/dashboard');
    } catch (err: any) {
      if (err.response?.data?.error === '2fa_required') {
        setNeeds2FA(true);
      } else {
        setError(err.response?.data?.error || 'Login failed');
      }
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: '100px auto', padding: 20, border: '1px solid #ccc', borderRadius: 8 }}>
      <h2 style={{ display: 'flex', alignItems: 'center', gap: 8 }}><LogIn /> Login</h2>
      {error && <div style={{ color: 'red', marginBottom: 10 }}>{error}</div>}
      <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
        <input type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} required style={{ padding: 8 }} />
        <input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} required style={{ padding: 8 }} />
        {needs2FA && (
          <input type="text" placeholder="2FA Passcode" value={passcode} onChange={e => setPasscode(e.target.value)} required style={{ padding: 8 }} />
        )}
        <button type="submit" style={{ padding: 10, cursor: 'pointer' }}>{needs2FA ? 'Verify & Login' : 'Login'}</button>
      </form>
      <div style={{ marginTop: 15 }}>
        Don't have an account? <Link to="/register">Register</Link>
      </div>
    </div>
  );
};
