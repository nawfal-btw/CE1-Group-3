# Problem Description
___
Different emergency services rely on different IT infrastructre each with a variaty of potentially proprietary or bespoke software. This offloads the complexity of coordination between departments to the dispatcher, introducing a human-link in the communication chain. A link which in some ways are more fallible to certain errors than computerized systems. This potential for error is exacerbated when the dispatcher acts both as a broker and coordinator.

# Solution Proposal
___
Brokering communication between emergency services is not a task that capitalizes on the dispatchers expertize and uniquely human ability to navigate complex situations.

Our product provides a unified communication interface using [insert chosen standard here]. This enables us to act as a centralized broker, transforming and distributing information across heterogenous systems, whilst providing each user with a graphical user interface to send and receive the exact information they need. This lessens the burden on the dispatcher and offloads the auxillary details of the situation to our system. Alleviating her of the menial responsibility of brokering and book-keeping, while having said information readily available at the press of a button.


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
	See [[Iteration 0]]
## Iteration 1 - 2026.10.16
> Implement canonical model and adapt API accordingly, deployment should be declarative such that a single (or few) commands initializes everything
	See [[Iteration 1]]

# Notes
k6: testing w/ human behavior
p95, p99 latency, 0% error
docker scratch images minimizes attack-surface
supply-chain attack from dependencies

Ask regarding where to host.
CLAAUDIA
Decide on which platofrm
Sokol avialability of compute

