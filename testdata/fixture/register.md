# Plumbline fixture register

A tiny C4 register used to prove the engine end to end. It exercises the full
requirement→architecture→code chain and, by construction, one of each gap the
engine must report: an uncovered requirement, a transitive defect, a broken
anchor and an orphan code-area.

## Features

### Authentication
`feat~authentication~1`
Status: approved

Users can authenticate before accessing the system.

Needs: req

### Key management
`feat~key-management~1`
Status: approved

Signing keys are managed on a schedule.

Needs: req

## Requirements

### Validate authentication request
`req~validate-auth-request~1`
Status: approved

Every authentication request is validated before access is granted.

Covers: feat~authentication~1
Needs: component

### Rotate signing keys
`req~rotate-signing-keys~1`
Status: approved

Signing keys are rotated on a fixed schedule. Deliberately left with no covering
component, to demonstrate an uncovered requirement — which in turn leaves the
feature above it a transitive defect.

Covers: feat~key-management~1
Needs: component

## Components

### Auth validator
`component~auth-validator~1`
Status: approved

Validates credentials and issues session tokens.

Covers: req~validate-auth-request~1
Needs: impl
