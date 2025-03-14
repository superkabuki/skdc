package skdc

// Pids holds collections of pids by type for skdc.Stream.
type Pids struct {
	PmtPids    []uint16
	PcrPids    []uint16
	Scte35Pids []uint16
	MaybePids  []uint16
}

func (pids *Pids) isPmtPid(pid uint16) bool {
	return IsIn(pids.PmtPids, pid)
}

func (pids *Pids) addPmtPid(pid uint16) {
	if !pids.isPmtPid(pid) {
		pids.PmtPids = append(pids.PmtPids, pid)
	}
}

func (pids *Pids) isPcrPid(pid uint16) bool {
	return IsIn(pids.PcrPids, pid)
}

func (pids *Pids) addPcrPid(pid uint16) {
	if !pids.isPcrPid(pid) {
		pids.PcrPids = append(pids.PcrPids, pid)
	}
}

func (pids *Pids) isScte35Pid(pid uint16) bool {
	return IsIn(pids.Scte35Pids, pid)
}

func (pids *Pids) addScte35Pid(pid uint16) {
	if !(pids.isScte35Pid(pid)) {
		pids.Scte35Pids = append(pids.Scte35Pids, pid)
	}
}
func (pids *Pids) delScte35Pid(pid uint16) {
	n := 0
	for _, val := range pids.Scte35Pids {
		if val != pid {
			pids.Scte35Pids[n] = val
			n++
		}
	}

	pids.Scte35Pids = pids.Scte35Pids[:n]
}

func (pids *Pids) isMaybePid(pid uint16) bool {
	return IsIn(pids.MaybePids, pid)
}

func (pids *Pids) addMaybePid(pid uint16) {
	if !(pids.isMaybePid(pid)) {
		pids.MaybePids = append(pids.MaybePids, pid)
	}
}
func (pids *Pids) delMaybePid(pid uint16) {
	n := 0
	for _, val := range pids.MaybePids {
		if val != pid {
			pids.MaybePids[n] = val
			n++
		}
	}

	pids.MaybePids = pids.MaybePids[:n]
}
