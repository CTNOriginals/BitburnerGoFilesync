package utils

type TBitMask uint

func (this TBitMask) Has(flag TBitMask) bool {
	return this&flag != 0
}
func (this *TBitMask) Set(flag TBitMask) {
	*this = *this | flag
}
func (this *TBitMask) Toggle(flag TBitMask) {
	*this = *this ^ flag
}
func (this *TBitMask) Clear(flag TBitMask) {
	*this = *this &^ flag
}

// Preforms an AND operation on this and flag
// and sets the result to this
func (this *TBitMask) Compare(flag TBitMask) {
	*this = *this & flag
}

func (this *TBitMask) SetIf(flag TBitMask, condition bool) {
	if !condition {
		return
	}

	this.Set(flag)
}
