# Problem Description
___
Different emergency services rely on different IT infrastructre each with a variaty of potentially proprietary or bespoke software. This offloads the complexity of coordination between departments to the dispatcher, introducing a human-link in the communication chain. A link which in some ways are more fallible to certain errors than computerized systems. This potential for error is exacerbated when the dispatcher acts both as a broker and coordinator.

# Solution Proposal
___
Brokering communication between emergency services is not a task that capitalizes on the dispatchers expertize and uniquely human ability to navigate complex situations.

Our product provides a unified communication interface using [insert chosen standard here]. This enables us to act as a centralized broker, transforming and distributing information across heterogenous systems, whilst providing each user with a graphical user interface to send and receive the exact information they need. This lessens the burden on the dispatcher and offloads the auxillary details of the situation to our system. Alleviating her of the menial responsibility of brokering and book-keeping, while having said information readily available at the press of a button.

# Requirements
___
See [[Requirements]] for how to formulate idiomatic requirements
## Functional 
- The system shall enable communication between emergency services with our platform as the centralized middleman.
- The system shall enable real-time tracking of mobile emergency resources during the emergency, with a latency < X.
- The system shall provide each user a role-specific dashboard
	- the role-specific dashboard shall contain all available and user-relevant information
		- the dashboard shall show GPS-view of their own position and the incident location
		- the dispatcher dashboard shall show GPS-view of all active mobile emegency resources
	- the role-specific dashboard shall be capable of transmitting information back to the system.
	- the role-specific utilities shall only be accessible by their respective roles
- The system shall provide the dispatcher with department internals such as: available vehicles, staff, etc.

## Non-Functional
### Performance
> Notes: performance is not mission-critical given human-involvement. Correctness is the metric of interest
- The business logic shall predefine permissible state-transitions effectively eliminating the possibility of certain logical bugs.
- system latency shall remain stable (not small) under variable loads (within reason: how much is "variable" - elaborate)
### Interface
> Notes: Actual CAD integration is out-of-scope, but most (?) follow ISO/other standards
- the system shall seemlessly integrate with CADs that follow standard [insert standard here]
- the system shall provide a flexible compatibility interface to integrate with CADs not following standard [insert same standard].
- the GUI shall function on chromium and gecko (firefox) based webbrowsers
> Maybe electron instead?
### Quality characteristic
- Inbound and outbound communication from server/clients shall be encrypted
- Communication shall happen over private tailnet
- tailnet shall only be reachable from from verified clients
> Additional security layer - Possible since every client is known prior

### Constraints
- The system shall be vendor-agnostic with the possibility to host it on-prem servers or private cloud providers
> Government entities requires certain levels of data-protection e.g. ownership of infra.
### Environmental
Dependent on the implementation of the client. 

# Architecture/TechStack
___ 
## Infrastructure
Infrastructure pertains to the distribution and deployment of the service that has no functional implication on the service. I.e. everything that goes unnoticed by the end-user.
#### Deployment
Deployment pertains to how the service is provided and continuously revised.
- CI/CD: GitHub Actions,
- Declarative builds: Ansible, Terraform
#### Distribution
Distribution pertains to how the service is executed and managed on the bare-metal
- Orchestration of services: Kubernetes
- Container images: Docker
- Virtualization/Isolation: VMware, KVM
- Managing bare-metal: OpenStack?

## Application
Application pertains to the layer that convert client inbound communication into information that the business logic understands and can execute upon, and converts and transmits resulting computations as outbound communication.
- API: Go
- Logging: Postman?
- Testing: []

## Business Logic
Business logic pertains to the functionality of the application also called the canonical model. Contains the business rules e.g. what state transitions are valid for a vehicle (Parked, En Route, Arrived, Returned); what units are operational and ready for dispatch, etc.

## Data
Data pertains to persistant and ephemeral storage of data and continous data-transmission in case of live-tracking
- Persistant storage: Relational DB - Postgresql
- Cold storage: Previous incidents - []
- Live-tracking: WebSocket
- Geospatial data: PostGIS

## Graphical User Interface
The user-facing interface used to interoperate with the API and WebSocket. Provides the user with the required information and communication with dispatch.
- React / Electron

# Implementation (in progress)






# Milestones
## Iteration 0 - 2026.09.30
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
## Iteration 1 - 2026.10.16
> Implement canonical model and adapt API accordingly, deployment should be declarative such that a single (or few) commands initializes everything
### Scope
Infrastructure (Kostas, Anakin)
- Host machine has kubernetes control-panel
	- 2 VMs identically configured with Ansible with kubernetes client
	- Make them talk
Application (Rasmus, Anakin)
- Adapt API according to Business Logic
- Container image for api (either scratch image or lightweight filesystem e.g. alpine)
- Authentication: TLS
Business Logic (Paula, Rasmus)
- 
Data (Paula, Kostas)
- DB should be single-node per EMS network, i.e. different db per locally conected EMS network
UI (LLM + Rasmus)
- TBD
Report: 


# Notes
k6: testing w/ human behavior
p95, p99 latency, 0% error
docker scratch images minimizes attack-surface
supply-chain attack from dependencies

