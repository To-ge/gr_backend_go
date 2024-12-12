package entity

type Status int

const (
	Unknown      Status = iota // 不明
	Disconnected               // 地上局が切断した
)
