// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/lal
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package hls

import (
	"github.com/q191201771/naza/pkg/nazabits"
	"github.com/q191201771/naza/pkg/nazalog"
)

// Program Map Table
// <iso13818-1.pdf> <2.4.4.8> <page 64/174>
// table_id                 [8b]  *
// section_syntax_indicator [1b]
// 0                        [1b]
// reserved                 [2b]
// section_length           [12b] **
// program_number           [16b] **
// reserved                 [2b]
// version_number           [5b]
// current_next_indicator   [1b]  *
// section_number           [8b]  *
// last_section_number      [8b]  *
// reserved                 [3b]
// PCR_PID                  [13b] **
// reserved                 [4b]
// program_info_length      [12b] **
// -----loop-----
// stream_type              [8b]  *
// reserved                 [3b]
// elementary_PID           [13b] **
// reserved                 [4b]
// ES_info_length_length    [12b] **
// --------------
// CRC32                    [32b] ****
type PMT struct {
	tid   uint8
	ssi   uint8
	sl    uint16
	pn    uint16
	vn    uint8
	cni   uint8
	sn    uint8
	lsn   uint8
	pp    uint16
	pil   uint16
	ppes  []PMTProgramElement
	crc32 uint32
}

type PMTProgramElement struct {
	st   uint8
	epid uint16
	esil uint16
}

func ParsePMT(b []byte) (pmt PMT) {
	br := nazabits.NewBitReader(b)
	pmt.tid = br.ReadBits8(8)
	pmt.ssi = br.ReadBits8(1)
	br.ReadBits8(3)
	pmt.sl = br.ReadBits16(12)
	len := pmt.sl - 13
	pmt.pn = br.ReadBits16(16)
	br.ReadBits8(2)
	pmt.vn = br.ReadBits8(5)
	pmt.cni = br.ReadBits8(1)
	pmt.sn = br.ReadBits8(8)
	pmt.lsn = br.ReadBits8(8)
	br.ReadBits8(3)
	pmt.pp = br.ReadBits16(13)
	br.ReadBits8(4)
	pmt.pil = br.ReadBits16(12)
	if pmt.pil != 0 {
		nazalog.Warn(pmt.pil)
		br.ReadBytes(uint(pmt.pil))
	}

	for i := uint16(0); i < len; i += 5 {
		var ppe PMTProgramElement
		ppe.st = br.ReadBits8(8)
		br.ReadBits8(3)
		ppe.epid = br.ReadBits16(13)
		br.ReadBits8(4)
		ppe.esil = br.ReadBits16(12)
		if ppe.esil != 0 {
			nazalog.Warn(ppe.esil)
			br.ReadBits32(uint(ppe.esil))
		}
		pmt.ppes = append(pmt.ppes, ppe)
	}

	return
}

func (pmt *PMT) searchPID(pid uint16) *PMTProgramElement {
	for _, ppe := range pmt.ppes {
		if ppe.epid == pid {
			return &ppe
		}
	}
	return nil
}