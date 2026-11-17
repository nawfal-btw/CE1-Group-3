// Dispatcher dashboard. For this use case the only action is creating a new incident.
export default function DashboardPage({ notification, onDismissNotification, onNewIncident }) {
  return (
    <main className="page">
      <h1>Dispatcher dashboard</h1>

      {notification && (
        <div className="notification success" role="status">
          <span>{notification}</span>
          <button type="button" className="link" onClick={onDismissNotification}>
            Dismiss
          </button>
        </div>
      )}

      <button type="button" className="primary" onClick={onNewIncident}>
        Create new incident
      </button>
    </main>
  )
}
