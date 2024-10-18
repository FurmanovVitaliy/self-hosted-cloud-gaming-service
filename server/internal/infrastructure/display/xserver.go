package display

type XServer struct {
	Card         string `json:"card"`
	DevicePath   string `json:"device_path"`
	ScreenNumber string `json:"screen"`
	Port         string `json:"port_name"`
	Connector    int    `json:"connector_id"`
	CRTS         int    `json:"crtc_id"`
	Plane        int    `json:"plane_id"`
	Used         bool
}

// Close implements srm.Display.
func (x *XServer) Close() {
	x.Used = false
}

// GetDisplayNumper implements srm.Display.
func (x *XServer) GetDisplayNumper() string {
	return x.ScreenNumber
}

// GetPlaneID implements srm.Display.
func (x *XServer) GetPlaneID() int {
	return x.Plane
}
