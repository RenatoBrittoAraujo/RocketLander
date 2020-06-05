package sim

// DetectGroundCollision returns true if ground collision happend
func DetectGroundCollision(r *Rocket) (collision int) {
	points := r.BoundingBox()
	for _, v := range points {
		if v.Y <= 0 {
			collision++
		}
	}
	return
}
