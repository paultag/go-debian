/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2025
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
 * THE SOFTWARE. }}} */

package dependency

var TUPLETABLE = string(`
eabi-uclibc-linux-arm		uclibc-linux-armel
base-uclibc-linux-<cpu>		uclibc-linux-<cpu>
eabihf-musl-linux-arm		musl-linux-armhf
base-musl-linux-<cpu>		musl-linux-<cpu>
eabihf-gnu-linux-arm		armhf
eabi-gnu-linux-arm		armel
abin32-gnu-linux-mips64r6el	mipsn32r6el
abin32-gnu-linux-mips64r6	mipsn32r6
abin32-gnu-linux-mips64el	mipsn32el
abin32-gnu-linux-mips64		mipsn32
abi64-gnu-linux-mips64r6el	mips64r6el
abi64-gnu-linux-mips64r6	mips64r6
abi64-gnu-linux-mips64el	mips64el
abi64-gnu-linux-mips64		mips64
spe-gnu-linux-powerpc		powerpcspe
x32-gnu-linux-amd64		x32
base-gnu-linux-<cpu>		<cpu>
base-gnu-kfreebsd-amd64		kfreebsd-amd64
base-gnu-kfreebsd-i386		kfreebsd-i386
base-gnu-kopensolaris-amd64	kopensolaris-amd64
base-gnu-kopensolaris-i386	kopensolaris-i386
base-gnu-hurd-amd64		hurd-amd64
base-gnu-hurd-i386		hurd-i386
base-bsd-dragonflybsd-amd64	dragonflybsd-amd64
base-bsd-freebsd-amd64		freebsd-amd64
base-bsd-freebsd-arm		freebsd-arm
base-bsd-freebsd-arm64		freebsd-arm64
base-bsd-freebsd-i386		freebsd-i386
base-bsd-freebsd-powerpc	freebsd-powerpc
base-bsd-freebsd-ppc64		freebsd-ppc64
base-bsd-freebsd-riscv		freebsd-riscv
base-bsd-openbsd-<cpu>		openbsd-<cpu>
base-bsd-netbsd-<cpu>		netbsd-<cpu>
base-bsd-darwin-amd64		darwin-amd64
base-bsd-darwin-arm		darwin-arm
base-bsd-darwin-arm64		darwin-arm64
base-bsd-darwin-i386		darwin-i386
base-bsd-darwin-powerpc		darwin-powerpc
base-bsd-darwin-ppc64		darwin-ppc64
base-sysv-aix-powerpc		aix-powerpc
base-sysv-aix-ppc64		aix-ppc64
base-sysv-solaris-amd64		solaris-amd64
base-sysv-solaris-i386		solaris-i386
base-sysv-solaris-sparc		solaris-sparc
base-sysv-solaris-sparc64	solaris-sparc64
base-tos-mint-m68k		mint-m68k
`)

// vim: foldmethod=marker
