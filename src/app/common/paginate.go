package common

import (
	"math"
)

const adjacent = 1

type Paginated struct {
	TotalItems int
	PerPage    int
	PageNumber int
	TotalPages int
	StartPoint int
	PageList   []int
	Flat       []int
	Left       []int
	Middle     []int
	Right      []int
	NextPage   int
	PrevPage   int
}

func GetPaginated(total int, pp int, pn int) *Paginated {

	total_pages := getTotalPages(total, pp)
	start_point := getStartPoint(pp, pn)
	page_list := getPageList(total_pages)
	next_page := getNextPage(pn, total_pages)
	prev_page := getPrevPage(pn)

	pure := total_pages / 3
	repu := (adjacent * 2) + 1

	var left []int = page_list[:2]
	var right []int = page_list[total_pages-2:]
	var mid []int
	var flat []int

	if pure >= repu {
		if pn <= (adjacent*2)+1 {
			left = page_list[:pn+adjacent*2]
		} else if pn >= total_pages-((adjacent*2)+1) {
			right = page_list[(pn-adjacent*2)-1:]
		} else {
			mid = page_list[(pn-adjacent)-1 : (pn + adjacent)]
		}
	} else {
		flat = page_list
	}

	return &Paginated{
		TotalItems: total,
		PerPage:    pp,
		PageNumber: pn,
		TotalPages: total_pages,
		StartPoint: start_point,
		PageList:   page_list,
		NextPage:   next_page,
		PrevPage:   prev_page,
		Left:       left,
		Middle:     mid,
		Right:      right,
		Flat:       flat,
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
