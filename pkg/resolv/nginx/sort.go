package nginx

func SortByOrders(slice *[]Parser, orders ...Order) {
	for _, order := range orders {
		switch {
		case order >= 1000:
			radixSort(slice, order)
		case order < 1000:
			insertionSort(slice, order)
		default:
			continue
		}
	}
}

func radixSort(slice *[]Parser, order Order) {
	maxBit := 0
	for _, sorter := range *slice {
		if maxBit < sorter.BitLen(order) {
			maxBit = sorter.BitLen(order)
		}
	}

	for i := maxBit - 1; i >= 0; i-- {
		radixSortC(slice, order, i)
	}
}

// radixSortC, insertion sort function.
func radixSortC(slice *[]Parser, order Order, bit int) {
	n := len(*slice)
	if n <= 1 {
		return
	}
	cache := map[Parser]byte{}
	for i := 1; i < n; i++ {
		tmp := (*slice)[i]
		d, tmpOK := cache[tmp]
		if !tmpOK {
			d = tmp.BitSize(order, bit)
			cache[tmp] = d
		}
		j := i - 1
		for ; j >= 0; j-- {
			c, ok := cache[(*slice)[j]]
			if !ok {
				c = (*slice)[j].BitSize(order, bit)
				cache[(*slice)[j]] = c
			}

			if c > d {
				(*slice)[j+1] = (*slice)[j]
			} else {
				break
			}
		}
		(*slice)[j+1] = tmp
	}
}

func insertionSort(slice *[]Parser, order Order) {
	n := len(*slice)
	if n <= 1 {
		return
	}
	cache := map[Parser]int{}
	for i := 1; i < n; i++ {
		tmp := (*slice)[i]
		d, tmpOK := cache[tmp]
		if !tmpOK {
			d = tmp.Size(order)
			cache[tmp] = d
		}
		j := i - 1
		for ; j >= 0; j-- {
			c, ok := cache[(*slice)[j]]
			if !ok {
				c = (*slice)[j].Size(order)
				cache[(*slice)[j]] = c
			}

			if c > d {
				(*slice)[j+1] = (*slice)[j]
			} else {
				break
			}
		}
		(*slice)[j+1] = tmp
	}
}
