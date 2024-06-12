package eventwatcher

type SID_NAME_USE uint32

const (
	// https://learn.microsoft.com/zh-cn/windows/win32/api/winnt/ne-winnt-sid_name_use
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
