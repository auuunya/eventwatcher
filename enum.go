package eventwatcher

type SID_NAME_USE uint32

const (
	SidTypeUser SID_NAME_USE = iota + 1
	SidTypeGroup
	SidTypeDomain
	SidTypeAlias
	SidTypeWellKnownGroup
	SidTypeDeletedAccount
	SidTypeInvalid
	SidTypeUnknown
	SidTypeComputer
	SidTypeLabel
	SidTypeLogonSession
)
