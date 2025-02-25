// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"github.com/xlplbo/win"
)

type PushButton struct {
	Button
}

func NewPushButton(parent Container) (*PushButton, error) {
	pb := new(PushButton)

	if err := InitWidget(
		pb,
		parent,
		"BUTTON",
		win.WS_TABSTOP|win.WS_VISIBLE|win.BS_PUSHBUTTON,
		0); err != nil {
		return nil, err
	}

	pb.Button.init()

	pb.GraphicsEffects().Add(InteractionEffect)
	pb.GraphicsEffects().Add(FocusEffect)

	return pb, nil
}

func (pb *PushButton) ImageAboveText() bool {
	return pb.hasStyleBits(win.BS_TOP)
}

func (pb *PushButton) SetImageAboveText(value bool) error {
	if err := pb.ensureStyleBits(win.BS_TOP, value); err != nil {
		return err
	}

	// We need to set the image again, or Windows will fail to calculate the
	// button control size correctly.
	return pb.SetImage(pb.image)
}

func (pb *PushButton) ensureProperDialogDefaultButton(hwndFocus win.HWND) {
	widget := windowFromHandle(hwndFocus)
	if widget == nil {
		return
	}

	if _, ok := widget.(*PushButton); ok {
		return
	}

	form := ancestor(pb)
	if form == nil {
		return
	}

	dlg, ok := form.(dialogish)
	if !ok {
		return
	}

	defBtn := dlg.DefaultButton()
	if defBtn == nil {
		return
	}

	if err := defBtn.setAndClearStyleBits(win.BS_DEFPUSHBUTTON, win.BS_PUSHBUTTON); err != nil {
		return
	}

	if err := defBtn.Invalidate(); err != nil {
		return
	}
}

func (pb *PushButton) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_GETDLGCODE:
		hwndFocus := win.GetFocus()
		if hwndFocus == pb.hWnd {
			form := ancestor(pb)
			if form == nil {
				break
			}

			dlg, ok := form.(dialogish)
			if !ok {
				break
			}

			defBtn := dlg.DefaultButton()
			if defBtn == pb {
				pb.setAndClearStyleBits(win.BS_DEFPUSHBUTTON, win.BS_PUSHBUTTON)
				return win.DLGC_BUTTON | win.DLGC_DEFPUSHBUTTON
			}

			break
		}

		pb.ensureProperDialogDefaultButton(hwndFocus)

	case win.WM_KILLFOCUS:
		pb.ensureProperDialogDefaultButton(win.HWND(wParam))
	}

	return pb.Button.WndProc(hwnd, msg, wParam, lParam)
}

func (pb *PushButton) CreateLayoutItem(ctx *LayoutContext) LayoutItem {
	return &pushButtonLayoutItem{
		buttonLayoutItem: buttonLayoutItem{
			idealSize: pb.idealSize(),
		},
	}
}

type pushButtonLayoutItem struct {
	buttonLayoutItem
}

func (*pushButtonLayoutItem) LayoutFlags() LayoutFlags {
	return GrowableHorz
}
