#include "textflag.h"

TEXT ·containsNum(SB), NOSPLIT, $0
    MOVQ $0, SI
    MOVQ sl+0(FP), BX  // &elem[0] address of the first element
    MOVQ sl+8(FP), CX  // len(elem)
    //MOVQ sl+16(FP), CX // cap(elem)
    MOVQ num+24(FP), DX
    INCQ CX
loop:
    DECQ CX
    JZ false
    CMPQ (BX), DX
    JE true
    ADDQ $8, BX
    JMP loop
true:
    MOVQ $1, ret+32(FP)
    RET
false:
    MOVQ $0, ret+32(FP)
    RET

TEXT ·containsStr(SB), NOSPLIT, $0
    MOVQ sl+0(FP), BX  // &elem[0] address of the first element
    MOVQ sl+8(FP), CX  // len(elem)
    //MOVQ sl+16(FP), CX // cap(elem)
    MOVQ num+24(FP), DX // second params
    INCQ CX
loop:
    DECQ CX
    JZ false
    CMPQ (BX), DX
    JE true
    ADDQ $16, BX
    JMP loop
true:
    MOVQ $1, ret+40(FP)
    RET
false:
    MOVQ $0, ret+40(FP)
    RET
