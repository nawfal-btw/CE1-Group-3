Questions:
How does operator actions translate into core functionality
- If operator allocates a resource, how does that look internally. 
- What guardrails are in place to deter or disable invalid actions

Operator Actions:
- Create an Incident
	- Create partial w/ required and optional fields
- Update an Incident
	- Multi-field update, Single-field update
- Attach Allocations to Incident
- View an Incident
- Forward an Incident
- Close an Incident
- Query Live Incidents
- Track Resources

In the end, this is just a giant CRUD app...

Basic model has been established, just adapt the API... but how