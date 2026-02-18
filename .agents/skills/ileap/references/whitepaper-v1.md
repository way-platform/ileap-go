# Demystifying the iLEAP Project and Technical Specifications V1

**June 2025**

**© Smart Freight Centre & SINE Foundation, 2025**

Suggested citation: Smart Freight Centre and SINE Foundation. iLEAP Whitepaper V1. 2025. Version Number: iLEAP-GUID-001

## Disclaimer

The views expressed in this publication are those of Smart Freight Centre and SINE Foundation teams, consultants and management, and do not necessarily reflect the views of the Board of Trustees of both organizations. Smart Freight Centre and SINE Foundation do not guarantee the accuracy of the data included in this publication and do not accept responsibility for consequence of their use.

## Acknowledgements

We extend our sincere gratitude to the more than 35 companies that dedicated their time, expertise, and resources to developing the iLEAP Technical Specifications throughout an intensive 18-month collaborative process. Smart Freight Centre and SINE Foundation deeply appreciate this commitment and the contributions that made these specifications possible.

## About Smart Freight Centre

Smart Freight Centre is an international non-profit organization focused on reducing greenhouse gas emissions from freight transportation. Smart Freight Centre’s vision is an efficient and zero emission global logistics sector. Smart Freight Centre’s mission is to collaborate with the organization’s global partners to quantify impacts, identify solutions, and propagate logistics decarbonization strategies. Smart Freight Centre’s goal is to guide the global logistics industry in tracking and reducing the industry’s greenhouse gas emissions by one billion tones by 2030 and to reach zero emissions by 2050 or earlier, consistent with a 1.5°C future.

## About SINE Foundation

SINE Foundation is a non-profit tech-for-good organization. Founded by academic experts and entrepreneurs, SINE’s vision is a world economy that operates within safe limits. To this end, SINE builds open-source software and collaborates with global organizations such as the World Business Council for Sustainable Development (WBCSD) or the Smart Freight Center. By establishing global standards for sustainability data exchange, the foundation aims to empower organizations to trust and act on sustainability data through its work, accelerating climate action and impact-driven investments.

## Contact

**Smart Freight Centre**
Keizersgracht 572, 1017 EM, Amsterdam, Netherlands
P.O. Box 11772, 1001 GT, Amsterdam, Netherlands
Tel office: +31 6 4695 4405
www.smartfreightcentre.org
info@smartfreightcentre.org

**SINE Foundation**
Bredowstraße 35a, Berlin, 10551, DE
https://sine.foundation/
hello@sine.foundation

## EXECUTIVE SUMMARY

Customers, investors, and regulators are increasingly demanding detailed and accurate emissions data across supply chains ("Scope 3 emissions"). Emissions data based on primary data – rather than industry averages – plays a particularly critical role, as it provides the clearest signal for the decarbonization of the economy. Specifically for the logistics industry, emissions calculated from primary data are not only more precise, but they also enable more targeted decarbonization decisions.

**Until now, access to logistics emissions data has been costly and difficult** as emissions reports are predominantly created and exchanged manually. This happens due to the lack of industry-wide standardization and the lack of automation through software. As a result, the interoperable flow of logistics emissions data is impeded at global scale and across all modalities.

**The iLEAP Project, with the publication of the iLEAP Technical Specifications V1, closes this digitization gap in the underlying accounting standards ISO 14083 and the GLEC Framework.** It enables *all* companies operating close to logistics value chains to exchange emissions data sovereignly and interoperable. It supports the logistics sector in completing the transition to a reporting regime based on primary data, bringing forward standardized data quality indicators (KPIs) such as the Primary Data Share. Lastly, with the definition of interoperable data transactions, reporting capabilities are simplified not only for larger companies but also for small and medium enterprises (SMEs).

Through the publication of this white paper, we aim to explain the strategic and technical considerations underpinning the iLEAP project. We also hope to inspire readers to explore and adopt our technical work in their own companies' processes and systems. By working together, we can further develop iLEAP, minimize decarbonization barriers, and accelerate the transition to a carbon-neutral future through standardization and digitization.

## THE CHALLENGE: EMISSIONS TRANSPARENCY IN GLOBAL LOGISTICS

### Growing Demand for carbon transparency based on primary data

The business imperative for carbon transparency has never been stronger. Customers, investors, and regulators are increasingly demanding accurate emissions data for products and services. Emerging compliance requirements such as the Corporate Sustainability Reporting Directive (CSRD) are requiring detailed reporting across value chains (Scope 3).

For the majority of product and service categories, 80% or more of the total carbon footprint is in Scope 3 [1]. In this context, emissions data calculated from actual *measurements*—that is, “primary data” rather than industry averages—has become essential, representing what the OECD calls “the clearest signal for the decarbonization of the economy [2]”.

In the logistics sector specifically, primary data offers several key benefits:

*   **Precise Reporting:** Enabling logistics companies to accurately measure and report their specific emissions profiles based on their unique vehicle fleets, routes, and handling procedures, and enabling shippers to more accurately report their Scope 3 emissions.
*   **Targeted improvements:** Identifying the most effective areas for emissions reduction within specific operations rather than relying on industry averages.
*   **Incentive alignment:** Creating appropriate rewards for continuous improvement as companies can quantify and demonstrate the benefits of cleaner vehicles, route optimization, and more efficient practices.

> “Along harmonization of GHG calculation methodologies, standardization is needed for the step of sending the results to the contractor of a transport service. In France, the mandatory reporting of GHG emissions from transport services providers to their customers, for each transport service, would certainly benefit from a standardized protocol like iLEAP.”
>
> **Marc Cottignies - Engineering Expert, ADEME**

[1] McKinsey Catalyst Zero analysis, based on 2022 CDP data.
[2] OECD (2025), “The carbon footprint of everything”, OECD Net Zero+ Policy Papers, No. 6, OECD Publishing, Paris, https://doi.org/10.1787/ae22f8e8-en

### Key Barriers To Progress

Despite the clear need, widespread access to high-quality primary emissions data in logistics is prevented by significant barriers:

1.  **High effort for primary data collection**: The global nature of logistics involves numerous subcontractors and partners, each of them often operating bespoke systems. This fragmentation, which is a common and necessary part of operations, creates significant barriers to data collection and exchange [3].
2.  **Data privacy and business concerns**: Primary emissions data, especially if it provides fuel use transparency, is seen as sensitive data potentially affecting the companies’ competitive position.
3.  **Insufficient technical support**: Most carriers are small and medium enterprises (SMEs) lacking the capabilities to collect and report emissions data.
4.  **Lack of technical interoperability and standards**: The logistics sector uses numerous IT systems without standardized approaches to logistics emissions data exchange. Despite the methodological interoperability that has been achieved with frameworks like ISO 14083 and Global Emissions Council (GLEC) Framework, data access processes remain predominantly manual. This results in high costs and data quality concerns.
5.  **Inconsistent emissions reporting formats**: Reporting processes differ between companies, with shippers often using custom formats and templates. This forces logistics companies to repeatedly prepare data differently while shippers bear the cost of manual reporting.
6.  **Paper-based processes**: Many companies are severely challenged to provide emissions data because their processes aren't sufficiently digitized yet.

> “The exchange of logistics emissions data between partners in a transport chain remains a hard problem that prevents reliable emissions reporting. iLEAP has finally laid the groundwork to make data sharing so much easier.”
>
> **Dimitrios Kourtesis - Founder, CO2Path**

[3] Also see the Smart Freight Centre’s project report from the “Data Access” project (2019-2022).

## iLEAP: SUPPORTING COLLECTIVE ACTION TOWARDS LOGISTICS EMISSIONS TRANSPARENCY

### Vision And Mission

iLEAP's vision is to accelerate the decarbonization and sustainability transition of global logistics by creating an ecosystem of interoperable software solutions. This ecosystem will make it possible to measure, verify, and continuously improve the environmental impact of every shipment.

To achieve this, iLEAP’s mission is to deliver open technical specifications and digital infrastructure designed to measure, report, and reduce logistics emissions globally. As a first step, iLEAP will enable all organizations operating close to logistics value chains to access and exchange trusted emissions data based on the GLEC Framework and ISO 14083 standards, thus supporting collective action toward sector decarbonization.

***Recognizing the industry’s maturity in relation to carbon accounting and the flow of logistics activity and emissions data, the iLEAP Technical Specifications support both primary and secondary data. The ambition is to empower the industry through easier access to primary data and lessened reliance on secondary data.***

### Design Principles

In alignment with its vision and mission statements, the development of the iLEAP Technical Specifications and foundational infrastructure is rooted in the following principles:

*   **Data confidentiality and sovereignty**: Protecting sensitive business information.
*   **Data quality and integrity**: Ensuring transparency and access to decision-grade data and information.
*   **Global scope**: Supporting geographic expansion, multiple transport modes, and diverse stakeholders.
*   **Openness**: Providing equal opportunity to participate in and benefit from iLEAP and its broader ecosystem, including full participation in standard development.
*   **Agility**: Fostering a structured approach of continuous and iterative testing, validation, and refinement by and for the members of the project.

### iLEAP Today

At its core, iLEAP represents:

1.  The informal **iLEAP Community** of committed companies working collaboratively to decarbonize the logistics sector through digitization:
    *   Spanning the entire value chain, i.e., transport service users (shippers), logistics service providers (transport service organizers), carriers (transport operators), tool providers, public sector entities, consultants, and NGOs.
    *   Steered by SFC and SINE to jointly define and develop the iLEAP Technical Specifications and supporting materials.
    *   Creating market signals for primary data by a) pooling demand from shippers and other stakeholders, and b) providing a catalog of conforming software tools and solutions.
    *   Supporting all interested companies, especially tool providers, in implementing iLEAP-conforming software through documentation and technical assistance.

2.  The **iLEAP Technical Specifications** that address key barriers to the use of primary data for carbon accounting (fragmentation, cost, data access challenges) by:
    *   Facilitating the interoperable flow of emissions data across the logistics value chains by specifying a protocol for machine-to-machine communication that works irrespective of company size, mode of transport or geography.
    *   Addressing sovereignty and data privacy concerns by relying solely on peer-to peer connectivity only, eliminating the need for a trusted third party.
    *   Deriving a data model from the GLEC Framework and ISO 14083, thereby streamlining the multitude of reporting (data) formats in the sector while further simplifying emissions data management for all.
    *   Achieving technical and semantic interoperability with the WBCSD PACT ecosystem, an initiative representing major shippers that can now seamlessly, autonomously access logistics footprint data through existing carbon management systems at a low cost.
    *   Balancing primary and secondary data usage through quality metrics like Primary Data Share [4].

3.  **Development of knowledge products, tools, and open-source software** to support the implementation:
    *   Specific guides and technical documentation for companies looking to join the project, apply iLEAP operationally, or implement the iLEAP Technical Specifications.
    *   Demo implementations of the iLEAP protocol and data model.
    *   Demo endpoints with synthetic data to simplify testing and development.
    *   An Automated Conformance Testing Tool to minimize the impact, efforts, and costs of iLEAP certification.
    *   Collaboration with other organizations, such as the Open Logistics Foundation, to deliver open-source integration components in dedicated working groups.
    *   Openly accessible feedback tracker, giving transparency over the full life cycle of the project’s feedback and related decision-making.

[4] See Pathfinder Framework Version 2 and Version 3 for further details.

### iLEAP Standardization Process

With the Version 1 release, the iLEAP project team has applied and refined a standardization process inspired by approaches used within the Internet Engineering Task Force (IETF), Integrating the Healthcare Environment (IHE), and Partnership for Carbon Transparency (PACT) communities. This process is unique in its ability to integrate diverse stakeholders and work specifically within the digital space, while also being pragmatic and agile.

However, this standardization approach goes beyond emissions transparency. It can be seen as a general blueprint for addressing business cases in the sector that can only be approached collaboratively. Thus, the iLEAP community is already well-positioned to expand the scope and impact of iLEAP across the global logistics sector. See the Future Roadmap chapter for further details on this.

***The iLEAP digitization framework is based on:***

*   ***A committed coalition of partners across the entire value chain who embrace iLEAP's vision, mission, and design principles.***
*   ***A consensus-driven development process that engages relevant stakeholders throughout all phases—definition, validation, testing, and publication—to establish robust interoperability standards.***
*   ***A business-first standardization methodology that begins with well-defined use cases addressing specific information requirements and workflows within the sector, achieving scalable interoperability and automation through structured "data transactions" (detailed in the following chapter).***

## THE iLEAP TECHNICAL SPECIFICATIONS

### Overview

The iLEAP Technical Specifications build on ISO 14083:2023 standard and the GLEC Framework version 3. They first introduce a data model that embodies core concepts from these foundational methodologies:

*   **Transport Operations Category (TOC)** and **Hub Operations Category (HOC)**: describing groups of operations with similar characteristics, with assigned emission intensity values.
*   **Transport Chain (TC)**: a sequence of transport and hub operations, also called Transport Chain Elements (TCEs), that are required to move cargo, enabling the calculation of logistics emissions on a per-shipment basis.

In Figure 1 below, there is a simplistic representation of a transport chain. The individual transport chain elements (from different modes of transport) can be identified and are numbered from 1-11. The main message is to understand how TOCs can be roughly abstracted and defined. The gray boxes, next to the transport chain elements, are depicting the varying combinations of vehicle(s) in a route or network. In practice, the stakeholder that works on collecting data to conduct emissions accounting for those TCEs shall keep in mind the route and fleet characteristics along with equipment and contract type. The summary of the possible distinctions of TOCs in order of appearance are:

*   TCE 1: single vehicle on a route
*   TCE 3: single vehicle on a network
*   TCE 7: Group of similar vehicles on a route
*   TCE 11: Group of similar vehicles on a network

The process of abstracting data and designing a TOC is detailed and is dependent on the choices of Transport Operators, Transport Service Organizers or from a contracted visibility platform that supports this. For more details on this topic, we prompt the reader to read the GLEC Framework v3 or the ISO14083:2023.

> *Figure 1: Transport / Hub Operation Category vs TCEs*

Based on these concepts and the resulting data model, iLEAP specifies so-called “Data Transactions”. These transactions define how exactly data is exchanged between IT systems in a fully automated manner. A representation of these transactions in found in Figure 2.

**Data Transaction 1** enables transparency of logistics emissions at the Transport Chain Element (TCE) level to eventually form Transport Chains (TCs).

**Data Transaction 2** enables transparency of emission intensity at Transport Operation Category (TOC) level and Hub Operation Category (HOC) level.

Recognizing the need to support companies at different maturity levels, iLEAP also defines a third transaction.

**Data Transaction 3** enables the exchange of Transport Activity Data, that is, data about the movement of goods or energy consumed, without emissions data. This data transaction was purposely designed for SMEs and software providers [5] who have not yet fully developed ISO 14083 and GLEC Framework capabilities or are not in the business of emissions calculations but would like to assist other stakeholders in doing so.

[5] Such as providers offering Telematics Services, Transport Management Systems, Warehouse Management Software or Charging Point Operators.

> *Figure 2: iLEAP Data Transactions*

All Data Transactions combined, form the “iLEAP Protocol”. The protocol serves the need to achieve interoperability between independent implementations of the iLEAP Technical Specifications, so that emissions data can flow in a seamless, automatic, and scalable manner.

> “We're proud to help the Smart Freight Centre create an industry standard that simplifies data interface processes for transport and logistics companies. This allows them to focus on what truly matters: reducing emissions.”
>
> **Philipp Huhn - Chief Data Officer, shipzero**

### Unique Value Proposition

iLEAP’s distinctiveness comes from combining a data exchange protocol with a data model to achieve end-to-end logistics emissions calculation through IT systems’ interoperability:

*   The **protocol** enables interoperable data flow between systems in a peer-to-peer approach.
*   The **data model** is built upon ISO 14083 and GLEC Framework semantics.

A particularly valuable feature is the inclusion of **data quality metrics** such as Primary Data Share [6], which indicates the proportion of primary data used to calculate total emissions or intensity values, or metrics indicating whether calculations used iLEAP-certified software.

A notable benefit of the iLEAP protocol is the inherent **standardization and consolidation of data flows and processes**, and through this a **significant reduction of complexity** for all parties. Even with complex transport chains or when complex subcontracting cases occur, only the immediate business partners ultimately need to execute the same two (or three) data transactions to achieve emissions transparency.

### Data Sovereignty and Openness

iLEAP supports companies in executing interoperable data flows while maintaining sovereignty:

*   Direct, immediate data exchange between parties.
*   No intermediary infrastructure or party required.
*   Both parties maintain control over data access and authentication.

As an open standard, iLEAP is designed for broad adoption within the sector:

*   Any individual or company can implement iLEAP free of charge and without restrictions. This includes software providers with commercial offerings.
*   Every organization can engage in conformance testing and implementation to assert its readiness to fully engage in the broader iLEAP ecosystem.

[6] See Pathfinder Framework Version 2 and Version 3 for further details.

With this approach, the iLEAP Technical Specifications also address geopolitical considerations. Prioritizing openness and connectivity allow each party to retain the necessary control over where they store data and how it flows. Consequently, iLEAP enhances companies' ability to swiftly adapt to evolving political and regulatory situations at minimal cost.

## THE ILEAP COMMUNITY

The iLEAP Community is an informal group of companies that operate close to logistics value chains. These companies came together to make collective decisions relating to the iLEAP Technical Specifications. Stakeholders include shippers or transport service users, transport service organizers, transport operators, tool providers, representatives from non-profits, or public authorities.

Building on the methods and lessons learned from previous digitization projects such as the SFC Data Access project, SFC Exchange Network, and WBCSD PACT, SFC and SINE designed the community to be open to as many participants as possible, including non-SFC members, while maintaining clear and explicit decision-making processes to facilitate consensus.

### Open And Collaborative Design Approach

SFC and SINE together with the iLEAP Community developed the iLEAP Technical Specifications through a stakeholder-focused process:

1.  Initial workshops to define the scope and business requirements for the iLEAP Technical Specifications.
2.  A series of workshops and working group meetings to decide on technical aspects.
3.  Ongoing feedback collection, especially from iLEAP Implementers.
4.  Lastly, a public consultation process leading to the publication of Version 1.0 of the iLEAP Technical Specifications.

In total, more than 35 independent industry stakeholders participated in this process, holding more than 50 workshops and working group meetings.

### Validation and Interoperability Testing

To validate and report back on the real-world applicability of the iLEAP Technical Specifications, SFC and SINE also opted for a community-driven approach, where independent software implementations test their interoperability with each other. The Automated Conformance Testing Tool, which supports this process, has been utilized over 500 times, demonstrating strong interest in conformance verification.

The Open Logistics Foundation joined this validation process, implementing iLEAP in open source software together with several of its key stakeholders as part of the Enabling Logistics Decarbonization Working Group.

### Why Organizations Choose iLEAP

Even before the publication of the Version 1.0 of the iLEAP Technical Specifications, multiple companies publicly committed to the operational use of the iLEAP Technical Specifications. Organizations choose iLEAP for emissions data exchange as it is:

1.  Founded on well-established methodologies, ISO 14083 and the GLEC Framework.
2.  Steered, sponsored and developed by two neutral organizations, SFC and SINE.
3.  Built for, by, and with companies operating close to logistics value chains.
4.  Addressing the gap on “how” to facilitate the flow of emissions and activity data.
5.  An open-source standard that is free to implement for all parties globally.

Some testimonials from companies are presented below:

> “Joining iLEAP strengthens our strategic advantage through quicker integrations and more seamless partnerships across the logistics sector.”
>
> **Gabriel Beslic - Chief Product Officer, Gryn Network**

> “The launch of iLEAP is a strong first step toward sharing logistics emission data in the chain. Now the real journey begins: making it work in day-to-day operations. That means improving data availability and quality, aligning processes and systems across companies, and building the trust and collaboration needed throughout the value chain to turn the new generated insights into action.”
>
> **Inge Tanke - Co-owner Sustainable Logistics, AllChiefs**

> “iLEAP lays the foundation for a more connected and accountable logistics ecosystem. It reduces integration costs, enhances emissions data quality, and makes it easier to deliver actionable insights to all supply chain stakeholders using our platform, from global cargo owners to carriers, all striving to decarbonize as a core part of their competitive strategy.”
>
> **Karsten Kopland - Head of Product Management / CPO, Kinver**

## CALL TO ACTION

**Join iLEAP**: By joining the project, you get full access to iLEAP’s digital public goods, including the iLEAP Technical Specifications. You can also collaborate with other iLEAP participants and learn about best practices for implementing iLEAP, whether you are an implementer or a “buyer” of iLEAP-conforming software.

**Contribute** to iLEAP: By joining iLEAP, you can materially contribute to the further development of iLEAP and its digital public goods, such as the iLEAP Technical Specifications, while gaining valuable market insights into logistics decarbonization.

**Implement** iLEAP: Everyone can implement iLEAP free of charge for all purposes, including commercial ones. As an implementer, you can immediately benefit from the iLEAP data model and protocol, which streamline implementation of the GLEC Framework and ISO14083 in multiple ways. In addition, you can engage with the iLEAP community and SFC to become iLEAP certified.

**Pilot** iLEAP: Pilot iLEAP by using iLEAP-conforming software, either built in-house or existing tools and solutions. You can then access emissions and intensity more easily, providing the basis for improved decarbonization planning and simulation, as well as greater transparency on low-carbon logistics services.

*Table 1: Benefits from implementing the iLEAP Technical Specifications and joining the community*

| Area & Benefits | Logistics Company (Transport Service Organizer or Operator) | Shipper | Tool Provider | Public Authority, Civil Sector, NGO |
| :---: | :--- | :---: | :---: | :--- |
| **Operational Use** | Lower reporting costs and efforts<br>Easier access to primary data<br>Secured customer relationships<br>Increased competitiveness<br>Increasing trust in and growing demand for low carbon services | Simplified access to emissions data at low cost<br>Gain transparency over primary data (share) and related data quality KPIs<br>One approach for emissions data at corporate- and at service-level (including PCFs)<br>Access to decision grade emissions intensity data for reducing Scope 3 emissions | Reduced GLEC FW and ISO 14083 implementation costs and efforts<br>Reputational gains from standards-based approach<br>Increased visibility through the iLEAP community<br>Lower costs and efforts to access and make available emissions data<br>Ecosystem interoperability, such as PCF calculations and PACT integration | Enhanced opportunities for data derived policy design and green procurement<br>Informed target setting<br>Cross-sector progress comparison |
| **Community Participation** | Reputational gains<br>Access to market-leading decarbonization knowledge<br>Access to all iLEAP digital goods<br>Participation in standards development and strategic decision-making | | | Collaboration with other authorities for harmonized approaches |
| **Standards Development** | Partake in shaping and defining iLEAP’s major standards and digital public goods<br>Assure and deliver a market standard ready for operational use | | | Contribution to B2G use cases within iLEAP<br>Foundation for emissions intensity databases |

### How To Join

***Logistics companies, shippers, tool providers, NGOs, as well as entities from the public and civil sectors, can join the iLEAP by sending an email to team@ileap.global or submitting their interest through the website.***

## FUTURE ROADMAP

The immediate focus is on further developing iLEAP’s governance, growing the base of participants, and enhancing the iLEAP Technical Specifications through community feedback and the development of new features.

### Governance Development

Considering the community's recent significant growth, a more formal governance approach will be necessary in the mid-term. This includes the formalization of participation rules, rights, and obligations, as well as the decision-making processes around working groups and technical specifications.

Acknowledging the importance of development speed of the industry, SFC and SINE intend to maintain a pragmatic approach to iLEAP’s governance with minimal "bureaucracy."

### Community Growth Strategy

The iLEAP community must further expand its global membership specifically with the following three stakeholder groups:

*   **Shippers**: further increase demand for iLEAP and iLEAP-ready software solutions.
*   **Tool providers and software solutions**: to make iLEAP accessible to as many stakeholders as possible.
*   **Small and Medium Enterprises Carriers**: to properly represent the interests of and assert the usability of iLEAP for the majority of logistics companies worldwide.

Lastly, SFC and SINE will extend their collaboration efforts with policymakers and government organizations to ensure alignment between iLEAP and national policies, regulatory frameworks, and digital ecosystems. With this combined approach, the two organizations aim for convergence and an increase in demand for iLEAP-compliant emissions data in the logistics sector worldwide.

### Enhancing SFC’s Decarbonization Initiatives with iLEAP

SFC is currently leading several initiatives to decarbonize the logistics sector. Through iLEAP, SFC will amplify and scale its efforts by applying the iLEAP standardization process and collaborating with iLEAP participants in the following areas:

*   **Book and Claim, Market Based Measures Program:** The iLEAP Technical Specifications will be enhanced with transparency for Book and Claim and related market-based measures.
*   **Clean Cargo:** iLEAP will work to harmonize Clean Cargo's established data collection methodologies with iLEAP's advanced data access, automation, and interoperability capabilities.

### iLEAP Certification Practices

Certification of iLEAP implementations in software will be integrated into SFC’s certification practices, in addition to its current ISO 14083 and GLEC Framework certification. SFC is committed to test and certify software against iLEAP and will apply integration and other automation software that is already applied within the iLEAP Community.

### Feature Expansion

The iLEAP community plans to extend the scope of its work to go *beyond* carbon reporting. Following the iLEAP vision statement, the community will bring forward:

1.  **Interoperability measures that reduce organizational barriers**: The iLEAP Technical Specifications reduce barriers at the technical and semantic level by enabling interoperability. The community aims to further reduce barriers at the organizational and legal domains as well, for instance by proposing model contractual terms [7].
2.  **Procurement-relevant data:** iLEAP already makes procurement-relevant data available through its Data Transaction 2. Following feedback from iLEAP participants, further enhancements to iLEAP are expected to provide greater access to this data.
3.  **Benchmarking**: Enabling members to compare iLEAP-related KPIs with peers while maintaining data confidentiality and compliance [8].
4.  **Knowledge Products**: Delivering documentation about the need for and use of iLEAP. Additionally, knowledge products supporting its application by providing highly targeted documentation and guidance for specific use cases will be developed.
5.  **Secondary Data as Software**: Providing regularly updated GLEC emission factors and intensity values in machine-readable formats, including databases, APIs, libraries, or CSV files.

**By joining iLEAP today, your organization can help shape these developments while positioning yourself at the forefront of logistics decarbonization.**

[7] Model Contractual Terms (MCTs) in the iLEAP context are proposed terms for contracts between companies in the logistics value chain, covering e.g. emissions data governance aspects.
Example of MCTs existing in the Data Act context: https://ec.europa.eu/transparency/expert-groups register/screen/meetings/consult?lang=en&meetingId=61683&fromExpertGroups=3840.

[8] Especially and including anti-trust regulation.

## GLOSSARY

| Term | Definition |
| :--- | :--- |
| **Community** | The group of contributors to the iLEAP project. |
| **GLEC** | Global Logistics Emissions Council. Also see this website |
| **iLEAP** | Integrating Logistics Emissions and Product Carbon Footprints. SFC and SINE co-sponsored this project. Depending on context, iLEAP can also mean the community or members of iLEAP |
| **Logistics company** | any company offering logistics services (transport service organizer or transport operator) |
| **PCF** | Product Carbon Footprint. Also see Pathfinder Framework V2 or later for further details |
| **Scope 3 (emissions)** | In this document, upstream Scope 3 emissions from services or products. Also see the GHG Protocol Product Standard, the Pathfinder Framework V2, and related standards such as ISO 14067 for further details. |
| **SINE, SINE Foundation** | The non-profit organization based in Berlin. |
| **SFC** | Smart Freight Centre organization based in Amsterdam |
| **SME** | Small and Medium-sized Enterprise |
| **Transport Operator** | see ISO 14083 for the definition, also known as “Carrier” |
| **Transport Service Organizer** | see ISO 14083 for the definition, also known as “Logistics Service Provider” |
| **Transport Service User** | see ISO 14083 for the definition, also known as “Shipper” |
