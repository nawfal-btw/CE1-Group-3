Build MVP and layer functionality on top
Functional core - Imperative shell

## Functional Core
Holds the current state for a given network of clients. Includes:
- Incidents


## Example Infrastructure
Clients

REST / WebSocket

Integration Layer

Application Layer - Imperative Shell
- Authentication
- Validation
- Transactions (DB)
- Orchestration
- Persistence

Domain command

Canonical Domain - Functional Core
- Business logic
- Invariants
- State and state change
- Decisions

Decision + events

Application Layer

## Work Split
- Infrastructure: (Adjacent: App| Diff:5/5) - Kostas, Rasmus
- Application: (Adjacent: Canon, Infra| Diff:4.5/5) - Salekin, Rasmus
- Data (Integration): (Adjacent: App (Infra) | Diff:2/5 (4/5)) - Paula
- Canonical Domain: (Adjacent: App, UI | Diff:3/5) - Paula, Rasmus
- User Interface: (Adjacent: Canon, App | Diff:4/5) - LLM

## Tech Stack
- Infrastructure: 
	- Kubernetes, containerd, etcd, Gateway (NGINX Ingress)
	- Networking: Load balancing,  DNS, TLS
	- Declarative builds: Ansible (Terraform, Nix)
	- Markup languages: YAML, TOML, Declarative builds (Terraform/Ansible)
	- CI/CD: Git (Actions)
- Data & (Minimal) Integration Layer: 
	- Persistant Storage: 
		- RelationalDB: Postgres
		- GIS: PostGIS
		- ~~Object storage: MongoDB~~
		- Cold storage: idk...
	- Cache storage
		- Redis
	- CAD Integration, Hospital integration, GIS, Parsing
- Application Layer: 
	- Go, Python 
	- REST API, WebSocket
	- JSON, HTTP
	- Networking, Authentication,
- Canonical Domain Model: 
	- Go is viable: decreases complexity (but i think oCaml is cool  )
	- Domain model: Invariants, State machines
	- Business logic: Policies, Decisions, Resource allocation
	- Domain knowledge (EMS, CAD...)
- UI: 
	- Domain knowledge
	- CSS, Tailwind, HTMX, Typescript, Python, React/Svelte etc...
	- WebSocket, REST, JSON

## MileStoneA
> A user can create an Incident containing location and view the incident in the UI
### Scope
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
### Outside scope
- No roles (dispatcher, Fire, EMS etc)
- No business logic i.e. state transitions AMB1 EN_ROUTE -> ARRIVED
- No auth
- No integration


- No roles (dispatcher, Fire, EMS etc)


