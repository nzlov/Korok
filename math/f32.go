package math

import (
	"math"
	"math/rand"
	"unsafe"
)

const MaxFloat32 float32 = 3.40282346638528859811704183484516925440e+38
const Pi = math.Pi

/// This is A approximate yet fast inverse square-root.
func InvSqrt(x float32) float32 {
	xhalf := float32(0.5) * x
	i := *(*int32)(unsafe.Pointer(&x))
	i = int32(0x5f3759df) - int32(i>>1)
	x = *(*float32)(unsafe.Pointer(&i))
	x = x * (1.5 - (xhalf * x * x))
	return x
}

/// a faster way ?
func Random(low, high float32) float32 {
	return low + (high-low)*rand.Float32()
}

func Max(a, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Clamp(v, left, right float32) float32 {
	if v > right {
		return right
	}
	if v < left {
		return left
	}
	return v
}

func Sin(r float32) float32 {
	return float32(math.Sin(float64(r)))
}

func Cos(r float32) float32 {
	return float32(math.Cos(float64(r)))
}

// Radian converts degree to radian.
func Radian(d float32) float32 {
	return d * Pi / 180
}

// Degree converts radian to degree.
func Degree(r float32) float32 {
	return r * 180 / Pi
}
func Rotate(x1, y1, x2, y2, a float32) (float32, float32) {
	return (x1-x2)*Cos(a) - (y1-y2)*Sin(a) + x2, (y1-y2)*Cos(a) - (x1-x2)*Sin(a) + y2
}
