package twilight

import (
	"math"
)

func daysSince2000(y, m, d int) float64 {
	return float64(367*(y) - ((7 * ((y) + (((m) + 9) / 12))) / 4) + ((275 * (m)) / 9) + (d) - 730530)
}

const (
	inv360   = float64(1.0 / 360.0)
	radToDeg = 180.0 / math.Pi
	degToRad = math.Pi / 180.0
)

// revolution will reduce the angle to within 0-360 degrees
func revolution(x float64) float64 {
	return x - 360.0*math.Floor(x*inv360)
}

// rev180 will reduce the angle to within 180-180 degrees
func rev180(x float64) float64 {
	return x - 360.0*math.Floor(x*inv360+0.5)
}

func sind(x float64) float64 {
	return math.Sin(x * degToRad)
}

func cosd(x float64) float64 {
	return math.Cos(x * degToRad)
}

func tand(x float64) float64 {
	return math.Tan(x * degToRad)
}

func atand(x float64) float64 {
	return radToDeg * math.Atan(x)
}

func asindh(x float64) float64 {
	return radToDeg * math.Asin(x)
}

func acosd(x float64) float64 {
	return radToDeg * math.Acos(x)
}

func atan2d(y, x float64) float64 {
	return radToDeg * math.Atan2(y, x)
}

func sunpos(d float64) (lon, r float64) {
	M := revolution(356.0470 + 0.9856002585*d)
	w := 282.9404 + 4.70935E-5*d
	e := 0.016709 - 1.151E-9*d

	E := M + e*radToDeg*sind(M)*(1.0+e*cosd(M))
	x := cosd(E) - e
	y := math.Sqrt(1.0-e*e) * sind(E)

	r = math.Sqrt(x*x + y*y)

	v := atan2d(y, x)
	lon = v + w

	for lon >= 360.0 {
		lon -= 360.0
	}

	return lon, r
}

func gmst0(d float64) float64 {
	return revolution((180.0 + 356.0470 + 282.9404) + (0.9856002585+4.70935E-5)*d)
}

func sunRADec(d float64) (ra, dec, r float64) {
	var lon float64
	lon, r = sunpos(d)

	x := r * cosd(lon)
	y := r * sind(lon)

	oblEcl := 23.4393 - 3.563E-7*d

	z := y * sind(oblEcl)
	y = y * cosd(oblEcl)

	ra = atan2d(y, x)
	dec = atan2d(z, math.Sqrt(x*x+y*y))

	return ra, dec, r
}

func sunRiseSet(year, month, day int, lon, lat float64, sunAlt float64, upperLimb bool) (float64, float64, SunriseStatus) {
	d := daysSince2000(year, month, day)

	sidTime := revolution(gmst0(d) + 180.0 + lon)

	sRA, sDec, sR := sunRADec(d)

	tSouth := 12.0 - rev180(sidTime-sRA)/15.0

	if upperLimb {
		sRadius := 0.2666 / sR
		sunAlt -= sRadius
	}

	cost := (sind(sunAlt) - sind(lat)*sind(sDec)) / (cosd(lat) * cosd(sDec))

	var s SunriseStatus
	var t float64

	switch {
	case cost >= 1.0:
		t = 0.0
		s = SunriseStatusBelowHorizon
	case cost <= -1.0:
		t = 12.0
		s = SunriseStatusAboveHorizon
	default:
		t = acosd(cost) / 15.0
		s = SunriseStatusOK
	}

	return tSouth - t, tSouth + t, s
}

func dayLen(year, month, day int, lon, lat float64, sunAlt float64, upperLimb bool) float64 {
	d := daysSince2000(year, month, day) + 0.5 - lon/360.0

	oblEcl := 23.4393 - 3.563E-7*d

	sLon, sR := sunpos(d)

	sinSDec := sind(oblEcl) * sind(sLon)
	cosSDec := math.Sqrt(1.0 - sinSDec*sinSDec)

	if upperLimb {
		sRadius := 0.2666 / sR
		sunAlt -= sRadius
	}

	cost := (sind(sunAlt) - sind(lat)*sinSDec) / (cosd(lat) * cosSDec)

	var t float64

	switch {
	case cost >= 1.0:
		t = 0.0
	case cost <= -1.0:
		t = 24.0
	default:
		t = (2.0 / 15.0) * acosd(cost)
	}

	return t
}
