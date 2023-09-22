package gossip

import "ece428_mp2/pkg/network/code"

type HeartbeatHandler struct {
	membership *MembershipManager
}

func NewHeartbeatHandler(manager *MembershipManager) *HeartbeatHandler {
	return &HeartbeatHandler{
		membership: manager,
	}
}

func (h *HeartbeatHandler) Handle(req *code.HeartbeatRequest) error {
	return h.merge(req)
}

func (h *HeartbeatHandler) merge(req *code.HeartbeatRequest) error {
	h.membership.MergeMembershipList(req.MemberShipList)
	return nil
}
