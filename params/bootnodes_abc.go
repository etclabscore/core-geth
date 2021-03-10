package params

// ABCBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the ABC network.
// They will be the first point of contact for nodes coming online
// to find peers.
var ABCBootnodes = []string{
	"enode://3e12c4c633157ae52e7e05c168f4b1aa91685a36ba33a0901aa8a83cfeb84c3633226e3dd2eaf59bfc83492139e1d68918bf5b60ba93e2deaedb4e6a2ded5d32@42.152.120.98:30303",
}

// Once ABC network has DNS discovery set up,
// this value can be configured.
// var ABCDNSNetwork = "enrtree://AJE62Q4DUX4QMMXEHCSSCSC65TDHZYSMONSD64P3WULVLSF6MRQ3K@example.network"
