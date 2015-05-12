/* Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. */

package dependency

/*
 */
type DepedencyIterator struct {
	dependency    Depedency
	relationIndex int
	arch          Arch
	next          *Possibility
}

func (iter *DepedencyIterator) ToSlice() (ret []*Possibility) {
	for {
		el := iter.Next()
		if el == nil {
			return
		}
		ret = append(ret, el)
	}
}

func (iter *DepedencyIterator) setNext() {
	iter.next = iter.peekNext()
}

/*
 */
func (iter *DepedencyIterator) peekNext() *Possibility {
	if len(iter.dependency.Relations) <= iter.relationIndex {
		return nil
	}

	relation := iter.dependency.Relations[iter.relationIndex]

	for _, possi := range relation.Possibilities {
		if len(possi.Arches.Arches) == 0 {
			return possi
		}
		/* OK, we need to do an arch constraint check */
		if possi.Arches.Matches(&iter.arch) {
			return possi
		}
	}

	return nil
}

func (iter *DepedencyIterator) HasNext() bool {
	return iter.next != nil
}

func (iter *DepedencyIterator) Next() *Possibility {
	el := iter.next
	if el != nil {
		iter.relationIndex += 1
	}
	iter.setNext()
	return el
}

/*
 *
 */
func (dep *Depedency) IterPossibilities(arch Arch) *DepedencyIterator {
	iter := DepedencyIterator{
		dependency:    *dep,
		relationIndex: 0,
		arch:          arch,
	}
	iter.setNext()
	return &iter
}
