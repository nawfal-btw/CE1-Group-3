# recap
Shared communication interface for emergency services
# Architecture
##  Before:
single-node kubernetes cluster on minikube
## Now
multi-node kubernetes cluster on 3 worker-nodes with kubectl managed with portainer.io
- 2 physical machines
    - machine 1: Control Plane + 1 worker-node (vm)
    - machine 2: 2 worker-nodes (vm)

# Domain
Established incident model
Established business rules for incident model
- Valid incident state transitions
# Application
API adapted to domain for "api/incident" endpoints
Basic CRUD operations with Domain verification
# Data
Established database schemas and basic queries
# UI
No work
# Report

