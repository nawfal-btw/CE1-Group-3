___
We have chosen to prioritize backend/API development for the foreseeable future. As we have established and implemented enough infrastructure to showcase the distributed nature of the application. There is still work to be done, but it will not be a bottleneck for some time. Is this pivot wise or should we develop all aspects in tandem?

Also, what are your requirements for the magnitude of distribution? does running it on e.g. 4 physical/virtual machines suffice, or is auto-scaling and similar critical to the acceptance criteria.

We have not found a use-case for certain technologies presented in Cloud Technologies. Should we emphasize technology coverage, even if it provides minimal utility?

___
We have segmented our project into 5 components of the overall system:
- Infrastructure 
- Application
- Domain/Business Logic
- User Interface
- Data
This separation has worked well from a development standpoint but should our report reflect the same separation?
___
We have chosen the following security measures for application:
- Given that all actors within the application are known in advance, we have opted to deploy it on a virtual private network, thus only those within the network can even reach the services. Furthermore, each device will also be equipped with a private-key for redundancy. Using passwords seems problematic due to the urgent nature and high-stress situation of the agents (Dispatcher, Operator, etc)
- Authorization will be implemented on a role-by-role basis, where privileges are propagated to the application layer to check if the requested action is permissible given the role.
	- Operators may only read/write to incidents of which they have been allocated
	- Dispatchers may only read/write to incidents of which they are assigned, or of which the incident owner has allowed read/write access.
- Thus the main attack-vector is unauthorized acquisition of a client-machine, but the exposure is minimized due to the *least privileges* setup
Is this adequate integration of the security aspect of the project.
___
Who is the target audience of the report i.e. what can be assumed to be known: Technical 










