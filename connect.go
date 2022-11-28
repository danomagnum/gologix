package main

func (plc *PLC) Connect() error {
	if plc.Size == 0 {
		plc.Size = 508
	}

	if ioi_cache == nil {
		ioi_cache = make(map[string]*IOI)
	}
	return plc.connect(plc.IPAddress)
}
