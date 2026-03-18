package source

type SourceType int

const (
	_ SourceType = iota
	DefSrc
	DefInstSrc
	DepSrc
	DepInstSrc
	HostSrc
	HostInstSrc
	HostDepSrc
	HostDepInstSrc
	LocSrc
	LocInstSrc
	LocDepSrc
	LocDepInstSrc
	EnvSrc
	VaultSrc
)

func (o SourceType) String() string {
	switch o {
	case DefSrc:
		return "default"
	case DefInstSrc:
		return "default/instance"
	case DepSrc:
		return "deployment"
	case DepInstSrc:
		return "deployment/instance"
	case HostSrc:
		return "host"
	case HostInstSrc:
		return "host/instance"
	case HostDepSrc:
		return "host/deployment"
	case HostDepInstSrc:
		return "host/deployment/instance"
	case LocSrc:
		return "local"
	case LocInstSrc:
		return "local/instance"
	case LocDepSrc:
		return "local/deployment"
	case LocDepInstSrc:
		return "local/deployment/instance"
	case EnvSrc:
		return "environment"
	case VaultSrc:
		return "vault"
	}

	return "unknown"
}
