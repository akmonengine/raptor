package raptor

import (
	"github.com/go-gl/mathgl/mgl32"
	"maps"
	"math"
	"math/rand"
	"slices"
	"time"
)

const (
	WORLD_SPACE = iota
	LOCAL_SPACE
)

const (
	STATUS_DOWN = iota
	STATUS_UP
)

type ParticleBuilderFnInterface func(emitter Emitter) Particle

type CurveFloat32 struct {
	keys   []float32
	Values map[float32]float32
}

type Emitter struct {
	Enabled  bool
	Duration float32 // Duration of the Emitter IF looping is false
	Looping  bool
	Delay    float32

	LifeTime          float32 // How long will a particle live
	LifeTimeVariation float32
	EmissionPerSecond int // number of particles per second
	Space             int // WORLD or LOCAL space

	Velocity          mgl32.Vec3
	VelocityVariation mgl32.Vec3
	Position          mgl32.Vec3
	PositionVariation mgl32.Vec3
	Rotation          CurveFloat32
	RotationVariation float32
	Scale             CurveFloat32
	ScaleVariation    float32
	Opacity           CurveFloat32
	OpacityVariation  float32

	GravityEffect float32

	ParticleBuilderFn ParticleBuilderFnInterface
	particles         []Particle
	status            int
	upAt              time.Time
	totalEmitted      int
	durationTimer     *time.Timer
	delayTimer        *time.Timer
}

type Particle struct {
	LifeTime      float32
	LifeRemaining float32
	Velocity      mgl32.Vec3
	Position      mgl32.Vec3
	Rotation      float32
	Scale         float32
	Opacity       float32

	initialRotationVariation float32
	initialScaleVariation    float32
	initialOpacityVariation  float32
}

func (emitter *Emitter) Start() {
	if emitter.Enabled && emitter.status == STATUS_DOWN {
		emitter.delayTimer = time.NewTimer(time.Duration(emitter.Delay * float32(time.Second)))
		emitter.durationTimer = &time.Timer{}

		if len(emitter.particles) == 0 {
			emitter.particles = make([]Particle, 0, emitter.EmissionPerSecond)
		}
	}
}

func (emitter *Emitter) Stop() {
	if emitter.status == STATUS_UP {
		emitter.totalEmitted = 0
		emitter.delayTimer.Stop()

		if emitter.durationTimer != nil {
			emitter.durationTimer.Stop()
		}
		emitter.status = STATUS_DOWN
	}
}

func (emitter *Emitter) GenerateParticles(maxParticles int, modelMatrix mgl32.Mat4, timeScale float64) {
	if emitter.status == STATUS_UP {
		elapsedTime := time.Since(emitter.upAt)
		quantityToGenerate := (elapsedTime.Seconds() * float64(emitter.EmissionPerSecond)) - float64(emitter.totalEmitted)
		quantityToGenerate = quantityToGenerate * timeScale
		for i := 0; i < int(math.Round(quantityToGenerate)); i++ {
			if len(emitter.particles) < maxParticles {
				emitter.addParticle(modelMatrix)
			}

			emitter.totalEmitted++
		}
	}
}

func (emitter *Emitter) Compute(elapsedTime float32, gravity float32) {
	if !emitter.Enabled || emitter.delayTimer == nil {
		return
	}

	select {
	case <-emitter.delayTimer.C:
		if !emitter.Looping {
			emitter.durationTimer = time.NewTimer(time.Duration(emitter.Duration * float32(time.Second)))
		}
		emitter.upAt = time.Now()
		emitter.status = STATUS_UP
	case <-emitter.durationTimer.C:
		emitter.Stop()
	default:
	}

	for i := 0; i < len(emitter.particles); i++ {
		gravityEffect := emitter.GravityEffect * gravity * elapsedTime
		emitter.particles[i].update(emitter, elapsedTime, gravityEffect)

		if emitter.particles[i].LifeRemaining <= 0 {
			emitter.particles = append(emitter.particles[:i], emitter.particles[i+1:]...)
			i-- // -1 as the slice just got shorter
		}
	}
}

func (emitter *Emitter) GetParticles() []Particle {
	return emitter.particles
}

func (emitter *Emitter) addParticle(modelMatrix mgl32.Mat4) {
	var particle Particle

	if emitter.ParticleBuilderFn != nil {
		particle = emitter.ParticleBuilderFn(*emitter)
	} else {
		particle = particleBuilderFn(*emitter)
	}

	if emitter.Space == WORLD_SPACE {
		particle.Position = mgl32.TransformCoordinate(particle.Position, modelMatrix)
		particle.Velocity = mgl32.Mat4ToQuat(modelMatrix).Rotate(particle.Velocity)
	}

	emitter.particles = append(emitter.particles, particle)
}

func particleBuilderFn(emitter Emitter) Particle {
	life := emitter.LifeTime + (emitter.LifeTimeVariation * (rand.Float32() - 0.5))
	position := emitter.Position
	position = position.Add(mgl32.Vec3{
		emitter.PositionVariation.X() * (rand.Float32() - 0.5),
		emitter.PositionVariation.Y() * (rand.Float32() - 0.5),
		emitter.PositionVariation.Z() * (rand.Float32() - 0.5),
	})

	velocity := emitter.Velocity
	velocity = velocity.Add(mgl32.Vec3{
		emitter.VelocityVariation.X() * (rand.Float32() - 0.5),
		emitter.VelocityVariation.Y() * (rand.Float32() - 0.5),
		emitter.VelocityVariation.Z() * (rand.Float32() - 0.5),
	})

	rotationVariation := emitter.RotationVariation * (rand.Float32() - 0.5)
	scaleVariation := emitter.ScaleVariation * (rand.Float32() - 0.5)
	opacityVariation := emitter.OpacityVariation * (rand.Float32() - 0.5)

	return Particle{
		LifeTime:      life,
		LifeRemaining: life,
		Velocity:      velocity,
		Position:      position,
		Rotation:      emitter.Rotation.lerp(0.0) + rotationVariation,
		Scale:         emitter.Scale.lerp(0.0) + scaleVariation,
		Opacity:       emitter.Opacity.lerp(0.0) + (opacityVariation * emitter.Opacity.lerp(0.0)),

		initialRotationVariation: rotationVariation,
		initialScaleVariation:    scaleVariation,
		initialOpacityVariation:  opacityVariation,
	}
}

func (particle *Particle) update(emitter *Emitter, time float32, gravityEffect float32) {
	life := particle.LifeRemaining / particle.LifeTime

	particle.LifeRemaining -= time
	particle.Velocity = particle.Velocity.Add(mgl32.Vec3{0.0, gravityEffect, 0.0})
	particle.Position = particle.Position.Add(particle.Velocity.Mul(time))
	particle.Rotation = emitter.Rotation.lerp(1.0-life) + particle.initialRotationVariation
	particle.Scale = emitter.Scale.lerp(1.0-life) + particle.initialScaleVariation
	particle.Opacity = emitter.Opacity.lerp(1.0-life) + (particle.Opacity * particle.initialOpacityVariation)
}

func NewCurveFloat32(values map[float32]float32) CurveFloat32 {
	keys := slices.Collect(maps.Keys(values))
	slices.Sort(keys)

	return CurveFloat32{
		keys:   keys,
		Values: values,
	}
}

func (curveFloat32 CurveFloat32) lerp(t float32) float32 {
	if t < 0 {
		return 0.0
	}
	if t > 1 {
		return 1.0
	}

	var previousK float32
	var value float32
	for _, k := range curveFloat32.keys {
		if t < k {
			factor := (k - t) / (k - previousK)
			value = curveFloat32.Values[previousK]*factor + (curveFloat32.Values[k] * (1 - factor))
			break
		} else {
			previousK = k
			value = curveFloat32.Values[k]
		}
	}

	return value
}
