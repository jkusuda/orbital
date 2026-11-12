package orbit

import (
	"math"
	"time"
)

type Elements struct {
	Epoch         float64 // JD (TDB) the elements are valid at
	Eccentricity  float64 // no units, 0 = circle, <1 = ellipse
	SemiMajorAxis float64 // AU
	Inclination   float64 // degrees
	AscendingNode float64 // degrees (longitude of ascending node, Ω)
	PerihelionArg float64 // degrees (argument of perihelion, ω)
	MeanAnomaly   float64 // degrees, at Epoch
}

type Vec3 struct{ X, Y, Z float64 }

func Position(el Elements, t time.Time) Vec3 {
	jd := julianDate(t)
	M := meanAnomaly(el, jd)
	E := solveKepler(M, el.Eccentricity)
	x, y := orbitalPlane(el.SemiMajorAxis, el.Eccentricity, E)
	return rotate(x, y, el)
}

func degToRad(d float64) float64 { return d * math.Pi / 180 }

// Internal steps
func julianDate(t time.Time) float64 { // step 1
	return float64(t.Unix())/86400 + 2440587.5
}

// https://en.wikipedia.org/wiki/Mean_anomaly (step 2)
func meanAnomaly(el Elements, jd float64) float64 {
	const gaussK = 0.01720209895                  // Gaussian gravitational constant
	n := gaussK / math.Pow(el.SemiMajorAxis, 1.5) // mean motion rad/day
	M := degToRad(el.MeanAnomaly) + n*(jd-el.Epoch)

	M = math.Mod(M, 2*math.Pi)
	if M < 0 {
		M += 2 * math.Pi
	}
	return M
}

// https://en.wikipedia.org/wiki/Kepler%27s_equation (step 3)
func solveKepler(M, e float64) float64 {
	E := M
	if e > 0.8 {
		E = math.Pi
	}
	for range 50 {
		next := E - (E-e*math.Sin(E)-M)/(1-e*math.Cos(E))
		if math.Abs(next-E) < 1e-12 {
			return next
		}
		E = next
	}
	return E
}

// https://ssd.jpl.nasa.gov/planets/approx_pos.html (step 4)
func orbitalPlane(a, e, E float64) (x, y float64) {
	x = a * (math.Cos(E) - e)
	y = a * math.Sqrt(1-e*e) * math.Sin(E)
	return x, y
}

// https://ssd.jpl.nasa.gov/planets/approx_pos.html (step 5)
func rotate(x, y float64, el Elements) Vec3 {
	i := degToRad(el.Inclination)
	O := degToRad(el.AscendingNode) // Ω
	w := degToRad(el.PerihelionArg) // ω

	cosO, sinO := math.Cos(O), math.Sin(O)
	cosw, sinw := math.Cos(w), math.Sin(w)
	cosi, sini := math.Cos(i), math.Sin(i)

	return Vec3{
		X: (cosw*cosO-sinw*sinO*cosi)*x + (-sinw*cosO-cosw*sinO*cosi)*y,
		Y: (cosw*sinO+sinw*cosO*cosi)*x + (-sinw*sinO+cosw*cosO*cosi)*y,
		Z: (sinw*sini)*x + (cosw*sini)*y,
	}
}

func EarthPosition(t time.Time) Vec3 {
	jd := julianDate(t)
	T := (jd - 2451545.0) / 36525 // centuries since J2000

	a := 1.00000261 + 0.00000562*T
	e := 0.01671123 - 0.00004392*T
	i := -0.00001531 - 0.01294668*T
	L := 100.46457166 + 35999.37244981*T // mean longitude
	peri := 102.93768193 + 0.32327364*T  // longitude of perihelion ϖ
	node := 0.0                          // Ω

	return Position(Elements{
		Epoch:         jd,
		Eccentricity:  e,
		SemiMajorAxis: a,
		Inclination:   i,
		AscendingNode: node,
		PerihelionArg: peri - node,
		MeanAnomaly:   L - peri,
	}, t)
}

func (v Vec3) Sub(o Vec3) Vec3 {
	return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}
