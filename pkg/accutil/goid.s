#include "textflag.h"
#include "./include/go_tls.h"

TEXT ·getgid(SB), NOSPLIT, $0-16
    get_tls(CX)
    MOVQ g(CX), AX
    MOVQ $152, BX
    LEAQ 0(AX)(BX*1), DX
    MOVQ (DX), AX
    MOVQ AX, ret+0(FP)
    RET

