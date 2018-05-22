package effect

import (
	"korok.io/korok/gfx"
	"korok.io/korok/math"
	"korok.io/korok/math/f32"
)

// SnowSimulator can simulate snow effect.
type SnowSimulator struct {
	Pool

	RateController
	LifeController
	VisualController

	velocity channel_v2
	rot      channel_f32
	rotDelta channel_f32

	// Configuration.
	Config struct {
		Duration, Rate float32
		Life           Var
		Size           Var
		Color          f32.Vec4
		Position       [2]Var
		Velocity       [2]Var
		Rotation       Var
	}
}

func NewSnowSimulator(cap int, w, h float32) *SnowSimulator {
	sim := SnowSimulator{Pool: Pool{cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)
	sim.AddChan(Rotation, RotationDelta)

	// config
	sim.Config.Duration = math.MaxFloat32
	sim.Config.Rate = 60
	sim.Config.Life = Var{10, 4}
	sim.Config.Color = f32.Vec4{1, 0, 0, 1}
	sim.Config.Size = Var{6, 6}
	sim.Config.Position[0] = Var{0, w}
	sim.Config.Position[1] = Var{h, 0}
	sim.Config.Velocity[0] = Var{-10, 20}
	sim.Config.Velocity[1] = Var{-50, 20}
	sim.Config.Rotation = Var{0, 360}

	return &sim
}

func (sim *SnowSimulator) Initialize() {
	sim.Pool.Initialize()

	sim.life = sim.Field(Life).(channel_f32)
	sim.size = sim.Field(Size).(channel_f32)
	sim.pose = sim.Field(Position).(channel_v2)
	sim.velocity = sim.Field(Velocity).(channel_v2)
	sim.color = sim.Field(Color).(channel_v4)
	sim.rot = sim.Field(Rotation).(channel_f32)
	sim.rotDelta = sim.Field(RotationDelta).(channel_f32)

	sim.RateController.Initialize(sim.Config.Duration, sim.Config.Rate)
}

func (sim *SnowSimulator) Simulate(dt float32) {
	if new := sim.Rate(dt); new > 0 {
		sim.NewParticle(new)
	}

	n := int32(sim.live)

	// update old particle
	sim.life.Sub(n, dt)

	// position integrate: p' = p + v * t
	sim.pose.Integrate(n, sim.velocity, dt)
	sim.rot.Integrate(n, sim.rotDelta, dt)

	// GC
	sim.GC(&sim.Pool)
}

func (sim *SnowSimulator) Size() (live, cap int) {
	return int(sim.live), sim.cap
}

func (sim *SnowSimulator) NewParticle(new int) {
	if (sim.live + new) > sim.cap {
		return
	}
	start := sim.live
	sim.live += new

	rot := Range{Var{0, 10}, Var{1, 10}}

	for i := start; i < sim.live; i++ {
		sim.life[i] = sim.Config.Life.Random()
		sim.color[i] = sim.Config.Color
		sim.size[i] = sim.Config.Size.Random()
		sim.rot[i], sim.rotDelta[i] = rot.RangeInit(1 / sim.life[i])

		f := sim.size[i] / (sim.Config.Size.Base + sim.Config.Size.Var)
		sim.color[i][3] = f

		px := sim.Config.Position[0].Random()
		py := sim.Config.Position[1].Random()
		sim.pose[i] = f32.Vec2{px, py}

		dx := sim.Config.Velocity[0].Random()
		dy := sim.Config.Velocity[1].Random()
		sim.velocity[i] = f32.Vec2{dx, dy}
	}
}

func (sim *SnowSimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	size := sim.size
	pose := sim.pose
	rotate := sim.rot

	// compute vbo
	for i := 0; i < sim.live; i++ {
		vi := i << 2
		h_size := size[i] / 2

		var (
			r  = math.Clamp(sim.color[i][0], 0, 1)
			g_ = math.Clamp(sim.color[i][1], 0, 1)
			b  = math.Clamp(sim.color[i][2], 0, 1)
			a  = math.Clamp(sim.color[i][3], 0, 1)
		)

		c := uint32(a*255)<<24 + uint32(b*255)<<16 + uint32(g_*255)<<8 + uint32(r*255)
		rg := tex.Region()
		rot := float32(0)
		if len(rotate) > i {
			rot = rotate[i]
		}

		// bottom-left
		buf[vi+0].X, buf[vi+0].Y = math.Rotate(pose[i][0]-h_size, pose[i][1]-h_size, pose[i][0], pose[i][1], rot)
		buf[vi+0].U = rg.X1
		buf[vi+0].V = rg.Y1
		buf[vi+0].RGBA = c

		// bottom-right
		buf[vi+1].X, buf[vi+1].Y = math.Rotate(pose[i][0]+h_size, pose[i][1]-h_size, pose[i][0], pose[i][1], rot)
		buf[vi+1].U = rg.X2
		buf[vi+1].V = rg.Y1
		buf[vi+1].RGBA = c

		// top-right
		buf[vi+2].X, buf[vi+2].Y = math.Rotate(pose[i][0]+h_size, pose[i][1]+h_size, pose[i][0], pose[i][1], rot)
		buf[vi+2].U = rg.X2
		buf[vi+2].V = rg.Y2
		buf[vi+2].RGBA = c

		// top-left
		buf[vi+3].X, buf[vi+3].Y = math.Rotate(pose[i][0]-h_size, pose[i][1]+h_size, pose[i][0], pose[i][1], rot)
		buf[vi+3].U = rg.X1
		buf[vi+3].V = rg.Y2
		buf[vi+3].RGBA = c
	}
}
