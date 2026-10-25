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
Note that this call only applies if the Emitter was not already started: it then waits for its Delay, and
emits for its Duration unless it loops. Both are counted on the time given to `Compute`, never on the wall
clock: pausing, stepping or scaling the time of your game does the same to its emitters.

If you need to stop an Emitter, you can call the stop function.
```go
emitter.Stop()
```
The Particles already born live on; `Start` starts the Emitter again, its delay first.

### Update
Each time your system needs to update the Particles, two calls are required:
To advance the Emitter and its Particles by the time elapsed, call the Compute method in your game loop:
```go
emitter.Compute(0.01, -9.8)
```
- The first argument is the time elapsed since the previous call, in seconds: your game's time, scaled
as your game scales it (x2.0 gives birth to twice the Particles in a second of the wall clock), 0 when
it is paused.
- The second argument is the gravity applied in your world. Terrestrial gravity is 9.8m/s².
But your world could have a different gravity force, feel free to play with this parameter.

Then give birth to the Particles the time owes:
```go
parentModelMatrix := mgl32.Ident4()
emitter.GenerateParticles(100, parentModelMatrix)
```
- The first argument is a hard limit of the maximum particles allowed to exist in an emitter.
It can be useful if your game engine sets this kind of configurable limits, for performance reasons.
- The second argument is the parent Model Matrix (i.e. the computed position), so that the particle
lives related to its parent's transform.

An Emitter of 100 Particles a second gives birth to 100 Particles in a second of the time given,
however it is cut into steps.

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
