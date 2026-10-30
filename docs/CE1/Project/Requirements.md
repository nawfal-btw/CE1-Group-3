___
See `Notes` for how to formulate idiomatic requirements
## Functional 
- The system shall function as a centralized broker, enabling standardized communication pertaining to incidents between emergency services
- The system shall persistently store all non-GIS communication and book keep any and all non-GIS changes to incidents.
- The system shall enable real-time tracking of mobile emergency resources during the emergency, with a latency < X & accuracy of Y.
- The system shall provide each user a role-specific dashboard
	- the role-specific dashboard shall contain only permitted information determined by the role and agent in quesiton
	- the role-specific dashboard shall contain only permitted actions determined by the role and agent in question
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


# Notes
Shall be
- Necesary: for viability of the product
- Appropriate: for the current scope of the project
- Unambigious: absolute - not up for interpretation
- Singular: pertain to only a single metric/feature
- Feasible: realistic completion within the constraints
- Verifiable: able to demonstrate completion of requirement
- Traceable: can be justified by reference to stakeholder 
- Correct: accurately reflects a need from which it originated
- Conforming: ...


