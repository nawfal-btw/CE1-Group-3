# Scope
- Infra
	- Debian VM
	- Single node K8s
	- 3 K8s pods  (w/ DB and API, UI)
	- reachable only from localhost. 
- Application
	- Simple REST API with:
		- POST /incidents
		- GET /incidents/{id}
			```json
			{
					"location": {
						"latitude": float32,
						"longitude": float32
					},
				}
			```
		- Validate req
		- create incident and store in postgres DB
		
- Data
	- Single Postgres table for incidents
		- id
		- latitude
		- longitude
		- created_at
- UI
	- Incident creation form
	- location input (in coordinates)
	- list/display incidents with ID, location, created at
	- Optionally map (MapLibre)
	- retrieve incident and show in UI