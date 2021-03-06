package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type StructTypeEditor struct {
	typ    *types.StructType
	fields []*StructFieldTypeEditor
}

func NewStructTypeEditor(typ *types.StructType) *StructTypeEditor {
	ed := &StructTypeEditor{
		typ:    typ,
		fields: make([]*StructFieldTypeEditor, len(typ.Fields)),
	}
	for i, f := range typ.Fields {
		ed.fields[i] = NewStructFieldTypeEditor(f)
	}
	return ed
}

func (s *StructTypeEditor) Type() types.Type { return s.typ }

func (s *StructTypeEditor) Layout(gtx C) D {
	maxFieldNameWidth := 0
	for _, f := range s.fields {
		if x := f.LayoutName(gtx); x > maxFieldNameWidth {
			maxFieldNameWidth = x
		}
	}
	fields := make([]layout.FlexChild, len(s.typ.Fields))
	for i, f := range s.fields {
		f := f
		fields[i] = layout.Rigid(func(gtx C) D {
			return f.Layout(gtx, maxFieldNameWidth)
		})
	}

	fieldsRec := Record(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, fields...)
	})

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "struct ").Layout),
		layout.Rigid(func(gtx C) D {
			width := gtx.Px(unit.Dp(16))
			height := fieldsRec.Dims.Size.Y + gtx.Px(unit.Dp(8))
			w := float32(width)
			h2 := float32(height) / 2
			path := clip.Path{}
			path.Begin(gtx.Ops)
			path.Move(f32.Pt(w, 0))
			path.Cube(f32.Pt(-w, 0), f32.Pt(0, h2), f32.Pt(-w, h2))
			path.Cube(f32.Pt(w, 0), f32.Pt(0, h2), f32.Pt(w, h2))
			paint.FillShape(gtx.Ops, theme.Fg, clip.Stroke{
				Path:  path.End(),
				Width: float32(gtx.Px(unit.Dp(1))),
			}.Op())
			return D{Size: image.Pt(width, height)}
		}),
		layout.Rigid(fieldsRec.Layout),
	)
}

type StructFieldTypeEditor struct {
	typ   *types.StructFieldType
	named widget.Editor
	typed *TypeEditor

	nameRec Recording
}

func NewStructFieldTypeEditor(typ *types.StructFieldType) *StructFieldTypeEditor {
	f := &StructFieldTypeEditor{
		typ: typ,
		named: widget.Editor{
			Alignment:  text.End,
			SingleLine: true,
		},
		typed: NewTypeEditor(&typ.Type),
	}
	f.named.SetText(typ.Name)
	return f
}

func (f *StructFieldTypeEditor) LayoutName(gtx C) int {
	f.nameRec = Record(gtx, material.Editor(theme, &f.named, "").Layout)
	return f.nameRec.Dims.Size.X
}

func (f *StructFieldTypeEditor) Layout(gtx C, nameWidth int) D {
	for _, e := range f.named.Events() {
		switch e := e.(type) {
		case widget.SubmitEvent:
			f.typ.Name = e.Text
		}
	}

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(layout.Spacer{Width: unit.Px(float32(nameWidth - f.nameRec.Dims.Size.X))}.Layout),
		layout.Rigid(f.nameRec.Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
		layout.Rigid(f.typed.Layout),
	)
}
