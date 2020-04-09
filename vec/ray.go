package vec

type Ray3D struct {
	Origin Vector3D
	Direction Vector3D
}

func NewRay3D(origin, direction Vector3D) Ray3D {
	return Ray3D{origin, direction}
}

func (r Ray3D) At(t float64) Vector3D {
	return r.Origin.Add(r.Direction.MultiplyScalar(t))
}

func (r Ray3D) Color() Vector3D {
	unit_direction := r.Direction.Unit()
	t := 0.5 * (unit_direction.Y + 1.0)
	return Unit().MultiplyScalar((1.0 - t)).Add(NewVector3D(0.5, 0.7, 1.0).MultiplyScalar(t))
}
