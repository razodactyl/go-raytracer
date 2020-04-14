package vec

import (
	"github.com/razodactyl/go-raytracer/util"
	"math"
)

type Camera struct {
	Origin Vector3D
	LowerLeftCorner Vector3D
	Horizontal Vector3D
	Vertical Vector3D
	U, V, W Vector3D
	LensRadius float64
}

func NewCamera(lookFrom Vector3D, lookAt Vector3D, vup Vector3D, vfov float64, aspect float64, aperture float64, focusDist float64) *Camera {
	origin := lookFrom
	lensRadius := aperture / 2.0

	theta := util.DegreesToRadians(vfov)
	halfHeight := math.Tan(theta/2)
	halfWidth := aspect * halfHeight

	w := lookFrom.Subtract(lookAt).Unit()
	u := vup.Cross(w).Unit()
	v := w.Cross(u)

	uMulHalfWidth := u.MultiplyScalar(halfWidth).MultiplyScalar(focusDist)
	vMulHalfHeight := v.MultiplyScalar(halfHeight).MultiplyScalar(focusDist)

	lowerLeftCorner := origin.Subtract(uMulHalfWidth).Subtract(vMulHalfHeight).Subtract(w.MultiplyScalar(focusDist))

	horizontal := uMulHalfWidth.MultiplyScalar(2)
	vertical := vMulHalfHeight.MultiplyScalar(2)

	return &Camera{
		LowerLeftCorner: lowerLeftCorner,
		Horizontal:      horizontal,
		Vertical:        vertical,
		Origin:          origin,
		U:               u,
		V:               v,
		W:               w,
		LensRadius:      lensRadius,
	}
}

func (c Camera) GetRay(s float64, t float64) Ray3D {
	rd := RandomInUnitDisk().MultiplyScalar(c.LensRadius)
	offset := c.U.MultiplyScalar(rd.X).Add(c.V.MultiplyScalar(rd.Y))

	sMulHorizontal := c.Horizontal.MultiplyScalar(s)
	tMulVertical := c.Vertical.MultiplyScalar(t)

	return NewRay3D(c.Origin.Add(offset), c.LowerLeftCorner.Add(sMulHorizontal).Add(tMulVertical).Subtract(c.Origin).Subtract(offset))
}
