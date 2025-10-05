# OpenTelemetry WASI for TypeScript

This is under active development and not currently working.

# How to generate bindings

https://github.com/bytecodealliance/jco

- Install jco:
    - `npm install @bytecodealliance/jco`

- Generate the types
    - `npx jco guest-types -o ./types ./wit`

# TODO

- Fill out the src implementation