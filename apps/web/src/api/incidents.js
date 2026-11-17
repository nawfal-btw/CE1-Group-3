// All communication with the incident API lives here, so pages never call fetch() directly.
// Paths start with /v1: Vite forwards them to the Go API (see vite.config.js).

// GET /v1/incident/options -> { version, types: [{code, label}], priorities: [{code, label}] }
export async function getIncidentOptions() {
  const response = await fetch('/v1/incident/options')
  if (!response.ok) {
    throw new Error(await errorMessage(response))
  }
  return response.json()
}

// POST /v1/incident -> the created incident (with its new id) on 201 Created.
export async function createIncident(incident) {
  const response = await fetch('/v1/incident', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(incident),
  })
  if (!response.ok) {
    throw new Error(await errorMessage(response))
  }
  return response.json()
}

// The API answers errors with plain text, e.g. "invalid incident\nunknown incident type 999".
async function errorMessage(response) {
  const text = (await response.text()).trim()
  return text || `Request failed with status ${response.status}`
}
