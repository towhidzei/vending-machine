package algo

import "sort"

func MinimumNumberOfElementsWhoseSumIs(assets []uint, sum uint) []uint {
	if sum <= 0 || len(assets) == 0 {
		return []uint{}
	}
	sort.SliceStable(assets, func(i, j int) bool {
		return assets[i] < assets[j]
	})
	// last index
	li := len(assets) - 1
	// biggest item
	bi := assets[li]
	if bi == sum {
		return []uint{sum}
	} else if bi > sum {
		return MinimumNumberOfElementsWhoseSumIs(assets[:li], sum)
	} else {
		q := sum / bi
		result := make([]uint, 0)
		for ; q > 0; q-- {
			result = append(result, bi)
		}
		r := sum % bi
		if r == 0 {
			return result
		}
		return append(result, MinimumNumberOfElementsWhoseSumIs(assets[:li], r)...)
	}
}
