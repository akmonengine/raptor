# Raptor
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/akmonengine/raptor)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/akmonengine/raptor)
[![Go Report Card](https://goreportcard.com/badge/github.com/akmonengine/raptor)](https://goreportcard.com/report/github.com/akmonengine/raptor)
![Tests](https://img.shields.io/github/actions/workflow/status/akmonengine/raptor/code_coverage.yml?label=tests)
![Codecov](https://img.shields.io/codecov/c/github/akmonengine/raptor)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/akmonengine/raptor)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-pr/akmonengine/raptor)


Raptor is a simple Go particle tool built for game development.
Emitters generate Particles, based on a simple yet powerful configuration.
Depending on this configuration, various effects can be achieved (e.g. smoke, fire, ambient twinkles).

## Dependencies
Raptor relies on https://github.com/go-gl/mathgl for mathematical computations.

## Basic Usage

### Emitter
First, you need to create an Emitter and configure it:
```go
emitter := raptor.Emitter{
    Enabled:           true,
    Looping:           true,
    LifeTime:          1.0,
    LifeTimeVariation: 0.0,
    EmissionPerSecond: 100,
    Velocity:          mgl32.Vec3{0.0, 1.0, 0.0},
    VelocityVariation: mgl32.Vec3{1.0, 0.2, 1.0},
    Position:          mgl32.Vec3{0.0, 0.0, 0.0},
    Rotation: CurveFloat32{
        Values: map[float32]float32{0.0: 0.0, 1.0: 180.0},
    },
    Scale: CurveFloat32{
        Values: map[float32]float32{0.0: 1.0, 1.0: 2.0},
    },
    Opacity: CurveFloat32{
        Values: map[float32]float32{0.0: 0.0, 0.1: 0.9, 1.0: 1.0},
    },
}
```
- The LifeTime property corresponds to the Particles.
- If Looping is set to true, the Emitter will keep generating new particles.
  Otherwise, the Emitter will automatically stop after the Duration property (in seconds).
- A Delay property helps to trigger after the start of the Emitter.
- The Space property (WORLD or LOCAL space) allows to position the particle using its Emitter
position (if WORLD, the Particle will follow a moving Emitter).
- GravityEffect is multiplied to the Gravity and applies a force on the Velocity of the Particle.

### Start & Stop
Once the Emitter is created, it needs to be started.
```go
emitter.Start()
```
Note that this call only applies and sets the status UP if the Emitter was not already started (i.e. status is DOWN).
If its timer is over, it does not automatically reset to the status DOWN.

If you need to stop an Emitter, you can call the stop function.
```go
emitter.Stop()
```
This call applies only if the status is UP. If a timer is set, it is not reset.

### Update
Each time your system needs to update the Particles, two calls are required:
To generate new particles, as the times goes on:
```go
parentModelMatrix := mgl32.Ident4()
emitter.GenerateParticles(100, parentModelMatrix, 1.0)
```
- The first argument is a hard limit of the maximum particles allowed to exist in an emitter.
It can be useful if your game engine sets this kind of configurable limits, for performance reasons.
- The second argument is the parent Model Matrix (i.e. the computed position), so that the particle
lives related to its parent's transform.
- The third argument is the time scale. For example if you decide to accelerate up to x2.0, the number
of generated Particles would need to be multiplied by x2.0.

You can then update the positions of the Particles, calling the Compute method.
This call would be applied in your game loop for example:
```go
emitter.Compute(0.01, -9.8)
```
- The first argument is the time elapsed between the previous call and the current time.
- The second argument is the gravity applied in your world. Terrestrial gravity is 9.8m/s².
But your world could have a different gravity force, feel free to play with this parameter.

### Fetch the Particles
You can retrieve all the Particles attached to an Emitter, as a slice:
```go
emitter.GetParticles()
```
All the Particles are active, so you can simply loop on this slice to render them.

## What is to come next ?
- Feel free to contribute to improve the performances of this module.

## Contributing Guidelines

See [how to contribute](CONTRIBUTING.md).

## Licence
This project is distributed under the [Apache 2.0 licence](LICENCE.md).
