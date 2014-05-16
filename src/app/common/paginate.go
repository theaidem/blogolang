package common

import (
	"math"
)

type Paginated struct {
	TotalItems int
	PerPage    int
	PageNumber int
	TotalPages int
	StartPoint int
	PageList   []int
	NextPage   int
	PrevPage   int
}

func GetPaginated(total int, pp int, pn int) *Paginated {

	total_pages := getTotalPages(total, pp)
	start_point := getStartPoint(pp, pn)
	page_list := getPageList(total_pages)
	next_page := getNextPage(pn, total_pages)
	prev_page := getPrevPage(pn)

	return &Paginated{
		TotalItems: total,
		PerPage:    pp,
		PageNumber: pn,
		TotalPages: total_pages,
		StartPoint: start_point,
		PageList:   page_list,
		NextPage:   next_page,
		PrevPage:   prev_page,
	}
}

func getTotalPages(total int, pp int) int {
	return int(math.Ceil(float64(total) / float64(pp)))
}

func getStartPoint(pp int, pn int) int {
	return (pn - 1) * pp
}

func getPageList(tt int) []int {
	var page_list []int
	for i := 1; i <= tt; i++ {
		page_list = append(page_list, i)
	}
	return page_list
}

func getNextPage(pn int, total_pages int) int {
	//var next_page int
	if next_page := pn + 1; next_page > total_pages {
		return 0
	} else {
		return next_page
	}
}

func getPrevPage(pn int) int {
	return pn - 1
}
