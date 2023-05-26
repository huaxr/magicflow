#include "textflag.h"

TEXT ·WrapF(SB), NOSPLIT, $0
    // call the enter func
    CALL ·enter(SB)
    // get func address
    MOVQ func+0(FP), DX
    MOVQ (DX), AX
    // call this func
    CALL AX

    CALL ·exit(SB)
    RET

TEXT ·wrapFN(SB), NOSPLIT, $0-40
    // no need enter cause doing this in l3
    // get func address
    MOVQ func+16(FP), DX
    MOVQ (DX), CX
    // call this func
    CALL CX
    NOP
    // record pc
    MOVQ pc+0(SP), AX
    MOVQ AX, ret+32(FP)          //   *     ^__^
    // get name address          //  \|/    (oo)\_______
    MOVQ name+0(FP), BX          //         (__)\       )\/\
    MOVQ BX, x_arg+0(SP)         //       *     ||----w |         ~~~~
    // +8 is the length of string//      \|/    ||     ||      *
    MOVQ name+8(FP), CX          //            ~--     --~    \|/
    MOVQ CX, x_arg+8(SP)         // ~~~~~ ~~~~~~ ~~~~~~~~~~ ~~~~
                                 //    ~~~~~~~    ~~~~~~~       ~~~~~~~~
    CALL ·exitMonitor(SB)
    NOP
    // reset pc
    MOVQ ret+32(FP), AX
    MOVQ AX, pc+0(SP)
    RET
    INT	$3

