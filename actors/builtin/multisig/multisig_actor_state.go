package multisig

import (
	abi "github.com/filecoin-project/specs-actors/actors/abi"
)

type MultiSigActorState struct {
	// Linear unlock
	InitialBalance abi.TokenAmount
	StartEpoch     abi.ChainEpoch
	UnlockDuration abi.ChainEpoch

	AuthorizedParties     []abi.ActorID
	NumApprovalsThreshold int64
	NextTxnID             TxnID
	PendingTxns           MultiSigTransactionHAMT
	PendingApprovals      MultiSigApprovalSetHAMT
}

func (st *MultiSigActorState) AmountLocked(elapsedEpoch abi.ChainEpoch) abi.TokenAmount {
	if elapsedEpoch >= st.UnlockDuration {
		return abi.TokenAmount(0)
	}

	lockedProportion := (st.UnlockDuration - elapsedEpoch) / st.UnlockDuration
	return abi.TokenAmount(uint64(st.InitialBalance) * uint64(lockedProportion))
}

func (st *MultiSigActorState) isAuthorizedParty(party abi.ActorID) bool {
	for _, ap := range st.AuthorizedParties {
		if party == ap {
			return true
		}
	}
	return false
}

// return true if MultiSig maintains required locked balance after spending the amount
func (st *MultiSigActorState) _hasAvailable(currBalance abi.TokenAmount, amountToSpend abi.TokenAmount, currEpoch abi.ChainEpoch) bool {
	if amountToSpend < 0 || currBalance < amountToSpend {
		return false
	}

	if currBalance-amountToSpend < st.AmountLocked(currEpoch-st.StartEpoch) {
		return false
	}

	return true
}