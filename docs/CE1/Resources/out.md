NENA Standard for Emergency Incident Data Object (EIDO) Abstract: As
agencies and regions move forward with implementing NG9-1-1 and IP based
emergency communications systems, it is critical that they adhere to a
standardized, industry neutral format for exchanging emergency Incident
information between disparate manufacturer's systems located within one
or more public safety agencies, and with other Incident stakeholders.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021 DSC Approval: 08/10/2021 PRC Approval: 09/29/2021
NENA Board of Directors Approval: 10/14/2021 ANSI Board of Standards
Review Approved: 10/19/2021 Next Scheduled Review Date: 10/20/2024
Prepared by: National Emergency Number Association (NENA) Agency Systems
Committee, EIDO JSON Working Group Published by NENA Printed in USA

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Executive Overview There are many Functional Elements (FEs) within an
NG9-1-1 system that are used to process emergency calls. Some of these
FEs may be within a specific agency, in another agency or elsewhere in
an Emergency Services IP Network (ESInet). In many cases, an emergency
call is related to, or results in the creation of an "Incident" as
defined in NENA STA-010.2-2016 \[3\]. As public safety communication
center personnel process emergency calls for service and their
associated Incidents, new information about the Incidents is obtained.
There are many sources available to communication center personnel for
obtaining new Incident information during call handling, Incident
creation, dispatch, Incident monitoring, and post Incident analysis
processes. Newly gathered information, as well as changes in Incident
status, must often be passed on to other FEs, other involved agencies,
and frequently to non-emergency entities authorized to receive emergency
Incident information. As agencies and regions move forward with
implementing NG9-1-1 and IP based emergency communications systems, it
is critical that they adhere to a standardized, industry-neutral format
for exchanging emergency Incident information between disparate
manufacturer's systems located within one or more public safety
agencies, and with other Incident stakeholders.

10/19/2021

Page 2 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Table of Contents NENA STANDARD FOR EMERGENCY INCIDENT DATA OBJECT
(EIDO) ......................................1 EXECUTIVE OVERVIEW
.......................................................................................................2
DOCUMENT TERMINOLOGY
........................................................................................................................
7 INTELLECTUAL PROPERTY RIGHTS (IPR) POLICY
............................................................................................
8 REASON FOR ISSUE/REISSUE
.....................................................................................................................
8 DATA ASSOCIATED WITH AN EMERGENCY INCIDENT
.....................................................11 DATA COMPONENTS
..................................................................................................................
11 DATA COMPONENT TABLE STRUCTURE
...........................................................................................
11 JSON REFERENCE DATA COMPONENT
............................................................................................
13 EMERGENCY INCIDENT DATA OBJECT DATA COMPONENT
.................................................................... 13
AGENT DATA COMPONENT
..........................................................................................................
16 AGENCY DATA COMPONENT
.........................................................................................................
18 SPLIT/MERGE DATA COMPONENT
.................................................................................................
20 LINK DATA COMPONENT
.............................................................................................................
22 INCIDENT DATA
COMPONENT.......................................................................................................
23 CALL DATA COMPONENT
.............................................................................................................
28 CALL BACK DATA COMPONENT
.....................................................................................................
33 UPDATED CALL BACK NUMBER DATA COMPONENT
.............................................................................
34 DISPATCH DATA COMPONENT
......................................................................................................
35 DISPOSITION DATA COMPONENT
..................................................................................................
37 NOTES DATA COMPONENT
..........................................................................................................
39 PERSON DATA COMPONENT
.........................................................................................................
40 VEHICLE DATA COMPONENT
........................................................................................................
41 LOCATION DATA COMPONENT
......................................................................................................
43 ADDITIONAL DATA DATA COMPONENT
...........................................................................................
47 EMERGENCY RESOURCE DATA COMPONENT
.....................................................................................
48 ALARMS AND SENSORS DATA COMPONENT
......................................................................................
54 IANA ACTIONS
..................................................................................................................56
AGENCY ROLE REGISTRY
............................................................................................................
56 INCIDENT TYPE -- COMMON REGISTRY
...........................................................................................
57

10/19/2021

Page 3 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 INCIDENT STATUS-COMMON
........................................................................................................
66 REPORT NUMBER
TYPE...............................................................................................................
67 COMMON DISPOSITION CODES
.....................................................................................................
67 PERSON ROLE
..........................................................................................................................
72 VEHICLE RELATIONSHIP TYPE
......................................................................................................
74 LOCATION TYPE
.......................................................................................................................
74 PRIMARY UNIT STATUS-COMMON
.................................................................................................
75 SECONDARY UNIT STATUS-COMMON
.............................................................................................
76 EMERGENCY RESOURCE TYPE-COMMON
..........................................................................................
79 RESOURCE ATTRIBUTE
...............................................................................................................
84 IMPACTS, CONSIDERATIONS, ABBREVIATIONS, TERMS, AND DEFINITIONS
............. 132 OPERATIONS IMPACTS SUMMARY
................................................................................................
132 TECHNICAL IMPACTS SUMMARY
..................................................................................................
132 SECURITY IMPACTS SUMMARY
....................................................................................................
132 RECOMMENDATION FOR ADDITIONAL DEVELOPMENT WORK
............................................................... 132
ANTICIPATED TIMELINE
............................................................................................................
133 COST FACTORS
......................................................................................................................
133 COST RECOVERY CONSIDERATIONS
.............................................................................................
133 ADDITIONAL IMPACTS (NON-COST
RELATED)..................................................................................
134 ABBREVIATIONS, TERMS, AND DEFINITIONS
..................................................................................
134 REFERENCES
..................................................................................................................
137 EXHIBITS
.......................................................................................................................
138 APPENDIX A REQUIRED NIEM DATA ELEMENT THAT NEED TO BE SUPPORTED FOR
EIDO COMPLIANCE
.........................................................................................................................
138 A.1

PERSON TYPE
............................................................................................................
138

A.2

VEHICLE TYPE
............................................................................................................
145

A.3

DATE TYPE
.................................................................................................................
148

A.4

LOCATION TYPE
.........................................................................................................
149

A.5

IMAGE TYPE
...............................................................................................................
150

A.6

LENGTH MEASURE TYPE
............................................................................................
151

A.7

INJURY TYPE
.............................................................................................................
152

A.8

ORGANIZATION TYPE
................................................................................................
153

10/19/2021

Page 4 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A.9

IDENTIFICATION TYPE
..............................................................................................
155

A.10

STATUS TYPE
.............................................................................................................
156

A.11

ADDRESS TYPE
..........................................................................................................
156

A.12

AREA TYPE
.................................................................................................................
158

A.13

CONTRACTINFORMATION TYPE
................................................................................
158

A.14

ENTITY TYPE
..............................................................................................................
159

A.15

CROSSSTREET TYPE
...................................................................................................
160

A.16

PHYSICALFEATURE TYPE
...........................................................................................
160

A.17

MEDICALCONDITION TYPE
.......................................................................................
161

APPENDIX B MAPPING BETWEEN VCARD AND NIEM PERSON TYPE
................................... 161 APPENDIX C OPENAPI SCHEMA
............................................................................................
164 ACKNOWLEDGEMENTS
..........................................................................................................
195

10/19/2021

Page 5 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 NENA STANDARD DOCUMENT NOTICE
This Standard Document (STA) is published by the National Emergency
Number Association (NENA) as an information source for 9-1-1 System
Service Providers, network interface vendors, system vendors,
telecommunication service providers, and 9-1-1 Authorities. As an
industry Standard, it provides for interoperability among systems and
services adopting and conforming to its specifications. NENA reserves
the right to revise this Standard Document for any reason including, but
not limited to: • • •

Conformity with criteria or standards promulgated by various agencies,
Utilization of advances in the state of the technical arts, Reflecting
changes in the design of equipment, network interfaces, or services
described herein. This document is an information source for the
voluntary use of communication centers. It is not intended to be a
complete operational directive. It is possible that certain advances in
technology or changes in governmental regulations will precede these
revisions. All NENA documents are subject to change as technology or
other influencing factors change. Therefore, this NENA document should
not be the only source of information used. NENA recommends that readers
contact their 9-1-1 System Service Provider (9-1-1 SSP) representative
to ensure compatibility with the 9-1-1 network, and their legal counsel,
to ensure compliance with current regulations. Patents may cover the
specifications, techniques, or network interface/system characteristics
disclosed herein. No license is granted, whether expressed or implied.
This document shall not be construed as a suggestion to any manufacturer
to modify or change any of its products, nor does this document
represent any commitment by NENA, or any affiliate thereof, to purchase
any product, whether or not it provides the described characteristics.
By using this document, the user agrees that NENA will have no liability
for any consequential, Incidental, special, or punitive damages arising
from use of the document. NENA's Committees have developed this
document. Recommendations for changes to this document may be submitted
to: National Emergency Number Association 1700 Diagonal Rd, Suite 500
Alexandria, VA 22314 202.466.4911 or commleadership@nena.org

10/19/2021

Page 6 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 NENA: The 9-1-1 Association
improves 9-1-1 through research, standards development, training,
education, outreach, and advocacy. Our vision is a public made safer and
more secure through universally-available state-of-the-art 9-1-1 systems
and better-trained 9-1-1 professionals. Learn more at nena.org. Document
Terminology This section defines keywords, as they should be interpreted
in NENA documents. The form of emphasis (UPPER CASE) shall be consistent
and exclusive throughout the document. Any of these words used in lower
case and not emphasized do not have special significance beyond normal
usage. 1. MUST, SHALL, REQUIRED: These terms mean that the definition is
a normative (absolute) requirement of the specification. 2. MUST NOT:
This phrase, or the phrase "SHALL NOT", means that the definition is an
absolute prohibition of the specification. 3. SHOULD: This word, or the
adjective "RECOMMENDED", means that there may exist valid reasons in
particular circumstances to ignore a particular item, but the full
implications must be understood and carefully weighed before choosing a
different course. 4. SHOULD NOT: This phrase, or the phrase "NOT
RECOMMENDED" means that there may exist valid reasons in particular
circumstances when the particular behavior is acceptable or even useful,
but the full implications should be understood and the case carefully
weighed before implementing any behavior described with this label. 5.
MAY: This word, or the adjective "OPTIONAL", means that an item is truly
optional. One vendor may choose to include the item because a particular
marketplace requires it or because the vendor feels that it enhances the
product while another vendor may omit the same item. An implementation
which does not include a particular option "must" be prepared to
interoperate with another implementation which does include the option,
though perhaps with reduced functionality. In the same vein, an
implementation which does include a particular option "must" be prepared
to interoperate with another implementation which does not include the
option (except, of course, for the feature the option provides). These
definitions are based on IETF RFC 2119 \[2\].

10/19/2021

Page 7 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Intellectual Property Rights
(IPR) Policy NOTE -- The user's attention is called to the possibility
that compliance with this standard may require use of an invention
covered by patent rights. By publication of this standard, NENA takes no
position with respect to the validity of any such claim(s) or of any
patent rights in connection therewith. If a patent holder has filed a
statement of willingness to grant a license under these rights on
reasonable and nondiscriminatory terms and conditions to applicants
desiring to obtain such a license, then details may be obtained from
NENA by contacting the Committee Resource Manager identified on NENA's
website at www.nena.org/ipr. Consistent with the NENA IPR Policy,
available at www.nena.org/ipr, NENA invites any interested party to
bring to its attention any copyrights, patents or patent applications,
or other proprietary rights that may cover technology that may be
required to implement this standard. Please address the information to:
National Emergency Number Association 1700 Diagonal Rd, Suite 500
Alexandria, VA 22314 202.466.4911 or commleadership@nena.org Reason for
Issue/Reissue NENA reserves the right to modify this document. Upon
revision, the reason(s) will be provided in the table below. Document
Number NENA-STA-021.1-2021

10/19/2021

Approval Date 10/19/2021

Reason For Issue/Reissue Initial Document

Page 8 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Document Number
NENA-STA-021.1a-2021

10/19/2021

Approval Date

Reason For Issue/Reissue

04/19/2022

Publishing error corrections: Section 2.5 page 15 agentRoleRegistryText
example values now match values from the "Agent Roles" registry. Section
2.10 page 29 answerDate description now say "date and time stamp"
instead of "date". Section 2.6 page 17 agencyType adds that values are
from IANA urn:emergency:service:responder registry and the subregistries
and provide examples.

Page 9 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Document Number
NENA-STA-021.1b-2021

Approval Date 12/17/2024

Reason For Issue/Reissue Publishing error corrections: Section 2: Added
normative text specifying that EIDO only convey Incident current state,
not history. Section 2.20: Changed the minimum for
resourceAttributeRegistryText from 0 to 1 because it is a Required data
element. Section 2.8: Set the maximum for involvedAgencyReference data
element in the link data component to 1 to match the normative schema.
Section 2.20: Added the sipIdentities data element to the call data
component that allow multiple values. It replace the sipIdentity header
from previous revision. Section 2.11: updatedCbn element maximum in
CallbackType component is now set to 1 to be compliant with the schema.
Section 2.19: additionalDataComponent element clarification to indicate
Additional Data in the provided-by section of an PIDF-LO should not be
copied in a Additional Data Data Component. Section 2.11:
deviceContactHeader element in Call Back Data Component is now optional
to match the minimum number of occurrence.

10/19/2021

Page 10 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Data Associated with an Emergency Incident All mention of registries in
this document reference IANA (https://www.iana.org/protocols)
"emergency" registry unless specified otherwise. An EIDO SHALL represent
the state of the Incident, as known by the sender, at the time it was
sent. Only known current state of the Incident is included. Every EIDO
is logged when sent or received. For history about the Incident,
entities will query the Logging Service. Data Components Chapter 2 (Data
Associated with an Emergency Incident) identifies the data elements
associated with an emergency Incident grouped into various data
components. Each data component, except Reference, has a common set of
elements which include: •

• • •

Identifier (\$id) -- required - An identifier that uniquely labels a
data component. For data components where there is an identifier such as
a Call Identifier, Agency Identifier or Agent Identifier, that
identifier is used for this member. For data components that do not have
such an identifier, the form of the Identifier is a URN formed by the
prefix "urn:emergency:uid:eidodatacomponent:", a unique string
containing alpha and/or numeric characters, the ":" character, and
Agency Identifier of the Agency that created the data component. For
example,
"urn:emergency:uid:eidodatacomponent:a56e556d871:idx.psap.allegheny.pa.us"
is a properly formatted eidodatacomponent. The identifier MUST be the
same for all instances of this specific data component for all EIDOs for
the incident, and unique from all other data components in the EIDO.
lastUpdateTimeStamp - required - Timestamp of the last time the data
component was modified. Must be in the W3C datetime format as specified
in NENA-STA010\[3\]. updatedByAgencyReference - mandatory -- Reference
to the agency that updated the data component last.
updatedByAgentReference - conditional -- Reference to the Agent that
updated the data component last. Must be present if the data component
was created or updated by an agent. Data Component Table Structure

The initial (heading) section of each data component contains the
following information blocks:

10/19/2021

Page 11 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 • • • • • •

Data Component -- the name of the data component (e.g., EIDO Header,
Agency Information, etc.) Data Component Use -- Identifies whether the
data component is required or optional in EIDO instances. Minimum Number
-- the minimum occurrences of the data component allowed in EIDO
instances. Maximum Number -- the maximum occurrences of the data
component allowed in EIDO instances. Child Of: Identifies the potential
parents of which a data component can be a child. See the parent
component to determine the relationship between the child and parent.
Data Component Description -- a general description of the purpose and
contents of the data component.

The above information blocks are followed by a table that identifies the
data elements included in the data components. Note that entire data
components are included as a complex data element within their parent
data component. In this case, the description identifies it as a complex
data element (data component) and defines the relationship between the
two data components. The following information is included for each data
element: • •

• • •

JSON Name -- Name of the data element. Use -- identifies whether the
data elements are required, optional or conditional. Data elements that
are conditional describe the conditions when they are required and when
they are optional. Required data elements can exist within optional data
components. Required data elements of an optional data component are
only required if the data component is included in an EIDO instance. For
example, not all EIDO instances will contain a Dispatch Information data
component since sufficient information to dispatch emergency resources
to the Incident are not yet available or assigned resource statuses have
not changed. However, if an EIDO instance contains a Dispatch
Information data component, that data component must always contain the
"Incident Type--Common" data element. Minimum number -- the minimum
occurrences of the data element that may be included in the data
component. Maximum Number -- the maximum occurrences of the data element
that may be included in the data component. Description -- a general
description of the data element. For complex data elements (data
components), the description identifies the relationship between the two
data components. 10/19/2021

Page 12 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Chapter 5 (References) includes
sources and references that can be used to obtain additional information
about related NENA standards, NIEM, and global justice terminology and
standards. Chapter 3 (IANA Actions) identifies registries that define
the domain of values that must be used for specific EIDO elements. JSON
Reference Data Component Data Component Use:

Optional component

Min: 0 Max: \* Child Of: Agent, Agency, Split/Merge, Link, Incident,
Call Data, Dispatch, Notes, Additional Data, Emergency Resource, Alarms
and Sensors Data Component Description: This Data Component is used by
other data components to reference data components present elsewhere in
the EIDO document. Table 2-1 JSON Reference Data Component

JSON Name \$ref

Use (required, optional, conditional) Required

Min Max Description 1 1 Data element containing the Identifier (\$id) of
the referenced data element.

Emergency Incident Data Object Data Component Data Component Use:

Required component

Min: 1 Max: 1 Child Of:

None

Data Component Description: This is the parent data component that
contains all the other data components to form an EIDO. This Data
Component must always be present. A functional element receiving an eido
for an incident where the state is already known and 10/19/2021

Page 13 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 that need to know what changed it
must compare the previously known state with the received EIDO. In this
data component, the identifier (\$id) is the Incident Tracking
Identifier. The Incident Tracking Identifier is a identifier assigned by
the first element in the ESInet. The form of an Incident Tracking
Identifier is defined in NENA-STA-0101\[3\]. Incident Tracking
Identifiers are globally unique and are associated with a single
emergency Incident. The Incident Tracking Identifier can be associated
with one or more emergency calls. It is carried through to any Incident
resulting from an emergency call. It may or may not be the same as the
local Incident ID. When EIDO is exchanged not as part of an Incident, it
contains an IncidentTrackingIdentifier as identified in STA-010 \[3\]
where "unique string" is "0". Example:
urn:emergency:uid:Incidentid:0:police.state.pa.us. Table 2-2 Emergency
Incident Data Object Component Use (required, optional, conditional)
Required

Min 1

issuingElementIdentific ation

Required

1

mergeComponent

Optional

0

JSON Name eidoVersion

Max Description 1 Version of the EIDO structure. For this version of the
standard the value MUST be "1.0". For non breaking change to the
structure the minor version number will be incremented. For breaking
change to the structure the major version number will be incremented. 1
An identifier in the format defined in NENA-STA-010\[3\].Identify the
element that created the EIDO. \* Complex data element (Data Component).
Contains merge and

1 NENA-STA-010 include text to handle the case of sending the EIDO by
reference or value in a transferred

SIP call.

10/19/2021

Page 14 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Use (required, optional,
conditional)

Min

Max

linkComponent

Optional

0

- 

incidentComponent

Optional

0

1

callComponent

Optional

0

- 

callbackComponent

Optional

0

- 

dispatchComponent

Optional

0

- 

notesComponent

Optional

0

- 

emergencyResourceCo mponent

Optional

0

- 

alarmsSensorCompone nt

Optional

0

- 

agencyComponent

Required

1

- 

JSON Name

10/19/2021

Description split information related to the Incident. Complex data
element (Data Component). Contains link information related to the
Incident. Complex data element (Data Component). Contains general
information about the Incident. Complex data element (Data Component).
Contains information about calls associated with the Incident. Complex
data element (Data Component). Contains information about how to call a
person. Complex data element (Data Component). Contains dispatch
information related to the Incident. Complex data element (Data
Component). Contains Incident notes and comments associated with the
Incident. Complex data element (Data Component). Indentifies emergency
resources involved with the Incident. Complex data element (Data
Component). Identifies Alarms/Sensors associated with the Incident.
Complex data element (Data Component). Identifies all agencies involved
with the Incident. Page 15 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Use (required, optional,
conditional) Conditional

Min 0

Optional

0

0

personComponent

Conditional: Required if a location is known about the Incident Optional

0

- 

vehicleComponent

Optional

0

- 

JSON Name agentComponent

additionalDataCompon ent

locationComponent

Max Description \* Complex data element (Data Component). Identifies all
agent involved with the Incident. Must be provided if an Agent is
involved in the Incident. \* Complex data element (Data Component). All
Additional Data related to the Incident except the objects found in the
`<providedby>`{=html} field of a PIDF-LO. \* Complex data element (Data
Component). All locations related to the Incident. Complex data element
(Data Component). Every person related to the Incident. Complex data
element (Data Component). Every vehicle related to the Incident.

Agent Data Component Data Component Use: Conditional Component Minimum
Number: 0 Maximum Number: \* Child Of: Emergency Incident Data Object
Data Component Description: This Data Component contains information
about agents (e.g., call takers, dispatchers, supervisors, responders,
etc.) that are involved in the Incident. There may be multiple Agent
Information data components in the case where both a call taker and
dispatcher are involved in an Incident, where multiple dispatch agencies
are associated with the same Incident, and similar situations.
10/19/2021

Page 16 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Rarely, as in the case of
automatically dispatched responses, the agent may be an automaton
(automated system) such as an Interactive Media Response (IMR).
Automatons that are actively involved in an Incident or call must be
assigned an Agent ID that follows the i3 naming conventions (see
NENA-STA-010 \[3\] for more information). In this data component, the
identifier (\$id) is the Agent Identifier. Table 2-3 Agent Data
Component

JSON Name agentJcard

Use (required, optional, conditional) Optional

Min Max Description 0 1 The agent's contact information.

agentWorkstationPosition Identification

Conditional: Required if agent workstation has a position
identification.

0

1 The workstation position ID within the agent's agency.

agentRoleRegistryText

Required

1

- The Role of the Agent -- Dispatching, Call Taking, etc. The acceptable
  roles are defined in the "Agent Roles" registry in NENA-STA-010 \[3\].

notesReference

Optional

0

- Reference to a Notes Data Component. Contains optional alphanumeric
  text further describing the agent.

10/19/2021

Page 17 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name agencyReference

Use (required, optional, conditional) Required

Min Max Description 1 1 Reference to an Agency Data Component.
Identifies the agency employing or contracting with the agent.

Agency Data Component Data Component Use: Required Minimum Number: 1
Maximum Number: \* Child Of: Emergency Incident Data Object Data
Component Description: Every EIDO will include at least one instance of
this data component in the agency component of the Emergency Incident
Data Object in order to identify the Agency creating the EIDO. Many
Incidents have one owner, a specific agency. Sometimes, ownership
changes from one owner to another. In some jurisdictions, there can be
more than one owner. Normally, ownership is passed from the current
owner to another, but there are circumstances where ownership is
unclear, and ownership must be claimed. The Agency Information data
component provides a mechanism for establishing the agency that owns the
Incident associated with the Incident tracking ID contained in the EIDO
Header data component or for removing current ownership from that
Incident. Agencies in this context are construed broadly to include such
organization as poison control centers, tow truck operators, hospitals
and utility compagnies. In this data component, the identifier (\$id) is
the Agency Identifier. The Agency Identifier Identifiy the agency. See
NENA-STA-010 \[3\] Agency Identifier section for the format and
requirements of the identifier data element within the identity
component. The agency whose identifier is found within this component
MUST have a record in the Service/Agency Locator as defined in
NENA-STA-010 \[3\] 10/19/2021

Page 18 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 2-4 Agency Data Component

JSON Name agencyJcard

Use (required, optional, conditional Mi Ma ) n x Description Optional 0
1 The agency's contact information.

agencyRoleDescriptionRegistryTe xt

Required

1

agencyType

Required

1

- The role of the agency in relation to the Incident. Valid roles are
  available in the "EIDOAgencyRole" registry and include:Dispatching,
  Dispatched, CallReceiving, and TransferredTo. See Section 3.1, below,
  for the registry description.
- One or more members of a list of available provider and agency types.
  The agency types are taken from the IANA
  urn:emergency:service:responder registry and the subregistries defined
  for that registry. Values exclude the "urn:emergency:service:responde
  r." prefix. Example values: "police.federal.fbi", "fire.forest",
  "poison_control"

notesReference

10/19/2021

Optional

0

- Reference to a Notes Data Component. Contains optional alphanumeric
  text about the agency in the context of the incident. Page 19 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name incidentOwningAgencyIndicator

Use (required, optional, conditional Mi Ma ) n x Description Optional 0
1 Boolean data element that, if true, indicates that the agency
associated with the Agency ID contained in this data component owns; or,
if false does not own, the Incident associated with the Incident
tracking ID in the EIDO data component. Once set to true, it should only
be set to false by the agency that originally set it true. Used to
enable the exchange and update of Incident ownership information. Only
one Agency Data Component can have this set to true in the EIDO.

Split/Merge Data Component Data Component Use: Optional Component
Minimum: 0 Maximum: \* Child Of: Emergency Incident Data Object Data
Component Description: An optional data component that is used to
indicate the existence of a merged Incident Tracking ID or to split an
Incident. The presence of a Split/Merge Information data component
indicates that another Incident Tracking ID has been merged with or is
being split from the Incident Tracking ID contained in the EIDO header.

10/19/2021

Page 20 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 2-5 Split/Merge Data
Component Use (required, optional, JSON Name conditional) Min Max
Description incidentTrackingIdentifier Required 1 1 The Incident
Tracking Identification of the Incident that is being merged with, or
split from the Incident represented by the Incident Tracking
Identifcation contained in the EIDO Header. See Incident Split/Merge
Indicator to determine the direction of the merge/split.
incidentMergeDirectionCode Required 1 1 The role of the
IncidentTrackingIdentifier in this data component. Indicator is: -
"REPLACED" the Incident Tracking Identifier in this data component
contains the retired Incident Tracking Identifier. Only applies to a
merge operation. -"REPLACING" the Incident Tracking Identifier in this
data component contains the new Incident Tracking Identifier. Only
applies to a merge operation. - "Split" The Incident Tracking Identifier
contained in this data component is split from the Incident Tracking
Identifier contained in the EIDO Header data component. The other data
components contained in the EIDO contain the data elements of the split
Incident.

10/19/2021

Page 21 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Link Data Component Data
Component Use: Optional Component Minimum: 0 Maximum: \* Child Of:
Emergency Incident Data Object Data Component Description: An optional
data component that is used to indicate the existence of linked calls
and Incidents. A Link data component indicates that an Incident has been
linked to the Incident Tracking Identifier in the EIDO header. Incidents
are linked when it is determined that while they are separate Incidents,
they are related in some way. When a link is declared, both Incident
Tracking Identifiers continue to be used to track the individual
Incidents. Incidents may be linked in a hierarchical relationship. For
more information on hierarchal Incidents see the Incident Tracking
Identifier section of NENA-STA-010 \[3\]. Table 2-6 Link Data Component
Use (required, optional, JSON Name conditional) Min Max Description
incidentTrackingIdentifier Required 1 1 The Incident Tracking Identifier
of the Incident that is being linked to the Incident represented by the
Incident Tracking Identifier contained in the EIDO Header. The nature of
the link is defined by the link Indicator below. transactionReasonText
Optional 0 1 Free format narrative description of the reason for the
link. involvedAgencyReference Optional 0 1 Other Agency involved with
the linked Incident if known. linkDirectionCode Required 1 1 The
direction of the link -- If the value of the Link Indicator is: -
"Parent" The Incident Tracking Identifier contained in this data
component is the parent of the Incident Tracking Identifier contained in
the EIDO header. 10/19/2021

Page 22 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

Description - "Child" The Incident Tracking Identifier contained in this
data component is the child of the Incident Tracking Identifier in the
EIDO header. - "Related" The Incident Tracking Identifier contained in
this data component is related to the Incident Tracking Identifier in
the EIDO header, without any parent-child relationship. - "UnLink" The
Incident Tracking Identifier contained in this data component is
unlinked from the Incident Tracking Identifier contained in the EIDO
Header data component.

Incident Data Component Data Component Use: Optional Component Minimum:
0 Maximum: 1 Child Of: Emergency Incident Data Object Data Component
Description: The Incident Information data component is optional and is
used to exchange general information about emergency Incidents gathered
by emergency agents, emergency responders, from reporting parties
(callers), and devices reporting emergency Incidents. Agencies
responding to the same Incident may have different versions of the
Incident Information data components for the Incident. Each agency's
version of the Incident Information data component contains information
that is relevant to their agency but may also contain shared information
that is common to both agencies. This module is used to exchange
Incident update information, as well as for exchanging initial Incident
creation information. For example, in high priority Incidents only
partial information may be exchanged between call takers and dispatchers
(i.e., the Incident's type and location) while additional information is
being collected. This data component is 10/19/2021

Page 23 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 used for the initial, high
priority exchange and the subsequent exchange containing the additional
information collected after the initial exchange was completed. When
multiple callers report a single Incident, this data component is used
to update involved agents and responders about new information gathered
from the other callers. The Incident Information Data Component is also
used to exchange general Incident information developed during dispatch
operations. Call takers, Dispatchers, Emergency Resources, and Emergency
Devices can enter information exchanged/carried by this data component.
Table 2-7 Incident Data Component Use (required, optional, JSON Name
conditional) Min Max Description incidentTypeInternalCode Optional 0 1
An alphanumeric code indicating the category of the Incident. This is
the internal code used by the local agencies involved in the Incident.
incidentTypeInternalText Optional 0 1 Textual description indicating the
category of the Incident. This is the internal text used by the local
agencies involved in the Incident. incidentTypeCommonRegistryText
Required 1 1 Incident type code that is available in the
EIDOIncidentType-Common" registry (see Section 3.2, below, for the
registry description) and that most closely corresponds to the Incident
Type internal code. APCO has developed an ANS set of globally unique
common Incident type codes (APCO ANS 2.103.2-2019) \[8\], which form the
basis for this registry. Each Agency should maintain a mapping of its
Internal Incident 10/19/2021

Page 24 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

incidentStatusInternalText

Optional

0

incidentStatusCommonRegistry Text

Optional

0

10/19/2021

Description Types (IncidentTypeInternal) to the list of Common Incident
Types (IncidentTypeCommon). The Common Incident Type should be selected
from this mapping when the EIDO is created to identify the Incident type
using a common code that is globally understood. \* An alphanumeric code
indicating the status of the Incident (active, closed, structure
cleared, etc.). This is the internal code used by the local agencies
involved in the Incident. \* Incident status code that is available in
the Incident Status Registriy (see Section 3.3, below, for the registry
description) and that most closely corresponds to the Incident
Status-Internal. Typically used to track significant changes in an
Incident's status. Each Agency should maintain a mapping of its internal
Incident status (IncidentStatusInternal) to the list of common Incident
status (IncidentStatusCommon). The common Incident status should be
selected from this mapping when an EIDO is created to Page 25 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

internalIncidentId

Optional

0

agentReference

Optional

0

locationReference

Optional (Required if available)

0

personReference

Optional

0

10/19/2021

Description identify the Incident status using a common code that is
globally understood.

1 The Internal Incident ID as an alphanumeric string assigned by the
agency involved in the Incident. Used To exchange Incident information
between systems using the same internal Incident IDs. Maintained for
conformance with legacy systems. \* Reference to Agent Data Components.
Identifies the agents (could be either an agent in a communication
center or an emergency responder) that are involved in this Incident. \*
Reference to an Location Data Component. Contains Incident location
information entered or updated by an agent receiving a call associated
with the Incident. \* Reference to an Person Data Component. Contains
person information entered or updated by an agent receiving a call
associated with the Incident.

Page 26 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name vehicleReference

documentIdentification

reportNumberType

incidentPriorityInternalText

10/19/2021

Use (required, optional, conditional) Min Max Description Optional 0 \*
Reference to an Vehicle Data Component. Contains vehicle information
entered or updated by an agent receiving a call associated with the
Incident. Optional 0 \* The Document Identification connects the
Incident to one or more associated follow-up reports and investigations.
Each responding agency may have its own Document Identification. Also
used by agency supervisor and other personnel to track the status of
reports. Conditional: 0 1 Report number type codes that MUST be are
available in the EIDO populated if Registries (see Section 3.4, Report
below, for the registry Number is description); may be New or present
Reopened. Optional

0

1 Priority of the Incident as alphanumeric text for the agency being
dispatched. This value may only be meaningful to the agency providing
the information and other closely cooperating agencies. Note, that
different responding agencies may assign different priorities to same
Incident; for example, a high priority fire Incident may be a medium
priority law enforcement Incident. Page 27 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name incidentCommonPriorityNumber

beatOrDispatchGroupText

disposition

notesReference

Use (required, optional, conditional) Min Max Description Optional 0 1
Globally understood numeric Incident priority that ranges from 0 to 10,
with 10 being the highest priority and 0 being the lowest priority. The
Local Priority, described above, should be mapped to this (Common
Priority) data element so that all involved and interested agencies can
determine the relative priority of the Incident. Optional 0 1 The beat
or dispatch group that contains the Incident. Note that each agency
involved in the Incident may have its own beat or dispatch group.
Optional 0 \* Complex data element. Contains Incident disposition
information entered or updated by a dispatch agent and/or an emergency
responder. Optional 0 \* Reference to a Notes Data Component. Contains
optional alphanumeric text further describing the Incident.

Call Data Component Data Component Use: Optional Component Minimum: 0
Maximum: \* Child Of: Emergency Incident Data Object

10/19/2021

Page 28 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Data Component Description: The
Call Information data component is optional and is used to exchange call
information about the Incident received and collected by the agent
identified in this data component. There can be more than one call about
an Incident, and thus more than one instance of this data component can
be in an EIDO. Some of the information in this data component is
contained in the call, some is additional data associated with the call,
and other information is collected by the agent. Every participant in
the call is referenced by either an agentReference or personReference.
In this data component, the identifier (\$id) is the Call identifier.
The Call identifier is automatically created by the first ESRP in the
first ESInet that handles a call. Call Identifiers are globally unique
and are only valid for a specific call. Table 2-8 Call Data Component

JSON Name

Use (required, optional, conditional)

Min Max

relatedCallIdentifier queueIdentifier

optional optional

0 0

standardPrimaryCallType

Required

1

standardSecondaryCallType

optional

0

10/19/2021

Description CallIdentifier of calls that are related to this \* call.
Eg: Callbacks 1 An identifier of a queue where the call is currently on.
1 Primary call type. Values limited to those in the LogEvent CallTypes
registry with "Primary" classification. 1 Secondary call type. Values
limited to those in the LogEvent CallTypes registry with "Secondary"
classification.

Page 29 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name localCallType

Use (required, optional, conditional) optional

Min Max Description 0 1 Locally defined call type.

direction

Required

1

additionalDataReference

optional

0

10/19/2021

1 The "direction" member has one of these values, "incoming", "outgoing"
or "internal", where "incoming" means received by the EIDO sender,
"outgoing" means sent by the EIDO sender and "internal" means both
caller and recipient are within the EIDO sender. \* Reference to an
Additional Data Data Component. Additional information about a call
received that is involved in or related to the Incident. There may be
multiple data providers for one call. Additional Information is defined
in RFC 7852 \[6\]

Page 30 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name callStartTimestamp

answerDate callStateRegistryText

10/19/2021

Use (required, optional, conditional) Required

optional Required

Min Max Description 1 1 Date and time stamp of when the call was
received. SHALL be represented by the "date-time" datatype described in
RFC 3339, Section 5.6 as specified in NENA-STA-010.3-202x \[3\]. 0 1 A
date and time stamp on which the call was answered by an agent. 1 1
Current call status (when the EIDO was created) from the available call
statuses in the "Call States" registry i.e., callBegin, callAnswered,
callEnd, partyAdd , etc. Values limited to those in the IANA CallStates
registry.

Page 31 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name agentReference

Use (required, optional, conditional) Conditional

callBackReference

Conditional

verstat

Conditional

sipIdentities

Conditional

10/19/2021

Min Max Description 0 \* Reference to an Agent Data Component.
Identifies the agent(s) that is involved with the call described in this
data component and/or entered additional information related to the
received call. Must always be present if one or more Agent is involved
in the call. 0 1 Reference to a call back data component. Identifies how
the caller described in this data component can be called back. Required
if available. 0 1 Value of the verstat parameter. Required if known. 0
\* Content of the SIP identity headers associated with the caller
identity. Required if known.

Page 32 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name locationReference

Use (required, optional, conditional) Conditional

notesReference

Optional

personReference

Conditional

Min Max Description 0 \* Reference to a Location Data Component.
Contains call location information received with the call or updated by
the agent receiving the call. Must be present if location is presented
with the call. 0 \* Reference to a Notes Data Component. Contains
optional alphanumeric text further describing the call. 0 \* Reference
to a Person Data Component. Contains information about the caller(s).
Must be present if the call is a SIP session.

Call Back Data Component Data Component Use: Optional Component Minimum:
0 Maximum: \* Child Of: Emergency Incident Data Object Data Component
Description: This data component is used to exchange information about
how to contact the Incident's reporting parties.

10/19/2021

Page 33 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 2-9 Call Back Data
Component

JSON Name callBackInformationUri

deviceContactHeader

updatedCbn

Use (required, optional, conditional) Min Max Description conditional 0
\* Information that enables agents and responders to call back. URI(s)
found in the P-Asserted-Identity header if provided else the URI in the
From header. Optional 0 1 Content of the device contact header.
Information that enables agents and responders to possibly reach (call
back) the device that initiated the call. Note, that this information is
only guaranteed to be valid during the call and for a few minutes after
it ends. Optional 0 1 Complex data element. Identifies additional
telephone numbers and SIP equivalents that can be used to contact the
individual that made the call described in this data component.

Updated Call Back Number Data Component Data Component Use: Optional
Component Minimum: 0 Maximum: \* Child Of: Call Back Data Component
Description: This data component is optional and is used to exchange
information about additional phone numbers that can be used to contact
the Incident's reporting parties.

10/19/2021

Page 34 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 2-10 Updated Call Back
Number Data Component Use (required, optional, JSON Name conditional)
Min Max Description updatedCbnIdentifierUri Required 1 1 This data
element (in the form of a URI) is used to track additional telephone
number or SIP equivalents that can be used to contact the reporting
party of the parent call. updatedCbnCallerDescription Optional 0 1
Descriptive (alphanumeric) text that provides additional information
about the updated call back number such as hours to use it, days to use
it, and the type of number (e.g., work, home, friend, etc.) Dispatch
Data Component Data Component Use: Optional Component Minimum: 0
Maximum: \* Child Of: Emergency Incident Data Object Data Component
Description: This Data Component contains dispatch related information.
It allows updates to be sent and received between Incident Handling FEs
and Dispatch FEs, between different Dispatch FEs that are working the
same Incident, and enables exchanging information provided directly by
emergency responders. It can also be used to provide dispatch related
status updates to involved agencies and authorized stakeholders.

10/19/2021

Page 35 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 2-11 Dispatch Data
Component Use (required, optional, JSON Name conditional) Min Max
Description agencyReference Optional 0 1 Reference to an Agency Data
Component. Agency that was dispatched through action performed in this
data component. Note, that if several agencies are dispatched (one fire
and one police, two fire agencies, etc.), there must be a separate
instance of the dispatch data component for each dispatched agency.
Agencies are globally unique. agentReference Optional 0 1 Reference to
an Agent Data Component. Identifies the agent that completed the
dispatch operation described in, and/or entered the information
contained in this data component. emergencyResourceReference Optional 0
\* Reference to an Emergency Resource Data Component. Contains
information about emergency responders assigned (dispatched) to the
Incident, as well as their status and location updates. This data
component is not required when emergency responders update their status
and/or location directly without requiring dispatcher entry.
notesReference Optional 0 \* Reference to a Notes Data Component.
Contains optional alphanumeric text further describing the dispatch.

10/19/2021

Page 36 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Disposition Data Component Data
Component Use: Optional Component Minimum: 0 Maximum: \* Child Of:
Incident Information and Responder Information Data Component
Description: Agency specific and standardized disposition codes assigned
to an Incident. Multiple disposition codes per Incident are supported.
Either a responder or a dispatcher can close an Incident and assign a
final disposition codes to it. Table 2-12 Disposition Data Component Use
(required, optional, JSON Name conditional) Min Max Description
dispositionCommonRegistryCode Required 1 1 An agency assigns a
disposition to an Incident when its participation in the Incident ends.
The disposition code indicates whether follow-up reports are required
and other information about the Incident such as whether it resulted
from a false or actual alarm. They are used to exchange the status and
follow up requirements of an Incident upon its closure. The disposition
codes are drawn from a registry containing common disposition codes for
Police, Fire EMS disciplines. See Section 3.4, below, for the
description of the registry.

10/19/2021

Page 37 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name dispositionPrimaryIndicator

dispositionCategoryText

dispositionDescriptionText

notesReference

10/19/2021

Use (required, optional, conditional) Min Max Description Optional 0 1 A
designation of whether the common disposition code is the primary
disposition code for the Incident. Note that multiple primary codes are
allowed, but some systems may not be able handle more than one primary
common disposition code. It is possible that no codes are marked as
primary. The value is Boolean, where True is a primary disposition code.
Optional 0 1 An agency-specific, alphanumeric code that indicates how
the Incident was closed. The Common Disposition Code, referenced above,
should be mapped to the closest value of this data element. Conditional:
0 1 Descriptive text describing the MUST be Disposition Code Internal.
populated if Disposition codes may be Disposition agency specific and
this field Code Type is explains the meaning of the presentOptional
internal disposition code. otherwise. Optional

0

- Reference to a Notes Data Component. Contains optional alphanumeric
  text further describing the disposition code.

Page 38 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Notes Data Component Data
Component Use: Optional Component Minimum: 0 Maximum: \* Child Of:
Emergency Incident Data Object Description: This Data Component is
typically populated by emergency service agents and responders and
occasionally by automated devices. There may be multiple notes from the
same agent and there may be notes from multiple agents and agencies. It
is difficult to determine when a note should be considered history, and
thus not sent in an EIDO. For that reason, if a Note is referenced by a
data component sent in the EIDO, then the note SHALL be included in the
EIDO and remains part of what is considered current state. If a system
permits notes to be deleted, then deleted notes SHALL NOT be sent in
subsequent EIDOs. Table 2-13 Notes Data Component Use (required,
optional, JSON Name conditional) Min Max Description notesActionComments
Required 1 1 Notes created by an agent entered in HTML fragments, as
supported by NIEM, and shall be limited to 16 MB. HTML is used to allow
multimedia data to be contained in the notes. Security issues may arise
from embedding scripts, images, and other references including
JavaScript in notes and the receiving system may ignore or filter out
such embedded information. agentReference Optional 0 1 Reference to an
Agent Data Component. Identifies the agent and agency that entered the
note contained in this data component. Defaults to the Agency
Information included in the EIDO header, if this data component is not
present. 10/19/2021

Page 39 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Person Data Component Data
Component Use: Optional Component Minimum: 0 Maximum: *Child Of:
Emergency Incident Data Object Data Component Description: This Data
Component is used to exchange information about people associated with
an Incident including: callers, suspects, victims, witnesses and other
individuals involved in the Incident. The information is provided by
reporting parties and emergency responders. Responders can enter the
information either directly through their mobile data computers or via
their assigned dispatchers. An instance of Person Data Component MUST be
created along with a Call Data Component at the beginning of an
emergency call. The Person Data Component must have the
personIncidentRoleRegistryText element set to "caller" and the
callBackReference element refers to the same Call Back Data Component
designated by the Call Data Component. If the emergency call included a
SubscriberInfo Additional Data block, the Additional Data Component
reference for that block and its associated ProviderInfo block SHOULD be
specified in an additionalDataReference element. At the end of the
emergency call, while the EIDO may no longer include a Call Data
Component, the caller's details will continue to be found in the Person
Data Component. Table 2-14 Person Data Component Use (required,
optional, JSON Name conditional) Min Max Description
personIncidentRoleRegistryText Required 1* Describes the relationship
(Caller, Victim, suspect, etc.) of a person to the Incident. Available
person types are contained in the Person Role registry. See section 3.6,
below for the registry description. Note that there could be multiple
relationships as when the reporting party is also the victim. 10/19/2021

Page 40 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name callIdentifier

locationReference notesReference

ncPersonComponent additionalDataReference callBackReference

Use (required, optional, conditional) Min Max Description conditional 0
\* Call Identifier of the calls this person was involved on. Must be
provided if this person was involved in a call. Must be an Identifier as
defined in the Call Identifier section of NENA-STA-010 \[3\]. Optional 0
\* Reference to an Location Data Component. Location of the person.
Optional 0 \* Reference to a Notes Data Component. Contains optional
alphanumeric text further describing the person. Optional 0 1 NIEM Core
Person Type, a data type for a human being. Optional 0 \* Additional
Data about a person. Conditional 0 1 Reference to a call back data
component. Identifies how the person described in this data component
can be called back. Required if available.

Vehicle Data Component Data Component Use: Optional Component Minimum: 0
Maximum: \* Child Of: Emergency Incident Data Object Data Component
Description: This Data Component is used to exchange information about
vehicles associated with an Incident including: suspect vehicles,
vehicles involved in accidents, and other vehicles involved with the
Incident. The information is provided by 10/19/2021

Page 41 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 reporting parties and emergency
responders. Responders can enter the information either directly through
their mobile data computers or via their assigned dispatchers. Note:
Vehicle telematics information is not located in this data component. It
is located in additional data associated with a call. Table 2-15 Vehicle
Data Component

JSON Name vehicleRelationshipType

vehicleRelationshipTimeStamp

additionalDataReference

10/19/2021

Use (required, optional, conditional) Min Max Description Required 1 \*
Describes the relationship (victim's vehicle, accident vehicle, suspect
vehicle, etc.) of a vehicle to the Incident. Available vehicle
relationship types are contained in the "EIDOVehicleRelationshipType"
registry. See Section 3.7, below, for the registry description. Required
1 1 The date and time that the relationship of the vehicle to the
Incident was established by this data component instance. SHALL be
represented by the "datetime" datatype described in RFC 3339, Section
5.6 as specified in NENA-STA-010. Optional 0 \* Reference to an
Additional Data Data Component. Additional information about a vehicle
that is involved in or related to the Incident. There may be multiple
data providers for one vehicle. Page 42 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

locationReference

Optional

0

notesReference

Optional

0

ncVehicleComponent

Optional

0

Description

- Reference to a Location Data Component. Location of the vehicle.
- Reference to a Notes Data Component. Contains optional alphanumeric
  text further describing the vehicle. 1 NIEM Core VehicleType. A data
  type for a conveyance designed to carry an operator, passengers and/or
  cargo, over land.

Location Data Component Data Component Use: Optional Component. Minimum:
0 Maximum: \* Child Of: Emergency Incident Data Object Data Component
Description: This Data Component represents a location description. The
type of location may be the caller's location, the Incident's location
or another type of location indicated by the Location Type field in the
data component. In order to dispatch emergency responders an "initial"
Incident location is required. It may be the same as the caller's
location, but it may also evolve as the Incident progresses. For
example, emergency responders are initially dispatched to the caller's
location, the caller verbally describes a different location for the
Incident, and finally the first responders arrive at the scene and
relates yet another location for the Incident. The Incident's 10/19/2021

Page 43 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 location may also be mobile such
as a caller reporting an Incident from a moving vehicle or a law
enforcement chase in progress. The location data component references or
includes a PIDF-LO as documented in NENASTA-010\[3\]. The EIDO
incorporates many data elements from NIEM, which in turn, incorporates
many different methods of conveying a location, including shapes,
coordinates, addresses and others. This is required for interoperability
with systems that are conformant with NIEM that are not NG9-1-1 systems.
Even when the components of a NIEM-conformant location have the same
information as a PIDF-LO, it will not be PIDF-LO and will not be
compatible with NG9-1-1 systems. Since there is not a single location
type in NIEM, there are many different NIEM-conformant location formats
that may be received by an NG9-1-1 system and may be found in an EIDO.
Location in NG9-1-1 is expressed as PIDF-LO, so whenever possible, any
conveyance of location in an EIDO is preferred to be PIDF-LO, even
though a NIEM-conformant location received or conveyed in an EIDO may
not be a PIDF-LO. This may require construction of a PIDF-LO when
NIEM-conformant location information is received by an NG9-1-1 entity.
When expressing a location in a civic format, the PIDF-LO by itself is
unable to describe either a Cross Street or an Intersection. In this
document: \* An Intersection is a location that is described as the
convergence of more than one street, regardless of the orientation of
the streets to each other. This includes perpendicular, T-shaped, Y
shaped and any other orientation of streets. Intersections have at least
two and may have more than two streets. The location is the Intersection
itself. \* A Cross Street is an additional piece of information about a
location, and describes the nearest streets that cross or terminate near
the location. The Cross Street can have a perpendicular, T-shaped,
Y-shaped or any other orientation of the Cross Street to the street the
location is on. The Cross Street is not the location: it is near the
location. There could be more than one Cross Street if a location is in,
or near the middle of a block, or the Cross Street has multiple street
names (such as when a street changes names on either side of the street
the location is on). To represent Intersections and Cross Streets, we
use additional PIDF-LOs. The PIDF-LO SHOULD have country, A1, A2, A3
and/or A4, MAY have A5 fields, and MUST have a street name.
Implementations MUST be prepared to receive a Cross Street or
Intersection PIDF that does not contain one or more of country, A1-A5.
If they are missing, they are assumed to be the same as the values in
the primary location PIDF. For an Intersection, address number MAY be
provided. If present it is interpreted as a block number. For a Cross
Street, a civic address MUST be provided in locationByValue or
10/19/2021

Page 44 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 locationByReferenceUrl. In this
case, address number must be provided unless two cross streets are
provided or the road ends past the cross street (e.g., cul de sac).
Cross streets and Intersection are represented as separate PIDF object
and not multiple location in a single PIDF.Cross Street and Intersection
are mutually exclusive. If one is provided, the other MUST NOT be
provided. Table 2-16 Location Data Component

JSON Name locationTypeDescriptionRegistry Text

Use (required, optional, conditional) Required

Min 1

locationByValue

Conditional: MUST be populated if Location By Reference is not
populated. Otherwise it MUST be omitted.

0

locationByReferenceUrl

Conditional: MUST be populated if Location By Value is not populated.
Otherwise it

0

10/19/2021

Max Description 1 Location type (Caller, Initial, CurrentIncident,
Staging, Investigation, Tower Location, Other) as defined in a registry.
See Section 3.8, below for the registry description. 1 The Location
Information data component must support all PIDF-LO data elements though
many of these elements may not be present in an EIDO. When populated
this data element must contain civic or geodetic location elements. 1 A
URL that can be dereferenced to obtain the location of the indicated
location type. The resulting dereference must support the PIDF- LO data
elements defined in Page 45 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) MUST be omitted.

Min

Max

locationDescriptionText

Optional

0

1

crossStreetByValue

Optional

0

- 

crossStreetByReferenceUrl

Optional

0

- 

intersectingStreetByValue

Optional

0

- 

intersectingStreetByReferenceUrl

Optional

0

- 

cellTowerSectorId

Optional

0

1

10/19/2021

Description Location By Value. This is particularly useful for
indicating the location of moving devices such as callers in moving
vehicles. The current location of the device can be de-referenced and
inserted into the Incident record. Optional text further describing the
location type. Note that the Location may be the Caller's location,
Incident's location or another type of location depending on the
Location Type field. The nearest Cross Street to the Incident's location
in PIDF-LO format. The URL of the nearest Cross Street to the Incident's
location. The nearest Intersection to the Incident's location in PIDF-LO
format. The URL of the nearest Intersection to the Incident's location.
Text field contain the ID of the nearest cell tower and the sector/face
of Page 46 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Use (required, optional,
conditional)

Min

Max

additionalDataReference

Optional

0

- 

notesReference

Optional

0

- 

JSON Name

Description the tower receiving the call. May be used with the "Provided
By" field of the PIDF-LO to identify the carrier if carrier specific
data is needed. Reference to Additional Data Data Component. Contains
additional data associated with a location. Reference to a Notes Data
Component. Contains optional alphanumeric text further describing the
location.

Additional Data Data Component Data Component Use: Optional Component
Minimum: 0 Maximum: \* Child Of: Emergency Incident Data Object Data
Component Description: This data component is optional and is used to
exchange information about individuals, call or location associated with
a call received by an agent handling the Incident. This data component
is used to track and exchange this information as defined in RFC7852.

10/19/2021

Page 47 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 2-17 Additional Data Data
Component Use (required, optional, JSON Name conditional) Min Max
Description additionalDataByReferenceUrl Conditional: MUST 0 1 The URL
is a link to additional be populated if information data component
additionalDataDetail received about the parent. is not populated. Allow
the dereferencing of Otherwise it MUST additional information be
omitted. associated with a call related to the Incident.
additionalDataByValue Conditional: MUST 0 1 If additional data
associated only be populated if with a call is sent by value
additionalDataURL these fields contain the is not populated.
information. Otherwise it MUST be omitted. urlPurpose Conditional: MUST
0 1 Purpose that accompanied the be populated if a URL. For example, the
purpose have been purpose parameter sent with sent along with the the
URL in the Call-Info URL headers. notesReference Optional 0 \* Reference
to a Notes Data Component. Contains notes and comments related to the
Additional Data that were entered by agents and emergency responders.
Emergency Resource Data Component Data Component Use: Optional Component
Minimum: 0 Maximum: \* Child Of: Emergency Incident Data Object
10/19/2021

Page 48 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Data Component Description: An
emergency resource can be a vehicle, a person (foot patrol), an
organizational unit such as a squad or strike team, and other emergency
responder configurations. An emergency resource is described by a unique
unit ID and unit type. There may be multiple Emergency Resource data
components where multiple emergency resources are dispatched to a single
Incident. When resources are assigned to an Incident by a dispatcher,
then the parent data component is Dispatch Information and the same
agency and agent that entered the information contained in the Dispatch
Information data component entered the information contained in the
Emergency Resource data component. However, when resources update their
status or change the Incident, then the parent data component is the
EIDO Header. Agent Information data components references of Emergency
Resource instances identify the individuals associated with the
emergency response unit; for example officer Jeff Smith and John Jones
are currently operating (riding in) engine number 52. In this data
component, the identifier (\$id) is the unit ID. Format of the
identifier is defined as resource@AgencyId eg:
engine52@riversideFD.riverside.ca. Table 2-18 Emergency Resource Data
Component Use (required, optional, JSON Name conditional) Min Max
Description emergencyResourceTypeCommonReg Required 1 1 A standard code
for istryText an emergency type (FireTruckLadder, PatrolUnit, etc.).
Restricted to values in EIDOEmergencyResource Type-Common registry (see
Section 3.11, below). emergencyResourceTypeInternalText Optional 0 1 A
local code for an emergency resource type. resourceAttributeRegistryText
Required 1 \* A standard code for an emergency resource attribute
10/19/2021

Page 49 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

emergencyResourceName

Optional

0

1

primaryUnitStatusRegistryText

Required

1

1

10/19/2021

Description (skill and equipment possessed by an emergency resource).
Restricted to values in EIDOResourceAttribute registry (see Section
3.12, below). Common name of the emergency resource. The common,
globally unique, status that sets the emergency resource's ability to be
assigned to an emergency Incident. An emergency resource can only have
one Primary Unit Status-Common at any given time. Available options for
Primary Unit StatusCommon are contained the
"EIDOPrimaryUnitStatusCommon" registry. See Section 3.9, below, for the
registry description. Page 50 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

secondaryUnitStatusRegistryText

10/19/2021

Use (required, optional, conditional) Min Max

Required

1

- 

Description Agencies should map their Unit Status-Internal to the most
appropriate combination of Primary Unit StatusCommon and Secondary Unit
Status-Common available in the registry. Globally unique, statuses that
further qualifies the Primary Unit Status by providing more detail about
the associated Primary status. Some systems may not be able to handle
Secondary statuses, which is acceptable, but not recommended. Available
options for Secondary Unit Status-Common are contained the
"EIDOSecondaryUnitStatus -Common" registry. See Section 3.10, Page 51 of
196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

notesReference

Optional

0

- 

unitStatusInternal

Optional

0

- 

10/19/2021

Description below, for the registry description. Agencies should map
their Unit Status- Internal to the most appropriate combination of
Primary Unit StatusCommon and Secondary Unit Status-Common available in
the registry. Reference to a Notes Data Component. Contains notes and
comments related to the status of an emergency responder (e.g., time
that status is expected to change, etc.) Local or internal status of a
response unit. May be meaningful only to the owning agency and possibly
to other closely affiliated agencies. Some systems may Page 52 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

agentReference

Required

1

- 

unitLocationReference

Optional

0

1

unitDestinationReference

Optional

0

1

10/19/2021

Description not be able to handle multiple Unit Statuses-Internal.
Reference to an Agent Data Component. Identifies the agents currently
staffing the emergency responder unit and the agency to which the unit
belongs. If the responding unit updated the Incident directly on their
FRC, then the agent role identifies which individual entered the
information as well as their role in the emergency response unit
(driver, passenger, etc.). Reference to an Location Data Component.
Valid location of the unit at the time indicated by the Date/Time Stamp.
Reference to an Location Data Component. Page 53 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name

Use (required, optional, conditional) Min Max

eta

Optional

0

1

personReference

Optional

0

- 

vehicleReference

Optional

0

- 

Description Location where the Emergency Resource is going. Timestamp of
when the unit is expected to reach it's destination Reference to a
Person Data Component. Contains person information entered or updated
directly by an emergency responder. Reference to a Vehicle Data
Component. Contains vehicle information entered or updated directly by
an emergency responder. Represent the vehicle(s) used by the Emergency
Resource.

Alarms and Sensors Data Component Data Component Use: Optional Component
Minimum: 0 Maximum: \* 10/19/2021

Page 54 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Child Of: Emergency Incident Data
Object Data Component Description: This Data Component is only used to
support the exchange of the legacy APCO CSAA ANS interface \[12\]. There
may not be an i3 call associated with this alarm. The Incident Record
Handling FE will have a direct interface that supports APCO CSAA ANS and
the FE will use it to automatically extract the relevant information and
create an Incident. This data component provides a link to the original
information received from the alarm company. Any other form of Alarm or
Sensor data will be contained in the Additional Data Associated with a
Call. Table 2-19 Alarms and Sensors Data Component Use (required,
optional, JSON Name conditional) Min Max Description
csaaAlarmInformation Conditional: MUST be 0 1 Original automated alarm
data that populated if alarmURL triggered the creation of the is not
populated . Incident. Otherwise it MUST be Read in NIEM schema Alarm and
omitted. Sensor data. For Alarms this would be the APCO/CSAA ANS
2.101.22014 \[12\] standard.

alarmUrl

Conditional: MUST be 0 populated if CSAAAlarmInformation is not
populated. Otherwise it MUST be omitted.

10/19/2021

1 Link to the automated alarm data that initiated the Incident. There
may be more than one transmission for a single Incident. Enable the
receiving agency to dereference and obtain the original alarm
information that triggered the Incident.

Page 55 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

JSON Name agentReference

Use (required, optional, conditional) Optional

notesReference

Optional

Min Max Description 0 1 Reference to an Agent Data Component. Identifies
the agent and agency that processed the Alarm/Sensor information
described in this data component. Defaults to the Agent Information
included in the EIDO header, if this data component is not present. 0 \*
Reference to a Notes Data Component. Contains optional alphanumeric text
further describing the Incident.

IANA Actions This section describes the registries defined for the Data
Components contained in Chapter 2 of the document. Each defined registry
is linked to one or more data elements specified in the data components
contained in Chapter 2 (Data Associated with an Emergency Incident).
Each registry is named after the data element that uses it and includes
references to the Data Components that contain the registry's data
element. Agency Role Registry The "Agency Role" data element is
described in Section 2.6 (Agency Data Component) of the document. IANA
is requested to add the following values to the "EIDO-AgencyRole"
registry. Table 3-1 Agency Role Registry Value Assisting

Literal Description An agency is assisting on, or being consulted
regarding an Incident. CallReceiving The Agency that received the
emergency call. 10/19/2021

Reference This document This document

Page 56 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value Dispatched

Literal Description The Agency that had emergency resources dispatched
to an Incident. Dispatching The Agency that dispatched emergency
resources to an Incident. Informational An agency is receiving
information regarding an Incident, but is not otherwise participating.
TransferredTo The Agency to which the emergency call is being
transferred.

Reference This document This document This document This document

Incident Type -- Common Registry The "Incident Type -- Common" data
element is described in Section 2.9 (Incident Data Component) and
Section 2.13 (Dispatch Data Component) of the document. IANA is
requested to add the following values to the "EIDO-IncidentType-Common"
registry. Table 3-2 Incident Type -- Common Registry Value NONEMERG 911H
911OPN ABDOM ABDUCTF ABDUCTNF ABUSE ACCIND ACTSHOOT ADMIN AIR AIRINC
ALIEN ALLERGIC ALMA ALMB ALMH 10/19/2021

Literal Description The Incident is not an emergency. 9-1-1 Hang Up
9-1-1 Open Line Abdominal Pain or Problems Abduction Family Abduction
Non-Family Abuse Industrial Accident Active Shooter Administrative
Aircraft Related Incidents Airport Incident (not Aircraft) Immigration
Allergic Reaction, Insect Stings Alarm Audible Commercial Alarm
Panic/Duress/Hold Up Alarm

Reference This document APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 Page 57 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value ALMR ALRBA ALRBOX ALRBS
ALRCO ALRCO2 ALRDEFIB ALRF ALRGAS ALRHEAT ALRKPAD ALRMED ALRNGAS ALRPANA
ALRPANS ALRPROPN ALRPULL ALRROBA ALRROBS ALRU ALRVEH ALRWTR ALS ANIMAL
ASSAULT ASSIST ATL BACK BARICADE BIO BITE 10/19/2021

Literal Description Resdential Alarm Burglar Alarm-Audible Box Alarm
Burglar Alarm - Silent Carbon Monoxide Fire Alarm - CO2 Automated
Implantable Cardiac Defibrillator (AICD) Alarm or Event Generic Fire
Alarm Fire Alarm - Gas Fire Alarm - Heat Fire Alarm - Keypad Medical
Alarm Fire Alarm - Natural Gas Panic Alarm Audible Panic Alarm Silent
Fire Alarm - Propane Fire Alarm - Pull Station Robbery Alarm Audible
Robbery Alarm Silent Alarm Unknown Vehicle Alarm Fire Alarm - Water Flow
Advanced Life Support Incident Animal With or Without Injury Assist to
Other Agency or Request for Manpower ATL Non-Traumatic Barricaded
Biological Threat Animal Bites or Attack by an Animal

Reference APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 Page
58 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value BLDCLPSE BLEEDING BLS BOMB
BRGCLPSE BURGLARY BURN CARDIAC CARJACK CDX CHEMICAL CHESPAIN CHOKE
CITASST CIVDIS CIVIL CNTRFT COURT CRIMPS CRIMTRAF CRWD CU1 CU2 CU3
CURFEW CUST DEATH DEATHNTF DIABETIC DISORD 10/19/2021

Literal Description Collapsed Building Includes Technical Rescue
Incidents Includes Rectal, Vaginal, Non-Traumatic Basic Life Support
Incident Bomb Threat Or Investigation Bridge Collapse Burglary Burns
Cardiac Related Event, Heart Attack, Chest Pain Carjack Civil Dispute
Includes Chemical and Environmental Incidents Chest Pain Non-Cardiac
Choking Public Assistance, Citizen Standby, Public Service Civil
Disobedience Civil Counterfeit Money Court Order Serve Impersonation
Criminal Traffic Large, Rowdy Crowd Civil Unrest Violent Civil Unrest
Non-Violent Civil Unrest Information Curfew Violation Custodial Issues
Obvious Death, Found Body Death Notification Diabetic Disorderly

Reference APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 Page 59 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value DISVEH DMPX DROWN DRUGPAR
DRUGS DUI DVA DVNA EDP ELEC ELECHZD ELUDING ENTRAP ESCAPE EXPLOS
EXPOSURE EXTORT EYE FAIRC FALL FAM FAPT FARSON FAS FASX FBRUSH

FELEC

10/19/2021

Literal Description Disabled Vehicle Person Dumping Rubbish, In Progress
Drowning Drug Paraphernalia Drugs DUI Domestic Violent Domestic
Non-Violent EDP, Emotional Disturbance, Abnormal Behavior Includes
Struck By Lightning Electrical Hazard Eluding Includes Accidental
Entrapments and Inaccessible Incidents Escape Explosion Heat or Cold
Exposure Extort Eye Problems or Injury Fire Aircraft Fall Family Offense
Fire Apartment Fire Arson Felonious Assault - No Suspect Felonious
Assault - Suspect In Area Includes Trees, Brush, Grass, Controlled Burns
with No Structural Exposure, Outdoor Fires Fire Electrical, Includes
Arcing Wires, Shorted Electrical Equipment, Light Ballast, Transformers.

Reference APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019

APCO 2.103.2-2019

Page 60 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value FFUEL FIGHT FIREWRKS
FMARINE

FORCENTR FOTHER FRAUD FREPORT FSMOKE FSTRU FTRAILER FTRAIN FTRASH FTRUCK
FUNK FVEH FVEH GASLEAK GASLEAKI GASLEAKO HARR HATE 10/19/2021

Literal Description Includes Gas Stations, Fuel Refinery Fight Fireworks
Fire Marine, Includes Marine Vessels & Boats On Waterways (For Boats On
Land Use Fire Vehicle) Includes Breech, Forcible Entry Required For
Structure, Residential Lock Outs Includes Other Types Not Listed Fraud,
Bad Checks, ID Theft, Counterfeit Money Response to take Fire Report
Only, No Active Fire, Incident Not In Progress Fire Smoke, Includes Fire
Investigations, Smoke Odors No Flames Seen Includes Sheds, Attached
Garages, House Fires Fire Trailer Any Type of Railway Fire, Locomotive,
Rail Car, Passenger Car, Etc. Fire Trash, Refuse Fires, Garbage Cans,
Dumpsters Includes Large Vehicles Such As Rv, Tractor Trailers, Dump
Trucks, Buses Fire Unknown Fire Vehicle Includes Small or Medium Size
Vehicles, Cars, SUVs, Cans, Pick-Up Trucks Includes Natural Gas, Liquid
Fuel or Propane Leaks Gas Leak Inside Gas Leak Outside Harassment Hate
Crime

Reference APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019

APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 Page 61 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value HAZ HAZMAT

HEADACHE HEMRG HLP1 HLP2 HOLD HOMICIDE HOSTAGE IMP INDECENT INFO
INTERNET INTOX INVEST JUV LARCENY LIQUOR MAJORINC MAYDAY MAYDAY MCI
MISBOATR MISCON MNTL MTLV MVA

10/19/2021

Literal Description Wires Down, Other Natural Hazards Hazardous
Materials Incident Can Include Large Fuel Spill (\>100 Gallons),
Chemical, Biological, Radiological or Nuclear Headache Hemorrhage
Police/Fire/EMS In Trouble Police/Fire/EMS Nonemergency Citizen/Security
Holding Suspect Homicide Hostage Impersonation Indecent Information
Internet Crimes Intoxicated Investigation Juvenile Offense Larceny
Liquor Violation Major Incident Mayday Immediate Assistance, Request,
Emergency Back-Up Request MCI - Typically Used As A Code Declared By
Incident Command Missing Boater Misconduct Mental Person Non-Violent
Mental Person Violent Motor Vehicle Accident, No Injury - Can Be
Received As Automated Crash Notification

Reference APCO 2.103.2-2019 APCO 2.103.2-2019

APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019

Page 62 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value MVAHR MVAINJY

MVAUNK MVCP NCIC NCRASLTWF NCRCPRF NCRELEV NCRFACTR

NCRMUAIDF NCRNDE NCRNONE NCRTASKFORCE NCRTESTF NCRTRANS NEGLECT NOISE OD
ODOR OPEN ORDERV PANDEMIC PARKING 10/19/2021

Literal Description Hit & Run Motor Vehicle Accident With Injury - can
be received rs Automated Crash Notification, includes Car, Pedestrian,
ATV Accidents MVA Unk Motor Vehicle Collision on Private Property NCIC
Hit Assault with Weapon Medical - Cardiopulmonary Resuscitation Rescue -
Elevator Approved, Field Infusion Of Blood Products Necessary For
Life-Saving Intervention Mutual Aid Fire Notification Event, Local Only,
No Units Dispatched None - No Code Found Task Force

Reference APCO 2.103.2-2019 APCO 2.103.2-2019

Test Call - Fire Transportation Incident, includes Planes, Trains,
Subway, Etc. Neglect Noise Overdose - Accidental, Poisoning, Accidental
Ingestion Odor Investigation - any type Unsecure-Bldg Order Violation
Pandemic Flu or Illness Outbreak, Epidemic Parking

APCO 2.103.2-2019 APCO 2.103.2-2019

APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019

APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019

APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
Page 63 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value PERFND PERMISS PERSUSP
PERUNCON PERWANT PHANDLE PREG PROPFND PROPLOST PROPSTLN PROW PUBWKS
PURSUIT RADIAT RECKDRV REPO RESCUE RESIST RESP RESPARR ROBBERY RUNAWAY
SAR SCH SCRM SEIZURE SEXOFFNS SHOOT SHOPLIFT

10/19/2021

Literal Description Found Person Missing Person Suspicious Person
Unconscious Person, includes Fainting and Near Fainting Wanted Person
Panhandling Pregnancy Problems, Childbirth, Miscarriage Found Property
Lost Property Stolen Property Prowler Public Works / Utilities Pursuit
Includes Nuclear Incidents Reckless Driver Repo Includes Lock-Ins, High
Angle Rescue Resisting Arrest Difficulty Breathing, Shortness Of Breath,
Respiratory Issue Respiratory Arrest Robbery Runaway Search And Rescue
School Offense Person Screaming Includes Convulsions Sex Offense
Shooting Shoplifting

Reference APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019

Page 64 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value SICK SLMP SMOKING SOL
STABBING STDBY STKX STROKE SUBJSTP SUICIDE SUICTHRT SUSP TERROR THAZ
THREATS TOW TRAIN TRAINDRL TRANSFER

TRAUMA TRAUMAP TRESPASS TROUBLE TRUANT TSTOP VAG VANDAL VEHABND
10/19/2021

Literal Description Sick Person Non-Specific, General Illness, Flu
Person Slumped In Vehicle Smoking Solicitors Stabbing Police, Fire Or
EMS Stalking With Suspect On Scene Stroke Or Cerebral Vascular Accident
(Cva) Subject Stop Suicide Attempt Or Threatening Suicide Suspicious
Incident,Circumstances, Unknown Problems Terrorist Activity Traffic
Hazard Threats Tow Incidents Other Than Fires or Derailments Related to
Trains Train Derail Inter-Facility Transfer, Palliative Care, also Used
For Transfer of Apparatus Between Jurisdictions. Traumatic Injury
including Head Injury Penetrating Trauma Trespass Unknown Type Trouble
Truant Traffic Stop Vagrant/Panhandler Vandalism Abandoned Veh

Reference APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019

APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019
APCO 2.103.2-2019 Page 65 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value VEHASST VEHLOCK VEHREC
VEHSTLN VEHSUSP VICE WARRANT WATER WEAPON WELFARE WILDLIFE WIRES
XTRICATE

Literal Description Motorist Assist Vehicle Lock Out Recovered Vehicle
Stolen Vehicle Suspicious Vehicle Vice Search Warrant All Water
Incidents, Diving & Scuba Accidents, Boat/Jet Ski Accident, Water Weapon
Violation Man Down Check, Welfare Check Wildlife Violations Wires Down
Confined Space Rescue, Vehicle Extrication

Reference APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO 2.103.2-2019 APCO
2.103.2-2019 APCO 2.103.2-2019

Incident Status-Common The "Incident Status-Common" data element is
described in Section 2.9 (Incident Information Data Component) of this
document.2 IANA is requested to add the following values to the
"EIDO-IncidentStatus-Common" registry. Table 3-3 Incident Status --
Common Registry Value Active Cancelled Closed

Literal Description The Incident is active. Incident cancelled. Incident
closed.

Reference This document This document This document

2 The goal of the incident statuses contained in the registry is to use
distinct entries that are readily

identifiable by agents receiving an EIDO. Multiple incident statuses may
be used to specify complex situations such as an active fire incident
that has had a structure cleared. In this case, the following incident
statuses could be assigned: "Active" and "StructureCleared." Incident
status codes may be used to indicate either when a particular event
occurred during the incident (patient contact) and/or when the
incident's status has changed (ResourcesOnscene).

10/19/2021

Page 66 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value FireUnderControl
NewLocation ReOpened ResourcesAssigned ResourcesEnroute ResourcesOnscene
StructureCleared

Literal Description Emergency responder has declared that a fire is
under control. The Incident's location has changed. Incident has been
re-opened. Incident has had at least one emergency resource assigned to
it. At least one emergency resource is en route to the Incident. At
least one emergency resource has arrived at the location (on scene) of
the Incident. Emergency responder has declared that the structure has
been cleared.

Reference This document This document This document This document This
document This document This document

Report Number Type The "Report Number Type" data element is described in
Section 2.9 (Incident Data Component) of this document. IANA is
requested to add the following values to the "EIDO-ReportNumberType"
registry. Table 3-4 Report Number Type Registry Value New Ongoing
Reopened

Literal Description The report number is new The report number has not
changed The report number was previously issued, the related Incident
was closed and the report and/or Incident are being activated again

Reference This document This document This document

Common Disposition Codes The "Disposition Type" data element is
described in Section 2.14 (Disposition Data component) of this document.
The initial values for this registry are adopted from APCO ANS
1.111.2-2018: Common Disposition Codes for Data Exchange.3

3 Available at
https://www.apcointl.org/download/apco_ans_1-111-2-2018-disposition-

codes/?wpdmdl=5997&ind=0. Retrieved 22 April 2021. 10/19/2021

Page 67 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 IANA is requested to add the
following values to the "EIDO-CommonDispositionCode" registry. Because
APCO ANS 1.111.2-2018 uses a two-digit integer as an incident status
code, IANA is also requested to hold values 46-100 for future versions
of the standard. Table 3-5 Common Disposition Codes Registry Value
Literal Description 01 Report Taken 02 No Report Taken 03 Arrest 04
Citation Issued (Verbal Or Written) 05 Assignment Completed/Settled By
Phone Contact 06 Broadcast/Attempt to Locate (ATL) 07 Civil 08 Follow-Up
Report Taken 09 Summons/Infraction Issued 10 Event Cancelled 11 Field
Interview 12 Property (Seized Or Found) 13 Information/Administrative 14
Gone On Arrival (GOA)/Unable To Locate (UATL) 15 Referred to partner
agency/Animal Control/Medics/Fire, etc 16 Transport/Escort Given 17
False Alarm caused by: Weather, Power Outage/Other Related Cause/Animal

Reference APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018

18 19 20 21 22 23

APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO
ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018

Assignment Completed/Settled By Contact Building Condemned Building
Secure Unit Cancelled enroute Cancelled - Duplicate Call Fire Control or
Extinguishment (NFIRS 10) 10/19/2021

APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO
ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018

Page 68 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 24 25 26 27

45 46

Search and Rescue (NFIRS 20) EMS and Transport (NFIRS 30) Hazardous
Condition (NFIRS 40) Fires, Rescues and Hazardous Conditions (NFIRS 50)
Systems and Services (NFIRS 60) Assistance (NFIRS 70) Information,
Investigation, and Enforcement (NFIRS 80) Fill-In, Standby (NFIRS 90)
Cancelled (Prior to Arrival at Scene) Patient Treated, Transported by
Law Enforcement Patient Treated, Transported by Private Vehicle
Standby-No Services or Support Provided Standby-Public Safety, Fire, or
EMS Operational Support Transport of Body Parts or Organs Only No
Patient Found Patient Dead at Scene-No EMS Transport Patient Dead at
Scene-Transported by EMS Patient Evaluated, No Treatment/Transport
Required Patient Refused Evaluation/Care Patient Treated, Released
(Without Transport) Patient Treated, Transferred Care to Another EMS
Patient Treated, Transported by EMS RESERVED

47

RESERVED

48

RESERVED

28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44

10/19/2021

APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO
ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS 1.111.2-2018 APCO ANS
1.111.2-2018 Future release of APCO ANS 1.111.x Future release of APCO
ANS 1.111.x Future release of APCO ANS 1.111.x Page 69 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 49

RESERVED

50

RESERVED

51

RESERVED

52

RESERVED

53

RESERVED

54

RESERVED

55

RESERVED

56

RESERVED

57

RESERVED

58

RESERVED

59

RESERVED

60

RESERVED

61

RESERVED

62

RESERVED

63

RESERVED

64

RESERVED

65

RESERVED

66

RESERVED

10/19/2021

Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x

Page 70 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 67

RESERVED

68

RESERVED

69

RESERVED

70

RESERVED

71

RESERVED

72

RESERVED

73

RESERVED

74

RESERVED

75

RESERVED

76

RESERVED

77

RESERVED

78

RESERVED

79

RESERVED

80

RESERVED

81

RESERVED

82

RESERVED

83

RESERVED

84

RESERVED

10/19/2021

Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x

Page 71 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 85

RESERVED

86

RESERVED

87

RESERVED

88

RESERVED

89

RESERVED

90

RESERVED

91

RESERVED

92

RESERVED

93

RESERVED

94

RESERVED

95

RESERVED

96

RESERVED

97

RESERVED

98

RESERVED

99

RESERVED

100

RESERVED

Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x
Future release of APCO ANS 1.111.x Future release of APCO ANS 1.111.x

Person Role The "Person Type" data element is described in Section 2.16
(Person Data Component) of this document. IANA is requested to add the
following values to the "EIDO-PersonRole" registry. 10/19/2021

Page 72 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 3-6 Person Role Registry
Value InvolvedPerson Caller Patient PatientAcquaintance PatientRelative
PersonOfInterest ReportingParty Suspect Victim VictimAcquaintance
VictimRelative Witness

10/19/2021

Literal Description Person described in the Person Information data
component is involved in the Incident. Used when no other relationship
is known. Person described in the Person Information data component is a
caller. Person described in the Person Information data component is
involved in the Incident as a patient. Person described in the Person
Information data component is an acquaintance of a patient involved in
the Incident. Person described in the Person Information data component
is a relative of a patient involved in the Incident. Person described in
the Person Information data component is involved in the Incident as a
person of interest. Person described in the Person Information data
component is involved in the Incident as a reporting party. Person
described in the Person Information data component is involved in the
Incident as a suspect. Person described in the Person Information data
component is involved in the Incident as a victim. Person described in
the Person Information data component is an acquaintance of a victim
involved in the Incident. Person described in the Person Information
data component is a relative of a victim involved in the Incident.
Person described in the Person Information data component is involved in
the Incident as a witness.

Reference This document This document This document This document This
document This document This document This document This document This
document This document This document

Page 73 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Vehicle Relationship Type The
"Vehicle Relationship Type" data element is described in Section 2.17
(Vehicle Data Component) of this document. IANA is requested to add the
following values to the "EIDO-VehicleRelationshipType" registry. Table
3-7 Vehicle Relationship Type Registry Value AccidentVehicle
InvolvedVehicle

SuspectVehicle VictimVehicle WitnessVehicle ResponseVehicle

Literal Description The vehicle described in the Vehicle Information
data component is involved in the Incident as an accident vehicle. The
vehicle described in the Vehicle Information data component is involved
in the Incident. Used when no other relationship is known. The vehicle
described in the Vehicle Information data component is involved in the
Incident as a suspect's vehicle. The vehicle described in the Vehicle
Information data component is involved in the Incident as a victim's
vehicle. The vehicle described in the Vehicle Information data component
is involved in the Incident as a witness's vehicle. The vehicle
described in the Vehicle Information data component may be involved in
the Incident response.

Reference This document This document

This document This document This document This document

Location Type The "Location Type" data element is described in Section
2.18 (Location Data Component) of this document. IANA is requested to
add the following values to the "EIDO-LocationType" registry. Table 3-8
Location Type Registry Value 10/19/2021

Literal Description

Reference Page 74 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

RoutingLocation Caller CurrentIncident

InitialIncident Investigation Other Staging TowerLocation UnitLocation

For a location referenced in a call represent the location that was used
to route the call. Can be used by other data component to indicate their
routing location. The Location Information data component contains the
caller's location. The Location Information data component contains the
current reported location of the Incident. When Incident Location is
mentioned, it refers to a Location where Location Type has this value.
The Location Information data component contains the initial Incident's
location. The Location Information data component contains the
Incident's investigation location. Location is unspecified. The Location
Information data component contains a staging location for emergency
responders assigned to the Incident. The Location Information data
component contains the location of a cell tower that processed the call.
The Location Information data component contains the location of an
emergency responder.

This document

This document This document

This document This document This document This document This document
This document

Primary Unit Status-Common The "Primary Unit Status-Common" data element
is described in Section 2.20 (Emergency Resource Data Component) of this
document.4 IANA is requested to add the following values to the
"EIDO-PrimaryUnitStatus-Common" registry.

4 The purpose of the common (global) primary unit statuses is to enable
an agency that needs a particular

type of resource to ascertain the availability of that resource in a
different agency and to either request that resource from the agency
that owns it (mutual aid) or to assign it to an incident (automatic
aid). In most situations, agreements will be structured between the two
agencies enabling automatic and mutual aid.

10/19/2021

Page 75 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 3-9 Primary Unit
Status-Common Registry Value Available ConditionallyAvailable
NotAvailable

Literal Description Emergency Unit is available for Dispatch. Emergency
Unit is assigned to an activity, but is available for dispatch or
reassignment.

Reference This document This document

Emergency Unit is not available for Dispatch and cannot be assigned to a
call.

This document

Secondary Unit Status-Common The "Secondary Unit Status-Common" data
element is described in Section 2.20 (Emergency Resource Data Component)
of this document. IANA is requested to add the following values to the
"EIDO-SecondaryUnitStatus-Common" registry. Table 3-10 Secondary Unit
Status-Common Registry Value Literal Description
AcknowledgedTransmission Emergency unit acknowledged receipt of a
dispatch/assignment. AlternateLocation Emergency unit is at an alternate
location when used as a standalone secondary status or is en route to,
transporting to, arrived at, etc. when used in combination with another
secondary unit statuscommon. Arrived Emergency unit arrived at the
Incident location or at some other location. Assigned Emergency unit has
been assigned to an Incident or to some other event. AssignmentCancelled
Emergency unit's assignment to an activity has been cancelled. This is
an observation useful for tracking unit history and Incident progress.
BacktoAssignedArea Emergency unit is back to patrolling or covering its
assigned area, beat, or district. This is an observation useful for
tracking unit history. 10/19/2021

Reference This document This document

This document This document This document

This document

Page 76 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value Backup Break CheckedIn

Cleared COP_or_POP Court CoveringAlternateArea

Delayed Departed

Dispatched EnRoute EquipmentIssues Event

10/19/2021

Literal Description Emergency unit is backing up another emergency unit
on an Incident. Emergency unit is on a break. Emergency unit checked in
with its dispatcher. This is an observation useful for tracking unit
history and Incident progress. Emergency unit cleared the Incident
location or some other location. Emergency unit is involved in Community
Oriented Policing or Problem Oriented Policing activities. Emergency
unit is assigned to Court. Emergency unit is patrolling, has moved up,
or is covering an alternate area, beat, station, or district when used
as a standalone secondary status or is en route to, arrived at, etc.
when used in combination with another secondary unit status-common.
Emergency unit is delayed from arriving at the Incident's location or
some other location. Emergency unit has departed a location, where
another, simultaneously assigned, secondary unit status-common describes
the destination (e.g., Departed and Court). Emergency unit has been
dispatched to an Incident or some other event. Emergency unit is en
route to an Incident location or some other location. Emergency unit is
experiencing equipment issues. Emergency unit is at an event (parade,
concert, etc.) when used as a standalone secondary status-common or is
en route

Reference This document This document This document

This document This document This document This document

This document This document

This document This document This document This document

Page 77 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value

Hospital

InService Investigation Location Meal Meeting

OffDuty OnDuty OnScene OutofService PatientContact

Post

ResponderInitiatedEvent 10/19/2021

Literal Description to, arrived at, etc. when used in combination with
another secondary unit status-common. Emergency unit is at the hospital
when used as a standalone secondary unit status-common or is en route
to, transporting to, arrived at, etc. when used in combination with
another secondary unit status-common. Emergency unit is in service.
Emergency unit is assigned to investigate an active or closed Incident.
A location other than the Incident that is used with en route, arrived,
etc. Emergency unit is at lunch, dinner, breakfast or some other meal.
Emergency unit is involved in a meeting when used as a standalone
secondary unit status-common or is en route to, arrived at, etc. when
used in combination with another secondary unit status-common. Emergency
Unit is off duty. Emergency unit is on duty. Emergency Unit is located
at the scene (location) of the Incident. Emergency unit is out of
service. Emergency responders made contact with a patient involved in
the Incident. This is an observation useful for tracking unit history
and Incident progress. Emergency unit is at a post when used as a
standalone secondary unit statuscommon or is en route to, arrived at,
etc. when used in combination with another secondary unit status-common.
Emergency unit is on a self initiated event that is not a traffic stop.

Reference

This document

This document This document This document This document This document

This document This document This document This document This document

This document

This document Page 78 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value RollCall

Roster ShiftPending Staging

Station

TrafficStop Training Transporting Unmanned

Literal Description Emergency unit is at Roll Call when used as a
standalone secondary unit statuscommon or is en route to, arrived at,
etc. when used in combination with another secondary unit status-common.
Emergency unit has automatically been activated, but is not yet
available and has not checked in. Emergency unit's end of shift is
pending. Emergency unit is at an Incident's staging location when used
as a standalone secondary unit status-common or is en route to, arrived
at, etc. when used in combination with another secondary unit
status-common. Emergency unit is at its headquarters, station, or
substation when used as a standalone secondary unit status-common or is
en route to, transporting to, arrived at, etc. when used in combination
with another secondary unit status-common. Emergency unit is on scene at
a self initiated traffic stop. Emergency unit and responders are
participating in a training activity. Emergency unit is transporting or
escorting a person or equipment to a location or destination. Emergency
unit is not adequately staffed.

Reference This document

This document This document This document

This document

This document This document This document This document

Emergency Resource Type-Common The "Emergency Resource Type-Common" data
element is described in Section 2.20 (Emergency Resource Data Component)
of this document. IANA is requested to add the following values to the
"EIDO-EmergencyResourceTypeCommon" registry 10/19/2021

Page 79 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 3-11 Emergency Resource
Type-Common Registry Value Literal Description AirAmbulanceFixedWing Air
Ambulance (Fixed-Wing) AirAmbulanceRotaryWing Air Ambulance
(Rotary-Wing) AirSearchTeam Air Search Team
AirborneCommunicationRelayTeamFixed Airborne Communications Relay Team
(Fixed-Wing) FixedWingAircraft Aircraft, Fixed Wing RotaryWingAircraft
Aircraft, Rotary Wing AirborneTransportTeamFixedWing Airborne Transport
Team (Fixed-Wing) AirCompressor Air Compressor and Air Refilling
AirConditioner-Heater Air Conditioner/Heater ATV All Terrain Vehicle
GroundALSAmbulance Ambulance, Advance Life Support (Ground)
LargeAnimalRescueTeam Large Animal (horse, cow, etc.) Rescue Team
LargeAnimalShelteringTeam Large Animal (horse, cow, etc.) Sheltering
Team LargeAnimalTransportTeam Large Animal (horse, cow, etc.) Transport
Team SmallAnimalRescueTeam Small Animal (cat, dog, etc.) Rescue Team
SmallAnimalShelteringTeam Small Animal (cat, dog, etc.) Sheltering Team
SmallAnimalTransportTeam Small Animal (cat, dog, etc.) Transport Team
AssaultHomicideUnit Assault / Homicide Unit GroundBLSAmbulance
Ambulance, Basic Life Support (Ground) BikeTeam Bike Team BombSquad Bomb
Squad/Explosives Team BrushFireEngineType3 Type 3 Brush Fire Engine
(120GPM, 500 gal tank, 1000' 10/19/2021

Reference This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document This document This document

Page 80 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value

BrushFireEngineType4

BrushFireEngineType5

BrushFireEngineType6

BrushFireEngineType7

Bus CanineSearchandRescue DisasterCollapsedStructureCanineTech
CaveSearchandRescue CivilServiceOfficer CoastGuard
CollapseSearchandRescue CommandLevel4 CommandLevel3 CommandLevel2

10/19/2021

Literal Description 1.5" or larger, 800' 1 " or larger, 3 crew members)
Type 4 Brush Fire Engine (70 GPM, 750 Gal tank, 300' 1.5" & 300? 1" or
larger, 3 crew members) Type 5 Brush Fire Engine (50 GPM, 500 Gal, 300'
1.5" & 300? 1" or larger, 3 crew members) Type 6 Brush Fire Engine (50
GPM, 200 Gal, 300' 1.5" & 300? 1" or larger, 2 crew members) Type 7
Brush Fire Engine (20 GPM, 125 Gal tank, 200 and 200' on 1.5 and 1"
hose, 2 crew members) Buses Canine Search and Rescue Team Disaster
Collapsed Structure Canine Search Technician Cave Search and Rescue Team
Civil Service Officer Coast Guard Search and Rescue Collapse Search and
Rescue Teams Initial Responding Resources that arrive on scene Field
Supervisor (patrol sergeant, detective sergeant, battalion chief, etc.)
Section, Precinct or Division Chief/Commander

Reference

This document

This document

This document

This document

This document This document This document This document This document
This document This document This document This document This document

Page 81 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value CommandLevel1
CommSupportTeam CommunityPolicingUnit Crane CrisisInterventionUnit
ResponderCrewTransport DetoxTransport UnmannedAerialVehicle
ElectronicBoard FireBoat FireTruckLadder FireFoamTender FoamTank
ForensicEvidence FuelTender GangEnforcement HazMatTeam FireHelicopter
HeavyRescue IncidentManagementTeam Investigator
LawEnforcementCommandUnit PatrolUnit SchoolOfficer LightRescue

10/19/2021

Literal Description Department/Agency chief, sheriff, deputy chief,
assistant chief, etc. Communications Support Team Community/Neighborhood
Policing Unit Crane Crisis Intervention Unit Responder Crew/Team
Transport Detox Transport Unmanned Aerial Vehicle (Drone) Electronic
Board Fire Boat Fire Truck - Aerial (Ladder or Platform) Foam Tender,
Firefighting Foam Tank Forensic Evidence Fuel Tender (Gasoline, Diesel,
AvGas, aka Gas Tanker) Gang Enforcement HazMat Entry Team Helicopters,
Firefighting Heavy Rescue Apparatus Incident Management Team
Investigator Law Enforcement Command Unit Law Enforcement Patrol Unit
Law Enforcement School Officer Light Rescue Apparatus (Jeep and Other
Vehicles)

Reference This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document

Page 82 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value MaintenanceandRepairTeam
ParkingEnforcement MobileCommunicationCenter MobileSupportVehicle
CrowdControlTeam MountedPatrol MountainSearchandRescueTeam
MedicalTransportVehicle ParkRanger DiveTeam RapidResponseSpecialist
RoadSweeper SafetyOfficer SnowBlower SnowCat StructureFireEngineType1

StructureFireEngineType2

SWAT

10/19/2021

Literal Description Maintenance and Repair Team - Light Equipment Public
Works Parking Enforcement Mobile Communications Center (Also referred to
as "Mobile EOC") Mobile Support Vehicle Law Enforcement Crowd Control
Team Mounted Patrol Mountain Search and Rescue Team Multi-Patient
Medical Transport Vehicle Park Ranger Public Safety Dive Team Rapid
Response Specialist Road Sweeper Safety Officer/Chief Snow Blower Snow
Cat Type 1 Structure Fire Engine (1000 GPM, 400 Gal tank, 1200' 2.5"
hose or larger, 400' 1.5" hose or larger (attack line), 200' 1" hose or
larger, ladder 20' or longer, 500GPM master stream, and 4 crew members)
Type 2 Structure Fire Engine (500GPM, 400 Gal Tank, hose 1000', 5400',
300' respectively, ladder 20' or longer, no master stream and 3 crew
members) SWAT/Tactical Teams

Reference This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document This document

This document

This document

Page 83 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Value SwiftwaterRescueTeam
VehicularCrimesTeam Veterinarian WaterRescueTeam TrafficUnit
TransitOfficer DumpTruck TowTruck PlowTruck TugBoat VictimAdvocate
DeWateringWaterPump WasteWaterPump FireWaterTender WaterTruck WheelDozer
WheelLoaderBackhoe WheelLoader WheelLoaderSkidSteer WoodChipper
WildlandTaskForce StrikeTeam HotShotTeam

Literal Description Swiftwater/Flood Search and Rescue Team Vehicular
Crimes Team Veterinarian Ocean/Lake/Other Water Body Rescue Team Traffic
Unit/Officer Transit Officer/Deputy Truck, Dump Truck, Tow (Wrecker)
Truck, Plow Tug Boat (General) Victim Advocate Unit Water Pumps,
De-Watering Water Pumps, Wastewater Water Tender, Firefighting (Tanker)
Truck, Water Wheel Dozer Wheel Loader Backhoe Wheel Loaders Wheel
Loaders, Skid Steer Wood Chipper Wild Land Task Force Strike Team Hot
Shot Team

Reference This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document This document This document This document This document This
document This document This document This document

Resource Attribute The "Resource Attribute" data element is described in
Section 2.20 (Emergency Resource Data Component) of this document. IANA
is requested to add the following values to the "EIDO-ResourceAttribute"
registry.

10/19/2021

Page 84 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Table 3-12 Resource Attribute
Registry Value AdvancedEMT AdvancedPRN ReconAircraft

SearchandRescueAircraft PursuitAircraft AerialPhotographyAircraft
TransportAircraft AirBags AirMedicalTransportParamedic
AirMedicalTransportPhysician AirMedicalTransportPilot
AirMedicalTransportNurse AirportFirefighter AmbulanceStrikeTeamLeader

10/19/2021

Literal Reference Description Advanced This document Emergency Medical
Technician Advanced Practice This document Registered Nurse (Nurse
Practitioner) Aircraft, This document Reconnaissance, manhunts and
surveillance, surveys Aircraft, Search This document and Rescue
Aircraft, Pursuit This document Aircraft, Aerial This document
Photography Aircraft, Transport This document Air Bags (special This
document equipment) Air Medical This document Transport Paramedic Air
Medical This document Transport Physician Air Medical This document
Transport Pilot Air Medical This document Transport Registered Nurse
Airport Firefighter This document Ambulance Strike This document Team /
Ambulance Task Force Leader Page 85 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value AnimalCaseManager AnimalControlSpecialist
AnimalDiseaseEpidemiologist AnimalHandlingSpecialist AnimalTechnician
Bicycle CarbonMonoxideDetector CivilFieldEngineer
CrisisInterventionSpecialist ConfinedSpaceRescueEquipBasic
ConfinedSpaceRescueEquipAdvanced ConcreteCutter DiveTeam
DrugRecognitionExpert DrugK9 ElectronicArrowBoard
ElectronicVariableMessageBoard EmergencyVehicleOperator
EmergencyVehicleOperatorHeavy

10/19/2021

Literal Reference Description Animal Case This document Manager Animal
Control This document Specialist Animal Disease This document
Epidemiologist Animal Handling This document Specialist Animal
Technician This document Bicycle This document Carbon Monoxide This
document Detector Civil/Field Engineer This document Crisis Intervention
This document Specialist Confined Space This document Rescue Equipment -
Basic Confined Space This document Rescue Equipment - Advanced Concrete
Cutter This document Certified Dive Team This document Member Drug
Recognition This document Expert Drug Canine This document Electronic
Arrow This document Board Electronic Variable This document Message
Board Emergency Vehicle This document Operator Emergency Vehicle This
document Operator - Heavy Page 86 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value EMSPhysician EMT ExplosiveK9 ExtricationEquipmentBasic
ExtricationEquipmentAdvanced FireApparatusDriver_Operator FireInspector
Firefighter FLIR Floodlights JawsOfLife Generators HazMatOfficer
HelicopterSearchandRescuePilot HelicopterSearchandRescueTech
HighAngleRescueEquipment HostageNegotiator IceRescueEquipment
AbkhazianInterpreter AchineseInterpreter

10/19/2021

Literal Reference Description EMS Physician This document Emergency
Medical This document Technician Explosive Canine This document
Extrication This document Equipment - Basic Extrication This document
Equipment Advanced Fire Apparatus This document Driver/Operator Fire
Inspector This document Firefighter This document FLIR Thermal This
document Imaging Floodlights This document Hurst Tool (Jaws of This
document Life) Generators This document HazMat Officer This document
Rotary Pilot This document Fixed Wing Pilot This document High Angle
Rescue This document specialist Hostage Negotiator This document Ice
Rescue This document Equipment Language This document Interpreter
-Abkhazian (abk) Language This document Interpreter -Achinese (ace)

Page 87 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value AcoliInterpreter AdangmeInterpreter AdygheAdygeiInterpreter

AfarInterpreter AfrihiliInterpreter AfrikaansInterpreter
Afro-AsiaticLanguagesInterpreter

AinuInterpreter AkanInterpreter AkkadianInterpreter AlbanianInterpreter

10/19/2021

Literal Description Language Interpreter -- Acoli (ach) Language
Interpreter -Adangme (ada) Language Interpreter -Adyghe Adygei (ady)
Language Interpreter -- Afar (aar) Language Interpreter -Afrihili (afh)
Language Interpreter -Afrikaans (afr) Language Interpreter --
AfroAsiatic Languages (afa) Language Interpreter -- Ainu (ain) Language
Interpreter -- Akan (aka) Language Interpreter -Akkadian (akk) Language
Interpreter -Albanian (alb (B))

Reference This document This document This document

This document This document This document This document

This document This document This document This document

Page 88 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value AleutInterpreter AlgonquianLanguagesInterpreter

AltaicLanguagesInterpreter AmharicInterpreter AngikaInterpreter
ApachelanguagesInterpreter

ArabicInterpreter AragoneseInterpreter AramaicInterpreter
ArapahoInterpreter ArawakInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter --
Aleut (ale) Language This document Interpreter -Algonquian Languages
(alg) Language This document Interpreter -- Altaic Languages (tut)
Language This document Interpreter -Amharic (amh) Language This document
Interpreter -Angika (anp) Language This document Interpreter -Apache
languages (apa) Language This document Interpreter -Arabic (ara)
Language This document Interpreter -Aragonese (arg) Language This
document Interpreter -Aramaic (arc) Language This document Interpreter
-Arapaho (arp) Language This document Interpreter -Arawak (arw) Page 89
of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value ArmenianInterpreter AromanianArumanianMacedo-RomanianInterpreter

AssameseInterpreter AsturianBableLeoneseAsturleoneseInterpreter

AthapascanLanguagesInterpreter

AustronesianLanguagesInterpreter

AvaricInterpreter AvestanInterpreter AwadhiInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Armenian (arm (B)) Language This document Interpreter -Aromanian
Arumanian Macedo-Romanian (rup) Language This document Interpreter
-Assamese (asm) Language This document Interpreter -Asturian Bable
Leonese Asturleonese (ast) Language This document Interpreter
-Athapascan Languages (ath) Language This document Interpreter
-Austronesian Languages (map) Language This document Interpreter -Avaric
(ava) Language This document Interpreter -Avestan (ave) Language This
document Interpreter -Awadhi (awa)

Page 90 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value AymaraInterpreter AzerbaijaniInterpreter BalineseInterpreter
BalticLanguagesInterpreter BaluchiInterpreter BambaraInterpreter
BamilekeLanguagesInterpreter

BandaLanguagesInterpreter

BantuLanguagesInterpreter

BasaInterpreter BashkirInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter -Aymara
(aym) Language This document Interpreter -Azerbaijani (aze) Language
This document Interpreter -Balinese (ban) Language This document
Interpreter -- Baltic Languages (bat) Language This document Interpreter
-Baluchi (bal) Language This document Interpreter -Bambara (bam)
Language This document Interpreter -Bamileke Languages (bai) Language
This document Interpreter -Banda Languages (bad) Language This document
Interpreter -Bantu Languages (bnt) Language This document Interpreter --
Basa (bas) Language This document Interpreter -Bashkir (bak) Page 91 of
196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value BasqueInterpreter BatakLanguagesInterpreter
BejaBedawiyetInterpreter BelarusianInterpreter BembaInterpreter
BengaliInterpreter BerberLanguagesInterpreter

BhojpuriInterpreter BihariLanguagesInterpreter BikolInterpreter
BiniEdoInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter -Basque
(baq (B)) Language This document Interpreter -- Batak Languages (btk)
Language This document Interpreter -- Beja Bedawiyet (bej) Language This
document Interpreter -Belarusian (bel) Language This document
Interpreter -Bemba (bem) Language This document Interpreter -Bengali
(ben) Language This document Interpreter -Berber Languages (ber)
Language This document Interpreter -Bhojpuri (bho) Language This
document Interpreter -- Bihari Languages (bih) Language This document
Interpreter -- Bikol (bik) Language This document Interpreter -- Bini
Edo (bin)

Page 92 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value BislamaInterpreter BlinBilinInterpreter
BlissymbolsBlissymbolicsBlissInterpreter

BosnianInterpreter BrajInterpreter BretonInterpreter BugineseInterpreter
BulgarianInterpreter BuriatInterpreter BurmeseInterpreter
CaddoInterpreter

10/19/2021

Literal Description Language Interpreter -Bislama (bis) Language
Interpreter -- Blin Bilin (byn) Language Interpreter -Blissymbols
Blissymbolics Bliss (zbl) Language Interpreter -Bosnian (bos) Language
Interpreter -- Braj (bra) Language Interpreter -Breton (bre) Language
Interpreter -Buginese (bug) Language Interpreter -Bulgarian (bul)
Language Interpreter -Buriat (bua) Language Interpreter -Burmese (bur
(B)) Language Interpreter -Caddo (cad)

Reference This document This document This document

This document This document This document This document This document
This document This document This document

Page 93 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value CatalanValencianInterpreter

CaucasianLanguagesInterpreter

CebuanoInterpreter CelticLanguagesInterpreter
CentralAmericanIndianLanguagesInterpreter

CentralKhmerInterpreter

ChagataiInterpreter ChamicLanguagesInterpreter

ChamorroInterpreter ChechenInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Catalan Valencian (cat) Language This document Interpreter -Caucasian
languages (cau) Language This document Interpreter -Cebuano (ceb)
Language This document Interpreter -- Celtic Languages (cel) Language
This document Interpreter -Central American Indian Languages (cai)
Language This document Interpreter -Central Khmer (khm) Language This
document Interpreter -Chagatai (chg) Language This document Interpreter
-Chamic Languages (cmc) Language This document Interpreter -Chamorro
(cha) Language This document Interpreter -Chechen (che) Page 94 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Literal Description
CherokeeInterpreter Language Interpreter -Cherokee (chr)
CheyenneInterpreter Language Interpreter -Cheyenne (chy)
ChibchaInterpreter Language Interpreter -Chibcha (chb)
ChichewaChewaNyanjaInterpreter Language Interpreter -Chichewa Chewa
Nyanja (nya) ChineseInterpreter Language Interpreter -Chinese (chi (B))
ChinookJargonInterpreter Language Interpreter -Chinook Jargon (chn)
ChipewyanDeneSulineInterpreter Language Interpreter -Chipewyan Dene
Suline (chp) ChoctawInterpreter Language Interpreter -Choctaw (cho)
ChuukeseInterpreter Language Interpreter -Chuukese (chk)
ChuvashInterpreter Language Interpreter -Chuvash (chv)
ClassicalNewariOldNewariClassicalNepalBhasaInterpret Language er
Interpreter -Classical Newari Value

10/19/2021

Reference This document This document This document This document

This document This document

This document

This document This document This document This document

Page 95 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

ClassicalSyriacInterpreter

CopticInterpreter CornishInterpreter CorsicanInterpreter CreeInterpreter
CreekInterpreter CreolesPidginsEnglishBasedInterpreter

CreolesPidginsFrenchBasedInterpreter

CreolesPidginsInterpreter

CreolesPidginsPortugueseBasedInterpreter 10/19/2021

Literal Reference Description Old Newari Classical Nepal Bhasa (nwc)
Language This document Interpreter -Classical Syriac (syc) Language This
document Interpreter -Coptic (cop) Language This document Interpreter
-Cornish (cor) Language This document Interpreter -Corsican (cos)
Language This document Interpreter -- Cree (cre) Language This document
Interpreter -- Creek (mus) Language This document Interpreter -Creoles
Pidgins English Based (cpe) Language This document Interpreter -Creoles
Pidgins French Based (cpf) Language This document Interpreter -Creoles
Pidgins (crp) Language This document Interpreter -Page 96 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

CrimeanTatarCrimeanTurkishInterpreter

CroatianInterpreter CushiticLanguagesInterpreter

CzechInterpreter DakotaInterpreter DanishInterpreter DargwaInterpreter
DelawareInterpreter DinkaInterpreter DivehiDhivehiMaldivianInterpreter

10/19/2021

Literal Reference Description Creoles Pidgins Portuguese Based (cpp)
Language This document Interpreter -Crimean Tatar Crimean Turkish (crh)
Language This document Interpreter -Croatian (hrv) Language This
document Interpreter -Cushitic Languages (cus) Language This document
Interpreter -Czech (cze (B)) Language This document Interpreter -Dakota
(dak) Language This document Interpreter -Danish (dan) Language This
document Interpreter -Dargwa (dar) Language This document Interpreter
-Delaware (del) Language This document Interpreter -- Dinka (din)
Language This document Interpreter -Page 97 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

DogribInterpreter DogriInterpreter DravidianLanguagesInterpreter

DualaInterpreter DutchFlemishInterpreter

DyulaInterpreter DzongkhaInterpreter EasternFrisianInterpreter
EfikInterpreter EkajukInterpreter

10/19/2021

Literal Reference Description Divehi Dhivehi Maldivian (div) Language
This document Interpreter -Dogrib (dgr) Language This document
Interpreter -- Dogri (doi) Language This document Interpreter -Dravidian
Languages (dra) Language This document Interpreter -- Duala (dua)
Language This document Interpreter -Dutch Flemish (dut (B)) Language
This document Interpreter -- Dyula (dyu) Language This document
Interpreter -Dzongkha (dzo) Language This document Interpreter -Eastern
Frisian (frs) Language This document Interpreter -- Efik (efi) Language
This document Interpreter -Ekajuk (eka)

Page 98 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value ElamiteInterpreter EnglishInterpreter ErzyaInterpreter
EsperantoInterpreter EstonianInterpreter EweInterpreter
EwondoInterpreter FangInterpreter FantiInterpreter FaroeseInterpreter
FijianInterpreter FilipinoPilipinoInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Elamite (elx) Language This document Interpreter -English (eng)
Language This document Interpreter -- Erzya (myv) Language This document
Interpreter -Esperanto (epo) Language This document Interpreter
-Estonian (est) Language This document Interpreter -- Ewe (ewe) Language
This document Interpreter -Ewondo (ewo) Language This document
Interpreter -- Fang (fan) Language This document Interpreter -- Fanti
(fat) Language This document Interpreter -Faroese (fao) Language This
document Interpreter -- Fijian (fij) Language This document Interpreter
-Filipino Pilipino (fil) Page 99 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value FinnishInterpreter FinnoUgrianLanguagesInterpreter

FonInterpreter FrenchInterpreter FriulianInterpreter FulahInterpreter
GaelicScottishGaelicInterpreter

GaInterpreter GalibiCaribInterpreter GalicianInterpreter
GandaInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Finnish (fin) Language This document Interpreter -- Finno Ugrian
Languages (fiu) Language This document Interpreter -- Fon (fon) Language
This document Interpreter -French (fre (B)) Language This document
Interpreter -Friulian (fur) Language This document Interpreter -- Fulah
(ful) Language This document Interpreter -Gaelic Scottish Gaelic (gla)
Language This document Interpreter -- Ga (gaa) Language This document
Interpreter -- Galibi Carib (car) Language This document Interpreter
-Galician (glg) Language This document Interpreter -Ganda (lug) Page 100
of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value GayoInterpreter GbayaInterpreter GeezInterpreter
GeorgianInterpreter GermanInterpreter GilberteseInterpreter
GondiInterpreter GorontaloInterpreter GothicInterpreter GreboInterpreter
GreekModernInterpreter

10/19/2021

Literal Description Language Interpreter -- Gayo (gay) Language
Interpreter -Gbaya (gba) Language Interpreter -- Geez (gez) Language
Interpreter -Georgian (geo (B)) Language Interpreter -German (ger (B))
Language Interpreter -Gilbertese (gil) Language Interpreter -Gondi (gon)
Language Interpreter -Gorontalo (gor) Language Interpreter -Gothic (got)
Language Interpreter -Grebo (grb) Language Interpreter -Greek Modern
(gre (B))

Reference This document This document This document This document This
document This document This document This document This document This
document This document

Page 101 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value GuaraniInterpreter GujaratiInterpreter GwichinInterpreter
HaidaInterpreter HaitianHaitianCreoleInterpreter

HausaInterpreter HawaiianInterpreter HebrewInterpreter HereroInterpreter
HiligaynonInterpreter HimachaliWesternPahariLanguagesInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Guarani (grn) Language This document Interpreter -Gujarati (guj)
Language This document Interpreter -Gwich'in (gwi) Language This
document Interpreter -- Haida (hai) Language This document Interpreter
-Haitian Haitian Creole (hat) Language This document Interpreter -Hausa
(hau) Language This document Interpreter -Hawaiian (haw) Language This
document Interpreter -Hebrew (heb) Language This document Interpreter
-Herero (her) Language This document Interpreter -Hiligaynon (hil)
Language This document Interpreter -Himachali Western Pahari Languages
(him) Page 102 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value HindiInterpreter HiriMotuInterpreter HittiteInterpreter
HmongMongInterpreter

HungarianInterpreter HupaInterpreter IbanInterpreter
IcelandicInterpreter IdoInterpreter IgboInterpreter
IjolanguagesInterpreter

10/19/2021

Literal Description Language Interpreter -- Hindi (hin) Language
Interpreter -- Hiri Motu (hmo) Language Interpreter -Hittite (hit)
Language Interpreter -Hmong Mong (hmn) Language Interpreter -Hungarian
(hun) Language Interpreter -- Hupa (hup) Language Interpreter -- Iban
(iba) Language Interpreter -Icelandic (ice (B)) Language Interpreter --
Ido (ido) Language Interpreter -- Igbo (ibo) Language Interpreter
-Ijolanguages (ijo)

Reference This document This document This document This document

This document This document This document This document This document
This document This document

Page 103 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value IlokoInterpreter InariSamiInterpreter IndicLanguagesInterpreter
Indo-EuropeanLanguagesInterpreter

IndonesianInterpreter IngushInterpreter InterlingueOccidentalInterpreter

InuktitutInterpreter InupiaqInterpreter IranianLanguagesInterpreter

IrishInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter --
Iloko (ilo) Language This document Interpreter -- Inari Sami (smn)
Language This document Interpreter -- Indic Languages (inc) Language
This document Interpreter -- IndoEuropean Languages (ine) Language This
document Interpreter -Indonesian (ind) Language This document
Interpreter -Ingush (inh) Language This document Interpreter
-Interlingue Occidental (ile) Language This document Interpreter
-Inuktitut (iku) Language This document Interpreter -Inupiaq (ipk)
Language This document Interpreter -Iranian Languages (ira) Language
This document Interpreter -- Irish (gle) Page 104 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value IroquoianLanguagesInterpreter

ItalianInterpreter JapaneseInterpreter JavaneseInterpreter
Judeo-ArabicInterpreter Judeo-PersianInterpreter KabardianInterpreter
KabyleInterpreter KachinJingphoInterpreter

KalaallisutGreenlandicInterpreter

KalmykOiratInterpreter

10/19/2021

Literal Description Language Interpreter -Iroquoian Languages (iro)
Language Interpreter -Italian (ita) Language Interpreter -Japanese (jpn)
Language Interpreter -Javanese (jav) Language Interpreter -Judeo-Arabic
(jrb) Language Interpreter -Judeo-Persian (jpr) Language Interpreter
-Kabardian (kbd) Language Interpreter -Kabyle (kab) Language Interpreter
-Kachin Jingpho (kac) Language Interpreter -Kalaallisut Greenlandic
(kal) Language Interpreter -Kalmyk Oirat (xal)

Reference This document

This document This document This document This document This document
This document This document This document

This document

This document

Page 105 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value KambaInterpreter KannadaInterpreter KanuriInterpreter
Karachay-BalkarInterpreter

Kara-KalpakInterpreter KarelianInterpreter KarenLanguagesInterpreter

KashmiriInterpreter KashubianInterpreter KawiInterpreter
KazakhInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter -Kamba
(kam) Language This document Interpreter -Kannada (kan) Language This
document Interpreter -Kanuri (kau) Language This document Interpreter
-Karachay-Balkar (krc) Language This document Interpreter -- KaraKalpak
(kaa) Language This document Interpreter -Karelian (krl) Language This
document Interpreter -Karen Languages (kar) Language This document
Interpreter -Kashmiri (kas) Language This document Interpreter
-Kashubian (csb) Language This document Interpreter -- Kawi (kaw)
Language This document Interpreter -Kazakh (kaz) Page 106 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value KhasiInterpreter KhoisanLanguagesInterpreter

KhotaneseSakanInterpreter

KikuyuGikuyuInterpreter KimbunduInterpreter KinyarwandaInterpreter
KirghizKyrgyzInterpreter KlingontlhIngan-HolInterpreter

KomiInterpreter KongoInterpreter KonkaniInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter --
Khasi (kha) Language This document Interpreter -Khoisan Languages (khi)
Language This document Interpreter -Khotanese Sakan (kho) Language This
document Interpreter -Kikuyu Gikuyu (kik) Language This document
Interpreter -Kimbundu (kmb) Language This document Interpreter
-Kinyarwanda (kin) Language This document Interpreter -Kirghiz Kyrgyz
(kir) Language This document Interpreter -KlingontlhInganHol (tlh)
Language This document Interpreter -- Komi (kom) Language This document
Interpreter -Kongo (kon) Language This document Interpreter -Konkani
(kok) Page 107 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value KoreanInterpreter KosraeanInterpreter KpelleInterpreter
KruLanguagesInterpreter KuanyamaKwanyamaInterpreter

KumykInterpreter KurdishInterpreter KurukhInterpreter KutenaiInterpreter
LadinoInterpreter LahndaInterpreter

10/19/2021

Literal Description Language Interpreter -Korean (kor) Language
Interpreter -Kosraean (kos) Language Interpreter -Kpelle (kpe) Language
Interpreter -- Kru Languages (kro) Language Interpreter -Kuanyama
Kwanyama (kua) Language Interpreter -Kumyk (kum) Language Interpreter
-Kurdish (kur) Language Interpreter -Kurukh (kru) Language Interpreter
-Kutenai (kut) Language Interpreter -Ladino (lad) Language Interpreter
-Lahnda (lah)

Reference This document This document This document This document This
document

This document This document This document This document This document
This document

Page 108 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value LambaInterpreter LandDayakLanguagesInterpreter

LaoInterpreter LatinInterpreter LatvianInterpreter LezghianInterpreter
LimburganLimburgerLimburgishInterpreter

LingalaInterpreter LithuanianInterpreter LojbanInterpreter
LowerSorbianInterpreter

10/19/2021

Literal Description Language Interpreter -Lamba (lam) Language
Interpreter -- Land Dayak Languages (day) Language Interpreter -- Lao
(lao) Language Interpreter -- Latin (lat) Language Interpreter -Latvian
(lav) Language Interpreter -Lezghian (lez) Language Interpreter
-Limburgan Limburger Limburgish (lim) Language Interpreter -Lingala
(lin) Language Interpreter -Lithuanian (lit) Language Interpreter
-Lojban (jbo) Language Interpreter --

Reference This document This document

This document This document This document This document This document

This document This document This document This document

Page 109 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

LoziInterpreter Luba-KatangaInterpreter Luba-LuluaInterpreter
LuisenoInterpreter LuleSamiInterpreter LundaInterpreter
Luo-KenyaTanzania-Interpreter

LushaiInterpreter LuxembourgishLetzeburgeschInterpreter

MacedonianInterpreter

10/19/2021

Literal Reference Description Lower Sorbian (dsb) Language This document
Interpreter -- Lozi (loz) Language This document Interpreter --
LubaKatanga (lub) Language This document Interpreter -- LubaLulua (lua)
Language This document Interpreter -Luiseno (lui) Language This document
Interpreter -- Lule Sami (smj) Language This document Interpreter -Lunda
(lun) Language This document Interpreter -Luo(KenyaTanzania ) (luo)
Language This document Interpreter -Lushai (lus) Language This document
Interpreter -Luxembourgish Letzeburgesch (ltz) Language This document
Interpreter -Macedonian (mac (B)) Page 110 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value MadureseInterpreter MagahiInterpreter MaithiliInterpreter
MakasarInterpreter MalagasyInterpreter MalayalamInterpreter
MalayInterpreter MalteseInterpreter ManchuInterpreter MandarInterpreter
MandingoInterpreter ManipuriInterpreter

10/19/2021

Literal Description Language Interpreter -Madurese (mad) Language
Interpreter -Magahi (mag) Language Interpreter -Maithili (mai) Language
Interpreter -Makasar (mak) Language Interpreter -Malagasy (mlg) Language
Interpreter -Malayalam (mal) Language Interpreter -Malay (may (B))
Language Interpreter -Maltese (mlt) Language Interpreter -Manchu (mnc)
Language Interpreter -Mandar (mdr) Language Interpreter -Mandingo (man)
Language Interpreter -Manipuri (mni)

Reference This document This document This document This document This
document This document This document This document This document This
document This document This document

Page 111 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value ManoboLanguagesInterpreter

ManxInterpreter MaoriInterpreter MapudungunMapucheInterpreter

MarathiInterpreter MariInterpreter MarshalleseInterpreter
MarwariInterpreter MasaiInterpreter MayanLanguagesInterpreter

MendeInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter -Manobo
Languages (mno) Language This document Interpreter -- Manx (glv)
Language This document Interpreter -- Maori (mao (B)) Language This
document Interpreter -Mapudungun Mapuche (arn) Language This document
Interpreter -Marathi (mar) Language This document Interpreter -- Mari
(chm) Language This document Interpreter -Marshallese (mah) Language
This document Interpreter -Marwari (mwr) Language This document
Interpreter -- Masai (mas) Language This document Interpreter -Mayan
Languages (myn) Language This document Interpreter -Mende (men) Page 112
of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value MikmaqMicmacInterpreter

MinangkabauInterpreter MirandeseInterpreter MohawkInterpreter
MokshaInterpreter MongoInterpreter MongolianInterpreter
Mon-KhmerLanguagesInterpreter

MossiInterpreter MundaLanguagesInterpreter

NahuatlLanguagesInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Mi'kmaq Micmac (mic) Language This document Interpreter -Minangkabau
(min) Language This document Interpreter -Mirandese (mwl) Language This
document Interpreter -Mohawk (moh) Language This document Interpreter
-Moksha (mdf) Language This document Interpreter -Mongo (lol) Language
This document Interpreter -Mongolian (mon) Language This document
Interpreter -- MonKhmer Languages (mkh) Language This document
Interpreter -- Mossi (mos) Language This document Interpreter -Munda
Languages (mun) Language This document Interpreter -Page 113 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

NauruInterpreter NavajoNavahoInterpreter

NdebeleNorthNdebeleInterpreter

NdebeleSouthNdebeleInterpreter

NdongaInterpreter NeapolitanInterpreter NepalBhasaNewariInterpreter

NepaliInterpreter NiasInterpreter Niger-KordofanianLanguagesInterpreter

10/19/2021

Literal Reference Description Nahuatl Languages (nah) Language This
document Interpreter -Nauru (nau) Language This document Interpreter
-Navajo Navaho (nav) Language This document Interpreter -Ndebele North
Ndebele (nde) Language This document Interpreter -Ndebele South Ndebele
(nbl) Language This document Interpreter -Ndonga (ndo) Language This
document Interpreter -Neapolitan (nap) Language This document
Interpreter -- Nepal Bhasa Newari (new) Language This document
Interpreter -Nepali (nep) Language This document Interpreter -- Nias
(nia) Language This document Interpreter -Page 114 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

Nilo-SaharanLanguagesInterpreter

NiueanInterpreter NKoInterpreter NogaiInterpreter NorseOldInterpreter
NorthernFrisianInterpreter

NorthernSamiInterpreter

NorwegianBokmalInterpreter

NorwegianInterpreter NubianLanguagesInterpreter

10/19/2021

Literal Reference Description Niger-Kordofanian Languages (nic) Language
This document Interpreter -- NiloSaharan Languages (ssa) Language This
document Interpreter -Niuean (niu) Language This document Interpreter --
N'Ko (nqo) Language This document Interpreter -- Nogai (nog) Language
This document Interpreter -- Norse Old (non) Language This document
Interpreter -Northern Frisian (frr) Language This document Interpreter
-Northern Sami (sme) Language This document Interpreter -Norwegian
Bokm?l (nob) Language This document Interpreter -Norwegian (nno)
Language This document Interpreter -Page 115 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

NyamweziInterpreter NyankoleInterpreter NyoroInterpreter
NzimaInterpreter OccitanInterpreter OjibwaInterpreter OriyaInterpreter
OromoInterpreter OsageInterpreter OssetianOsseticInterpreter

OtomianLanguagesInterpreter

10/19/2021

Literal Description Nubian Languages (nub) Language Interpreter
-Nyamwezi (nym) Language Interpreter -Nyankole (nyn) Language
Interpreter -Nyoro (nyo) Language Interpreter -Nzima (nzi) Language
Interpreter -Occitan (oci) Language Interpreter -Ojibwa (oji) Language
Interpreter -- Oriya (ori) Language Interpreter -Oromo (orm) Language
Interpreter -Osage (osa) Language Interpreter -Ossetian Ossetic (oss)
Language Interpreter --

Reference

This document This document This document This document This document
This document This document This document This document This document

This document

Page 116 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

PahlaviInterpreter PalauanInterpreter PaliInterpreter
PampangaKapampanganInterpreter

PangasinanInterpreter PanjabiPunjabiInterpreter

PapiamentoInterpreter PapuanLanguagesInterpreter

PediSepediNorthernSothoInterpreter

10/19/2021

Literal Description Otomian Languages (oto) Language Interpreter
-Pahlavi (pal) Language Interpreter -Palauan (pau) Language Interpreter
-- Pali (pli) Language Interpreter -Pampanga Kapampangan (pam) Language
Interpreter -Pangasinan (pag) Language Interpreter -Panjabi Punjabi
(pan) Language Interpreter -Papiamento (pap) Language Interpreter
-Papuan Languages (paa) Language Interpreter -- Pedi Sepedi Northern
Sotho (nso)

Reference

This document This document This document This document

This document This document

This document This document

This document

Page 117 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value PersianInterpreter PhilippineLanguagesInterpreter

PhoenicianInterpreter PohnpeianInterpreter PolishInterpreter
PortugueseInterpreter PrakritLanguagesInterpreter

PushtoPashtoInterpreter

QuechuaInterpreter RajasthaniInterpreter RapanuiInterpreter

10/19/2021

Literal Description Language Interpreter -Persian (per (B)) Language
Interpreter -Philippine Languages (phi) Language Interpreter -Phoenician
(phn) Language Interpreter -Pohnpeian (pon) Language Interpreter -Polish
(pol) Language Interpreter -Portuguese (por) Language Interpreter
-Prakrit Languages (pra) Language Interpreter -Pushto Pashto (pus)
Language Interpreter -Quechua (que) Language Interpreter -Rajasthani
(raj) Language Interpreter -Rapanui (rap)

Reference This document This document

This document This document This document This document This document

This document

This document This document This document

Page 118 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value RarotonganCookIslandsMaoriInterpreter

RomanianMoldavianMoldovanInterpreter

RomanshInterpreter RomanyInterpreter RundiInterpreter RussianInterpreter
SalishanLanguagesInterpreter

SamaritanAramaicInterpreter

SamiLanguagesInterpreter SamoanInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Rarotongan Cook Islands Maori (rar) Language This document Interpreter
-Romanian Moldavian Moldovan (rum (B)) Language This document
Interpreter -Romansh (roh) Language This document Interpreter -Romany
(rom) Language This document Interpreter -- Rundi (run) Language This
document Interpreter -Russian (rus) Language This document Interpreter
-Salishan Languages (sal) Language This document Interpreter -Samaritan
Aramaic (sam) Language This document Interpreter -- Sami Languages (smi)
Language This document Interpreter -Samoan (smo) Page 119 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value SandaweInterpreter SangoInterpreter SanskritInterpreter
SantaliInterpreter SardinianInterpreter SasakInterpreter
ScotsInterpreter SelkupInterpreter SemiticLanguagesInterpreter

SerbianInterpreter SererInterpreter

10/19/2021

Literal Description Language Interpreter -Sandawe (sad) Language
Interpreter -Sango (sag) Language Interpreter -Sanskrit (san) Language
Interpreter -Santali (sat) Language Interpreter -Sardinian (srd)
Language Interpreter -Sasak (sas) Language Interpreter -- Scots (sco)
Language Interpreter -Selkup (sel) Language Interpreter -Semitic
Languages (sem) Language Interpreter -Serbian (srp) Language Interpreter
-- Serer (srr)

Reference This document This document This document This document This
document This document This document This document This document

This document This document

Page 120 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value ShanInterpreter ShonaInterpreter SichuanYiNuosuInterpreter

SicilianInterpreter SidamoInterpreter SignLanguagesInterpreter
SiksikaInterpreter SindhiInterpreter SinhalaSinhaleseInterpreter

Sino-TibetanLanguagesInterpreter

SiouanLanguagesInterpreter

10/19/2021

Literal Description Language Interpreter -- Shan (shn) Language
Interpreter -Shona (sna) Language Interpreter -Sichuan YiNuosu (iii)
Language Interpreter -Sicilian (scn) Language Interpreter -Sidamo (sid)
Language Interpreter -- Sign Languages (sgn) Language Interpreter
-Siksika (bla) Language Interpreter -Sindhi (snd) Language Interpreter
-Sinhala Sinhalese (sin) Language Interpreter -- SinoTibetan Languages
(sit) Language Interpreter --

Reference This document This document This document

This document This document This document This document This document
This document

This document

This document

Page 121 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

SkoltSamiInterpreter Slave(Athapascan)Interpreter
SlavicLanguagesInterpreter SlavicSlavonicOldBulgarianInterpreter

SlovakInterpreter SlovenianInterpreter SogdianInterpreter
SomaliInterpreter SonghaiLanguagesInterpreter

SoninkeInterpreter SorbianLanguagesInterpreter 10/19/2021

Literal Reference Description Siouan Languages (sio) Language This
document Interpreter -- Skolt Sami (sms) Language This document
Interpreter -- Slave - Athapascan (den) Language This document
Interpreter -- Slavic Languages (sla) Language This document Interpreter
-- Slavic Slavonic Old Bulgarian (chu) Language This document
Interpreter -Slovak (slo (B)) Language This document Interpreter
-Slovenian (slv) Language This document Interpreter -Sogdian (sog)
Language This document Interpreter -Somali (som) Language This document
Interpreter -Songhai Languages (son) Language This document Interpreter
-Soninke (snk) Language This document Interpreter -Page 122 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

SothoSouthernInterpreter SouthAmericanIndianLanguagesInterpreter

SouthernAltaiInterpreter SouthernSamiInterpreter

SpanishCastilianInterpreter

SrananTongoInterpreter StandardMoroccanTamazightInterpreter

SukumaInterpreter SumerianInterpreter

10/19/2021

Literal Reference Description Sorbian Languages (wen) Language This
document Interpreter -- Sotho Southern (sot) Language This document
Interpreter -South American Indian Languages (sai) Language This
document Interpreter -Southern Altai (alt) Language This document
Interpreter -Southern Sami (sma) Language This document Interpreter
-Spanish Castilian (spa) Language This document Interpreter -Sranan
Tongo (srn) Language This document Interpreter -Standard Moroccan
Tamazight (zgh) Language This document Interpreter -Sukuma (suk)
Language This document Interpreter -Sumerian (sux)

Page 123 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value SundaneseInterpreter SusuInterpreter SwahiliInterpreter
SwatiInterpreter SwedishInterpreter
SwissGermanAlemannicAlsatianInterpreter

SyriacInterpreter TagalogInterpreter TahitianInterpreter
TaiLanguagesInterpreter TajikInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Sundanese (sun) Language This document Interpreter -- Susu (sus)
Language This document Interpreter -Swahili (swa) Language This document
Interpreter -- Swati (ssw) Language This document Interpreter -Swedish
(swe) Language This document Interpreter -- Swiss German Alemannic
Alsatian (gsw) Language This document Interpreter -Syriac (syr) Language
This document Interpreter -Tagalog (tgl) Language This document
Interpreter -Tahitian (tah) Language This document Interpreter -- Tai
Languages (tai) Language This document Interpreter -- Tajik (tgk)

Page 124 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value TamashekInterpreter TamilInterpreter TatarInterpreter
TeluguInterpreter TerenoInterpreter TetumInterpreter ThaiInterpreter
TibetanInterpreter TigreInterpreter TigrinyaInterpreter TimneInterpreter
TivInterpreter

10/19/2021

Literal Reference Description Language This document Interpreter
-Tamashek (tmh) Language This document Interpreter -- Tamil (tam)
Language This document Interpreter -- Tatar (tat) Language This document
Interpreter -Telugu (tel) Language This document Interpreter -Tereno
(ter) Language This document Interpreter -Tetum (tet) Language This
document Interpreter -- Thai (tha) Language This document Interpreter
-Tibetan (tib (B)) Language This document Interpreter -- Tigre (tig)
Language This document Interpreter -Tigrinya (tir) Language This
document Interpreter -Timne (tem) Language This document Interpreter --
Tiv (tiv) Page 125 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value TlingitInterpreter TokelauInterpreter TokPisinInterpreter
TongaNyasaInterpreter Tonga-TongaIslands-Interpreter

TsimshianInterpreter TsongaInterpreter TswanaInterpreter
TumbukaInterpreter TupiLanguagesInterpreter TurkishInterpreter

10/19/2021

Literal Description Language Interpreter -Tlingit (tli) Language
Interpreter -Tokelau (tkl) Language Interpreter -- Tok Pisin (tpi)
Language Interpreter -Tonga Nyasa (tog) Language Interpreter -Tonga
(Tonga Islands) (ton) Language Interpreter -Tsimshian (tsi) Language
Interpreter -Tsonga (tso) Language Interpreter -Tswana (tsn) Language
Interpreter -Tumbuka (tum) Language Interpreter -- Tupi Languages (tup)
Language Interpreter -Turkish (tur)

Reference This document This document This document This document This
document

This document This document This document This document This document
This document

Page 126 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value TurkmenInterpreter TuvaluInterpreter TuvinianInterpreter
TwiInterpreter UdmurtInterpreter UgariticInterpreter
UighurUyghurInterpreter

UkrainianInterpreter UmbunduInterpreter UpperSorbianInterpreter

UrduInterpreter

10/19/2021

Literal Description Language Interpreter -Turkmen (tuk) Language
Interpreter -Tuvalu (tvl) Language Interpreter -Tuvinian (tyv) Language
Interpreter -- Twi (twi) Language Interpreter -Udmurt (udm) Language
Interpreter -Ugaritic (uga) Language Interpreter -Uighur Uyghur (uig)
Language Interpreter -Ukrainian (ukr) Language Interpreter -Umbundu
(umb) Language Interpreter -Upper Sorbian (hsb) Language Interpreter --
Urdu (urd)

Reference This document This document This document This document This
document This document This document

This document This document This document

This document

Page 127 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value UzbekInterpreter VaiInterpreter VendaInterpreter
VietnameseInterpreter VolapukInterpreter VoticInterpreter
WakashanLanguagesInterpreter

WalloonInterpreter WarayInterpreter WashoInterpreter WelshInterpreter
WesternFrisianInterpreter 10/19/2021

Literal Description Language Interpreter -Uzbek (uzb) Language
Interpreter -- Vai (vai) Language Interpreter -Venda (ven) Language
Interpreter -Vietnamese (vie) Language Interpreter -Volap?k (vol)
Language Interpreter -- Votic (vot) Language Interpreter -Wakashan
Languages (wak) Language Interpreter -Walloon (wln) Language Interpreter
-Waray (war) Language Interpreter -Washo (was) Language Interpreter
-Welsh (wel (B)) Language Interpreter --

Reference This document This document This document This document This
document This document This document

This document This document This document This document This document

Page 128 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value

WolaittaWolayttaInterpreter

WolofInterpreter XhosaInterpreter YakutInterpreter YaoInterpreter
YapeseInterpreter YiddishInterpreter YorubaInterpreter
YupikLanguagesInterpreter ZandeLanguagesInterpreter

10/19/2021

Literal Reference Description Western Frisian (fry) Language This
document Interpreter -Wolaitta Wolaytta (wal) Language This document
Interpreter -- Wolof (wol) Language This document Interpreter -Xhosa
(xho) Language This document Interpreter -- Yakut (sah) Language This
document Interpreter -- Yao (yao) Language This document Interpreter
-Yapese (yap) Language This document Interpreter -Yiddish (yid) Language
This document Interpreter -Yoruba (yor) Language This document
Interpreter -- Yupik Languages (ypk) Language This document Interpreter
-Zande Languages (znd)

Page 129 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value ZapotecInterpreter ZazaDimiliDimlKirdkiKirmanjkiInterpreter

ZenagaInterpreter ZhuangChuangInterpreter

ZuluInterpreter ZuniInterpreter LoJack CrowdControlSpecialist
MountainSearchandRescueTech NarcoticsIdentificationKitOperator Paramedic
PortablePumpFloating PortablePumpSubmersible 10/19/2021

Literal Reference Description Language This document Interpreter
-Zapotec (zap) Language This document Interpreter -- Zaza Dimili Diml
Kirdki Kirmanjki (zza) Language This document Interpreter -Zenaga (zen)
Language This document Interpreter -Zhuang Chuang (zha) Language This
document Interpreter -- Zulu (zul) Language This document Interpreter --
Zuni (zun) LoJack? Stolen This document Vehicle Recovery System Crowd
Control This document Specialist Mountain Search This document and/or
Rescue Technical Specialist Narcotics This document Identification Kit
Operator Paramedic This document Portable Pump -This document Floating
Portable Pump -This document Submersible Page 130 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Value PublicInformationOfficer Telecommunicator RadiologicTech
RadiationDetector RegisteredNurse RescueRope SCBA Shotgun
StructuralCollapaseRescueTech StructuralEngineer SWAT
SwiftwaterRescueEquipment SwiftwaterRescueTech TaserTrained
ThermalCamera VacuumWater WaterCanineSearchTech WildLandFireTools

10/19/2021

Literal Reference Description Public Information This document Officer
Public Safety This document Telecommunicator Radiologic This document
Technician Victoreen Radiation This document Detector Registered Nurse
This document Rescue Rope This document Self-Contained This document
Breathing Apparatus Shotgun This document Structural Collapse This
document Rescue Technician Structural Engineer This document SWAT
Officer This document Swift Water Rescue This document Equipment Swift
Water/Flood This document Rescue Technical Specialist Taser Trained This
document Thermal Imaging This document Camera Vacuum / Water This
document Water Canine This document Search Technician Wildland Fire This
document Fighting Tools

Page 131 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Impacts, Considerations, Abbreviations, Terms, and Definitions The FEs
exchanging the data defined in the EIDO may be physically or virtually
connected to each other. The FEs may belong to one or more disparate
manufacturer's systems located within the same public safety agency or
within different agencies. Operations Impacts Summary A standardized
format for electronically exchanging emergency Incident information will
provide stakeholders with many operational benefits. These benefits are
similar to those realized when agencies implement a local or regional
"CAD to CAD" exchange. Use of the EIDO at a national, regional, and
local level and within communication centers themselves is critical to
the implementation of NG9-1-1 (i3) compliant systems \[3\] . Technical
Impacts Summary As with the implementation of any technical standard,
the final EIDO will have significant impact. Initially, all FE's
involved in the exchange of emergency Incident information (e.g. call
handling, logging, dispatch, etc.) will need to be modified to comply
with EIDO transactions. Minimal impacts are expected on the ESInet or
other E9-1-1 IP based networks however they will need to be compatible
with the EIDO structure in order to carry EIDOs from one FE to another.
Security Impacts Summary EIDOs carry confidential information and they
must be transmitted over secure transports such as TLS (Transport Layer
Security) protected TCP (Transmission Control Protocol), and only
to/from i3 authenticated FEs and systems \[3\] . EIDOs accessible to
authenticated FEs and other systems must have their contents filtered to
contain only data authorized to be transmitted to those systems by the
data owner's policy. Furthermore, FEs and systems that pass EIDOs or
their contents along to other FEs and systems must filter those EIDO
instances based on the authentication of the receiving FEs and other
systems to contain only data authorized to be transmitted to those FEs
and systems by the data owner's policy. Recommendation for Additional
Development Work Future versions may expand the EIDO to include new data
components and data elements, provide additional allowable registry
values, or modified to support additional emergency Incident related
exchanges such as the transfer of a patient's medical diagnostics,
administered procedures, and medical status between transporting
ambulances and receiving hospitals. 10/19/2021

Page 132 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Table 4-1 Additonal Development Work Section(s) 2 2.17, Appendix A 2,
Appendix A Appendix A 2 TBD Apppendix A TBD 2 3.1 2.7,2.8 2.9 2.x

Reference to future work Add radio channels/push to talk information
Include NIEM Commercial vehicle type Adding NIEM Injury treatment type
Conformance with future NIEM JSON version Add weapon information.
Consider using JSON-LD context for NIEM compliance. Add a floor type
that include floor integers. Add a mapping between EIDO and NIEM
em:EmergencyIncidentType Add the concept of automaton identifier Make
TransferredTo Agency role transient Add reason code for split/merge and
Link components Add information about how an incident is originally
reported. Add URL for media related to an Incident.

Anticipated Timeline EIDO is an integral piece of NG9-1-1 and must be
considered during implementation planning. Early implementations are
anticipated considering the migration to NG9-1-1 in Canada. Cost Factors
Significant consideration was given to the cost impact of proposing yet
another "change" to 9-1-1 stakeholders. Standardizing the format of
Incident information exchange, however, is critical for NG9-1-1's true
potential to be realized \[3\]. Absent this standard, Incident
information data exchanges will continue to rely on limited, proprietary
implementations that are costly and offer limited capability to share
Incident information both within and between communication centers. It
should be noted that implementation of the EIDO IEPD ANS, as well as
most NG9-1-1 systems, may require updates to existing public safety
systems resident at communication centers. Cost Recovery Considerations
Normal business practices shall be assumed to be the cost recovery
mechanism. 10/19/2021

Page 133 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Additional Impacts (non-cost
related) The information or requirements contained in this document are
not expected to have additional impacts, based on the analysis of the
authoring group. Abbreviations, Terms, and Definitions See NENA Master
Glossary of 9-1-1 Terminology, NENA-ADM-000 \[1\], for a complete
listing of terms used in NENA documents. All abbreviations used in this
document are listed below, along with any new or updated terms and
definitions.

10/19/2021

Page 134 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Term or Abbreviation (Expansion)
Data Component Data Element

Definition / Description A structure composed of multiple Data Elements.
A unit of information composed of a name and value. The value can be a
number, a boolean, a character string or a Data Component.

A JSON-based (JavaScript Object Notation) object that is used to share
emergency Incident information between and EIDO (Emergency Incident
among authorized entities and systems. NENA has adopted Data Object) the
JSON-based EIDO for sharing Incident information among authorized
NG9-1-1 entities and systems.

ESInet (Emergency Services IP Network)

An ESInet is a managed IP network that is used for emergency services
communications, and which can be shared by all public safety agencies.
It provides the IP transport infrastructure upon which independent
application platforms and core services can be deployed, including, but
not restricted to, those necessary for providing NG9-1-1 services.
ESInets may be constructed from a mix of dedicated and shared
facilities. ESInets may be interconnected at local, regional, state,
federal, national and international levels to form an IP-based
inter-network (network of networks). The term ESInet designates the
network, not the services that ride on the network. See NGCS (Next
Generation Core Services).

FE (Functional Element)

A set of software features that may be combined with hardware interfaces
and operations on those interfaces to accomplish a defined task.

FRC (Field Responder Client)

A device or software installed on a device such as computer, tablet or
smartphone used by field responders to access public safety systems,
such as to receive dispatch instructions or perform database lookups.
FRCs are often colloquially referred to as "MDT" or "Mobile Data
Terminals".

10/19/2021

Page 135 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Term or Abbreviation (Expansion)

Definition / Description

IMR (Interactive Media Response)

An automated service used to play announcements, record responses and
interact with callers using any or all of audio, video and text.

MDT (Mobile Data Terminal)

A computerized device used in a vehicle to exchange information between
an end user and a communications center. MDTs are a specific type of FRC
(Field Responder Client).

Report Number

A number associated with an emergency Incident that once generated
indicates that a follow up report or investigation will be associated
with the Incident. May also be known as: case number, fire Incident
number, department report numbers and other variations depending on the
types of service agencies involved and local customs. Typically, the
report number is sequential within a year and also identifies the agency
for which it is issued. Within an agency, the report number is globally
unique. More than one report number may be associated with a single
Incident.

SIP (Session Initiation Protocol)

A protocol specified by the IETF (RFC 3261) that defines a method for
establishing multimedia sessions over the Internet. Used as the call
signaling protocol in VoIP, NENA i2 and NENA i3.

A communications protocol linking different computer TCP (Transmission
Control platforms across networks. TCP/IP functions at the 3rd and
Protocol) 4th levels of the open system integration model. TLS
(Transport Layer Security)

An IETF standard protocol that provides authentication, integrity
protection, and confidentiality using mutually negotiated cryptographic
capabilities and digital certificates on a connection such as TCP.

URL (Uniform Resource Locator)

A URL is a type of URI specifically used for describing and navigating
to a resource (e.g. http://www.nena.org)

10/19/2021

Page 136 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

References

\[11\]

National Emergency Number Association. Master Glossary of 9-1-1
Terminology. NENA-ADM-000.23-2020. Arlington, VA: NENA, approved January
20, 2020. Internet Engineering Task Force. Key words for use in RFCs to
Indicate Requirement Levels. S. Bradner. RFC 2119, March 1997. National
Emergency Number Association. NENA Detailed Functional and Interface
Standards for the NENA i3 Solution. NENA-STA-010.2-2016 (originally
08-003). Arlington, VA: NENA, approved September 10, 2016. National
Emergency Number Association. NENA i3 Standard for Next Generation
9-1-1. NENA-STA-010.3-202Y. Arlington, VA: NENA (forthcoming). NIEM
Program Management Office. Introduction to the National Information
Exchange Model (NIEM), version 0.3. Issued February 12, 2007.
https://reference.niem.gov/niem/guidance/introduction/0.3/niem-introduction0.3.pdf
Internet Engineering Task Force. Additional Data Related to an Emergency
Call. R. Gellens, B. Rosen, H. Tschofenig, R. Marshall, J. Winterbottom.
RFC 7852 Internet Engineering Task Force. A Uniform Resource Name (URN)
for Emergency and Other Well-Known Services. H. Schulzrinne. RFC 5031,
January 2008. Association of Public-Safety Communications Officials.
Public Safety Communications Common Incident Types for Data Exchange.
APCO ANS 2.103.22019. Daytona Beach, FL: APCO, approved 2019. Internet
Engineering Task Force. Presence Information Data Format (PIDF). H.
Sugano, S. Fujimoto, G. Klyne, A. Bateman, W. Carr, and J. Peterson. RFC
3863, August 2004. Internet Engineering Task Force. A Presence-based
GEOPRIV Location Object Format. J. Peterson. RFC 4119, December 2005.
Internet Engineering Task Force. GEOPRIV Presence Information Data
Format

\[12\]

March 2009. Association of Public-Safety Communications Officials. Alarm
Monitoring Companyto

\[1\] \[2\]

\[3\] \[4\] \[5\]

\[6\] \[7\] \[8\]

\[9\]

\[10\]

Location Object (PIDF-LO) Usage Clarification, Considerations, and
Recommendations. J. Winterbottom, M. Thomson, and H. Tschofenig. RFC
5491,

Public Safety Answering Point (PSAP) Computer-Aided Dispatch (CAD)
Automated Secure Alarm Protocol (ASAP). APCO/CSAA ANS 2.101.2-2014.
Daytona Beach, FL: APCO, approved August 5, 2014.

10/19/2021

Page 137 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Exhibits Not Applicable.

Appendix A Required NIEM data element that need to be supported for EIDO
compliance The Person data component and the Vehicle data component both
contains their respective NIEM data components. The NIEM PersonType and
VehicleType also contains other NIEM data components. In order for
implementations to be interoperable, a subset of these NIEM data
elements and data components must be supported by every EIDO
implementation. In this context support means that the data element must
be preserved when transmitting, receiving, storing or retrieving an
EIDO. To be compliant with this standard any EIDO implementation must
support the following. A.1

Person Type Table A-1 Person Type

JSON Name ageMeasure

Type Description Comment nc:TimeMeasu A measurement reType of the age of
a person. birthDate nc:PersonBirthDate nc:DateType A date a person was
born. birthLocation nc:PersonBirthLocation nc:LocationTy A location pe
where a person was born. bloodTypeCode j:PersonBloodTypeCode
ncic:BLTCodeT A blood group ype and RH factor of a person. buildCode
j:PersonBuildCode ndex:PersonB A person's uildCodeType physique or body
shape. citizenshipIso3166A iso_3166:Coun A country that lpha2Code
nc:PersonCitizenshipISO31 tryAlpha2Code assigns rights, 66Alpha2Code
Type duties, and privileges to a 10/19/2021

NIEM Data element nc:PersonAgeMeasure

Page 138 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

clothing

complexion

description digitalImage

disguise

ethnicityCode eyeColorCode eyewearCode

10/19/2021

NIEM Data element

Type

Description Comment person because of the birth or naturalization of the
person in that country. j:PersonClothing j:ClothingType An article of
clothing, dress, or attire for a person. nc:PersonComplexionText
nc:TextType An appearance clear, or condition of freckled, the skin of a
wrinkled person. nc:PersonDescriptionText nc:TextType A description of a
person. nc:PersonDigitalImage nc:ImageType A photograph or image of a
person in a digital format. nc:PersonDisguiseDescript nc:TextType A
description of wig, mask, ionText something a glasses, person wears to
uniform conceal or mislead others as to the true appearance or identity
of that person. j:PersonEthnicityCode ucr:EthnicityC A cultural odeType
lineage of a person. j:PersonEyeColorCode ncic:EYECodeT A color of the
ype eyes of a person. j:PersonEyewearCode ndex:PersonEy A description of
ewearCodeTyp glasses or other e Page 139 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Description Comment eyewear a person wears. facialHairCode
j:PersonFacialHairCode ndex:PersonFa A kind of facial cialHairCodeTy
hair of a person. pe generalAppearance nc:PersonGeneralAppeara
nc:TextType A description of nceDescriptionText the way a person looks
and is presented overall. hairAppearance nc:PersonHairAppearance
nc:TextType An overall Text appearance of the hair of a person.
hairCategory nc:PersonHairCategoryTex nc:TextType A kind of hair of t a
person, such as wavy or straight. hairColorCode j:PersonHairColorCode
ncic:HAIRCode A color of the Type hair of a person. hairLengthCode
j:PersonHairLengthCode ndex:PersonH A length of hair airLengthCode of a
person. Type hairStyleCode j:PersonHairStyleCode ndex:PersonH A style or
cut of airStyleCodeTy hair worn by a pe person. handedness
nc:PersonHandednessText nc:TextType A hand with left, right, which a
person ambidextro is more adept us using. height nc:PersonHeightMeasure
nc:LengthMea A measurement sureType of the height of a person. injury
nc:PersonInjury nc:InjuryType A form of physical harm or damage
10/19/2021

NIEM Data element

Type

Page 140 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

jewelry languageIsEnglish

license

isLiving

unionStatusText

mentalState name 10/19/2021

NIEM Data element

Type

Description Comment sustained by a person. nc:PersonJewelryDescripti
nc:TextType A description of onText adornments a person wears.
nc:PersonLanguageEnglish niemTrue if a person Indicator xs:boolean
understands and speaks English; false otherwise.
nc:PersonLicenseIdentifica nc:Identificati An identification This may
tion onType that references be granted a license to certify a
certification or professiona registration of a l person for some
occupation purpose. or skill. nc:PersonLivingIndicator niemTrue if a
person This xs:boolean is alive; false if indicator a person is may be
dead. useful if a person is known to be dead but the death date is
unknown. nc:PersonUnionStatusText nc:TextType A legal status of a union
between two people. nc:PersonMentalStateText nc:TextType A mental
condition of a person. nc:PersonName nc:PersonNam A combination eType of
names and/or Page 141 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Description Comment titles by which a person is known.
nationalIdentificatio nc:PersonNationalIdentific nc:Identificati An
identification n ation onType that references a person within a country
but is not based on fingerprint. nationalityCode
nc:PersonNationalityISO31 iso_3166:Coun A country of a 66Alpha2Code
tryAlpha2Code person's Type citizenship or a country in which a person
is deemed a national. physicalFeature nc:PersonPhysicalFeature
nc:PhysicalFea A prominent or tureType easily identifiable aspect of a
person. primaryLanguage nc:PersonPrimaryLanguag nc:PersonLang A capacity
of a native e uageType person for a language. language with This can
which that also be person has the used for strongest bi-lingual
familiarity. or multilingual people. raceCode j:PersonRaceCode
ncic:RACECod A classification eType of a person based on factors such as
geographical locations and genetics. 10/19/2021

NIEM Data element

Type

Page 142 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name sexCode

Type Description Comment ncic:SEXCodeT A gender or sex ype of a person.
sexualOrientation nc:PersonSexualOrientatio nc:TextType A target gender
nText of the sexual interest of a person. skinToneCode
j:PersonSkinToneCode ncic:SKINCode A color or tone Type of the skin of a
person. ssnIdentification nc:PersonSSNIdentificatio nc:Identificati A
unique In the US n onType identification and reference to a Canada,
living person; this is a 9assigned by a digit country. numeric
identifier. ssn, SSN, social security number, social insurance number
weight nc:PersonWeightMeasure nc:WeightMea A measurement sureType of the
weight of a person. nationalityIso3166A nc:PersonNationalityISO31
iso_3166:Coun A country of a lpha2Code 66Alpha2Code tryAlpha2Code
person's Type citizenship or a country in which a person is deemed a
national. accent nc:PersonAccentText nc:TextType A manner of
pronunciation; a way of pronouncing 10/19/2021

NIEM Data element j:PersonSexCode

Page 143 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Description Comment words that may indicate the place of origin or
social background of the speaker. medicalCondition
nc:PersonMedicalCondition nc:MedicalCon A state of health ditionType for
a person, on-going or present. medicalDescription
nc:PersonMedicalDescripti nc:TextType A description of onText the
overall health of a person. medicationRequired nc:PersonMedicationRequi
nc:TextType A medication redText and dosage required for a person. mood
nc:PersonMoodDescription nc:TextType A description of temper, Text a
state of humor feeling of a person. physicalDisability
nc:PersonPhysicalDisability nc:TextType A physical Text disability of a
person. emergencyContact nc:PersonEmergencyConta nc:ContactInfo A means
of ctInformation rmationType contacting someone in the event of an
emergency.

10/19/2021

NIEM Data element

Type

Page 144 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A.2

Vehicle Type Table A-2 Vehicle Type

JSON Name condition

Type Description Comment nc:TextType A state or new, used, appearance of
damaged an item. description nc:ItemDescriptionText nc:TextType A
description of an item. owner nc:ItemOwner nc:EntityType An entity which
property owns a owner property item. ownerAppliedID
nc:ItemOwnerAppliedID niemAn identifier A number on xs:string applied to
an a car to item by the identify it owner. within a fleet of cars.
serialIdentification nc:ItemSerialIdentification nc:Identificati An
onType identification inscribed on or attached to a part, collection of
parts, or complete unit by the manufacturer. category
nc:ItemCategoryText nc:TextType A kind of property item.
colorPrimaryCode j:ConveyanceColorPrimary ncic:VCOCode A single, Code
Type upper-most, front-most, or majority color of a vehicle.
colorSecondaryCod j:ConveyanceColorSecond ncic:VCOCode A lower-most e
aryCode Type or rear-most color of a twotone vehicle or a lesser color
10/19/2021

Data element nc:ItemConditionText

Page 145 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Data element

Type

Description Comment of a multicolored vehicle. image nc:ItemImage
nc:ImageTyp A binary e representation of an image of an item. makeName
nc:ItemMakeName nc:ProperNa A name of the meTextType manufacturer that
produced an item. modelName nc:ItemModelName nc:ProperNa A name of a
meTextType specific design or kind of item made by a manufacturer.
modelYear nc:ItemModelYearDate niemA year in which xs:gYear an item was
manufactured or produced. styleCode j:VehicleStyleCode ncic:VSTCode A
style of a Type vehicle. registrationPlateIde j:ConveyanceRegistrationP
nc:Identificati An ntification lateIdentification onType identification
on a metal plate fixed to a conveyance. isMotorized

nc:ConveyanceMotorizedIn niemdicator xs:boolean

isTowed

nc:ConveyanceTowedIndic niemator xs:boolean

10/19/2021

True if a conveyance is powered by a motor; false otherwise. True if a
conveyance is towed; false otherwise. Page 146 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name useAnsiD20Code isCMV

doorQuantity identification

makeCode modelCode

isRental

towerAssociation

vina

10/19/2021

Data element Type Description Comment j:VehicleUseANSID20Code
aamva_d20:V A manner or ehicleUseCod way in which a eType vehicle is
used. nc:VehicleCMVIndicator niemTrue if a xs:boolean vehicle is a
commercial motor vehicle; false otherwise. nc:VehicleDoorQuantity niemA
number of xs:nonNegati doors on a veInteger vehicle.
nc:VehicleIdentification nc:Identificati A unique VIN, VIN onType
identification number, for a specific registration, vehicle. plate
j:VehicleMakeCode ncic:VMACod A manufacturer eType of a vehicle.
j:VehicleModelCode ncic:VMOCod A specific eType design or class of
vehicle made by a manufacturer. nc:ItemRentalIndicator niemTrue if an
item Applies to xs:boolean is rented; false short term otherwise. rental
service only. nc:VehicleTowerAssociatio nc:ItemEntity An association n
AssociationTy between a pe vehicle and the entity that towed it.
j:VehicleVINAText nc:TextType A Vehicle Identification Number Analysis;
a combination of a vehicle make Page 147 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

isWanted

A.3

Data element

j:VehicleWantedIndicator

Type

niemxs:boolean

Description Comment and model information. True if a vehicle is being
searched for by law enforcement; false otherwise.

Date Type Table A-3 Date Type

JSON Name date dateTime fiscalYear

yearMonth marginOfError

10/19/2021

Data element nc:Date nc:DateTime

Type Description Comment niem-xs:date A full date. niemA full date and
xs:dateTime time. nc:FiscalYearDate nc:FiscalYearD A year of a ateType
twelve month period that does not necessarily correspond to the calendar
year. nc:YearMonthDate niemA year and xs:gYearMont month. h
nc:DateMarginOfErrorDura niemA subjective tion xs:duration assessment of
the uncertainty of an estimated point by bounding an elements value with
an estimated Page 148 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 margin of error. A.4

Location Type Table A-4 Locataion Type

JSON Name address

Data element Type Description Comment nc:Address nc:Address A postal
location to which Type paper mail can be directed. addressGrid
nc:AddressGrid nc:Address A location identified by a unit GridType of a
grid system overlaid on an area. area nc:LocationArea nc:AreaTyp A
location identified by e geographic boundaries. wgs84LocationCyl
mo:WGS84Locati mo:WGS84 A location identified by a inder onCylinder
LocationCyl cylinder oriented vertically and inderType centered on a
point described with WGS84 coordinates. If it is appropriate for the
radius and half-height properties to represent an error value (for
example, because the event is a laser-designated target), then the true
event location follows a normal distribution such that the cylinder
defines the one-sigma (p≈0.67)(almost equal to) deviation. (A cylinder
with twice the volume would be the two-sigma (p≈0.95)(almost equal to)
deviation, etc.) Otherwise the cylinder encloses the full physical
extent of the event. categoryCodeText it:LocationCatego nc:TextTyp A
category of operation ryCodeText e performed at a given Location
10/19/2021

Page 149 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name categoryText

Data element Type Description Comment nc:LocationCateg nc:TextTyp A kind
or functional description oryText e of a location. contact
nc:LocationConta nc:Contact A set of contact information ctInformation
Informatio for a location. nType addressCrossStre nc:AddressCross
nc:CrossStr A location identified by two or NIEM uses et Street eetType
more streets which intersect. the term CrossStreet as a location but
this document uses it to denote proximity to a location. description
nc:LocationDescr nc:TextTyp A description of a location. iptionText e
locationLandmark nc:LocationLand nc:TextTyp A distinguishing physical
markText e feature at a location. name nc:LocationName nc:ProperN A name
of a location. ameTextTy pe floorNumber nc:LocationFloor nc:TextTyp A
floor number of a location. NumberText e uncertainty mo:LocationUnce
nc:AreaTyp An area of uncertainty about a rtaintyArea e location. A.5

Image Type Table A-5 Image Type

JSON Name Data element id nc:BinaryID isAvailable

Type niemxs:string nc:BinaryAvailableInd niemicator xs:boolean

10/19/2021

Description Comment An identifier that references a binary object. True
if a binary is available; false if it is not. Page 150 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name Data element Type
Description captureDate nc:BinaryCaptureDat nc:DateType A date on which
a binary e object is captured or created.

Comment This is the date an image was taken, an audio file was recorded,
...

description

nc:BinaryDescription nc:TextType A description of a binary Text object.
format nc:BinaryFormatText nc:TextType A file format or content type of
a binary object. uri nc:BinaryURI niemA URL or file reference of a
xs:anyURI binary object. size nc:BinarySizeValue nc:NonNegati A size of
a binary object in veDecimalTyp kilobytes. e categoryText
nc:BinaryCategoryTe nc:TextType A kind of object that has mug shot, xt
been encoded. driver license picture, audio confession height
nc:ImageHeightValue niemA height of an image in xs:nonNegati pixels.
veInteger width nc:ImageWidthValue niemA width of an image in
xs:nonNegati pixels. veInteger A.6

Length Measure Type Table A-6 Length Measure Type

JSON Name Data element Type doubleValue nc:MeasureDoubleVal niemue
xs:double 10/19/2021

Description A double measurement value.

Comment

Page 151 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name Data element Type
Description Comment doubleRange nc:MeasureDoubleRa nc:DoubleRa A double
measurement nge ngeType range. textList nc:MeasureTextList nc:StringList
A list of measurement values Type , all using the same measurement
method/device and of the same units. unitText nc:MeasureUnitText
nc:TextType A unit that qualifies the measurement value. method
nc:MeasureMethodTe nc:TextType A method used to make a xt measurement.
unitCode nc:LengthUnitCode unece:Lengt A unit of measure of a length
hCodeType value. A.7

Injury Type Table A-7 Injury Type

JSON Name Data element Type Description Comment categoryCode
j:InjuryCategoryCode nc:TextType A kind of bodily harm or injury. date
nc:InjuryDate nc:DateType A date on which an injury In EIDO occurred.
InjuryDate should include DateRange/En dDate and DateRange/St artDate
description nc:InjuryDescriptionT nc:TextType A description of an
injury. ext locationCode j:InjuryLocationNDEX nc:TextType A place on the
body of a Code person where an injury occurred. severity
nc:InjurySeverityText nc:TextType A degree of the seriousness or
intensity of an injury.

10/19/2021

Page 152 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A.8

Organization Type Table A-8 Organizaation Type

JSON Name activity

Data element Type Description nc:OrganizationActi nc:TextType An
activity that an vityText organization is known or thought to be
involved with. branchName nc:OrganizationBra nc:TextType A name of the
chapter or nchName branch by which an organization is known within a
larger group of organizations. dayContact nc:OrganizationDa nc:ContactIn
A means of contacting an yContactInformatio formationTy organization
during n pe daytime hours. doingBusinessA nc:OrganizationDoi nc:TextType
A name an organization sName ngBusinessAsName uses for conducting
business. emergencyCont nc:OrganizationEm nc:ContactIn A means of
contacting an act ergencyContactInf formationTy organization in the
event ormation pe of an emergency. establishedDate nc:OrganizationEst
nc:DateType A date an organization ablishedDate was started.
eveningContact nc:OrganizationEve nc:ContactIn A means of contacting an
ningContactInform formationTy organization during ation pe evening or
early night hours. identification nc:OrganizationIde nc:Identificat An
identification that ntification ionType references an organization.

10/19/2021

Comment Ex: law enforcement, supervision

For a school, this would be a school identifier, for a lien holder, this
would be a lien holder identifier, for a court, this would be a

Page 153 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Data element

Type

Description

Comment court identifier.

localIdentificati nc:OrganizationLoc nc:Identificat An identification
assigned on alIdentification ionType at a local level to an
organization. location nc:OrganizationLoc nc:LocationT A location of an
ation ype organization. name nc:OrganizationNa nc:TextType A name of an
me organization. nightContact nc:OrganizationNig nc:ContactIn A means of
contacting an htContactInformati formationTy organization during lateon
pe night hours. otherIdentificati nc:OrganizationOth nc:Identificat An
alternate identification on erIdentification ionType assigned to an
organization. primaryContact nc:OrganizationPri nc:ContactIn A preferred
means of maryContactInform formationTy contacting an ation pe
organization. principalOfficial nc:OrganizationPri nc:PersonTy A chief
or high ranking ncipalOfficial pe executive of an organization. status
nc:OrganizationSta nc:StatusTyp A status of an Examples tus e
organization. include: active, inactive subUnit nc:OrganizationSu
nc:Organizat A division of an Ex: bUnit ionType organization.
department, group taxIdentification nc:OrganizationTax nc:Identificat A
tax identification Ex: Federal Identification ionType assigned to an
Employer organization. Identification Number, FEIN, an Employer
Identification Number, EIN 10/19/2021

Page 154 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name abbreviation
description A.9

Data element Type Description Comment nc:OrganizationAb nc:TextType An
abbreviation, acronym, Ex: FBI, NCIC breviationText or code for an
organization name. nc:OrganizationDe nc:TextType A description of an
scriptionText organization.

Identification Type Table A-9 Identification Type

JSON Name id

Data element Type Description Comment nc:IdentificationID niemAn
identifier. xs:string categoryCode j:PersonIDCategory ncic:MNUCo A kind
of identifier Code deType assigned to a person. categoryDescri
nc:IdentificationCat nc:TextType A description of a kind of ption
egoryDescriptionTe identification. xt effectiveDate nc:IdentificationEff
nc:DateType A date an identification This may or ectiveDate takes
effect. may not be the issue date. expirationDate nc:IdentificationExp
nc:DateType A date after which an irationDate identification is no
longer valid. jurisdiction nc:IdentificationJuri nc:Jurisdictio An area,
region, or unit For EIDO only sdiction nType where a unique
nc:Jurisdiction identification is issued. Text in JurisdictionTyp e must
be supported source nc:IdentificationSo nc:TextType A person,
organization, urceText or locale which issues an identification. status
nc:IdentificationSta nc:StatusTyp A status of an Ex: valid, tus e
identification. expired, revoked, suspended, 10/19/2021

Page 155 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Data element

Type

Description

Comment replaced, duplicate, lost

A.10 Status Type Table A-10 Status Type JSON Name status

Data element nc:StatusText

Type Description Comment nc:TextType A status or condition of something
or someone. date nc:StatusDate nc:DateType A date a status was set,
effected, or reported. description nc:StatusDescriptio nc:TextType A
description of a status nText or condition of something or someone.
issuerIdentifica nc:StatusIssuerIde nc:Identificat An identification of
a tion ntification ionType person or organization which assigns a
status. A.11 Address Type Table A-11 Address Type JSON Name fullText

Data element Type Description nc:AddressFullText nc:TextType A complete
address.

Comment Usually includes street, city, state, and zip code information.
Can also include other mailing specifics such as PO Box, apartment
number, or suite number.

recipientName nc:AddressRecipien nc:TextType A name of a person, tName
organization, or other 10/19/2021

Page 156 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Data element

Description Comment recipient to whom physical mail may be sent.
buildingName nc:AddressBuilding nc:TextType A name of a specific Name
building at an address to distinguish it from other buildings at the
same site. deliveryPointId nc:AddressDelivery niemAn identifier of a
single PointID xs:string place or unit at which mail is delivered.
privateMailbox nc:AddressPrivateM nc:TextType A private mailbox within
mail stop code ailboxText a company. secondaryUnit nc:AddressSeconda
nc:TextType A particular unit within a Ex: apartment ryUnitText larger
unit or grouping at number, suite a location. number routeName
nc:LocationRouteN nc:TextType A name and number of a ame postal route.
street nc:LocationStreet nc:StreetTyp A road, thoroughfare or e highway.
city nc:LocationCityNam nc:ProperNa A name of a city or town. e
meTextType countyCode nc:LocationCountyC census:USCo A county, parish,
ode untyCodeTyp vicinage, or other such e geopolitical subdivision of a
state. state nc:LocationState nc:StateTyp A state, commonwealth, Only e
province, or other such nc:LocationSta geopolitical subdivision of
teFIPS5-2Code a country. is required to be supported in StateType
country nc:LocationCountry nc:CountryT A country, territory, Only ype
dependency, or other nc:LocationCo such geopolitical untryISO3166
subdivision of a location. Alpha2Code is required to be 10/19/2021

Type

Page 157 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Data element

Type

Description

Comment supported in CountryType

A.12 Area Type Table A-12 Area Type JSON Name Data element Type
Description Comment wgs84Location mo:WGS84Location mo:WGS84El A area
region described Ellipse Ellipse lipseType by an ellipse specified by a
point, major axis, minor axis and rotation, using WGS84, meters, and
decimal degrees. wgs84Location mo:WGS84Location mo:WGS84E An area region
described ExternalPolygo ExternalPolygon xternalPolyg by a polygon with
no n onType interior region, using WGS84 coordinates. A.13
ContractInformation Type Table A-13 ContractInformation Type JSON Name
emailId

Data element Type nc:ContactEmailID niemxs:string

Description Comment An electronic mailing address by which a person or
organization may be contacted. instantMesseng nc:ContactInstantM
nc:InstantMe A user account for an er essenger ssengerType instant
messaging program by which a person or organization may be contacted.
mailingAddress nc:ContactMailingA nc:AddressTy A postal address by
ddress pe which a person or organization may be contacted. radio
nc:ContactRadio nc:ContactRa A method of contacting dioType a person or
10/19/2021

Page 158 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 JSON Name

Data element

Type

Description Comment organization by messages over a radio. telephoneNum
nc:ContactTelephon nc:Telephone A telephone number for ber eNumber
NumberType a telecommunication device by which a person or organization
may be contacted. entity nc:ContactEntity nc:EntityType An entity that
may be contacted by using the given contact information.
entityDescriptio nc:ContactEntityDe nc:TextType A description of the
This could be a n scriptionText entity being contacted. title or
function of the person with the given contact information.'
informationDes nc:ContactInformati nc:TextType A description of the
cription onDescriptionText contact information. responder
nc:ContactRespond nc:PersonTyp A third party person er e who answers a
call and connects or directs the caller to the intended person. A.14
Entity Type Table A-14 Entity Type JSON Name Data element Type
Description organization nc:EntityOrganizatio nc:Organizatio An
organization n nType capable of bearing legal rights and
responsibilities. person nc:EntityPerson nc:PersonType A person capable
of bearing legal rights and responsibilities. 10/19/2021

Comment

Page 159 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A.15 CrossStreet Type Table A-15
Cross Street Type JSON Name description

relativeLocation

Data element Type Description Comment nc:CrossStreetDescri nc:TextType A
description of NIEM uses the ptionText a street term CrossStreet
intersection. as a location but this document uses it to denote
proximity to a location. nc:CrossStreetRelati nc:RelativeLo A location
of veLocation cationType something relative to a street intersection.

A.16 PhysicalFeature Type Table A-16 PhysicalFeature Type JSON Name
generalCategory categoryCode description image location

10/19/2021

Data element Type Description Comment nc:PhysicalFeatureGe nc:TextType A
general kind Ex: scar, mark, neralCategoryText of physical tattoo,
missing feature. limb j:PhysicalFeatureCate nc:TextType A specific kind
goryCode of physical feature. nc:PhysicalFeatureDe nc:TextType A
description of scriptionText a physical feature. nc:PhysicalFeatureIm
nc:ImageTy A digital image age pe of a physical feature.
nc:PhysicalFeatureLo nc:TextType A location on a cationText person's
body of a physical feature.

Page 160 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A.17 MedicalCondition Type Table
A-17 MedicalCondition Type JSON Name conditionText

cause

description isPresent

severity

dateRange

Data element Type Description Comment nc:MedicalCondition nc:TextType A
state of Ex: scar, mark, Text health, ontattoo, missing going or limb
present. nc:MedicalCondition nc:TextType A trigger that Ex: drug, food,
CauseText can initiate the allergen onset of a medical condition.
nc:MedicalCondition nc:TextType A description of DescriptionText a
medical condition. nc:MedicalCondition niemTrue if a PresentIndicator
xs:boolean medical condition currently exists; false otherwise.
nc:MedicalCondition nc:TextType A degree to SeverityText which a medical
condition is affecting a person. nc:MedicalCondition nc:DateRan A date
range for DateRange geType the start and end of a medical condition.

Appendix B Mapping between vCard and NIEM Person Type Table B-1 Mapping
between vCard and NIEM Person Type vCard Data Properties General
Properties: 10/19/2021

NIEM Data Element Page 161 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 vCard Data Properties BEGIN END
SOURCE KIND

NIEM Data Element No Mapping No Mapping No Mapping No Mapping

XML Identification Properties: FN N

No Mapping

NICKNAME PHOTO BDAY ANNIVERSARY GENDER ADR

TEL EMAIL IMPP LANG TZ GEO TITLE ROLE ORG 10/19/2021

nc:PersonFullName nc:PersonSurName nc:PersonGivenName
nc:PersonMiddleName nc:PersonNamePrefixText nc:PersonNameSuffixText
im:PersonNickName nc:PersonDigitalImage nc:PersonBirthDate No Mapping
j:PersonSexCode nc:AddressPrivateMailboxText nc:AddressSecondaryUnitText
nc:LocationStreet nc:LocationCityName nc:LocationCityName
nc:LocationPostalCode nc:LocationCountry nc:FullTelephoneNumber nc:Email
nc:ContactInstantMessenger nc:PersonPrimaryLanguage No Mapping
geo:LocationGeospatialPoint em:JobTitleOrRoleName em:JobTitleOrRoleName
nc:OrganizationName Page 162 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 vCard Data Properties MEMBER

NIEM Data Element nc:OrganizationBranchName
nc:OrganizationPrimaryContactInformation

RELATED

No Mapping except for emergency contact:
nc:PersonEmergencyContactInformation nc:PersonDescriptionText
nc:PersonDescriptionText No Mapping No Mapping No Mapping
nc:PersonOtherIdentification No Mapping nc:URLID No Mapping
nc:ContactWebsiteURI nc:ContactWebsiteURI nc:ContactWebsiteURI
nc:ContactWebsiteURI

CATEGORIES NOTE PRODID REV SOUND UID CLIENTPIDMAP URL VERSION KEY FBURL
CALADRURI CALURI

10/19/2021

Page 163 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021

Appendix C OpenAPI Schema The schemas defined in this document are
described using OpenAPI V3, with JSON-ld exchange schemas for NIEM
conformance. This section contains the OpenAPI interface definition as a
YAML object. Only the component section of this schema is relevant to
the EIDO JSON object described in this document. The normative reference
can be found in NENA git hub repository available here
https://github.com/NENA911/EIDO-JSON. openapi: 3.0.2 info: title:
Emergency Incident Data Object (EIDO) description: \>A standard format
and content definition for exchanging emergency incident related data.
contact: name: NENA email: commleadership@nena.org version: '1.0.0'
servers: - url: 'http://localhost/' - url: 'https://localhost/' paths:
'/EIDO/{IncidentId}': get: tags: - Dereference summary: Dereference an
EIDO reference operationId: GetEIDO parameters: - name: IncidentId in:
path required: true description: The id of the incident represented in
the EIDO

10/19/2021

Page 164 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 schema: type: string responses:
'200': description: OK content: application/json: schema: \$ref:
'#/components/schemas/EmergencyIncidentDataObjectType' components:
schemas: ReferenceType: title: Reference required: - $ref
properties:$ref: type: string description: The \$id value of the
referenced data component PrologueType: title: Prologue required: - $id
- lastUpdateTimeStamp
- updatedByAgencyReference
type: object
properties:$id: type: string description: Globally unique identifier of
the data component lastUpdateTimeStamp: type: string

10/19/2021

Page 165 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 description: Timestamp of the
last time the data component was modified format: date-time
updatedByAgencyReference: $ref: '#/components/schemas/ReferenceType'
description:
Owning Agency Identifier. This includes private and public
providers. Agencies are globally unique.
updatedByAgentReference:$ref: '#/components/schemas/ReferenceType'
description: A globally unique identifier of the agent or automaton that
updated the data component last description: A common prologue on each
data component AgencyType: title: Agency Type description: A data type
for an organization that owns or is otherwise related to an emergency
incident. allOf: - \$ref: '#/components/schemas/PrologueType' - type:
object required: - agencyRoleDescriptionRegistryText - agencyType
properties: agencyJcard: type: string description: The agency's contact
information. agencyRoleDescriptionRegistryText: type: array

10/19/2021

Page 166 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 description: Shall only contains
values present in the Agency Role registry items: type: string
incidentOwningAgencyIndicator: type: boolean description: True if an
agency owns the emergency incident; false otherwise agencyType: type:
array items: type: string description: 'One or more members of a list of
available provider and agency types. The agency types are taken from the
IANA urn:emergency:service:responder registry and the subregistries
defined for that registry. Values exclude the
"urn:emergency:service:responder." prefix. Example values:
"police.federal.fbi", "fire.forest", "poison_control" including: Law
Enforcement, Fire, EMS, Consolidated Dispatch, Ambulance Company, etc.
Agency Types are defined in the IANA urn:emergency:service:responder
registry' notesReference: type: array description: Reference to a Notes
Data Component. Contains optional alphanumeric text further describing
the agency. items: \$ref: '#/components/schemas/ReferenceType'
AgentType: title: Agent Type description: A data type for information
about agents (e.g., call takers, dispatchers, supervisors, responding
emergency resources, etc.) and automated systems acting as agents that
are involved in an incident. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object required:

10/19/2021

Page 167 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 - agentRoleRegistryText -
agencyReference properties: agentJcard: type: string description: The
agent's contact information. agentWorkstationPositionIdentification:
type: string description: The workstation position ID within the agent's
agency. agentRoleRegistryText: type: array items: type: string
description: Shall only contains values present in the Agent Role
registry agencyReference: $ref: '#/components/schemas/ReferenceType'
description: Linked element type is AgencyType
notesReference:
type: array
description:
Reference to a Notes Data Component. Contains optional
alphanumeric text further describing the agent.
items:$ref: '#/components/schemas/ReferenceType'
IncidentMergeDirectionCodeSimpleType: title: Incident Merge Direction
Code Simple Type type: string description: A data type for a code list
for a direction of a merge or split. example: REPLACING

10/19/2021

Page 168 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 enum: - REPLACING - REPLACED -
SPLIT MergeInformationType: title: MergeInformationType description: A
data type for information about a merged or split incident. allOf: -
$ref: '#/components/schemas/PrologueType'
- type: object
required:
- incidentMergeDirectionCode
- incidentTrackingIdentifier
properties:
incidentTrackingIdentifier:
type: string
description:
The Incident Tracking Identification of the Incident
that is being merged with, or split from the Incident.
incidentMergeDirectionCode:$ref:
'#/components/schemas/IncidentMergeDirectionCodeSimpleType'
LinkDirectionCodeSimpleType: title: Link Direction Code Simple Type
type: string description: A data type for a code list for a kind of link
performed. example: PARENT enum: - PARENT - CHILD - RELATED - UNLINK
LinkInformationType:

10/19/2021

Page 169 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 title: Link Information Type
description: A data type for information about a linked incident.
allOf: - $ref: '#/components/schemas/PrologueType'
- type: object
required:
- incidentTrackingIdentifier
- linkDirectionCode
properties:
incidentTrackingIdentifier:
type: string
description: A globally unique emergency incident identifier.
transactionReasonText:
type: string
description: Free format narrative description of the reason for
the link
involvedAgencyReference:$ref: '#/components/schemas/ReferenceType'
description: Agencies involved with the linked incident. Must be
populated if the agency sending the EIDO is not involved with the linked
incident. linkDirectionCode: \$ref:
'#/components/schemas/LinkDirectionCodeSimpleType' NotesType: title:
Notes Type description: A data type for notes and comments associated
with an incident. allOf: - \$ref: '#/components/schemas/PrologueType' -
type: object required: - notesActionComments

10/19/2021

Page 170 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 properties: notesActionComments:
type: string description: A note about an incident created by an agent;
in HTML5 format. agentReference: \$ref:
'#/components/schemas/ReferenceType' description: Linked element type is
AgentType AdditionalDataType: title: Additional Data Type description: A
data type for additional information associated with an incident call,
location, or person. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object properties:
additionalDataByReferenceUrl: type: string description: A URL that
references additional data associated with a call.
additionalDataByValue: type: string description: Additional data
associated with a call. urlPurpose: type: string description: Purpose
that accompanied the URL. For example, the purpose parameter sent with
the URL in the Call-Info headers. notesReference: type: array
description: Notes and comments associated with an incident. Linked
elements type are NotesType

10/19/2021

Page 171 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 items: \$ref:
'#/components/schemas/ReferenceType' LocationInformationType: title:
Location Information Type description: A data type for a location
associated with an incident. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object required: -
locationTypeDescriptionRegistryText properties:
locationTypeDescriptionRegistryText: type: string description: Shall
only contains values present in the Location type registry
locationByValue: type: string description: A type for an external
adapter for a Presence Information Data Format-Location Object
(PIDF-LO). locationDescriptionText: type: string description: Optional
text further describing the location type. locationByReferenceUrl: type:
string description: A URL where a location associated with an incident
can be dereferenced. format: uri crossStreetByValue: type: array
description:

10/19/2021

Page 172 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A nearest cross street to an
incident's location in PIDF-LO format. items: type: string
crossStreetByReferenceUrl: type: array description: A URL where a
nearest cross street to an incident's location can be dereferenced.
items: type: string intersectingStreetByValue: type: array description:
An additional street, that with the location forms an intersection, in
PIDF-LO format. items: type: string intersectingStreetByReferenceUrl:
type: array description: A URL where an additional street, that with the
location forms an intersection, can be dereferenced. items: type: string
cellTowerSectorId: type: string description: The nearest cell tower and
sector/face of the tower receiving a call associated with an incident.
additionalDataReference: type: array description:

10/19/2021

Page 173 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Additional information associated
with an incident call, location, or person. Linked elements type are
AdditionalDataType items: $ref: '#/components/schemas/ReferenceType'
notesReference:
type: array
description:
Notes and comments associated with an incident. Linked elements
type are NotesType
items:$ref: '#/components/schemas/ReferenceType' PersonInformationType:
title: Person Information Type description: A data type for information
about a person associated with an incident. allOf: -
$ref: '#/components/schemas/PrologueType'
- type: object
required:
- personIncidentRoleRegistryText
properties:
ncPersonComponent:$ref:
'NIEMComponents/nc.yaml#/components/schemas/ncPersonType'
personIncidentRoleRegistryText: type: array description: Shall only
contains values present in the Person Role registry items: type: string
callIdentifier: type: array description: Call Identifier of the calls
this person was involved on.

10/19/2021

Page 174 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 items: type: string
additionalDataReference: type: array description: Additional information
associated with an incident call, location, or person. Linked elements
type are AdditionalDataType items:
$ref: '#/components/schemas/ReferenceType'
locationReference:
type: array
description:
A location associated with an Person. Linked elements type is
LocationInformationType
items:$ref: '#/components/schemas/ReferenceType' notesReference: type:
array description: Reference to a Notes Data Component. Contains
optional alphanumeric text further describing the person. items:
$ref: '#/components/schemas/ReferenceType'
callBackReference:$ref: '#/components/schemas/ReferenceType'
description: Reference to a call back data component. Identifies how the
person described in this data component can be called back.
DispositionType: title: Disposition Type description:

10/19/2021

Page 175 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A data type for incident
disposition information entered or updated by a dispatch agent and/or a
responding emergency resource. allOf: -
$ref: '#/components/schemas/PrologueType'
- type: object
required:
- dispositionCommonRegistryCode
properties:
dispositionCommonRegistryCode:
type: string
description:
Shall only contains values present in the Common Disposition
Code registry
dispositionPrimaryIndicator:
type: boolean
description:
True if the common disposition is the primary disposition for an
agency for an incident; false otherwise.
dispositionCategoryText:
type: string
dispositionDescriptionText:
type: string
notesReference:
type: array
description:
Reference to a Notes Data Component. Contains optional
alphanumeric text further describing the incident.
items:$ref: '#/components/schemas/ReferenceType' VehicleInformationType:
title: Vehicle Information Type allOf:

10/19/2021

Page 176 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 -
$ref: '#/components/schemas/PrologueType'
- type: object
required:
- vehicleRelationshipType
- vehicleRelationshipTimeStamp
properties:
ncVehicleComponent:$ref:
'NIEMComponents//nc.yaml#/components/schemas/ncVehicleType'
vehicleRelationshipType: type: array description: Information about a
vehicle associated with an incident. Shall only contains values present
in the Vehicle Relationship Type registry items: type: string
vehicleRelationshipTimeStamp: type: string description: The date and
time that the relationship of the vehicle to the incident was
established format: date-time additionalDataReference: type: array
description: Additional information associated with a vehicle. Linked
elements type is AdditionalDataType items: \$ref:
'#/components/schemas/ReferenceType' locationReference: type: array
description:

10/19/2021

Page 177 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A location associated with an
vehicle. Linked elements type is LocationInformationType items:
$ref: '#/components/schemas/ReferenceType'
notesReference:
type: array
description:
Reference to a Notes Data Component. Contains optional
alphanumeric text further describing the vehicle.
items:$ref: '#/components/schemas/ReferenceType'
IncidentInformationType: title: Incident Information Type description: A
data type for general information about emergency incidents gathered by
emergency agents, responding emergency resources, from reporting parties
(callers), and devices reporting emergency incidents. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object required: -
incidentTypeCommonRegistryText properties: incidentTypeInternalCode:
type: string description: An alphanumeric code indicating the category
of the incident. This is the internal code used by the local agencies
involved in the incident. incidentTypeInternalText: type: string
description:

10/19/2021

Page 178 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Textual description indicating
the category of the incident. This is the internal text used by the
local agencies involved in the incident. incidentTypeCommonRegistryText:
type: string description: Shall only contains values present in the
Incident Type registry incidentStatusInternalText: type: array
description: An internal code used by a local agency for a status of an
incident. items: type: string incidentStatusCommonRegistryText: type:
array description: Shall only contains values present in the Incident
Status registry items: type: string internalIncidentId: type: string
agentReference: type: array description: Linked elements type is
AgentType items: \$ref: '#/components/schemas/ReferenceType'
locationReference: type: array description: A location associated with
an incident. Linked elements type is

10/19/2021

Page 179 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 LocationInformationType items:
$ref: '#/components/schemas/ReferenceType'
personReference:
type: array
description:
Information about a person associated with an incident. Linked
elements type is PersonInformationType
items:$ref: '#/components/schemas/ReferenceType' vehicleReference: type:
array description: Information about a vehicle associated with an
incident. Linked elements type is VehicleInformationType items: \$ref:
'#/components/schemas/ReferenceType' documentIdentification: type: array
items: type: string reportNumberType: type: string
incidentPriorityInternalText: type: string description: An internal
value used by a local agency for the priority of an incident.
incidentCommonPriorityNumber: type: integer description: A globally
unique priority of an incident. format: int32

10/19/2021

Page 180 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 beatOrDispatchGroupText: type:
string description: A beat or dispatch group that contains an incident.
dispositionComponent: type: array description: Incident disposition
information entered or updated by a dispatch agent and/or a responding
emergency resource. items: $ref: '#/components/schemas/DispositionType'
notesReference:
type: array
description:
Reference to a Notes Data Component. Contains optional
alphanumeric text further describing the incident.
items:$ref: '#/components/schemas/ReferenceType' UpdatedCBNType: title:
Updated CBN Type description: A data type for additional telephone
numbers and SIP equivalents that can be used to contact the individual
that made a call associated with an incident. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object required: -
updatedCbnIdentifierUri properties: updatedCbnIdentifierUri: type:
string description:

10/19/2021

Page 181 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 A telephone number or SIP
equivalent of a reporting party described in the parent call. format:
uri updatedCbnCallerDescription: type: string description: Additional
descriptive text about an updated caller call back number such as the
days and hours to use it and the type of number (e.g., work, home,
friend, etc.). CallInformationType: title: CallInformationType
description: A data type for call information about an incident received
and collected by an agent. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object required: -
callStateRegistryText - standardPrimaryCallType - direction -
callStartTimestamp properties: relatedCallIdentifier: type: array
description: CallIdentifier of calls that are related to this call. Eg
Callbacks items: type: string queueIdentifier: type: string

10/19/2021

Page 182 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 description: An identifier of a
queue where the call is currently on. standardPrimaryCallType: type:
string description: Values limited to those in the LogEvent CallTypes
registry with "Primary" classification. standardSecondaryCallType: type:
string description: Values limited to those in the LogEvent CallTypes
registry with "Secondary" classification. localCallType: type: string
description: Locally defined call type. direction: type: string enum: -
incoming - outgoing - internal description: "incoming" means received by
the EIDO sender, "outgoing" means sent by the EIDO sender and "internal"
means both caller and recipient are within the EIDO sender.
additionalDataReference: type: array description: Additional information
associated with an incident call, location, or person. Linked elements
type is AdditionalDataType

10/19/2021

Page 183 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 items:
$ref: '#/components/schemas/ReferenceType'
callStartTimestamp:
type: string
description: Date and time stamp of when the call was received.
format: date-time
answerDate:
type: string
description: A date and time stamp on which the call was answered
by an agent.
format: date-time
callStateRegistryText:
type: string
description: Current call status (when the EIDO was created) from
the available call statuses in the “Call States” registry i.e., callBegin,
callAnswered, callEnd, partyAdd , etc. Values limited to those in the IANA
CallStates registry.
callBackReference:$ref: '#/components/schemas/ReferenceType'
description: Reference to a call back data component. Identifies how the
caller described in this data component can be called back. verstat:
type: string description: Value of the verstat parameter. Required if
known. sipIdentity: type: string description: Content of the SIP
identity header. Required if known. agentReference: type: array
description: Linked elements type is AgentType

10/19/2021

Page 184 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 items:
$ref: '#/components/schemas/ReferenceType'
locationReference:
type: array
description:
A location associated with an incident. Linked elements type is
LocationInformationType
items:$ref: '#/components/schemas/ReferenceType' personReference: type:
array description: Information about a person associated with an
incident. Linked elements type is PersonInformationType items:
$ref: '#/components/schemas/ReferenceType'
notesReference:
type: array
description:
Contains optional alphanumeric text further describing the call.
Linked elements type are NotesType
items:$ref: '#/components/schemas/ReferenceType' CallbackType: title:
Callback Type description: This data component is used to exchange
information about how to contact the Incident's reporting parties.
allOf: - \$ref: '#/components/schemas/PrologueType' - type: object
properties:

10/19/2021

Page 185 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 callBackInformationUri: type:
array description: Information that enables agents and responders to
call back. URI(s) found in the P-Asserted-Identity header if provided
else

the URI in the From header.

items: type: string deviceContactHeader: type: string description:
Content of the device contact header. Information that enables agents
and responders to possibly reach (call back) the device that initiated
the call. updatedCbn: description: Identifies additional telephone
numbers and SIP equivalents that can be used to contact. \$ref:
'#/components/schemas/UpdatedCBNType' EmergencyResourceType: title:
Emergency Resource Type description: A data type for a resource
(vehicle, person \[foot patrol\], organizational unit such as a squad or
strike team, and other emergency responder configurations) that responds
to an emergency incident. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object required: -
agentReference - emergencyResourceTypeCommonRegistryText -
primaryUnitStatusRegistryText - resourceAttributeRegistryText

10/19/2021

Page 186 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 - secondaryUnitStatusRegistryText
properties: eta: type: string description: Timestamp of when the unit is
expected to reach it's destination format: date-time
emergencyResourceTypeCommonRegistryText: type: string description: Shall
only contains values present in the Emergency Resource Type registry
emergencyResourceTypeInternalText: type: string description: A local
code for an emergency resource type. resourceAttributeRegistryText:
type: array description: Shall only contains values present in the
Resource Attribute registry items: type: string emergencyResourceName:
type: string description: Common name of the emergency resource.
primaryUnitStatusRegistryText: type: string description: Shall only
contains values present in the Primary Unit Status registry
secondaryUnitStatusRegistryText:

10/19/2021

Page 187 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 type: array description: Shall
only contains values present in the Secondary Unit Status registry
items: type: string notesReference: type: array description: Reference
to a Notes Data Component. Contains optional alphanumeric text further
describing the incident. items:
$ref: '#/components/schemas/ReferenceType'
unitStatusInternal:
type: array
description: A local code for a status of an emergency resource.
items:
type: string
agentReference:
type: array
description:
Information about agents (e.g., call takers, dispatchers,
supervisors, responding emergency resources, etc.) and automated
systems acting as agents that are involved in an incident.
Linked elements type is AgentType
items:$ref: '#/components/schemas/ReferenceType' unitLocationReference:
$ref: '#/components/schemas/ReferenceType'
description: Linked elements type is LocationInformationType
unitDestinationReference:$ref: '#/components/schemas/ReferenceType'

10/19/2021

Page 188 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 description: Linked elements type
is LocationInformationType personReference: type: array description:
Information about a person associated with an incident. Linked elements
type is PersonInformationType items:
$ref: '#/components/schemas/ReferenceType'
vehicleReference:
type: array
description:
Information about a vehicle associated with an incident. Linked
elements type is VehicleInformationType
items:$ref: '#/components/schemas/ReferenceType'
DispatchInformationType: title: Dispatch Information Type description: A
data type for dispatch information related to an incident. allOf: -
$ref: '#/components/schemas/PrologueType'
- type: object
properties:
agencyReference:$ref: '#/components/schemas/ReferenceType' description:
Reference to an AgencyType. Agency that was dispatched through action
performed in this data component. agentReference: \$ref:
'#/components/schemas/ReferenceType' description: Linked elements type
is AgentType. Identifies the agent that completed the dispatch operation
described

10/19/2021

Page 189 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 emergencyResourceReference: type:
array description: A resource (vehicle, person \[foot patrol\],
organizational unit such as a squad or strike team, and other emergency
responder configurations) that responds to an emergency incident. Linked
elements type is EmergencyResourceType items:
$ref: '#/components/schemas/ReferenceType'
notesReference:
type: array
description:
Contains optional alphanumeric text further describing the
dispatch.
Linked elements type are NotesType
items:$ref: '#/components/schemas/ReferenceType' AlarmsSensorsType:
title: AlarmsSensorsType description: A data type for alarms and sensors
information associated with an incident. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object properties:
csaaAlarmInformation: type: string description: Central Station Alarm
Association (CSAA) data associated with an incident. alarmUrl:

10/19/2021

Page 190 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 type: string description: A URL
that references Central Station Alarm Association (CSAA) data associated
with an incident. format: uri agentReference:
$ref: '#/components/schemas/ReferenceType'
description: Linked elements type is AgentType
notesReference:
type: array
description:
Contains optional alphanumeric text further describing the
dispatch.
Linked elements type are NotesType
items:$ref: '#/components/schemas/ReferenceType'
EmergencyIncidentDataObjectType: title: Emergency Incident Data Object
Type description: A data type for an information exchange containing
emergency incident related data. allOf: - \$ref:
'#/components/schemas/PrologueType' - type: object required: -
eidoVersion - issuingElementIdentification - agencyComponent properties:
eidoVersion: type: string description: Version of the EIDO structure.

10/19/2021

Page 191 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 enum: - '1.0'
issuingElementIdentification: type: string description: An identifier of
a functional element that created the EIDO instance. mergeComponent:
type: array description: Information about a merged or split incident.
items: $ref: '#/components/schemas/MergeInformationType'
linkComponent:
type: array
description: Information about a linked incident.
items:$ref: '#/components/schemas/LinkInformationType'
incidentComponent: $ref: '#/components/schemas/IncidentInformationType'
description:
General information about the emergency incident gathered by
emergency agents, responding emergency resources, from reporting
parties (callers), and devices reporting emergency incidents.
callComponent:
type: array
description:
Call information about an incident received and collected by an
agent.
items:$ref: '#/components/schemas/CallInformationType'
callbackComponent: type: array

10/19/2021

Page 192 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 description: Contains information
about how to call a person. items:
$ref: '#/components/schemas/CallbackType'
dispatchComponent:
type: array
description: Dispatch information related to an incident.
items:$ref: '#/components/schemas/DispatchInformationType'
alarmsSensorsComponent: type: array description: Alarms and sensors
information associated with an incident. items:
$ref: '#/components/schemas/AlarmsSensorsType'
agencyComponent:
type: array
description: All agencies referenced in this document are listed
here
items:$ref: '#/components/schemas/AgencyType' agentComponent: type:
array description: All agents referenced in this document are listed
here items: $ref: '#/components/schemas/AgentType'
notesComponent:
type: array
description: All notes referenced in this document are listed here
items:$ref: '#/components/schemas/NotesType' additionalDataComponent:

10/19/2021

Page 193 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 type: array description: All
Additional Data related to the Incident except the objects found in the
`<provided-by>`{=html} field of a PIDF-LO. items:
$ref: '#/components/schemas/AdditionalDataType'
locationComponent:
type: array
description:
All location informations referenced in this document are listed
here
items:$ref: '#/components/schemas/LocationInformationType'
personComponent: type: array description: All person informations
referenced in this document are listed here items:
$ref: '#/components/schemas/PersonInformationType'
vehicleComponent:
type: array
description:
All vehicle informations referenced in this document are listed
here
items:$ref: '#/components/schemas/VehicleInformationType'
emergencyResourceComponent: type: array description: All emergency
resources information referenced in this document are listed here items:

10/19/2021

Page 194 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 \$ref:
'#/components/schemas/EmergencyResourceType' '@context': type: object
description: JSON-LD context properties: '@vocab': type: string enum: -
'./JSONLD_Contexts/EmergencyIncidentDataObjectType.jsonld'

ACKNOWLEDGEMENTS The National Emergency Number Association (NENA) Agency
Systems Committee, EIDO JSON Working Group developed this document. NENA
Board of Directors Approval Date: 10/19/2021 NENA recognizes the
following industry experts and their employers for their contributions
to the development of this document.

Members Employer Michael Smith, Agency Systems Equature/DSS Corp.
Committee Co-Chair Rick Blackwell, ENP, Agency Systems Greenville
County, SC Committee Co-Chair Amy McDowell, ENP, Working Group
CoGreenville County, SC Chair Etienne Paquet, Working Group Co-Chair
Komutel Tony Dunsworth, Note-Taker City of Alexandria, VA Brandon Abley,
ENP NENA Jamie Algatt Rapid Deploy Roy Artin E-Comm E9-1-1 Vancouver
Bernard Brabant Consultant Guy Caron Bell Canada Rakhee Duneja Toronto
Police Service ON Simon Farrow Stancil Québec 9-1-1 Agency Pierre
Foucault Ryan Fouts Rapid Deploy Rebecca Green, ENP City of
Murfreesboro, TN 10/19/2021

Page 195 of 196

© Copyright 2021 National Emergency Number Association, Inc.

NENA Standard for Emergency Incident Data Object (EIDO)
NENA-STA-021.1b-2021, October 19, 2021 Members Diane Harris Ryan Keeling
James Kinney Alain Labrosse Darrell Marsh, P.Eng, PMP Roger Marshall
Richard Maw Steven McDowall Steve McMurrer, ENP George Michallas Dan
Mongrain Richard Muscat Katy Myers, ENP Justin Prahar Theresa Reese
Brian Rosen Jerry Schlesinger, PMP Brooks Shannon Johanne Tanguay Matt
Tenold Henry Unger Ravi Valavandan Michael Vislocky Lisa Wirtanen
Michael Salonish

Employer Zetron Bell Canada INdigitial Montreal Police QC CA SPVM Nova
Scotia Emergency Management Office NS Comtech Telecommunications
Corporation RAM Solutions Carbyne, Inc Fairfax County, VA
MicroAutomation, Inc Motorola Solutions, Inc Bexar Metro 9-1-1 Network
District, TX Norcomm Public Saefty Communications, Inc ECats Ericsson,
Inc Consultant Intelligent Systems, Inc Rapid Deploy Montreal Police QC
CA SPVM City of Lynwood, WA Pulsiam City of Calgary AB CA Network
Orange, Inc AT&T DataTech911

Special Acknowledgements: Delaine Arnold, ENP, Committee Resource
Manager, has facilitated the production of this document through the
prescribed approval process. The EIDO JSON Working Group Working Group
is part of the NENA Development Group that is led by: • Wendi Rooney,
ENP, and Jim Shepard, ENP, Development Steering Council CoChairs •
Brandon Abley, Technical Issues Director • April Heinze, ENP, PSAP
Operations Director 10/19/2021

Page 196 of 196

© Copyright 2021 National Emergency Number Association, Inc.


