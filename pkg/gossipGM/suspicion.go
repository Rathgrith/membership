package gossipGM

import (
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network/code"
	"strings"
	"time"
)

func (m *MembershipManager) MarkMembersSuspectedIfNotUpdated(TSuspicion time.Duration, TConfirm time.Duration) {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()

	currentTime := time.Now()

	for k, v := range m.membershipList {
		// if k is current hostname, skip
		if strings.HasPrefix(k, m.selfHostName) {
			continue
		}
		timeElapsed := currentTime.Sub(v.LocalTime)
		if timeElapsed > TSuspicion && v.StatusCode == code.Alive { // If member is alive and time elapsed exceeds TSuspicion
			v.StatusCode = code.Suspected // Mark as suspected
			v.LocalTime = currentTime
			m.membershipList[k] = v
			logutil.Logger.Infof("mark %v as suspected", k)
			go m.ReportSuspectedMember(k, v)
			go m.ReadyReportConfirm(k, TConfirm)
		}

	}
}

func (m *MembershipManager) ReportSuspectedMember(memberID string, memberInfo *code.MemberInfo) {
	suspectRequest := code.SuspensionRequest{
		TargetID: memberID,
		InfoType: code.Suspect,
	}

	m.mu.Lock()
	suspectRequest.IncarnationNumber = m.IncarnationNumberTrack[memberID]
	m.forwardRequestBuf = append(m.forwardRequestBuf, &suspectRequest)
	m.mu.Unlock()
}

func (m *MembershipManager) ReportConfirmFailedMember(memberID string) {
	confirmFailedRequest := code.SuspensionRequest{
		TargetID: memberID,
		InfoType: code.ConfirmFailed,
	}

	m.mu.Lock()
	confirmFailedRequest.IncarnationNumber = m.IncarnationNumberTrack[memberID]
	m.forwardRequestBuf = append(m.forwardRequestBuf, &confirmFailedRequest)
	m.mu.Unlock()
}

func (m *MembershipManager) ReportSelfAlive() {
	aliveRequest := code.SuspensionRequest{
		TargetID: m.selfID,
		InfoType: code.InformAlive,
	}

	m.mu.Lock()
	m.IncarnationNumberTrack[m.selfID] += 1
	aliveRequest.IncarnationNumber = m.IncarnationNumberTrack[m.selfID]
	m.forwardRequestBuf = append(m.forwardRequestBuf, &aliveRequest)
	m.mu.Unlock()
}

func (m *MembershipManager) HandleSuspicionRequest(req *code.SuspensionRequest) {
	if req.InfoType == code.Suspect {
		// override rule
		if m.IncarnationNumberTrack[req.TargetID] >= req.IncarnationNumber {
			return
		}

		if req.TargetID == m.selfID {
			m.ReportSelfAlive()
			return
		}

		// mark suspected
		m.listMutex.Lock()
		m.membershipList[req.TargetID].StatusCode = code.Suspected
		m.listMutex.Unlock()
		logutil.Logger.Infof("mark %v as suspected", req.TargetID)
	}

	if req.InfoType == code.InformAlive {
		if m.IncarnationNumberTrack[req.TargetID] >= req.IncarnationNumber {
			return
		}

		// mark alive
		m.listMutex.Lock()
		m.membershipList[req.TargetID].StatusCode = code.Alive
		m.listMutex.Unlock()
		logutil.Logger.Infof("mark %v as alive", req.TargetID)
	}

	if req.InfoType == code.ConfirmFailed {
		if m.IncarnationNumberTrack[req.TargetID] >= req.IncarnationNumber {
			return
		}

		// mark failed
		m.listMutex.Lock()
		delete(m.membershipList, req.TargetID)
		m.listMutex.Unlock()
		logutil.Logger.Infof("mark %v as failed", req.TargetID)
	}

	m.IncarnationNumberTrack[req.TargetID] = req.IncarnationNumber
	m.mu.Lock()
	m.forwardRequestBuf = append(m.forwardRequestBuf, req)
	m.mu.Unlock()
}

func (m *MembershipManager) ReadyReportConfirm(targetID string, TConfirm time.Duration) {
	timer := time.NewTimer(TConfirm)
	<-timer.C

	if m.membershipList[targetID] == nil || m.membershipList[targetID].StatusCode != code.Suspected {
		return
	}

	m.listMutex.Lock()
	delete(m.membershipList, targetID)
	m.listMutex.Unlock()
	logutil.Logger.Infof("mark %v as failed", targetID)
	m.ReportConfirmFailedMember(targetID)
}

func (m *MembershipManager) GetAllForwardSuspicionRequest() []*code.SuspensionRequest {
	ret := make([]*code.SuspensionRequest, len(m.forwardRequestBuf))
	m.mu.Lock()
	copy(ret, m.forwardRequestBuf)
	m.forwardRequestBuf = make([]*code.SuspensionRequest, 0)
	m.mu.Unlock()
	return ret
}
