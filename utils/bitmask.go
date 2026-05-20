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
