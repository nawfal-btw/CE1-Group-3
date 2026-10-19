# Scope
Infrastructure (Kostas, Anakin)
- Host machine has kubernetes control-panel
	- 2 VMs identically configured with Ansible with kubernetes client
	- Make them talk
Application (Rasmus, Anakin)
- Adapt API according to Business Logic
- Container image for api (either scratch image or lightweight filesystem e.g. alpine, distroless)
- Authentication: TLS
business logic (paula, rasmus)
- 
Data (Paula, Kostas)
- DB should be single-node per EMS network, i.e. different db per locally conected EMS network
UI (LLM + Rasmus)
- TBD
Report: 
- Methodology
	- Iterative development
	- Requirements
	- TechStack
		- Vague
- Theory / SOA
- Problem description / Context & Motivation 

# Infrastructure
Running two virtual machines with Debian 13. Each of these machines are a worker node in the Kubernetes cluster. The kubernetes control plane is running on the host machine. An Ansible script configures the OS of the worker nodes for reproducability and identically configured machines. 

Terraform was another candidate technology but requires an API to allocate compute (i think). 

# Application

# Business Logic
Most of the current implementation is agnostic as to the actual functionality of the application. The logic of the app is predominantly determining the canonical model, that is:
- what information is received from the various peripheral systems (CADs etc)
- what information is propagated to the various peripheral systems
- how is the information represented internally.
- how is the information represented externally to the dispatcher
- what actions are exposed to the dispatcher
	- how are valid actions constrained

Business logic is modelled around the Incident type. It should contain or reference all information that pertains to said incident. 
```go
type Incident struct {
	IncidentID
	DispatcherID
	
	AllocatedResources
	Status
}
```

# Data

# UI
