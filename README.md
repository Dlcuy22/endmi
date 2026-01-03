# ENDMI  
**A Simple Golang Project Manager**

ENDMI is a lightweight project management tool designed to streamline Golang project creation, experimentation, and extension. It focuses on speed, structure, and safety, especially when testing temporary code or prototypes.

## Overview

Managing Golang projects often involves repetitive setup, manual folder creation, and cleanup after experiments. ENDMI removes that friction by providing templated project generation and a built-in temporary workspace concept, allowing developers to focus on writing code instead of managing files.


## Key Features

### 1. Project Creation Using Templates
- Quickly generate Golang projects from predefined templates
- Maintain consistent project structure across teams
- Reduce setup time for new services, tools, or experiments

### 2. Extensible Architecture
- Add custom extensions to support additional workflows
- Adapt ENDMI to different project styles or internal standards
- Keep the core simple while allowing advanced customization

### 3. Temporary Code Workspace (TempCode)
- Write and test Golang code without creating permanent folders
- Temporary projects are treated as disposable by default
- Ideal for:
  - Prototyping
  - Experimenting with APIs or libraries
  - Testing small code snippets safely

This prevents unused test folders from polluting your workspace while encouraging experimentation.

## Use Cases

- Bootstrapping new Golang projects
- Rapid prototyping and proof-of-concept development
- Testing ideas without committing to project structure
- Maintaining clean repositories and local environments

## Philosophy

ENDMI follows three core principles:

- **Simplicity**: Minimal setup, predictable behavior
- **Safety**: Temporary code should stay temporary unless promoted
- **Productivity**: Less project management, more development

## Status

This project is under active development. Features and APIs may evolve as the tool matures.


## License

MIT
---

## Contributing

Contributions, suggestions, and extensions are welcome.  
Please open an issue or submit a pull request to participate.

