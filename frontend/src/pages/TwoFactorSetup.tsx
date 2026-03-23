import { useState } from 'react';
import { authApi } from '../api/auth';
import { ShieldCheck, ShieldAlert } from 'lucide-react';
import { Link } from 'react-router-dom';

export const TwoFactorSetup = () => {
  const [url, setUrl] = useState('');
  const [passcode, setPasscode] = useState('');
  const [status, setStatus] = useState<'' | 'success' | 'error'>('');
  const [errorMsg, setErrorMsg] = useState('');

  const fetchQRCode = async () => {
    try {
      const res = await authApi.setup2FA();
      setUrl(res.data.url);
      setStatus('');
    } catch (err) {
      console.error(err);
      setErrorMsg('Failed to initialize 2FA');
    }
  };

  const handleVerify = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await authApi.verify2FA({ passcode });
      setStatus('success');
    } catch (err: any) {
      setStatus('error');
      setErrorMsg(err.response?.data?.error || 'Invalid passcode');
    }
  };

  return (
    <div style={{ maxWidth: 500, margin: '50px auto', padding: 20, border: '1px solid #ccc', borderRadius: 8 }}>
      <h2>Two-Factor Authentication</h2>
      
      {status === 'success' ? (
        <div style={{ textAlign: 'center', padding: 20 }}>
          <ShieldCheck size={48} color="green" />
          <h3 style={{ color: 'green' }}>2FA Successfully Enabled!</h3>
          <Link to="/dashboard">Return to Dashboard</Link>
        </div>
      ) : (
        <>
          <p>Enhance your account security by enabling TOTP-based 2FA.</p>
          {!url ? (
            <button onClick={fetchQRCode} style={{ padding: '8px 16px' }}>Generate Secret</button>
          ) : (
            <div>
              <div style={{ padding: 15, background: '#f5f5f5', wordBreak: 'break-all', marginBottom: 15 }}>
                <code>{url}</code>
                <div style={{ fontSize: '0.8em', color: '#666', marginTop: 5 }}>(Copy this URL into your Authenticator app, real QR code logic not rendered here to avoid bulky deps)</div>
              </div>
              
              {status === 'error' && <div style={{ color: 'red', marginBottom: 10 }}><ShieldAlert size={16} /> {errorMsg}</div>}
              
              <form onSubmit={handleVerify} style={{ display: 'flex', gap: 10 }}>
                <input type="text" placeholder="Enter 6-digit code" value={passcode} onChange={e => setPasscode(e.target.value)} required style={{ padding: 8, flex: 1 }} />
                <button type="submit" style={{ padding: 8 }}>Verify & Enable</button>
              </form>
            </div>
          )}
        </>
      )}
    </div>
  );
};
