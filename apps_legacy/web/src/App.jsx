import { useState } from 'react'
import { createIncident, getAllIncidents, getIncidentById } from './api/incidents'
import './App.css'

const DEFAULT_API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api'

function App() {
  const [apiBaseUrl, setApiBaseUrl] = useState(DEFAULT_API_BASE_URL)
  const [incidentId, setIncidentId] = useState('1')
  const [newIncident, setNewIncident] = useState({ latitude: '', longitude: '' })
  const [response, setResponse] = useState(null)
  const [error, setError] = useState('')

  async function runRequest(requestFn) {
    setError('')
    setResponse(null)

    try {
      const data = await requestFn()
      setResponse(data)
    } catch (requestError) {
      setError(requestError.message)
    }
  }

  function handleCreateIncident(event) {
    event.preventDefault()

    runRequest(() =>
      createIncident(apiBaseUrl, {
        latitude: newIncident.latitude,
        longitude: newIncident.longitude,
      }),
    )
  }

  return (
    <main>
      <h1>Incident API Prototype</h1>

      <label>
        API Base URL
        <input
          value={apiBaseUrl}
          onChange={(event) => setApiBaseUrl(event.target.value)}
          placeholder="/api"
        />
      </label>

      <section>
        <h2>GET /incident</h2>
        <button type="button" onClick={() => runRequest(() => getAllIncidents(apiBaseUrl))}>
          Fetch all incidents
        </button>
      </section>

      <section>
        <h2>GET /incident/{'{id}'}</h2>
        <div className="row">
          <input
            value={incidentId}
            onChange={(event) => setIncidentId(event.target.value)}
            placeholder="Incident ID"
          />
          <button
            type="button"
            onClick={() => runRequest(() => getIncidentById(apiBaseUrl, incidentId))}
          >
            Fetch incident by ID
          </button>
        </div>
      </section>

      <section>
        <h2>POST /incident</h2>
        <form className="stack" onSubmit={handleCreateIncident}>
          <input
            value={newIncident.latitude}
            onChange={(event) =>
              setNewIncident((current) => ({ ...current, latitude: event.target.value }))
            }
            placeholder="Latitude"
            required
          />
          <input
            value={newIncident.longitude}
            onChange={(event) =>
              setNewIncident((current) => ({ ...current, longitude: event.target.value }))
            }
            placeholder="Longitude"
            required
          />
          <button type="submit">Create incident</button>
        </form>
      </section>

      <section>
        <h2>Response</h2>
        {error ? <pre className="error">{error}</pre> : <pre>{JSON.stringify(response, null, 2)}</pre>}
      </section>
    </main>
  )
}

export default App
