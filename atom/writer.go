
package atom

import (
	"io"
	"log"
)

func WriteBytes(w io.Writer, b []byte) (err error) {
	_, err = w.Write(b)
	return
}

func WriteUInt(w io.Writer, val uint, n int) (err error) {
	var b [8]byte
	for i := n-1; i >= 0; i-- {
		b[i] = byte(val)
		val >>= 8
	}
	return WriteBytes(w, b[0:n])
}

func WriteInt(w io.Writer, val int, n int) (err error) {
	return WriteUInt(w, uint(val), n)
}

func WriteFixed(w io.Writer, val Fixed, n int) (err error) {
	return WriteUInt(w, uint(val), n)
}

func WriteTimeStamp(w io.Writer, ts TimeStamp, n int) (err error) {
	return WriteUInt(w, uint(ts), n)
}

func WriteString(w io.Writer, val string, n int) (err error) {
	wb := make([]byte, n)
	sb := []byte(val)
	copy(wb, sb)
	return WriteBytes(w, wb)
}

func WriteDummy(w io.Writer, n int) (err error) {
	return WriteBytes(w, make([]byte, n))
}

type Writer struct {
	io.WriteSeeker
	sizePos int64
}

func (self *Writer) WriteEmptyInt(n int) (pos int64, err error) {
	if pos, err = self.Seek(0, 1); err != nil {
		return
	}
	if err = WriteInt(self, 0, n); err != nil {
		return
	}
	return
}

func (self *Writer) RefillInt(pos int64, val int, n int) (err error) {
	var curPos int64
	if curPos, err = self.Seek(0, 1); err != nil {
		return
	}
	if _, err = self.Seek(pos, 0); err != nil {
		return
	}
	if err = WriteInt(self, val, n); err != nil {
		return
	}
	if _, err = self.Seek(curPos, 0); err != nil {
		return
	}
	return
}

func (self *Writer) Close() (err error) {
	var curPos int64
	if curPos, err = self.Seek(0, 1); err != nil {
		return
	}
	if err = self.RefillInt(self.sizePos, int(curPos - self.sizePos), 4); err != nil {
		return
	}
	if false {
		log.Println("writeback", self.sizePos, curPos, curPos-self.sizePos)
	}
	return
}

func WriteAtomHeader(w io.WriteSeeker, cc4 string) (res *Writer, err error) {
	self := &Writer{WriteSeeker: w}

	if self.sizePos, err = self.WriteEmptyInt(4); err != nil {
		return
	}
	if err = WriteString(self, cc4, 4); err != nil {
		return
	}

	res = self
	return
}

