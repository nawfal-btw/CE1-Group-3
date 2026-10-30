const DEFAULT_HEADERS = {
  'Content-Type': 'application/json',
}

async function requestJson(path, options = {}) {
  const response = await fetch(path, options)
  const responseText = await response.text()

  if (!response.ok) {
    throw new Error(responseText || `Request failed with status ${response.status}`)
  }

  if (!responseText) {
    return null
  }

  try {
    return JSON.parse(responseText)
  } catch {
    return responseText
  }
}

function buildEndpoint(baseUrl, endpoint) {
  const normalizedBase = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl
  return `${normalizedBase}${endpoint}`
}

export function getAllIncidents(baseUrl) {
  return requestJson(buildEndpoint(baseUrl, '/incident'))
}

export function getIncidentById(baseUrl, id) {
  return requestJson(buildEndpoint(baseUrl, `/incident/${id}`))
}

export function createIncident(baseUrl, payload) {
  return requestJson(buildEndpoint(baseUrl, '/incident'), {
    method: 'POST',
    headers: DEFAULT_HEADERS,
    body: JSON.stringify(payload),
  })
}
