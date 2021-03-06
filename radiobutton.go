// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	. "github.com/lxn/go-winapi"
)

type RadioButtonGroup struct {
	buttons       []*RadioButton
	checkedButton *RadioButton
}

func (rbg *RadioButtonGroup) Buttons() []*RadioButton {
	buttons := make([]*RadioButton, len(rbg.buttons))
	copy(buttons, rbg.buttons)
	return buttons
}

func (rbg *RadioButtonGroup) CheckedButton() *RadioButton {
	return rbg.checkedButton
}

type radioButtonish interface {
	radioButton() *RadioButton
}

type RadioButton struct {
	Button
	group *RadioButtonGroup
	value interface{}
}

func NewRadioButton(parent Container) (*RadioButton, error) {
	rb := new(RadioButton)

	if count := parent.Children().Len(); count > 0 {
		if prevRB, ok := parent.Children().At(count - 1).(radioButtonish); ok {
			rb.group = prevRB.radioButton().group
		}
	}
	var groupBit uint32
	if rb.group == nil {
		groupBit = WS_GROUP
		rb.group = new(RadioButtonGroup)
	}

	if err := InitChildWidget(
		rb,
		parent,
		"BUTTON",
		groupBit|WS_TABSTOP|WS_VISIBLE|BS_AUTORADIOBUTTON,
		0); err != nil {
		return nil, err
	}

	rb.Button.init()

	rb.MustRegisterProperty("CheckedValue", NewProperty(
		func() interface{} {
			if rb.Checked() {
				return rb.value
			}

			return nil
		},
		func(v interface{}) error {
			checked := v == rb.value
			if checked {
				rb.group.checkedButton = rb
			}
			rb.SetChecked(checked)

			return nil
		},
		rb.CheckedChanged()))

	rb.group.buttons = append(rb.group.buttons, rb)

	return rb, nil
}

func (rb *RadioButton) radioButton() *RadioButton {
	return rb
}

func (*RadioButton) LayoutFlags() LayoutFlags {
	return 0
}

func (rb *RadioButton) MinSizeHint() Size {
	defaultSize := rb.dialogBaseUnitsToPixels(Size{50, 10})
	textSize := rb.calculateTextSizeImpl("n" + widgetText(rb.hWnd))

	// FIXME: Use GetThemePartSize instead of GetSystemMetrics?
	w := textSize.Width + int(GetSystemMetrics(SM_CXMENUCHECK))
	h := maxi(defaultSize.Height, textSize.Height)

	return Size{w, h}
}

func (rb *RadioButton) SizeHint() Size {
	return rb.MinSizeHint()
}

func (rb *RadioButton) Group() *RadioButtonGroup {
	return rb.group
}

func (rb *RadioButton) Value() interface{} {
	return rb.value
}

func (rb *RadioButton) SetValue(value interface{}) {
	rb.value = value
}

func (rb *RadioButton) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint32(wParam)) {
		case BN_CLICKED:
			prevChecked := rb.group.checkedButton
			rb.group.checkedButton = rb

			if prevChecked != rb {
				if prevChecked != nil {
					prevChecked.setChecked(false)
				}

				rb.setChecked(true)
			}
		}
	}

	return rb.Button.WndProc(hwnd, msg, wParam, lParam)
}
