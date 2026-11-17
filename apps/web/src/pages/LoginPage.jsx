// Login screen. For now the button only moves to the dashboard.
// TODO(auth team): send credentials to the backend and continue only if they are accepted.
export default function LoginPage({ onLogin }) {
  return (
    <main className="page narrow">
      <h1>Incident Dispatch</h1>
      <p className="muted">Dispatcher login</p>
      <button type="button" className="primary" onClick={onLogin}>
        Log in
      </button>
    </main>
  )
}
