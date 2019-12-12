package npm

type Package struct {
	Name   string
	Weekly int
	Total  int
}

type PackageHandler func(*Package, error)

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

type packageResponse struct {
	Package   string `json:"package"`
	Downloads []struct {
		Downloads int `json:"downloads"`
	} `json:"downloads"`
}

type packageResponseHandler func(*packageResponse, error)

func converterFunc(f PackageHandler) packageResponseHandler {
	return func(r *packageResponse, err error) {
		if err != nil {
			f(nil, err)
			return
		}
		p := &Package{
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
		f(p, nil)
	}
}
