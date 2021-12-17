package main

type DNSHeader struct {
	TransactionID byte
	Flags         byte
}
type DNSBase struct {
}
type DNSRecords struct {
}

type DNS struct {
	DNSHeader
	DNSBase
	DNSRecords
}
