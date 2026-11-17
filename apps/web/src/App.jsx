import { useState } from 'react'
import LoginPage from './pages/LoginPage.jsx'
import DashboardPage from './pages/DashboardPage.jsx'
import NewIncidentPage from './pages/NewIncidentPage.jsx'

// App decides which page is visible and holds the notification shown on the dashboard.
// The pages never change the screen themselves: they call the functions App gives them.
export default function App() {
  const [page, setPage] = useState('login') // 'login' | 'dashboard' | 'new-incident'
  const [notification, setNotification] = useState(null)

  function handleIncidentCreated(incident) {
    setNotification(`Incident created (id ${incident.id})`)
    setPage('dashboard')
  }

  if (page === 'login') {
    return <LoginPage onLogin={() => setPage('dashboard')} />
  }

  if (page === 'new-incident') {
    return (
      <NewIncidentPage
        onCreated={handleIncidentCreated}
        onCancel={() => setPage('dashboard')}
      />
    )
  }

  return (
    <DashboardPage
      notification={notification}
      onDismissNotification={() => setNotification(null)}
      onNewIncident={() => {
        setNotification(null)
        setPage('new-incident')
      }}
    />
  )
}
