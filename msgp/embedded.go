package msgp

// InterceptField intercept the next object bytes with field.
func (m *Reader) InterceptField(field []byte) (object []byte, err error) {
	b, err := m.next()
	if err != nil {
		return
	}
	object = make([]byte, 0, len(b)+32)
	object = AppendMapHeader(object, 1)
	object = AppendString(object, UnsafeString(field))
	object = append(object, b...)
	return object, nil
}

func (m *Reader) next() (b []byte, err error) {
	var (
		v uintptr // bytes
		o uintptr // objects
		p []byte
	)

	// we can use the faster
	// method if we have enough
	// buffered data
	if m.R.Buffered() >= 5 {
		p, err = m.R.Peek(5)
		if err != nil {
			return nil, err
		}
		v, o, err = getSize(p)
		if err != nil {
			return nil, err
		}
	} else {
		v, o, err = getNextSize(m.R)
		if err != nil {
			return nil, err
		}
	}

	// 'v' is always non-zero
	// if err == nil
	b, err = m.R.Next(int(v))
	if err != nil {
		return nil, err
	}
	var bb []byte
	// for maps and slices, skip elements
	for x := uintptr(0); x < o; x++ {
		bb, err = m.next()
		if err != nil {
			return nil, err
		}
		b = append(b, bb...)
	}
	return b, nil
}

// InterceptField intercept the next object bytes with field.
func InterceptField(field []byte, b []byte) (object, last []byte, err error) {
	last, err = Skip(b)
	if err != nil {
		return
	}
	object = make([]byte, 0, len(b)+32)
	object = AppendMapHeader(object, 1)
	object = AppendString(object, UnsafeString(field))
	object = append(object, b[:len(b)-len(last)]...)
	return object, last, nil
}
