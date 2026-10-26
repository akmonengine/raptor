package raptor

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

// emissionEmitter is an emitter of 100 particles a second, living 10 s: none dies during a test.
func emissionEmitter(looping bool, delay, duration float32) *Emitter {
	return &Emitter{
		Enabled: true, Looping: looping, Delay: delay, Duration: duration,
		LifeTime: 10, EmissionPerSecond: 100,
		Scale:   NewCurveFloat32(map[float32]float32{0: 1}),
		Opacity: NewCurveFloat32(map[float32]float32{0: 1}),
	}
}

// emit runs the emitter by the steps given (seconds), each a birth then a step, and counts the
// particles alive.
func emit(emitter *Emitter, steps ...float32) int {
	for _, step := range steps {
		emitter.Compute(step, -9.8)
		emitter.GenerateParticles(100000, mgl32.Ident4())
	}
	return len(emitter.GetParticles())
}

// repeat is n steps of seconds each.
func repeat(n int, seconds float32) []float32 {
	steps := make([]float32, n)
	for i := range steps {
		steps[i] = seconds
	}
	return steps
}

// A second of the time the emitter is given, at 100 particles a second, gives birth to 100 particles
// however it is cut into steps — the wall clock plays no part.
func TestASecondGivesItsParticlesWhateverItsSteps(t *testing.T) {
	for name, steps := range map[string][]float32{
		"one step":         {1},
		"ten steps":        repeat(10, 0.1),
		"a hundred steps":  repeat(100, 0.01),
		"uneven steps":     {0.3, 0.7},
		"sixty a second":   repeat(60, 1.0/60),
		"a thousand steps": repeat(1000, 0.001),
	} {
		emitter := emissionEmitter(true, 0, 0)
		emitter.Start()
		if got := emit(emitter, steps...); got != 100 {
			t.Errorf("%s: %d particles born in a second, want 100", name, got)
		}
	}
}

// A delay of 0.5 s gives birth to nothing before 0.5 s of the time given, and then to what the time past
// it owes.
func TestTheDelayCountsTheTimeGiven(t *testing.T) {
	emitter := emissionEmitter(true, 0.5, 0)
	emitter.Start()
	if got := emit(emitter, repeat(49, 0.01)...); got != 0 {
		t.Fatalf("0.49 s into a delay of 0.5 s, %d particles born, want none", got)
	}
	if got := emit(emitter, repeat(51, 0.01)...); got != 50 {
		t.Fatalf("1 s into a delay of 0.5 s, %d particles born, want the 50 of 0.5 s", got)
	}
}

// A non-looping emitter of 0.5 s gives birth to the particles of 0.5 s, then to none; its particles live
// on.
func TestTheDurationStopsTheEmission(t *testing.T) {
	emitter := emissionEmitter(false, 0, 0.5)
	emitter.Start()
	if got := emit(emitter, repeat(30, 0.03)...); got != 50 {
		t.Fatalf("0.9 s of an emitter of 0.5 s, %d particles born, want 50", got)
	}
	if got := emit(emitter, 1); got != 50 {
		t.Fatalf("after its duration, %d particles, want the 50 born", got)
	}
}

// Without time — a game paused —, nothing is born, however long the wall clock runs.
func TestNoTimeGivesNoParticle(t *testing.T) {
	emitter := emissionEmitter(true, 0, 0)
	emitter.Start()
	if got := emit(emitter, repeat(100, 0)...); got != 0 {
		t.Fatalf("a hundred steps of no time give %d particles, want none", got)
	}
}

// A step of emission allocates nothing once the emitter holds its particles.
func TestAStepOfEmissionAllocatesNothing(t *testing.T) {
	emitter := emissionEmitter(true, 0, 0)
	emitter.LifeTime = 0.5
	emitter.Start()
	emit(emitter, repeat(200, 0.01)...)
	if allocs := testing.AllocsPerRun(100, func() { emit(emitter, 0.01) }); allocs != 0 {
		t.Fatalf("a step of emission allocated %v times, want 0", allocs)
	}
}
