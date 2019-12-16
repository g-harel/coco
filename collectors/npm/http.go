package npm

import (
	"fmt"
	"net/http"

	"github.com/g-harel/coco/internal/httpc"
)

func fetchOwner(owner string, page int) (*ownerResponse, error) {
	res := &ownerResponse{}
	_, err := httpc.Get(
		fmt.Sprintf("https://www.npmjs.com/~%v?page=%v", owner, page),
		http.Header{"x-spiferack": []string{"1"}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch owner %v page %v: %v", owner, page, err)
	}
	return res, nil
}

func fetchPackage(name string) (*pkgResponse, error) {
	res := &pkgResponse{}
	_, err := httpc.Get(
		fmt.Sprintf("https://www.npmjs.com/package/%v", name),
		http.Header{"x-spiferack": []string{"1"}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch owner package %v: %v", name, err)
	}
	return res, nil
}
