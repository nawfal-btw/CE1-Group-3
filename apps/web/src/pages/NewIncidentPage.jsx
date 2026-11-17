import { useEffect, useState } from 'react'
import { createIncident, getIncidentOptions } from '../api/incidents.js'

// Every input starts empty. All values are kept as text while the user types;
// they are converted to numbers only when the request is built.
const EMPTY_FORM = {
  type: '',
  priority: '',
  description: '',
  address: '',
  locationDescription: '',
  latitude: '',
  longitude: '',
}

export default function NewIncidentPage({ onCreated, onCancel }) {
  const [form, setForm] = useState(EMPTY_FORM)
  const [options, setOptions] = useState(null) // the lists for the dropdowns, from the backend
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  // Load the dropdown options once, when the page opens.
  useEffect(() => {
    getIncidentOptions()
      .then(setOptions)
      .catch((err) => setError(`Could not load form options: ${err.message}`))
  }, [])

  // One handler for every input: the input's "name" says which field of the form to update.
  function handleChange(event) {
    const { name, value } = event.target
    setForm((current) => ({ ...current, [name]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault() // stop the browser from reloading the page
    setError('')
    setSubmitting(true)
    try {
      const incident = await createIncident(toRequestBody(form))
      onCreated(incident)
    } catch (err) {
      // Keep what the user typed and show the backend's reason (e.g. a 400 from the domain).
      setError(err.message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="page">
      <h1>New incident</h1>

      {error && (
        <div className="notification error" role="alert">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="form">
        <label>
          Type
          <select name="type" value={form.type} onChange={handleChange} required>
            <option value="">Select a type…</option>
            {options?.types.map((entry) => (
              <option key={entry.code} value={entry.code}>
                {entry.label}
              </option>
            ))}
          </select>
        </label>

        <label>
          Priority
          <select name="priority" value={form.priority} onChange={handleChange} required>
            <option value="">Select a priority…</option>
            {options?.priorities.map((entry) => (
              <option key={entry.code} value={entry.code}>
                {entry.label}
              </option>
            ))}
          </select>
        </label>

        <label>
          Description
          <textarea name="description" value={form.description} onChange={handleChange} rows={3} />
        </label>

        <fieldset>
          <legend>Location</legend>
          <label>
            Address
            <input name="address" value={form.address} onChange={handleChange} />
          </label>
          <label>
            Location description
            <input name="locationDescription" value={form.locationDescription} onChange={handleChange} />
          </label>
          <div className="row">
            <label>
              Latitude
              <input name="latitude" type="number" step="any" value={form.latitude} onChange={handleChange} />
            </label>
            <label>
              Longitude
              <input name="longitude" type="number" step="any" value={form.longitude} onChange={handleChange} />
            </label>
          </div>
        </fieldset>

        <div className="actions">
          <button type="button" onClick={onCancel} disabled={submitting}>
            Cancel
          </button>
          <button type="submit" className="primary" disabled={submitting || !options}>
            {submitting ? 'Creating…' : 'Create'}
          </button>
        </div>
      </form>
    </main>
  )
}

// Builds the JSON the API expects. Empty inputs are left out, so the backend treats them
// as "not set" instead of receiving empty strings.
function toRequestBody(form) {
  const body = {
    type: Number(form.type),
    priority: Number(form.priority),
    location: {},
  }
  if (form.description.trim()) body.description = form.description.trim()
  if (form.address.trim()) body.location.address = form.address.trim()
  if (form.locationDescription.trim()) body.location.locationDescription = form.locationDescription.trim()
  if (form.latitude !== '') body.location.latitude = Number(form.latitude)
  if (form.longitude !== '') body.location.longitude = Number(form.longitude)
  return body
}
