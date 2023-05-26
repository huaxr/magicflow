#include "textflag.h"

TEXT ·LockSwap(SB),NOSPLIT,$0-17
	MOVQ	ptr+0(FP), BX
	MOVL	old+8(FP), AX
	MOVL	new+12(FP), CX
	LOCK
	CMPXCHGL	CX, 0(BX)
	SETEQ	ret+16(FP)
	RET

// int32 add, the return size = 20  $0-20
// int64 add, the return size = 24  $0-24
// if we not set, this works as same
TEXT ·LockAdd(SB), NOSPLIT, $0
	MOVQ	ptr+0(FP), BX
	MOVL	delta+8(FP), AX
	MOVL	AX, CX
	LOCK
	XADDL	AX, 0(BX)
	ADDL	CX, AX
	MOVL	AX, ret+16(FP)
	RET

