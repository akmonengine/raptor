package raptor

import (
	"github.com/go-gl/mathgl/mgl32"
	"maps"
	"slices"
	"testing"
	"time"
)

func BenchmarkParticleBuilderFn(b *testing.B) {
	emitter := Emitter{
		Enabled:           true,
		Duration:          1.0,
		Looping:           true,
		Delay:             0.0,
		LifeTime:          0.5,
		LifeTimeVariation: 0.1,
		Space:             WORLD_SPACE,
		EmissionPerSecond: 100,
		Velocity:          mgl32.Vec3{0.1, 0.5, 0.2},
		VelocityVariation: mgl32.Vec3{0.02, 0.2, 0.1},
		Position:          mgl32.Vec3{0.0, 0.0, 0.0},
		PositionVariation: mgl32.Vec3{0.0: 0.1},
		Rotation:          CurveFloat32{},
		RotationVariation: 0.1,
		Scale: CurveFloat32{
			keys:   slices.Collect(maps.Keys(map[float32]float32{0.0: 0.1})),
			Values: map[float32]float32{0.0: 0.1},
		},
		ScaleVariation: 0.1,
		Opacity: CurveFloat32{
			keys:   slices.Collect(maps.Keys(map[float32]float32{0.0: 0.1})),
			Values: map[float32]float32{0.0: 0.1},
		},
		OpacityVariation: 0.0,
		GravityEffect:    0.2,
	}

	for b.Loop() {
		particleBuilderFn(emitter)
	}

	b.ReportAllocs()
}

func BenchmarkEmitter_GenerateParticles(b *testing.B) {
	emitter := Emitter{
		Enabled:           true,
		Looping:           true,
		LifeTime:          1.0,
		LifeTimeVariation: 0.0,
		EmissionPerSecond: 500000,
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
		GravityEffect: 1.0,
	}
	emitter.Start()
	emitter.Compute(1.0, -9.8)
	time.Sleep(time.Second * 1)

	modelMatrix := mgl32.Ident4()
	for b.Loop() {
		emitter.GenerateParticles(1000000, modelMatrix, 1.0)
		emitter.particles = emitter.particles[:0]
		time.Sleep(time.Second * 1)
	}

	b.ReportAllocs()
}

func BenchmarkEmitter_Compute(b *testing.B) {
	emitter := Emitter{
		Enabled:           true,
		Looping:           true,
		LifeTime:          1.0,
		LifeTimeVariation: 0.0,
		EmissionPerSecond: 500000,
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
		GravityEffect: 1.0,
	}
	emitter.Start()

	modelMatrix := mgl32.Ident4()
	for b.Loop() {
		t := time.NewTimer(time.Second * 10)
		elapsedTime := time.Now()

	loop:
		for {
			emitter.GenerateParticles(1000000, modelMatrix, 1.0)
			emitter.Compute(float32(time.Since(elapsedTime).Seconds()), -9.8)
			elapsedTime = time.Now()

			select {
			case <-t.C:
				break loop
			default:
			}
		}
	}

	b.ReportAllocs()
}

func TestCurveFloat32_lerp(t *testing.T) {
	type args struct {
		t float32
	}
	tests := []struct {
		name         string
		curveFloat32 CurveFloat32
		args         args
		want         float32
	}{
		// TODO: Add test cases.
		{name: "lerp01", curveFloat32: NewCurveFloat32(map[float32]float32{0.0: 0.0, 1.0: 1.0}),
			args: args{t: 0.0}, want: 0.0},
		{name: "lerp02", curveFloat32: NewCurveFloat32(map[float32]float32{0.0: 0.0, 1.0: 1.0}),
			args: args{t: 1.0}, want: 1.0},
		{name: "lerp03", curveFloat32: NewCurveFloat32(map[float32]float32{0.0: 0.0, 1.0: 1.0}),
			args: args{t: -0.5}, want: 0.0},
		{name: "lerp04", curveFloat32: NewCurveFloat32(map[float32]float32{0.0: 0.0, 1.0: 1.0}),
			args: args{t: 1.3}, want: 1.0},
		{name: "lerp05", curveFloat32: NewCurveFloat32(map[float32]float32{0.0: 0.0, 0.8: 1.0}),
			args: args{t: 0.5}, want: 0.625},
		{name: "lerp06", curveFloat32: NewCurveFloat32(map[float32]float32{0.0: 0.0, 0.5: 1.0, 1.0: 0.0}),
			args: args{t: 0.7}, want: 0.6},
		{name: "lerp07", curveFloat32: NewCurveFloat32(map[float32]float32{0.0: 0.0, 0.8: 0.5, 1.0: 1.0}),
			args: args{t: 0.6}, want: 0.375},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.curveFloat32.lerp(tt.args.t); got != tt.want {
				t.Errorf("lerp() = %v, want %v", got, tt.want)
			}
		})
	}
}
