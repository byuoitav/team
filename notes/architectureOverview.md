# AV API


## Overview
Major players: 

1. AV-Control API
1. Configuration Database
1. Event System
1. Monitoring Services [planned]
1. Deployment Pipeline

Other
1. Pi Setup

## AV-Control API
The AV control api exists primarily to provide an a) abstracted and b) RESTful interface to control a room

Major players in the AV-Control API

1. AV-API 
1. Supporting Microservices 
1. Configuration Database

Important Concepts

1. [RESTful Request Body](https://github.com/byuoitav/av-api/wiki/PUT-Body-Definition)
    * The abstracted definition of *controllable* elements of a room
1. [Actions](https://github.com/byuoitav/av-api/wiki/Logical-Constructs-Internal-to-the-API-Service#action-structure)
    * Actions that will set the state of a device. Take the form of a microservice, endpoint, and paramters
    * Endpoints are parameterized
    * RPC endpoints
    * Also contain STATUS events to be published on successful execution.
1. [Command Evaluators](https://github.com/byuoitav/av-api/wiki/Logical-Constructs-Internal-to-the-API-Service#action-structure)
    * Command Evaluators take a *RESTful Request Body* and generate *Actions*
1. Status Evaluators
    * Status Evaluators generate Status Actions, which contain the information necessary to retrieve the state of devices
1. [Room configuration](https://github.com/byuoitav/av-api/wiki/Logical-Constructs-Internal-to-the-API-Service#action-structure)
    * Room configurations map which evaluators are used to evaluate requests against a room. 
    * Contain a *Reconciler* which is used to define logic on how different commands (generated from different evaluators) interact within a request/room. 


#### Command Evaluators
Look at flow. Main points:

1. Evaluate -> generates actions
1. Validate -> validates actions
1. Incompatable actions -> May be used in reconcilers to define "incomptable actions"

The idea is that you define logic that determines how one aspect of a room in an evaluator, if you need to change ONE aspect of a room, you have to create a new configuration, but may reuse most of the evaluators. 

Really. Look at flow doc. Linked above.

Overview: Think of it like an assembly line.

```
Request -Put Body-> Evaluators.Evaluate/Validate -Actions-> Gather responses -Actions-> Reconcile -Actions-> ExecuteActions -Response-> Return Response
```

#### Status Evaluators 

1. Generate Actions -> Determines how to retrieve the status of a given attribute for a device, contains information about where that information belongs in the evental response
1. Evaluate Response -> Can perform mapping on the responses from the microservices i.e. remapping ports to devices

#### Supporting Microservices
* Define how to talk to individual devices (Protocols) 
* Can expose lots of functionality: AV-Control API is mostly concerned with the *Standard Endpoints*, which are RPC (Lowest common demoninator)
* Dumb, don't know about any other microservices/Database

#### Standard Endpoints 
Look in the DB. But they pretty much look like:

`/:address/power/on`
`/:address/input/:port`
`/:address/input/:in/:out`

Where the elements preceeded with a `:` are *parameters*
Are defined in the DB, each command that a device may issue is linked to an endpoint and microservice. This means that two different devices don't have to use the same endpoint to issue the same command. In fact standard endpoints don't really have to be standard, nor do command names need to be standard, but by convention they are. 


## Configuration Database 
This contains all the information necessary to 
