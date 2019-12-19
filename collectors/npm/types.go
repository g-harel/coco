package npm

// Pkg is extracted package download data.
type pkg struct {
	Name   string
	Weekly int
	Total  int
}

// PkgHandler accepts and handles package download data.
type pkgHandler func(*pkg, error)

// OwnerResponse represents the response data for a request
// for owner data.
type ownerResponse struct {
	Packages struct {
		Total   int `json:"total"`
		Objects []struct {
			Name string `json:"name"`
		} `json:"objects"`
	} `json:"packages"`
	Pagination struct {
		PerPage int `json:"perPage"`
		Page    int `json:"page"`
	}
}

// PkgResponse represents the response data for a request
// for package data.
type pkgResponse struct {
	Package   string `json:"package"`
	Downloads []struct {
		Downloads int `json:"downloads"`
	} `json:"downloads"`
}

// Convert converts between the HTTP response and extracted
// package download data.
func convert(r *pkgResponse) *pkg {
	p := &pkg{
		Name:   r.Package,
		Weekly: 0,
		Total:  0,
	}
	if len(r.Downloads) > 0 {
		p.Weekly = r.Downloads[len(r.Downloads)-1].Downloads
		for i := 0; i < len(r.Downloads); i++ {
			p.Total += r.Downloads[i].Downloads
		}
	}
	return p
}
