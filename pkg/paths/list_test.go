/*
 * This file is part of PathsHelper library.
 *
 * Copyright 2018 Arduino AG (http://www.arduino.cc/)
 *
 * PathsHelper library is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 */

package paths

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListConstructors(t *testing.T) {
	list0 := NewPathList()
	require.Len(t, list0, 0)

	list1 := NewPathList("test")
	require.Len(t, list1, 1)
	require.Equal(t, "[test]", fmt.Sprintf("%s", list1))

	list3 := NewPathList("a", "b", "c")
	require.Len(t, list3, 3)
	require.Equal(t, "[a b c]", fmt.Sprintf("%s", list3))

	require.False(t, list3.Contains(New("d")))
	require.True(t, list3.Contains(New("a")))
	require.False(t, list3.Contains(New("d/../a")))

	require.False(t, list3.ContainsEquivalentTo(New("d")))
	require.True(t, list3.ContainsEquivalentTo(New("a")))
	require.True(t, list3.ContainsEquivalentTo(New("d/../a")))

	list4 := list3.Clone()
	require.Equal(t, "[a b c]", fmt.Sprintf("%s", list4))
	list4.AddIfMissing(New("d"))
	require.Equal(t, "[a b c d]", fmt.Sprintf("%s", list4))
	list4.AddIfMissing(New("b"))
	require.Equal(t, "[a b c d]", fmt.Sprintf("%s", list4))
	list4.AddAllMissing(NewPathList("a", "e", "i", "o", "u"))
	require.Equal(t, "[a b c d e i o u]", fmt.Sprintf("%s", list4))
}

func TestListSorting(t *testing.T) {
	list := NewPathList(
		"pointless",
		"spare",
		"carve",
		"unwieldy",
		"empty",
		"bow",
		"tub",
		"grease",
		"error",
		"energetic",
		"depend",
		"property")
	require.Equal(t, "[pointless spare carve unwieldy empty bow tub grease error energetic depend property]", fmt.Sprintf("%s", list))
	list.Sort()
	require.Equal(t, "[bow carve depend empty energetic error grease pointless property spare tub unwieldy]", fmt.Sprintf("%s", list))
}

func TestListFilters(t *testing.T) {
	list := NewPathList(
		"aaaa",
		"bbbb",
		"cccc",
		"dddd",
		"eeff",
		"aaaa/bbbb",
		"eeee/ffff",
		"gggg/hhhh",
	)

	l1 := list.Clone()
	l1.FilterPrefix("a")
	require.Equal(t, "[aaaa]", fmt.Sprintf("%s", l1))

	l2 := list.Clone()
	l2.FilterPrefix("b")
	require.Equal(t, "[bbbb aaaa/bbbb]", fmt.Sprintf("%s", l2))

	l3 := list.Clone()
	l3.FilterOutPrefix("b")
	require.Equal(t, "[aaaa cccc dddd eeff eeee/ffff gggg/hhhh]", fmt.Sprintf("%s", l3))

	l4 := list.Clone()
	l4.FilterPrefix("a", "b")
	require.Equal(t, "[aaaa bbbb aaaa/bbbb]", fmt.Sprintf("%s", l4))

	l5 := list.Clone()
	l5.FilterPrefix("test")
	require.Equal(t, "[]", fmt.Sprintf("%s", l5))

	l6 := list.Clone()
	l6.FilterOutPrefix("b", "c", "h")
	require.Equal(t, "[aaaa dddd eeff eeee/ffff]", fmt.Sprintf("%s", l6))

	l7 := list.Clone()
	l7.FilterSuffix("a")
	require.Equal(t, "[aaaa]", fmt.Sprintf("%s", l7))

	l8 := list.Clone()
	l8.FilterSuffix("a", "h")
	require.Equal(t, "[aaaa gggg/hhhh]", fmt.Sprintf("%s", l8))

	l9 := list.Clone()
	l9.FilterSuffix("test")
	require.Equal(t, "[]", fmt.Sprintf("%s", l9))

	l10 := list.Clone()
	l10.FilterOutSuffix("a")
	require.Equal(t, "[bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]", fmt.Sprintf("%s", l10))

	l11 := list.Clone()
	l11.FilterOutSuffix("a", "h")
	require.Equal(t, "[bbbb cccc dddd eeff aaaa/bbbb eeee/ffff]", fmt.Sprintf("%s", l11))

	l12 := list.Clone()
	l12.FilterOutSuffix("test")
	require.Equal(t, "[aaaa bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]", fmt.Sprintf("%s", l12))

	l13 := list.Clone()
	l13.FilterOutSuffix()
	require.Equal(t, "[aaaa bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]", fmt.Sprintf("%s", l13))

	l14 := list.Clone()
	l14.FilterSuffix()
	require.Equal(t, "[]", fmt.Sprintf("%s", l14))

	l15 := list.Clone()
	l15.FilterOutPrefix()
	require.Equal(t, "[aaaa bbbb cccc dddd eeff aaaa/bbbb eeee/ffff gggg/hhhh]", fmt.Sprintf("%s", l15))

	l16 := list.Clone()
	l16.FilterPrefix()
	require.Equal(t, "[]", fmt.Sprintf("%s", l16))

	l17 := list.Clone()
	l17.Filter(func(p *Path) bool {
		return p.Base() == "bbbb"
	})
	require.Equal(t, "[bbbb aaaa/bbbb]", fmt.Sprintf("%s", l17))
}
